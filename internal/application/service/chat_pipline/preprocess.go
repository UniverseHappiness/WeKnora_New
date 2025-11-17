package chatpipline

import (
    "bytes"
    "context"
    "html/template"
    "regexp"
    "strings"
    "unicode"

    "github.com/Tencent/WeKnora/internal/config"
    "github.com/Tencent/WeKnora/internal/logger"
    "github.com/Tencent/WeKnora/internal/models/chat"
    "github.com/Tencent/WeKnora/internal/types"
    "github.com/Tencent/WeKnora/internal/types/interfaces"
    "github.com/yanyiwu/gojieba"
)

// PluginPreprocess Query preprocessing plugin
type PluginPreprocess struct {
    config    *config.Config
    modelService interfaces.ModelService
    jieba     *gojieba.Jieba
    stopwords map[string]struct{}
}

// Regular expressions for text cleaning
var (
	multiSpaceRegex = regexp.MustCompile(`\s+`)                                 // Multiple spaces
	urlRegex        = regexp.MustCompile(`https?://\S+`)                        // URLs
	emailRegex      = regexp.MustCompile(`\b[\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,}\b`) // Email addresses
	punctRegex      = regexp.MustCompile(`[^\p{L}\p{N}\s]`)                     // Punctuation marks
)

// NewPluginPreprocess Creates a new query preprocessing plugin
func NewPluginPreprocess(
    eventManager *EventManager,
    modelService interfaces.ModelService,
    config *config.Config,
    cleaner interfaces.ResourceCleaner,
) *PluginPreprocess {
	// Use default dictionary for Jieba tokenizer
	jieba := gojieba.NewJieba()

	// Load stopwords from built-in stopword library
	stopwords := loadStopwords()

    res := &PluginPreprocess{
        config:       config,
        modelService: modelService,
        jieba:        jieba,
        stopwords:    stopwords,
    }

	// Register resource cleanup function
	if cleaner != nil {
		cleaner.RegisterWithName("JiebaPreprocessor", func() error {
			res.Close()
			return nil
		})
	}

	eventManager.Register(res)
	return res
}

// Load stopwords
func loadStopwords() map[string]struct{} {
	// Directly use some common stopwords built into Jieba
	commonStopwords := []string{
		"的", "了", "和", "是", "在", "我", "你", "他", "她", "它",
		"这", "那", "什么", "怎么", "如何", "为什么", "哪里", "什么时候",
		"the", "is", "are", "am", "I", "you", "he", "she", "it", "this",
		"that", "what", "how", "a", "an", "and", "or", "but", "if", "of",
		"to", "in", "on", "at", "by", "for", "with", "about", "from",
		"有", "无", "好", "来", "去", "说", "看", "想", "会", "可以",
		"吗", "呢", "啊", "吧", "的话", "就是", "只是", "因为", "所以",
	}

	result := make(map[string]struct{}, len(commonStopwords))
	for _, word := range commonStopwords {
		result[word] = struct{}{}
	}
	return result
}

// ActivationEvents Register activation events
func (p *PluginPreprocess) ActivationEvents() []types.EventType {
	return []types.EventType{types.PREPROCESS_QUERY}
}

