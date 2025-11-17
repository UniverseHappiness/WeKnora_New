# EventManager.Trigger 实现与事件流水线详解

## 入口与上下文
- 触发点：`internal/application/service/session.go:406` 在会话服务按预设事件序列循环触发：
  - 入口函数：`KnowledgeQAByEvent` `internal/application/service/session.go:378-429`
  - 调用：`err := s.eventManager.Trigger(ctx, event, chatManage)` `internal/application/service/session.go:406`
- 事件序列来源：`types.Pipline["rag_stream"]` `internal/types/chat_manage.go:98-116`
  - 顺序：`rewrite_query → preprocess_query → chunk_search → entity_search → chunk_rerank → chunk_merge → filter_top_k → into_chat_message → chat_completion_stream → stream_filter`

## EventManager 核心实现
- 接口定义：插件需实现 `Plugin` `internal/application/service/chat_pipline/chat_pipline.go:11-21`
- 管理器结构：`EventManager` 维护 `listeners` 与 `handlers` 两个映射 `internal/application/service/chat_pipline/chat_pipline.go:23-29`
- 初始化：`NewEventManager` 构造空映射 `internal/application/service/chat_pipline/chat_pipline.go:31-37`
- 注册：`Register(plugin)`
  - 将插件加入对应事件的 `listeners` 列表
  - 通过 `buildHandler` 为该事件重建职责链 handler `internal/application/service/chat_pipline/chat_pipline.go:39-51`
- 职责链构造：`buildHandler(plugins)` 采用“后向折叠”将 `next` 逐层包裹，形成链式调用 `internal/application/service/chat_pipline/chat_pipline.go:53-68`
  - 链执行顺序与注册顺序一致：先注册的插件先执行，其 `next` 指向后注册的插件
- 触发：`Trigger(ctx, eventType, chatManage)` 查表调用对应事件的链式 `handler`，无则返回 `nil` `internal/application/service/chat_pipline/chat_pipline.go:70-78`

## 错误传播与兜底
- 插件错误类型：`PluginError` 及预置错误如 `ErrSearchNothing` `internal/application/service/chat_pipline/chat_pipline.go:80-125`
- 触发方的处理：
  - 检索空：匹配 `ErrSearchNothing`，下发兜底流 `NewFallbackChan`，并填充 `ChatResponse` `internal/application/service/session.go:409-413`
  - 其它错误：记录日志与追踪，返回底层 `err.Err` `internal/application/service/session.go:416-423`
- 兜底流实现：`NewFallbackChan`/`NewFallback` 构造一次性完成的流式响应 `internal/application/service/chat_pipline/chat_pipline.go:142-161`

## 插件注册与依赖注入
- 在容器构建期注册所有插件到同一个 `EventManager`：`internal/container/container.go:102-117`
  - 先提供 `NewEventManager` `internal/container/container.go:103`
  - 逐个 `Invoke(NewPluginXxx)` 完成插件创建与 `Register`

## 事件与插件映射
- 追踪插件：覆盖所有关键事件用于打点 `internal/application/service/chat_pipline/tracing.go:25-36`
- 重写查询：`PluginRewrite` 响应 `rewrite_query` `internal/application/service/chat_pipline/rewrite.go:49-52`
- 预处理查询：`PluginPreprocess` 响应 `preprocess_query` `internal/application/service/chat_pipline/preprocess.go:82-84`
- 语义检索：`PluginSearch` 响应 `chunk_search` `internal/application/service/chat_pipline/search.go:33-36`
- 图谱检索：`PluginSearchEntity` 响应 `entity_search` `internal/application/service/chat_pipline/search_entity.go:34-36`
- 重排：`PluginRerank` 响应 `chunk_rerank` `internal/application/service/chat_pipline/rerank.go:30-33`
- 合并：`PluginMerge` 响应 `chunk_merge` `internal/application/service/chat_pipline/merge.go:22-25`
- TopK 过滤：`PluginFilterTopK` 响应 `filter_top_k` `internal/application/service/chat_pipline/filter_top_k.go:20-23`
- 生成消息：`PluginIntoChatMessage` 响应 `into_chat_message` `internal/application/service/chat_pipline/into_chat_message.go:28-31`
- 流式生成：`PluginChatCompletionStream` 响应 `chat_completion_stream` `internal/application/service/chat_pipline/chat_completion_stream.go:29-32`
- 流过滤：`PluginStreamFilter` 响应 `stream_filter` `internal/application/service/chat_pipline/stream_filter.go:19-22`

