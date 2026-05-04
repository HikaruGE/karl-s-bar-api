# Karl's Bar API 文档

## 基础信息

- **服务器地址**: `http://localhost:9527`
- **内容类型**: `application/json`
- **认证方式**: JWT Bearer Token

---

## 目录

1. [健康检查](#健康检查)
2. [鸡尾酒相关](#鸡尾酒相关)
   - [获取鸡尾酒列表](#获取鸡尾酒列表)
   - [获取单个鸡尾酒](#获取单个鸡尾酒)
3. [认证相关](#认证相关)
   - [用户注册](#用户注册)
   - [用户登录](#用户登录)
4. [收藏相关](#收藏相关)
   - [添加收藏](#添加收藏)
   - [获取收藏列表](#获取收藏列表)
   - [删除收藏](#删除收藏)

---

## 健康检查

### 获取服务器状态

**端点**: `GET /cheers`

**描述**: 检查服务器是否在线

**认证**: 否

**请求示例**:

```bash
curl -X GET http://localhost:9527/cheers
```

**响应示例**:

```json
{
  "message": "pong"
}
```

**状态码**:

- `200` - 服务器正常

---

## 鸡尾酒相关

### 获取鸡尾酒列表

**端点**: `GET /cocktails`

**描述**: 获取所有鸡尾酒列表

**认证**: 否

**查询参数**: 无

**请求示例**:

```bash
curl -X GET http://localhost:9527/cocktails
```

**响应示例**:

```json
{
  "cocktails": [
    {
      "id": "1",
      "name": "Mojito",
      "category": "Cocktail",
      "description": "Refreshing Cuban cocktail with rum, mint, and lime",
      "image": "https://images.pexels.com/photos/3407750/...",
      "ingredients": [
        "60 ml White Rum",
        "22 ml Fresh Lime Juice",
        "15 ml Simple Syrup",
        "8-10 Fresh Mint Leaves",
        "Club Soda",
        "Ice",
        "Lime Slice"
      ],
      "instructions": "1. Add mint leaves and simple syrup to a tall glass...",
      "abv": 12,
      "servingSize": "355 ml"
    }
  ]
}
```

**状态码**:

- `200` - 成功
- `500` - 服务器错误

---

### 获取单个鸡尾酒

**端点**: `GET /cocktails/:id`

**描述**: 根据 ID 获取单个鸡尾酒详情

**认证**: 否

**路径参数**:

- `id` (string, required) - 鸡尾酒 ID（例如："1", "2", "3"）

**请求示例**:

```bash
curl -X GET http://localhost:9527/cocktails/1
```

**响应示例**:

```json
{
  "id": "1",
  "name": "Mojito",
  "category": "Cocktail",
  "description": "Refreshing Cuban cocktail with rum, mint, and lime",
  "image": "https://images.pexels.com/photos/3407750/...",
  "ingredients": [
    "60 ml White Rum",
    "22 ml Fresh Lime Juice",
    "15 ml Simple Syrup",
    "8-10 Fresh Mint Leaves",
    "Club Soda",
    "Ice",
    "Lime Slice"
  ],
  "instructions": "1. Add mint leaves and simple syrup to a tall glass...",
  "abv": 12,
  "servingSize": "355 ml"
}
```

**状态码**:

- `200` - 成功
- `404` - 鸡尾酒不存在
- `500` - 服务器错误

---

## 认证相关

### 用户注册

**端点**: `POST /auth/register`

**描述**: 用新的邮箱和密码注册用户

**认证**: 否

**请求体**:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**请求示例**:

```bash
curl -X POST http://localhost:9527/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**响应示例** (成功):

```json
{
  "message": "user created",
  "email": "test@example.com"
}
```

**错误响应示例**:

```json
{
  "error": "email already registered"
}
```

**状态码**:

- `200` - 注册成功
- `400` - 邮箱已注册或请求格式错误
- `500` - 服务器错误

---

### 用户登录

**端点**: `POST /auth/login`

**描述**: 使用邮箱和密码登录，获取 JWT token

**认证**: 否

**请求体**:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**请求示例**:

```bash
curl -X POST http://localhost:9527/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**响应示例** (成功):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**错误响应示例**:

```json
{
  "error": "invalid email or password"
}
```

**状态码**:

- `200` - 登录成功
- `400` - 请求格式错误
- `401` - 邮箱或密码错误
- `500` - 服务器错误

---

## 收藏相关

### 添加收藏

**端点**: `POST /favorite`

**描述**: 添加鸡尾酒到收藏列表

**认证**: 是 (需要 JWT token)

**请求头**:

```
Authorization: Bearer <token>
```

**请求体**:

```json
{
  "cocktailId": "1"
}
```

**请求示例**:

```bash
curl -X POST http://localhost:9527/favorite \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "cocktailId": "1"
  }'
```

**响应示例** (成功):

```json
{
  "message": "ok"
}
```

**错误响应示例**:

```json
{
  "error": "already favorited"
}
```

**状态码**:

- `200` - 添加成功
- `400` - 已收藏或请求格式错误
- `401` - 未认证
- `500` - 服务器错误

---

### 获取收藏列表

**端点**: `GET /favorite`

**描述**: 获取当前用户的收藏列表

**认证**: 是 (需要 JWT token)

**请求头**:

```
Authorization: Bearer <token>
```

**请求示例**:

```bash
curl -X GET http://localhost:9527/favorite \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**:

```json
[
  {
    "cocktailId": "1",
    "createdAt": "2026-04-26T10:30:00Z"
  },
  {
    "cocktailId": "2",
    "createdAt": "2026-04-26T10:35:00Z"
  }
]
```

**状态码**:

- `200` - 获取成功
- `401` - 未认证
- `500` - 服务器错误

---

### 删除收藏

**端点**: `DELETE /favorite/:cocktailId`

**描述**: 从收藏列表中删除鸡尾酒

**认证**: 是 (需要 JWT token)

**请求头**:

```
Authorization: Bearer <token>
```

**路径参数**:

- `cocktailId` (string, required) - 要删除的鸡尾酒 ID

**请求示例**:

```bash
curl -X DELETE http://localhost:9527/favorite/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例** (成功):

```json
{
  "message": "ok"
}
```

**状态码**:

- `200` - 删除成功
- `401` - 未认证
- `500` - 服务器错误

---

## 错误处理

### 常见错误状态码

| 状态码 | 说明                       |
| ------ | -------------------------- |
| 200    | 请求成功                   |
| 400    | 请求格式错误或业务逻辑错误 |
| 401    | 未认证或认证失败           |
| 404    | 资源不存在                 |
| 500    | 服务器内部错误             |

### 错误响应格式

所有错误响应都遵循以下格式：

```json
{
  "error": "错误信息描述"
}
```

---

## CORS

此 API 启用了 CORS，允许来自任何来源的请求。

支持的 HTTP 方法:

- GET
- POST
- PUT
- DELETE
- OPTIONS

---

## 认证说明

### JWT Token

- **获取方式**: 调用 `/auth/login` 端点登录
- **有效期**: 24 小时
- **使用方式**: 在请求头中添加 `Authorization: Bearer <token>`
- **签名算法**: HS256

### 刷新 Token

如果 token 过期，需要重新调用 `/auth/login` 登录获取新 token。