// OnEvent Process events
func (p *PluginPreprocess) OnEvent(ctx context.Context, eventType types.EventType, chatManage *types.ChatManage, next func() *PluginError) *PluginError {
    if chatManage.RewriteQuery == "" {
        return next()
    }
    logger.GetLogger(ctx).Info("开始查询预处理")
    logger.GetLogger(ctx).Infof("Starting query preprocessing, original query: %s", chatManage.RewriteQuery)

    cleaned := p.cleanText(chatManage.RewriteQuery)

    useLLM := p.config != nil && p.config.Conversation != nil &&
        p.config.Conversation.KeywordsExtractionPrompt != "" &&
        p.config.Conversation.KeywordsExtractionPromptUser != ""

    if useLLM && p.modelService != nil {
        systemTmpl, err := template.New("keywordsSystem").Parse(p.config.Conversation.KeywordsExtractionPrompt)
        if err == nil {
            userTmpl, err2 := template.New("keywordsUser").Parse(p.config.Conversation.KeywordsExtractionPromptUser)
            if err2 == nil {
                var systemContent, userContent bytes.Buffer
                _ = systemTmpl.Execute(&systemContent, map[string]interface{}{
                    "Query": cleaned,
                })
                _ = userTmpl.Execute(&userContent, map[string]interface{}{
                    "Query": cleaned,
                })

                chatModel, err3 := p.modelService.GetChatModel(ctx, chatManage.ChatModelID)
                if err3 == nil && chatModel != nil {
                    thinking := false
                    resp, err4 := chatModel.Chat(ctx, []chat.Message{
                        {Role: "system", Content: systemContent.String()},
                        {Role: "user", Content: userContent.String()},
                    }, &chat.ChatOptions{Temperature: 0.1, MaxCompletionTokens: 64, Thinking: &thinking})
                    if err4 == nil && strings.TrimSpace(resp.Content) != "" {
                        content := strings.TrimSpace(resp.Content)
                        content = strings.TrimPrefix(content, "Output:")
                        content = strings.TrimPrefix(content, "输出:")
                        content = strings.ReplaceAll(content, "，", ",")
                        parts := strings.Split(content, ",")
                        keywords := make([]string, 0, len(parts))
                        seen := map[string]struct{}{}
                        for _, part := range parts {
                            w := strings.TrimSpace(part)
                            if w == "" {
                                continue
                            }
                            if _, ok := seen[w]; ok {
                                continue
                            }
                            seen[w] = struct{}{}
                            keywords = append(keywords, w)
                        }
                        if len(keywords) > 0 {
                            chatManage.ProcessedQuery = strings.Join(keywords, " ")
                        }
                    }
                }
            }
        }
    }

    if strings.TrimSpace(chatManage.ProcessedQuery) == "" {
        segments := p.segmentText(cleaned)
        filteredSegments := p.filterStopwords(segments)
        chatManage.ProcessedQuery = strings.Join(filteredSegments, " ")
    }

    logger.GetLogger(ctx).Infof("Query preprocessing complete, processed query: %s", chatManage.ProcessedQuery)

    return next()
}

// cleanText Basic text cleaning
func (p *PluginPreprocess) cleanText(text string) string {
	// Remove URLs
	text = urlRegex.ReplaceAllString(text, " ")

	// Remove email addresses
	text = emailRegex.ReplaceAllString(text, " ")

	// Remove excessive spaces
	text = multiSpaceRegex.ReplaceAllString(text, " ")

	// Remove punctuation marks
	text = punctRegex.ReplaceAllString(text, " ")

	// Trim leading and trailing spaces
	text = strings.TrimSpace(text)

	return text
}

// segmentText Text tokenization
func (p *PluginPreprocess) segmentText(text string) []string {
	// Use Jieba tokenizer for tokenization, using search engine mode
	segments := p.jieba.CutForSearch(text, true)
	return segments
}

// filterStopwords Filter stopwords
func (p *PluginPreprocess) filterStopwords(segments []string) []string {
	var filtered []string

	for _, word := range segments {
		// If not a stopword and not blank, keep it
		if _, isStopword := p.stopwords[word]; !isStopword && !isBlank(word) {
			filtered = append(filtered, word)
		}
	}

	// If filtering results in empty list, return original tokenization results
	if len(filtered) == 0 {
		return segments
	}

	return filtered
}

// isBlank Check if a string is blank
func isBlank(str string) bool {
	for _, r := range str {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// Ensure resources are properly released
func (p *PluginPreprocess) Close() {
	if p.jieba != nil {
		p.jieba.Free()
		p.jieba = nil
	}
}

// ShutdownHandler Returns shutdown function
func (p *PluginPreprocess) ShutdownHandler() func() {
	return func() {
		p.Close()
	}
}
