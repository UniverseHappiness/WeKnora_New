toolName: search_codebase
            
status: success
          
query: Where in this project is the question → retrieval → answer pipeline implemented? Look for modules handling query parsing, vector search/DB retrieval, and response generation.
          

toolName: search_codebase
            
status: success
          
query: Find where RAG pipeline events (EventType) are executed: state machine, processing function, chat pipeline service orchestrating preprocess_query, chunk_search, rerank, merge, into_chat_message, chat_completion_stream.
          

        
            
toolName: search_codebase
            
status: success
          
query: Find HTTP handlers/controllers that trigger knowledge QA or chat pipeline, e.g., endpoints like /api/qa, /knowledge, session service usage, calls to KnowledgeQAByEvent or EventManager.Trigger.
          

        
            
toolName: search_codebase
            
status: success
          
query: Locate implementation of FILTER_TOP_K event handler in chat pipeline, likely PluginTracing.FilterTopK or dedicated plugin file.
          

        
            
toolName: search_by_regex
            
status: success
          
query: type EventManager
          

        
            
toolName: view_files
            
status: success
          
            
filePath: /home/wu/myproject/hello/WeKnora/internal/application/service/chat_pipline/chat_pipline.go
          

        
            
toolName: search_codebase
            
status: success
          
query: Find implementation of INTO_CHAT_MESSAGE event handler constructing prompt/messages for chat model.
          

        
            
toolName: search_codebase
            
status: success
          
query: Find implementation of KnowledgeBaseService.HybridSearch and underlying search logic for keyword/vector retrieval.
          

        
            
toolName: view_files
            
status: success
          
            
filePath: /home/wu/myproject/hello/WeKnora/internal/router/router.go
          

        
            
toolName: search_codebase
            
status: success
          
query: Find prepareMessagesWithHistory implementation used in chat pipeline.
          

        
            
toolName: search_codebase
            
status: success
          
query: Where are chat pipeline plugins registered into EventManager? Look for NewEventManager usage and Register(plugin) calls and SessionService construction injecting EventManager.
          

        
            
toolName: search_by_regex
            
status: success
          
