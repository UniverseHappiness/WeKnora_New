# 文件上传日志分析

- 入口路径：`POST /api/v1/knowledge-bases/:id/knowledge/file`（`internal/router/router.go:114`）
- 处理器：`CreateKnowledgeFromFile`（`internal/handler/knowledge.go:85-152`）
- 服务层：`knowledgeService.CreateKnowledgeFromFile`（`internal/application/service/knowledge.go:98-259`）
- 存储层：`FileService.SaveFile`（本地：`internal/application/service/file/local.go:31-79`；MinIO：`internal/application/service/file/minio.go:53-77`）
- 异步解析：`processDocument`/`processDocumentFromURL`→`processChunks`（`internal/application/service/knowledge.go:657-765`、`768-816`、`850-872`）

## 典型成功序列
- 请求进入处理器开始日志：`Start creating knowledge from file`（`internal/handler/knowledge.go:87`）
- 上传文件成功：`File upload successful, filename: <name>, size: <KB>`（`internal/handler/knowledge.go:104`）
- 记录创建动作：`Creating knowledge, knowledge base ID: <kbID>, filename: <name>`（`internal/handler/knowledge.go:105`）
- 元数据解析成功（若提供）：`Received file metadata: <map>`（`internal/handler/knowledge.go:116`）
- 服务层开始：`Start creating knowledge from file`（`internal/application/service/knowledge.go:102`）
- 拉取知识库配置：`Getting knowledge base configuration`（`internal/application/service/knowledge.go:109`）
- 文件类型校验日志：`Checking file type: <filename>`（`internal/application/service/knowledge.go:147`）
- 计算哈希：`Calculating file hash`（`internal/application/service/knowledge.go:154`）
- 查重：`Checking if file exists, tenant ID: <id>`（`internal/application/service/knowledge.go:163`）
- 写入知识记录：
  - `Creating knowledge record`（`internal/application/service/knowledge.go:210`）
  - `Saving knowledge record to database`（`internal/application/service/knowledge.go:228`）
- 保存文件：`Saving file, knowledge ID: <id>`（`internal/application/service/knowledge.go:234`）
  - 本地存储：
    - `Starting to save file locally`（`internal/application/service/file/local.go:34`）
    - `Creating directory: <path>`（`internal/application/service/file/local.go:40`）
    - `Generated file path: <path>`（`internal/application/service/file/local.go:50`）
    - `Copying file content`（`internal/application/service/file/local.go:71`）
    - `File saved successfully: <path>`（`internal/application/service/file/local.go:77`）
  - MinIO 存储：
    - 生成对象名并上传，成功后返回 `minio://<bucket>/<object>`（`internal/application/service/file/minio.go:53-77`）
- 更新知识记录文件路径：`Updating knowledge record with file path`（`internal/application/service/knowledge.go:243`）
- 启动异步解析：`Starting asynchronous document processing`（`internal/application/service/knowledge.go:250`）
- 处理器返回成功：`Knowledge created successfully, ID: <id>, title: <title>`（`internal/handler/knowledge.go:147`）

## 常见错误与对应日志
- 上传失败：
  - `File upload failed`（`internal/handler/knowledge.go:99`）→ 400
- 元数据解析失败：
  - `Failed to parse metadata`（`internal/handler/knowledge.go:112`）→ 400
- `enable_multimodel` 解析失败：
  - `Failed to parse enable_multimodel`（`internal/handler/knowledge.go:124`）→ 400
- 获取知识库配置失败：
  - `Failed to get knowledge base: <err>`（`internal/application/service/knowledge.go:111-113`）→ 500
- 图片多模态配置不完整（COS/VLM/MinIO）：
  - `COS configuration incomplete...`（`internal/application/service/knowledge.go:127-133`）→ 400
  - `VLM configuration incomplete...`（`internal/application/service/knowledge.go:138-141`）→ 400
- 文件类型非法：
  - `Invalid file type`（`internal/application/service/knowledge.go:149-151`）→ 400 / 业务错误 `ErrInvalidFileType`
- 哈希计算失败：
  - `Failed to calculate file hash: <err>`（`internal/application/service/knowledge.go:157-158`）→ 500
- 文件重复：
  - 服务层：`File already exists: <filename>`（`internal/application/service/knowledge.go:175`）
  - 处理器冲突返回：`Detected duplicate file: ...`（`internal/handler/knowledge.go:70-79`）→ 409，携带已存在文档信息
- 存储配额超限：
  - `Storage quota exceeded`（`internal/application/service/knowledge.go:187-188`）→ 400 / 业务错误 `NewStorageQuotaExceededError`
- 文件名非法：
  - `Invalid filename: <name>`（`internal/application/service/knowledge.go:205-207`）→ 400 / 校验错误
- 写库/更新失败：
  - `Failed to create knowledge record...`（`internal/application/service/knowledge.go:230-231`）→ 500
  - `Failed to update knowledge with file path...`（`internal/application/service/knowledge.go:245-246`）→ 500
- 保存文件失败：
  - 本地：`Failed to create directory/open/copy...`（`internal/application/service/file/local.go:42/56/73`）→ 500
  - MinIO：`failed to upload file to MinIO: <err>`（`internal/application/service/file/minio.go:72-74`）→ 500

## 异步解析阶段日志
- 入口与追踪：
  - `processDocument enableMultimodel: <bool>`（`internal/application/service/knowledge.go:661`）
  - `processDocument trace id: <traceId>`（`internal/application/service/knowledge.go:677`）
- 图片但未启用多模态：
  - `processDocument image without enable multimodel`（`internal/application/service/knowledge.go:679-686`）→ `ParseStatus=failed`
- 读取文件失败：
  - `processDocument open file failed`（`internal/application/service/knowledge.go:700-707`）→ `failed`
  - `processDocument read file failed`（`internal/application/service/knowledge.go:752-759`）→ `failed`
- 文档分块（DocReader）：
  - 发送 `ReadFromFile`，失败日志同上（`internal/application/service/knowledge.go:724-751`）
- 分块处理：
  - 获取嵌入模型失败：`processChunks get embedding model failed`（`internal/application/service/knowledge.go:865-869`）
  - 获取摘要模型失败：`processChunks get summary model failed`（`internal/application/service/knowledge.go:874-878`）
  - 统计与处理日志：`chunk_count` 属性、图片子块统计等（`internal/application/service/knowledge.go:857-887`）

## 存储后端差异
- 本地存储：详细目录／文件复制日志（`internal/application/service/file/local.go:34-77`）
- MinIO：对象名生成与 `PutObject` 上传日志（`internal/application/service/file/minio.go:53-77`）

## 快速定位建议
- 按“处理器→服务→存储→异步解析”顺序查看日志；每个阶段都包含明确的 `Info/Debug/Error` 级别输出（例如 `internal/handler/knowledge.go:87/104/147`）
- 冲突（409）优先在处理器查看重复文件日志，再对服务层的查重与更新时间日志进行比对（`internal/application/service/knowledge.go:163-182`）
- 图片多模态相关问题同时检查对象存储与 VLM 配置日志（`internal/application/service/knowledge.go:121-144`）
- 异步解析失败时，结合 trace id 在 `processDocument` → DocReader → `processChunks` 全链路跟踪（`internal/application/service/knowledge.go:677/724/850+`）
