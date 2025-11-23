# API参考

<cite>
**本文档引用的文件**   
- [main.go](file://cmd/server/main.go)
- [router.go](file://internal/router/router.go)
- [auth.go](file://internal/handler/auth.go)
- [knowledgebase.go](file://internal/handler/knowledgebase.go)
- [session.go](file://internal/handler/session.go)
- [message.go](file://internal/handler/message.go)
- [model.go](file://internal/handler/model.go)
- [knowledge.go](file://internal/handler/knowledge.go)
- [tenant.go](file://internal/handler/tenant.go)
- [config.yaml](file://config/config.yaml)
- [interfaces/user.go](file://internal/types/interfaces/user.go)
- [chat.go](file://internal/types/chat.go)
- [knowledgebase.go](file://internal/types/knowledgebase.go)
- [session.go](file://internal/types/session.go)
</cite>

## 目录
1. [简介](#简介)
2. [认证与用户管理](#认证与用户管理)
3. [知识库管理](#知识库管理)
4. [文档操作](#文档操作)
5. [会话与消息](#会话与消息)
6. [模型配置](#模型配置)
7. [租户管理](#租户管理)
8. [错误处理](#错误处理)

## 简介
WeKnora_New后端服务提供了一套完整的RESTful API，用于构建智能知识问答系统。API基于Gin框架实现，采用JWT进行身份验证，支持知识库管理、文档处理、会话交互、模型配置等功能。所有API端点均位于`/api/v1`路径下，通过HTTP方法、路径参数、查询参数和请求体来定义操作。

API设计遵循REST原则，使用标准HTTP状态码表示响应结果。请求和响应数据格式均为JSON。系统通过租户（Tenant）机制实现多租户隔离，每个用户属于特定租户，资源访问受租户权限控制。

**Section sources**
- [main.go](file://cmd/server/main.go#L1-L105)
- [router.go](file://internal/router/router.go#L1-L283)

## 认证与用户管理
认证模块提供用户注册、登录、令牌管理等功能，所有需要认证的API端点都需要在请求头中包含有效的JWT令牌。

### 用户注册
创建新用户账户。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/auth/register`  
**请求头**: 
- `Content-Type: application/json`

**请求体**:
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**响应状态码**:
- `201 Created`: 用户注册成功
- `400 Bad Request`: 请求参数无效
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Registration successful",
  "user": {
    "id": "string",
    "username": "string",
    "email": "string",
    "tenant_id": 0,
    "created_at": "string",
    "updated_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [auth.go](file://internal/handler/auth.go#L36-L80)
- [interfaces/user.go](file://internal/types/interfaces/user.go#L11-L12)

### 用户登录
验证用户凭据并获取访问令牌。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/auth/login`  
**请求头**: 
- `Content-Type: application/json`

**请求体**:
```json
{
  "email": "string",
  "password": "string"
}
```

**响应状态码**:
- `200 OK`: 登录成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 认证失败
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "access_token": "string",
  "refresh_token": "string",
  "expires_in": 86400
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [auth.go](file://internal/handler/auth.go#L82-L128)
- [interfaces/user.go](file://internal/types/interfaces/user.go#L13-L14)

### 令牌刷新
使用刷新令牌获取新的访问令牌。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/auth/refresh`  
**请求头**: 
- `Content-Type: application/json`

**请求体**:
```json
{
  "refreshToken": "string"
}
```

**响应状态码**:
- `200 OK`: 令牌刷新成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 刷新令牌无效或已过期
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "access_token": "string",
  "refresh_token": "string"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [auth.go](file://internal/handler/auth.go#L175-L211)
- [interfaces/user.go](file://internal/types/interfaces/user.go#L33-L34)

### 令牌验证
验证访问令牌的有效性。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/auth/validate`  
**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 令牌有效
- `400 Bad Request`: 请求头格式错误
- `401 Unauthorized`: 令牌无效或已过期
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Token is valid",
  "user": {
    "id": "string",
    "username": "string",
    "email": "string",
    "tenant_id": 0
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [auth.go](file://internal/handler/auth.go#L299-L343)
- [interfaces/user.go](file://internal/types/interfaces/user.go#L31-L32)

### 获取当前用户信息
获取当前认证用户的信息。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/auth/me`  
**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取用户信息
- `401 Unauthorized`: 未授权或令牌无效
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "string",
      "username": "string",
      "email": "string",
      "tenant_id": 0,
      "created_at": "string",
      "updated_at": "string"
    },
    "tenant": {
      "id": 0,
      "name": "string",
      "description": "string",
      "created_at": "string",
      "updated_at": "string"
    }
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [auth.go](file://internal/handler/auth.go#L213-L251)
- [interfaces/user.go](file://internal/types/interfaces/user.go#L37-L38)

### 修改密码
修改当前用户的密码。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/auth/change-password`  
**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "old_password": "string",
  "new_password": "string"
}
```

**响应状态码**:
- `200 OK`: 密码修改成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权或旧密码错误
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [auth.go](file://internal/handler/auth.go#L253-L297)
- [interfaces/user.go](file://internal/types/interfaces/user.go#L25-L26)

### 用户登出
使当前用户的访问令牌失效。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/auth/logout`  
**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 登出成功
- `400 Bad Request`: 请求头格式错误
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Logout successful"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [auth.go](file://internal/handler/auth.go#L130-L173)
- [interfaces/user.go](file://internal/types/interfaces/user.go#L35-L36)

## 知识库管理
知识库是文档和知识的容器，支持创建、查询、更新、删除和混合搜索等操作。

### 创建知识库
创建一个新的知识库。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/knowledge-bases`  
**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "chunking_config": {
    "chunk_size": 0,
    "chunk_overlap": 0,
    "separators": ["string"],
    "enable_multimodal": false
  },
  "image_processing_config": {
    "model_id": "string"
  },
  "embedding_model_id": "string",
  "summary_model_id": "string",
  "rerank_model_id": "string",
  "vlm_model_id": "string",
  "vlm_config": {
    "model_name": "string",
    "base_url": "string",
    "api_key": "string",
    "interface_type": "string"
  },
  "cos_config": {
    "secret_id": "string",
    "secret_key": "string",
    "region": "string",
    "bucket_name": "string",
    "app_id": "string",
    "path_prefix": "string",
    "provider": "string"
  },
  "extract_config": {
    "text": "string",
    "tags": ["string"],
    "nodes": [
      {
        "name": "string",
        "attributes": ["string"]
      }
    ],
    "relations": [
      {
        "node1": "string",
        "node2": "string",
        "type": "string"
      }
    ]
  }
}
```

**响应状态码**:
- `201 Created`: 知识库创建成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "name": "string",
    "description": "string",
    "tenant_id": 0,
    "chunking_config": {},
    "image_processing_config": {},
    "embedding_model_id": "string",
    "summary_model_id": "string",
    "rerank_model_id": "string",
    "vlm_model_id": "string",
    "vlm_config": {},
    "cos_config": {},
    "extract_config": {},
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledgebase.go](file://internal/handler/knowledgebase.go#L67-L95)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 获取知识库列表
获取当前租户的所有知识库。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/knowledge-bases`  
**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取知识库列表
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "tenant_id": 0,
      "chunking_config": {},
      "image_processing_config": {},
      "embedding_model_id": "string",
      "summary_model_id": "string",
      "rerank_model_id": "string",
      "vlm_model_id": "string",
      "vlm_config": {},
      "cos_config": {},
      "extract_config": {},
      "created_at": "string",
      "updated_at": "string",
      "deleted_at": "string"
    }
  ]
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledgebase.go](file://internal/handler/knowledgebase.go#L147-L148)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 获取知识库详情
获取指定知识库的详细信息。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/knowledge-bases/{id}`  
**路径参数**:
- `id`: 知识库ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取知识库详情
- `400 Bad Request`: 知识库ID为空
- `401 Unauthorized`: 未授权
- `403 Forbidden`: 无权访问该知识库
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "name": "string",
    "description": "string",
    "tenant_id": 0,
    "chunking_config": {},
    "image_processing_config": {},
    "embedding_model_id": "string",
    "summary_model_id": "string",
    "rerank_model_id": "string",
    "vlm_model_id": "string",
    "vlm_config": {},
    "cos_config": {},
    "extract_config": {},
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledgebase.go](file://internal/handler/knowledgebase.go#L149-L150)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 更新知识库
更新现有知识库的信息。

**HTTP方法**: `PUT`  
**URL路径**: `/api/v1/knowledge-bases/{id}`  
**路径参数**:
- `id`: 知识库ID

**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "config": {
    "chunking_config": {
      "chunk_size": 0,
      "chunk_overlap": 0,
      "separators": ["string"],
      "enable_multimodal": false
    },
    "image_processing_config": {
      "model_id": "string"
    }
  }
}
```

**响应状态码**:
- `200 OK`: 知识库更新成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `403 Forbidden`: 无权访问该知识库
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "name": "string",
    "description": "string",
    "tenant_id": 0,
    "chunking_config": {},
    "image_processing_config": {},
    "embedding_model_id": "string",
    "summary_model_id": "string",
    "rerank_model_id": "string",
    "vlm_model_id": "string",
    "vlm_config": {},
    "cos_config": {},
    "extract_config": {},
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledgebase.go](file://internal/handler/knowledgebase.go#L151-L152)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 删除知识库
删除指定的知识库。

**HTTP方法**: `DELETE`  
**URL路径**: `/api/v1/knowledge-bases/{id}`  
**路径参数**:
- `id`: 知识库ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 知识库删除成功
- `400 Bad Request`: 知识库ID为空
- `401 Unauthorized`: 未授权
- `403 Forbidden`: 无权访问该知识库
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Knowledge base deleted successfully"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledgebase.go](file://internal/handler/knowledgebase.go#L153-L154)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 混合搜索
在指定知识库中执行混合向量和关键词搜索。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/knowledge-bases/{id}/hybrid-search`  
**路径参数**:
- `id`: 知识库ID

**查询参数**:
- `query_text`: 搜索查询文本
- `top_k`: 返回结果数量
- `keyword_threshold`: 关键词召回阈值
- `vector_threshold`: 向量召回阈值

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 搜索成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `403 Forbidden`: 无权访问该知识库
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": [
    {
      "id": "string",
      "knowledge_id": "string",
      "content": "string",
      "metadata": {},
      "score": 0,
      "source": "string"
    }
  ]
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledgebase.go](file://internal/handler/knowledgebase.go#L155-L156)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 拷贝知识库
拷贝一个知识库到另一个知识库。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/knowledge-bases/copy`  
**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "source_id": "string",
  "target_id": "string"
}
```

**响应状态码**:
- `200 OK`: 知识库拷贝成功
- `400 Bad Request`: 请求参数无效
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Knowledge base copy successfully"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledgebase.go](file://internal/handler/knowledgebase.go#L270-L297)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

## 文档操作
文档操作包括从文件或URL创建知识、获取知识详情、删除知识、更新知识等。

### 从文件创建知识
从上传的文件创建知识条目。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/knowledge-bases/{id}/knowledge/file`  
**路径参数**:
- `id`: 知识库ID

**请求头**: 
- `Authorization: Bearer <token>`

**请求体** (multipart/form-data):
- `file`: 要上传的文件
- `metadata`: 文件元数据（JSON格式）
- `enable_multimodel`: 是否启用多模态处理

**响应状态码**:
- `200 OK`: 知识创建成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `403 Forbidden`: 无权访问该知识库
- `409 Conflict`: 文件已存在
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "knowledge_base_id": "string",
    "title": "string",
    "type": "string",
    "source": "string",
    "status": "string",
    "metadata": {},
    "chunking_config": {},
    "image_processing_config": {},
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string",
  "code": "string"
}
```

**Section sources**
- [knowledge.go](file://internal/handler/knowledge.go#L84-L152)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 从URL创建知识
从指定URL创建知识条目。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/knowledge-bases/{id}/knowledge/url`  
**路径参数**:
- `id`: 知识库ID

**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "url": "string",
  "enable_multimodel": false
}
```

**响应状态码**:
- `201 Created`: 知识创建成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `403 Forbidden`: 无权访问该知识库
- `409 Conflict`: URL已存在
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "knowledge_base_id": "string",
    "title": "string",
    "type": "string",
    "source": "string",
    "status": "string",
    "metadata": {},
    "chunking_config": {},
    "image_processing_config": {},
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string",
  "code": "string"
}
```

**Section sources**
- [knowledge.go](file://internal/handler/knowledge.go#L154-L197)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 获取知识列表
获取知识库中的知识列表。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/knowledge-bases/{id}/knowledge`  
**路径参数**:
- `id`: 知识库ID

**查询参数**:
- `page`: 页码
- `page_size`: 每页数量

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取知识列表
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `403 Forbidden`: 无权访问该知识库
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": [
    {
      "id": "string",
      "knowledge_base_id": "string",
      "title": "string",
      "type": "string",
      "source": "string",
      "status": "string",
      "metadata": {},
      "chunking_config": {},
      "image_processing_config": {},
      "created_at": "string",
      "updated_at": "string",
      "deleted_at": "string"
    }
  ],
  "total": 0,
  "page": 0,
  "page_size": 0
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledge.go](file://internal/handler/knowledge.go#L117-L118)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 获取知识详情
获取指定知识的详细信息。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/knowledge/{id}`  
**路径参数**:
- `id`: 知识ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取知识详情
- `400 Bad Request`: 知识ID为空
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "knowledge_base_id": "string",
    "title": "string",
    "type": "string",
    "source": "string",
    "status": "string",
    "metadata": {},
    "chunking_config": {},
    "image_processing_config": {},
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledge.go](file://internal/handler/knowledge.go#L200-L226)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 删除知识
删除指定的知识条目。

**HTTP方法**: `DELETE`  
**URL路径**: `/api/v1/knowledge/{id}`  
**路径参数**:
- `id`: 知识ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 知识删除成功
- `400 Bad Request`: 知识ID为空
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Deleted successfully"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledge.go](file://internal/handler/knowledge.go#L272-L297)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 更新知识
更新知识条目的信息。

**HTTP方法**: `PUT`  
**URL路径**: `/api/v1/knowledge/{id}`  
**路径参数**:
- `id`: 知识ID

**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "id": "string",
  "knowledge_base_id": "string",
  "title": "string",
  "type": "string",
  "source": "string",
  "status": "string",
  "metadata": {},
  "chunking_config": {},
  "image_processing_config": {}
}
```

**响应状态码**:
- `200 OK`: 知识更新成功
- `400 Bad Request`: 请求参数无效
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Knowledge chunk updated successfully"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledge.go](file://internal/handler/knowledge.go#L401-L431)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 下载知识文件
下载与知识条目关联的原始文件。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/knowledge/{id}/download`  
**路径参数**:
- `id`: 知识ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 文件下载成功
- `400 Bad Request`: 知识ID为空
- `500 Internal Server Error`: 服务器内部错误

**响应**: 文件流

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledge.go](file://internal/handler/knowledge.go#L300-L345)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

### 批量获取知识
批量获取多个知识条目。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/knowledge/batch`  
**查询参数**:
- `ids`: 知识ID列表（逗号分隔）

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取知识列表
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": [
    {
      "id": "string",
      "knowledge_base_id": "string",
      "title": "string",
      "type": "string",
      "source": "string",
      "status": "string",
      "metadata": {},
      "chunking_config": {},
      "image_processing_config": {},
      "created_at": "string",
      "updated_at": "string",
      "deleted_at": "string"
    }
  ]
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [knowledge.go](file://internal/handler/knowledge.go#L352-L399)
- [knowledgebase.go](file://internal/types/knowledgebase.go#L16-L49)

## 会话与消息
会话管理用户与系统的交互，支持创建会话、发送消息、加载历史消息等功能。

### 创建会话
创建一个新的对话会话。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/sessions`  
**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "knowledge_base_id": "string",
  "session_strategy": {
    "max_rounds": 0,
    "enable_rewrite": false,
    "fallback_strategy": "string",
    "fallback_response": "string",
    "embedding_top_k": 0,
    "keyword_threshold": 0,
    "vector_threshold": 0,
    "rerank_model_id": "string",
    "rerank_top_k": 0,
    "rerank_threshold": 0,
    "summary_model_id": "string",
    "summary_parameters": {
      "max_tokens": 0,
      "repeat_penalty": 0,
      "top_k": 0,
      "top_p": 0,
      "frequency_penalty": 0,
      "presence_penalty": 0,
      "prompt": "string",
      "context_template": "string",
      "no_match_prefix": "string",
      "temperature": 0,
      "seed": 0,
      "max_completion_tokens": 0
    },
    "no_match_prefix": "string"
  }
}
```

**响应状态码**:
- `201 Created`: 会话创建成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "title": "string",
    "description": "string",
    "tenant_id": 0,
    "knowledge_base_id": "string",
    "max_rounds": 0,
    "enable_rewrite": false,
    "fallback_strategy": "string",
    "fallback_response": "string",
    "embedding_top_k": 0,
    "keyword_threshold": 0,
    "vector_threshold": 0,
    "rerank_model_id": "string",
    "rerank_top_k": 0,
    "rerank_threshold": 0,
    "summary_model_id": "string",
    "summary_parameters": {},
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string",
    "messages": []
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [session.go](file://internal/handler/session.go#L82-L223)
- [session.go](file://internal/types/session.go#L48-L79)

### 获取会话详情
获取指定会话的详细信息。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/sessions/{id}`  
**路径参数**:
- `id`: 会话ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取会话详情
- `400 Bad Request`: 会话ID为空
- `401 Unauthorized`: 未授权
- `404 Not Found`: 会话不存在
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "title": "string",
    "description": "string",
    "tenant_id": 0,
    "knowledge_base_id": "string",
    "max_rounds": 0,
    "enable_rewrite": false,
    "fallback_strategy": "string",
    "fallback_response": "string",
    "embedding_top_k": 0,
    "keyword_threshold": 0,
    "vector_threshold": 0,
    "rerank_model_id": "string",
    "rerank_top_k": 0,
    "rerank_threshold": 0,
    "summary_model_id": "string",
    "summary_parameters": {},
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string",
    "messages": []
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [session.go](file://internal/handler/session.go#L225-L259)
- [session.go](file://internal/types/session.go#L48-L79)

### 获取会话列表
获取当前租户的所有会话。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/sessions`  
**查询参数**:
- `page`: 页码
- `page_size`: 每页数量

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取会话列表
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": [
    {
      "id": "string",
      "title": "string",
      "description": "string",
      "tenant_id": 0,
      "knowledge_base_id": "string",
      "max_rounds": 0,
      "enable_rewrite": false,
      "fallback_strategy": "string",
      "fallback_response": "string",
      "embedding_top_k": 0,
      "keyword_threshold": 0,
      "vector_threshold": 0,
      "rerank_model_id": "string",
      "rerank_top_k": 0,
      "rerank_threshold": 0,
      "summary_model_id": "string",
      "summary_parameters": {},
      "created_at": "string",
      "updated_at": "string",
      "deleted_at": "string",
      "messages": []
    }
  ],
  "total": 0,
  "page": 0,
  "page_size": 0
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [session.go](file://internal/handler/session.go#L261-L294)
- [session.go](file://internal/types/session.go#L48-L79)

### 更新会话
更新现有会话的信息。

**HTTP方法**: `PUT`  
**URL路径**: `/api/v1/sessions/{id}`  
**路径参数**:
- `id`: 会话ID

**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "tenant_id": 0,
  "knowledge_base_id": "string",
  "max_rounds": 0,
  "enable_rewrite": false,
  "fallback_strategy": "string",
  "fallback_response": "string",
  "embedding_top_k": 0,
  "keyword_threshold": 0,
  "vector_threshold": 0,
  "rerank_model_id": "string",
  "rerank_top_k": 0,
  "rerank_threshold": 0,
  "summary_model_id": "string",
  "summary_parameters": {},
  "messages": []
}
```

**响应状态码**:
- `200 OK`: 会话更新成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `404 Not Found`: 会话不存在
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "title": "string",
    "description": "string",
    "tenant_id": 0,
    "knowledge_base_id": "string",
    "max_rounds": 0,
    "enable_rewrite": false,
    "fallback_strategy": "string",
    "fallback_response": "string",
    "embedding_top_k": 0,
    "keyword_threshold": 0,
    "vector_threshold": 0,
    "rerank_model_id": "string",
    "rerank_top_k": 0,
    "rerank_threshold": 0,
    "summary_model_id": "string",
    "summary_parameters": {},
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string",
    "messages": []
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [session.go](file://internal/handler/session.go#L296-L348)
- [session.go](file://internal/types/session.go#L48-L79)

### 删除会话
删除指定的会话。

**HTTP方法**: `DELETE`  
**URL路径**: `/api/v1/sessions/{id}`  
**路径参数**:
- `id`: 会话ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 会话删除成功
- `400 Bad Request`: 会话ID为空
- `401 Unauthorized`: 未授权
- `404 Not Found`: 会话不存在
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Session deleted successfully"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [session.go](file://internal/handler/session.go#L351-L383)
- [session.go](file://internal/types/session.go#L48-L79)

### 生成会话标题
为会话生成一个标题。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/sessions/{session_id}/generate_title`  
**路径参数**:
- `session_id`: 会话ID

**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "messages": [
    {
      "session_id": "string",
      "role": "string",
      "content": "string",
      "request_id": "string",
      "created_at": "string",
      "is_completed": false
    }
  ]
}
```

**响应状态码**:
- `200 OK`: 标题生成成功
- `400 Bad Request`: 请求参数无效
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": "string"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [session.go](file://internal/handler/session.go#L391-L428)
- [session.go](file://internal/types/session.go#L48-L79)

### 知识问答
在会话中进行知识问答。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/knowledge-chat/{session_id}`  
**路径参数**:
- `session_id`: 会话ID

**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "query": "string"
}
```

**响应状态码**:
- `200 OK`: 流式响应开始
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**响应**: Server-Sent Events (SSE) 流

**事件类型**:
- `message`: 包含响应片段或引用信息

**响应数据**:
```json
{
  "id": "string",
  "response_type": "answer|references",
  "content": "string",
  "done": false,
  "knowledge_references": [
    {
      "id": "string",
      "knowledge_id": "string",
      "content": "string",
      "metadata": {},
      "score": 0,
      "source": "string"
    }
  ]
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [session.go](file://internal/handler/session.go#L686-L797)
- [chat.go](file://internal/types/chat.go#L32-L44)

### 继续流式响应
继续接收流式响应。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/sessions/continue-stream/{session_id}`  
**路径参数**:
- `session_id`: 会话ID

**查询参数**:
- `message_id`: 消息ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 流式响应继续
- `400 Bad Request`: 参数无效
- `401 Unauthorized`: 未授权
- `404 Not Found`: 会话或消息不存在
- `500 Internal Server Error`: 服务器内部错误

**响应**: Server-Sent Events (SSE) 流

**事件类型**:
- `message`: 包含响应片段或引用信息

**响应数据**:
```json
{
  "id": "string",
  "response_type": "answer|references",
  "content": "string",
  "done": false,
  "knowledge_references": [
    {
      "id": "string",
      "knowledge_id": "string",
      "content": "string",
      "metadata": {},
      "score": 0,
      "source": "string"
    }
  ]
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [session.go](file://internal/handler/session.go#L492-L684)
- [chat.go](file://internal/types/chat.go#L32-L44)

### 知识搜索
在知识库中搜索知识，不进行LLM总结。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/knowledge-search`  
**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "query": "string",
  "knowledge_base_id": "string"
}
```

**响应状态码**:
- `200 OK`: 搜索成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": [
    {
      "id": "string",
      "knowledge_id": "string",
      "content": "string",
      "metadata": {},
      "score": 0,
      "source": "string"
    }
  ]
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [session.go](file://internal/handler/session.go#L441-L490)
- [chat.go](file://internal/types/chat.go#L32-L44)

## 模型配置
模型管理模块用于创建、查询、更新和删除模型配置。

### 创建模型
创建一个新的模型配置。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/models`  
**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "name": "string",
  "type": "string",
  "source": "string",
  "description": "string",
  "parameters": {
    "max_tokens": 0,
    "temperature": 0,
    "top_p": 0,
    "frequency_penalty": 0,
    "presence_penalty": 0,
    "repeat_penalty": 0,
    "seed": 0,
    "max_completion_tokens": 0,
    "prompt": "string",
    "context_template": "string",
    "no_match_prefix": "string"
  },
  "is_default": false
}
```

**响应状态码**:
- `201 Created`: 模型创建成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "tenant_id": 0,
    "name": "string",
    "type": "string",
    "source": "string",
    "description": "string",
    "parameters": {},
    "is_default": false,
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [model.go](file://internal/handler/model.go#L41-L89)
- [model.go](file://internal/types/model.go#L15-L45)

### 获取模型列表
获取当前租户的所有模型。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/models`  
**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取模型列表
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": [
    {
      "id": "string",
      "tenant_id": 0,
      "name": "string",
      "type": "string",
      "source": "string",
      "description": "string",
      "parameters": {},
      "is_default": false,
      "created_at": "string",
      "updated_at": "string",
      "deleted_at": "string"
    }
  ]
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [model.go](file://internal/handler/model.go#L128-L157)
- [model.go](file://internal/types/model.go#L15-L45)

### 获取模型详情
获取指定模型的详细信息。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/models/{id}`  
**路径参数**:
- `id`: 模型ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取模型详情
- `400 Bad Request`: 模型ID为空
- `401 Unauthorized`: 未授权
- `404 Not Found`: 模型不存在
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "tenant_id": 0,
    "name": "string",
    "type": "string",
    "source": "string",
    "description": "string",
    "parameters": {},
    "is_default": false,
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [model.go](file://internal/handler/model.go#L91-L126)
- [model.go](file://internal/types/model.go#L15-L45)

### 更新模型
更新现有模型的配置。

**HTTP方法**: `PUT`  
**URL路径**: `/api/v1/models/{id}`  
**路径参数**:
- `id`: 模型ID

**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "parameters": {
    "max_tokens": 0,
    "temperature": 0,
    "top_p": 0,
    "frequency_penalty": 0,
    "presence_penalty": 0,
    "repeat_penalty": 0,
    "seed": 0,
    "max_completion_tokens": 0,
    "prompt": "string",
    "context_template": "string",
    "no_match_prefix": "string"
  },
  "is_default": false
}
```

**响应状态码**:
- `200 OK`: 模型更新成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `404 Not Found`: 模型不存在
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": "string",
    "tenant_id": 0,
    "name": "string",
    "type": "string",
    "source": "string",
    "description": "string",
    "parameters": {},
    "is_default": false,
    "created_at": "string",
    "updated_at": "string",
    "deleted_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [model.go](file://internal/handler/model.go#L168-L229)
- [model.go](file://internal/types/model.go#L15-L45)

### 删除模型
删除指定的模型。

**HTTP方法**: `DELETE`  
**URL路径**: `/api/v1/models/{id}`  
**路径参数**:
- `id`: 模型ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 模型删除成功
- `400 Bad Request`: 模型ID为空
- `401 Unauthorized`: 未授权
- `404 Not Found`: 模型不存在
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Model deleted"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [model.go](file://internal/handler/model.go#L231-L265)
- [model.go](file://internal/types/model.go#L15-L45)

## 租户管理
租户管理模块用于创建、查询、更新和删除租户。

### 创建租户
创建一个新的租户。

**HTTP方法**: `POST`  
**URL路径**: `/api/v1/tenants`  
**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "name": "string",
  "description": "string"
}
```

**响应状态码**:
- `201 Created`: 租户创建成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": 0,
    "name": "string",
    "description": "string",
    "created_at": "string",
    "updated_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [tenant.go](file://internal/handler/tenant.go#L33-L71)
- [tenant.go](file://internal/types/tenant.go#L11-L20)

### 获取租户详情
获取指定租户的详细信息。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/tenants/{id}`  
**路径参数**:
- `id`: 租户ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取租户详情
- `400 Bad Request`: 租户ID无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": 0,
    "name": "string",
    "description": "string",
    "created_at": "string",
    "updated_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [tenant.go](file://internal/handler/tenant.go#L73-L110)
- [tenant.go](file://internal/types/tenant.go#L11-L20)

### 更新租户
更新现有租户的信息。

**HTTP方法**: `PUT`  
**URL路径**: `/api/v1/tenants/{id}`  
**路径参数**:
- `id`: 租户ID

**请求头**: 
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

**请求体**:
```json
{
  "id": 0,
  "name": "string",
  "description": "string"
}
```

**响应状态码**:
- `200 OK`: 租户更新成功
- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "id": 0,
    "name": "string",
    "description": "string",
    "created_at": "string",
    "updated_at": "string"
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [tenant.go](file://internal/handler/tenant.go#L112-L157)
- [tenant.go](file://internal/types/tenant.go#L11-L20)

### 删除租户
删除指定的租户。

**HTTP方法**: `DELETE`  
**URL路径**: `/api/v1/tenants/{id}`  
**路径参数**:
- `id`: 租户ID

**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 租户删除成功
- `400 Bad Request`: 租户ID无效
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "message": "Tenant deleted successfully"
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [tenant.go](file://internal/handler/tenant.go#L159-L195)
- [tenant.go](file://internal/types/tenant.go#L11-L20)

### 获取租户列表
获取所有租户的列表。

**HTTP方法**: `GET`  
**URL路径**: `/api/v1/tenants`  
**请求头**: 
- `Authorization: Bearer <token>`

**响应状态码**:
- `200 OK`: 成功获取租户列表
- `401 Unauthorized`: 未授权
- `500 Internal Server Error`: 服务器内部错误

**成功响应体**:
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 0,
        "name": "string",
        "description": "string",
        "created_at": "string",
        "updated_at": "string"
      }
    ]
  }
}
```

**错误响应体**:
```json
{
  "success": false,
  "error": "string",
  "details": "string"
}
```

**Section sources**
- [tenant.go](file://internal/handler/tenant.go#L197-L226)
- [tenant.go](file://internal/types/tenant.go#L11-L20)

## 错误处理
API使用标准的HTTP状态码和统一的错误响应格式来处理错误情况。

### 常见错误状态码
- `400 Bad Request`: 请求参数无效或缺失
- `401 Unauthorized`: 未提供认证信息或认证失败
- `403 Forbidden`: 无权访问指定资源
- `404 Not Found`: 请求的资源不存在
- `409 Conflict`: 请求与现有资源冲突（如重复创建）
- `500 Internal Server Error`: 服务器内部错误

### 错误响应格式
所有错误响应都遵循统一的格式：
```json
{
  "success": false,
  "error": "错误消息",
  "details": "详细错误信息"
}
```

### 常见错误场景
1. **认证失败**: 未提供有效的JWT令牌或令牌已过期
2. **权限不足**: 用户尝试访问不属于其租户的资源
3. **参数校验失败**: 请求参数不符合验证规则
4. **资源不存在**: 请求的ID对应的资源不存在
5. **重复资源**: 尝试创建已存在的资源

**Section sources**
- [errors.go](file://internal/errors/errors.go)
- [middleware/error_handler.go](file://internal/middleware/error_handler.go)