## 典型 rag_stream 执行轨迹（链式 next）
- `rewrite_query`：`PluginRewrite.OnEvent` 基于历史与模板重写查询 `internal/application/service/chat_pipline/rewrite.go:55-118`
- `preprocess_query`：`PluginPreprocess.OnEvent` 清洗分词生成 `ProcessedQuery` `internal/application/service/chat_pipline/preprocess.go:86-109`
- `chunk_search`：`PluginSearch.OnEvent` 调用知识库混合检索，可能追加历史命中 `internal/application/service/chat_pipline/search.go:42-73`
- `entity_search`：`PluginSearchEntity.OnEvent` 基于抽取实体做图谱扩展检索 `internal/application/service/chat_pipline/search_entity.go:40-73`
- `chunk_rerank`：`PluginRerank.OnEvent` 以不同查询候选依次重排、阈值过滤，空则抛 `ErrSearchNothing` `internal/application/service/chat_pipline/rerank.go:35-91`
- `chunk_merge`：`PluginMerge.OnEvent` 按知识来源分组合并、处理重叠与图片信息 `internal/application/service/chat_pipline/merge.go:27-118`
- `filter_top_k`：`PluginFilterTopK.OnEvent` 对 `MergeResult`/`RerankResult`/`SearchResult` 取前 K `internal/application/service/chat_pipline/filter_top_k.go:25-53`
- `into_chat_message`：`PluginIntoChatMessage.OnEvent` 用上下文模板生成用户内容，安全校验与图片信息融合 `internal/application/service/chat_pipline/into_chat_message.go:33-75`
- `chat_completion_stream`：`PluginChatCompletionStream.OnEvent` 取模型与消息，发起流式通道 `internal/application/service/chat_pipline/chat_completion_stream.go:34-62`
- `stream_filter`：`PluginStreamFilter.OnEvent` 包裹流做前缀过滤或透传 `internal/application/service/chat_pipline/stream_filter.go:29-65`

## Tracing 包裹与响应拼装
- 每个事件进入 `PluginTracing.OnEvent`，转发到具体打点函数后调用 `next()` `internal/application/service/chat_pipline/tracing.go:50-74`
- 关键打点示例：
  - 搜索：记录阈值与检索结果 JSON `internal/application/service/chat_pipline/tracing.go:76-95`
  - 重排：记录模型与过滤结果 `internal/application/service/chat_pipline/tracing.go:97-117`
  - 流式：在 goroutine 中累计答案、统计耗时与速率 `internal/application/service/chat_pipline/tracing.go:176-214`

## 执行顺序与扩展方式
- 顺序规则：同一事件下“注册顺序即执行顺序”，通过反向折叠保证 `A → B → C` 链式 `internal/application/service/chat_pipline/chat_pipline.go:53-68`
- 扩展新事件/插件：
  - 新增插件类型并实现 `ActivationEvents` 与 `OnEvent`
  - 在容器中 `Invoke(NewPluginYourX)` 完成注册 `internal/container/container.go:104-117`
  - 由上层服务在合适的 `eventList` 中加入该事件并触发

## 调用方的控制流
- 逐事件循环触发，遇空检索兜底返回，遇其它错误立即终止 `internal/application/service/session.go:403-417`
- 成功则继续直至末尾，最终返回合并后的检索结果与模型流通道 `internal/application/service/session.go:425-429` 与 `internal/application/service/session.go:374-376`