query: SSEvent\(
          

        
            
toolName: view_files
            
status: success
          
            
filePath: /home/wu/myproject/hello/WeKnora/internal/handler/session.go
          

        
            
toolName: search_by_regex
            
status: success
          
query: prepareMessagesWithHistory
          

        
            
toolName: view_files
            
status: success
          
            
filePath: /home/wu/myproject/hello/WeKnora/internal/application/service/chat_pipline/common.go
          

        
            
toolName: search_by_regex
            
status: success
          
query: NewEventManager\(
          

        
            
toolName: search_codebase
            
status: success
          
query: Find wiring code where EventManager is instantiated and plugins created: likely in a module initialization or dependency injection setup (dig container). Search for Register(res) calls or NewPluginX used in a provider.
          

        
            
toolName: view_files
            
status: success
          
            
filePath: /home/wu/myproject/hello/WeKnora/internal/container/container.go
          

        
            
toolName: view_files
            
status: success
          
            
filePath: /home/wu/myproject/hello/WeKnora/internal/application/service/chat_pipline/extract_entity.go
          

        
            
toolName: view_files
            
status: success
          
            
filePath: /home/wu/myproject/hello/WeKnora/internal/application/service/chat_pipline/search.go
          

        
            
toolName: search_codebase
            
status: success
          
query: Find server-side stream manager to manage SSE streaming state (RegisterStream, UpdateStream, CompleteStream).
          

        
            
toolName: view_files
            
status: success
          
            
filePath: /home/wu/myproject/hello/WeKnora/internal/application/service/chat_pipline/rewrite.go
          

        
            
toolName: view_files
            
status: success
          
            
filePath: /home/wu/myproject/hello/WeKnora/internal/application/service/chat_pipline/preprocess.go
          

        
            
toolName: search_codebase
            
status: success
          
query: Find summarization/LLM call code for CHAT_COMPLETION (non-stream) and response update in SessionService or handler; check chat_pipline/chat_completion.go.
          

        
            
toolName: search_by_regex
            
status: success
          
query: type SummaryConfig
          

        
            
toolName: view_files
            
status: success
          
            
filePath: /home/wu/myproject/hello/WeKnora/internal/types/session.go
          
**整体链路概览**
- 入口路由将请求交到会话处理器，随后调用会话服务启动“事件流水线”，依次完成改写、预处理、检索、重排、合并、成文、模型生成与流式过滤，最终以 SSE 推送答案与引用
- 事件流水线由可插拔插件实现，按模式选择事件序列；默认问答走 `rag_stream` 模式

**HTTP 入口**
- 路由注册 `POST /api/v1/knowledge-chat/:session_id` 映射到会话处理器的知识问答方法 `internal/router/router.go:188-200`
- 会话处理器接收请求、写消息记录、调用会话服务进行知识问答并以 SSE 输出 `internal/handler/session.go:686-798`
  - 创建用户消息与助手消息占位，并注册流管理器以支持继续拉流
  - 调用 `SessionService.KnowledgeQA` 获取“知识引用 + 模型流式通道”并发往前端

**事件流水线**
- 事件类型与模式定义（`rag_stream` 顺序）：`rewrite_query → preprocess_query → chunk_search → entity_search → chunk_rerank → chunk_merge → filter_top_k → into_chat_message → chat_completion_stream → stream_filter` `internal/types/chat_manage.go:72-118`
- 会话服务按模式触发事件，链式执行插件；出现“检索为空”走兜底 `internal/application/service/session.go:377-430`
- 事件管理器注册与触发机制（职责链）：`internal/application/service/chat_pipline/chat_pipline.go:23-78`
- 所有插件在容器启动时注册到事件管理器：`internal/container/container.go:102-117`

**提问阶段**
- 会话服务构造本次问答的上下文与参数（阈值、TopK、模型、Prompt、模板等），并选用 `rag_stream` 模式触发事件 `internal/application/service/session.go:307-343`

**改写与预处理**
- 改写当前问题，基于历史对话和模板生成更利于检索的语句，同时整理近几轮历史用于后续提示词 `internal/application/service/chat_pipline/rewrite.go:55-109`
- 预处理将改写后的文本做清洗、分词、停用词过滤，得到“关键词序列”作为另一种检索查询 `internal/application/service/chat_pipline/preprocess.go:86-109`

**检索阶段**
- 混合检索（关键词 + 向量），先用改写后的自然语言查询，再用预处理后的关键词序列查询；合并去重为检索结果集 `internal/application/service/chat_pipline/search.go:43-95`
- 混检服务组装向量与关键词检索参数、调用检索引擎，然后批量拉取知识与分块并构建搜索结果 `internal/application/service/knowledgebase.go:296-334,375-419,524-549`
- Postgres 关键词检索（ParadeDB FTS）与向量检索（pgvector）实现 `internal/application/repository/retriever/postgres/repository.go:145-185,206-242`
- 若启用图检索，先抽取实体再对图谱查询，补充新的分块引用 `internal/application/service/chat_pipline/extract_entity.go:50-95`、`internal/application/service/chat_pipline/search_entity.go:32-73`

**重排与合并**
- 使用重排模型对检索结果按相关性打分排序；优先用改写查询，必要时用关键词查询作为补偿 `internal/application/service/chat_pipline/rerank.go:34-69`
- 按知识来源对分块分组，按原文位置排序并做邻近合并，得到用于生成的上下文片段 `internal/application/service/chat_pipline/merge.go:21-68`
- 过滤 Top-K（在合并、重排或检索结果上择一生效）`internal/application/service/chat_pipline/filter_top_k.go:19-53`

**成文与提示**
- 将选定片段与图片信息（Caption/OCR）融合，按上下文模板渲染为最终用户消息内容 `internal/application/service/chat_pipline/into_chat_message.go:22-59,91-191`
- 组装对话消息（系统提示 + 近两轮历史 Q/A + 当前用户内容）供模型生成 `internal/application/service/chat_pipline/common.go:46-67`

**生成与流式输出**
- 获取聊天模型与参数，启动流式生成，得到服务侧流通道 `internal/application/service/chat_pipline/chat_completion_stream.go:28-62`
- 追踪插件在流式过程中记录响应内容与耗时指标 `internal/application/service/chat_pipline/tracing.go:175-214`
- 流式过滤器用于处理“无匹配前缀”、必要时输出兜底响应并收尾 `internal/application/service/chat_pipline/stream_filter.go:20-84`

**SSE 回传与消息落盘**
- 处理器先推送引用，再将模型流的增量片段以 SSE 推送；同时累积保存到助手消息并更新流管理器（用于断线重连继续拉流）`internal/handler/session.go:762-798`
- 内存流管理器记录与更新流状态（注册、追加、完成）`internal/stream/memory_manager.go:33-97`
- 继续拉流接口从缓存增量补齐并以 SSE 续传 `internal/handler/session.go:491-684`

**关键数据结构**
- 本次问答状态体 `ChatManage`：包含查询、检索/重排/合并结果、历史、模板渲染后的用户内容、模型响应通道等 `internal/types/chat_manage.go:22-36`
- 流片段 `StreamResponse`：类型（答案/引用）、内容与完成标记 `internal/types/chat.go:27-44`
- 会话配置 `SummaryConfig`：提示词、上下文模板、温度、TopP、惩罚、令牌上限等 `internal/types/session.go:20-45`

**可配置与扩展点**
- 会话粒度的检索与生成参数（阈值、TopK、模型、模板）在会话实体中设置并由服务注入 `internal/types/session.go:58-72`
- 系统级改写提示与上下文模板在配置文件中维护，可结合业务调整 `config/config.yaml:42-73`
- 检索后端可切换 Postgres/Elasticsearch，容器初始化按环境变量决定注册哪些引擎 `internal/container/container.go:265-334`
- 图检索开关与连接由环境变量控制，未开启则跳过实体抽取与图查询 `internal/application/service/chat_pipline/extract_entity.go:55-58`

**调试与定位建议**
- 按“入口→服务→事件→插件”顺序查看日志；每个事件触发与完成都有日志与追踪记录 `internal/application/service/session.go:377-430`
- 若检索为空会记录并走兜底，可从 HybridSearch 参数与阈值着手调整 `internal/application/service/chat_pipline/search.go:81-95`
- 生成阶段问题（提示与模板）可检查 `ContextTemplate/Prompt` 与 `prepareMessagesWithHistory` 的最终消息 `internal/application/service/chat_pipline/common.go:46-67`

这就是从“提问 → 检索 → 回答”的完整可观测链路，涉及入口路由与处理器、会话服务触发的事件流水线、可插拔检索与生成插件、以及 SSE 流的回传与持久化流程。您若需要我针对某一环节做参数调整或加日志定位，我可以直接帮您改到位并跑通验证。