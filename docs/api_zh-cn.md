# 🚀 GoAccountHub API 参考
**__由 AI 生成，如有错误，请在 issue 中反馈。__**

## 🌐Language
[English](api.md) | [简体中文](api_zh-cn.md)

**GoAccountHub (GAH)** 的 REST API —— 一个用 Go 编写的用户中心，可存储用户元数据，并支持多个子用户（在本项目中称为「角色」）。

- **基础 URL：** `http://<host>:<port>`
- **API 版本前缀：** `/api/v1`
- **内容类型：** `application/json`（请求体均为 JSON，包括 `DELETE` / `PUT`）
- **鉴权模型：** 基于 Cookie（管理员使用 `admin_token`，第三方应用使用 `app_key`）+ 应用用户登录后返回的 `token`

## 📑 目录

- [约定](#-约定)
- [Cookie 与鉴权模型](#-cookie-与鉴权模型)
- [权限模型](#-权限模型)
- [快速参考](#-快速参考)
- [🌐 根路由](#-根路由)
- [🛡️ 管理员 API](#️-管理员-api)
- [📱 应用 API](#-应用-api)

## 🧭 约定

| 主题 | 规则 |
| --- | --- |
| 成功 | HTTP `200`，响应体为 `{"code": 200, ...}` |
| 错误 | 非 2xx 状态码，响应体为 `{"code": <status>, "error": "<reason>"}`（部分端点使用 `"message"` 而非 `"error"`） |
| 布尔条件 | 搜索条件中为空 / `false` 的字段会被忽略 |
| 分页 | `begin_table_id` 从 **1** 开始；`length = -1` 表示**返回全部**；`offset = (begin_table_id - 1) * length` |
| 对象字段 | `User` / `Character` / `Admin` 对象序列化时使用 **PascalCase** 键名（如 `Username`、`UUHash`、`MetaData`）；`Permission` 使用 **snake_case**（如 `can_add_admin`） |

## 🍪 Cookie 与鉴权模型

| Cookie | 设置方 | 使用者 | 含义 |
| --- | --- | --- | --- |
| `admin_token` | `POST /api/v1/admin/login` | 所有 `/api/v1/*` 管理员端点 | 管理员会话令牌 |
| `admin_name` | `POST /api/v1/admin/login` | 仅前端展示 | 当前管理员用户名 |
| `app_key` | 外部（线下发放） | 所有 `/api/v1/app/*` 端点，**两个 `metadata` 端点除外** | 第三方应用密钥 |

> `is_remember = true` → Cookie / 会话有效期 **30 天**；否则为 **1 小时**（管理员令牌为会话 Cookie）。

## 🔑 权限模型

普通管理员携带一个 `permission` 对象：

```json
{
  "can_add_admin": false,
  "can_delete_admin": false,
  "can_edit_admin": false,
  "can_get_admin": false,
  "can_operate_user": false,
  "can_operate_character": false
}
```

- 内置的 `root` 管理员存储于配置文件中（而非数据库），**绕过所有权限检查**。
- 权限字段缺失或为 `false` → HTTP `403` `Permission Denied`。
- `admin_token` 无效 / 过期 → HTTP `401`。

| 权限 | 可访问的端点 |
| --- | --- |
| `can_get_admin` | `POST /admin/count`、`POST /admin/get`、`POST /admin/range` |
| `can_add_admin` | `POST /admin/add` |
| `can_delete_admin` | `DELETE /admin/delete` |
| `can_edit_admin` | `PUT /admin/edit` |
| `can_operate_user` | 所有 `/api/v1/user/*` 管理员端点 |
| `can_operate_character` | 所有 `/api/v1/character/*` 管理员端点 |

## 📋 快速参考

| 方法 | 路径 | 🍪 Cookie | 🔑 权限 |
| --- | --- | --- | --- |
| `GET` | `/` | ❌ | — |
| `POST` | `/api/v1/admin/login` | ❌ | — |
| `DELETE` | `/api/v1/admin/logout` | ✅ `admin_token` | — |
| `POST` | `/api/v1/info` | ✅ `admin_token` | — |
| `POST` | `/api/v1/admin/check_token` | ✅ `admin_token` | — |
| `POST` | `/api/v1/admin/info` | ✅ `admin_token` | — |
| `POST` | `/api/v1/admin/count` | ✅ `admin_token` | `can_get_admin` |
| `POST` | `/api/v1/admin/add` | ✅ `admin_token` | `can_add_admin` |
| `DELETE` | `/api/v1/admin/delete` | ✅ `admin_token` | `can_delete_admin` |
| `PUT` | `/api/v1/admin/edit` | ✅ `admin_token` | `can_edit_admin` |
| `POST` | `/api/v1/admin/get` | ✅ `admin_token` | `can_get_admin` |
| `POST` | `/api/v1/admin/range` | ✅ `admin_token` | `can_get_admin` |
| `POST` | `/api/v1/user/count` | ✅ `admin_token` | `can_operate_user` |
| `POST` | `/api/v1/user/add` | ✅ `admin_token` | `can_operate_user` |
| `DELETE` | `/api/v1/user/delete` | ✅ `admin_token` | `can_operate_user` |
| `PUT` | `/api/v1/user/edit` | ✅ `admin_token` | `can_operate_user` |
| `POST` | `/api/v1/user/get` | ✅ `admin_token` | `can_operate_user` |
| `POST` | `/api/v1/user/range` | ✅ `admin_token` | `can_operate_user` |
| `POST` | `/api/v1/character/count` | ✅ `admin_token` | `can_operate_character` + 多角色开关 |
| `POST` | `/api/v1/character/add` | ✅ `admin_token` | `can_operate_character` + 多角色开关 |
| `DELETE` | `/api/v1/character/delete` | ✅ `admin_token` | `can_operate_character` + 多角色开关 |
| `PUT` | `/api/v1/character/edit` | ✅ `admin_token` | `can_operate_character` + 多角色开关 |
| `POST` | `/api/v1/character/get` | ✅ `admin_token` | `can_operate_character` + 多角色开关 |
| `POST` | `/api/v1/character/range` | ✅ `admin_token` | `can_operate_character` + 多角色开关 |
| `POST` | `/api/v1/app/user/login` | ✅ `app_key` | — |
| `DELETE` | `/api/v1/app/user/logout` | ✅ `app_key` | — |
| `POST` | `/api/v1/app/user/add` | ✅ `app_key` | — |
| `DELETE` | `/api/v1/app/user/delete` | ✅ `app_key` | — |
| `PUT` | `/api/v1/app/user/edit` | ✅ `app_key` | — |
| `POST` | `/api/v1/app/user/get` | ✅ `app_key` | — |
| `POST` | `/api/v1/app/user/range` | ✅ `app_key` | — |
| `POST` | `/api/v1/app/character/add` | ✅ `app_key` | — |
| `DELETE` | `/api/v1/app/character/delete` | ✅ `app_key` | — |
| `PUT` | `/api/v1/app/character/edit` | ✅ `app_key` | — |
| `POST` | `/api/v1/app/character/get` | ✅ `app_key` | — |
| `POST` | `/api/v1/app/character/range` | ✅ `app_key` | — |
| `POST` | `/api/v1/app/user/get/metadata` | ❌ | — |
| `POST` | `/api/v1/app/character/get/metadata` | ❌ | — |

---

## 🌐 根路由

### `GET` `/`

健康检查 / 欢迎检查。

- 🍪 **是否需要 Cookie：** 否
- 🔑 **所需权限：** —

**响应**

```json
{ "code": 200, "message": "Welcome to GoAccountHub" }
```

---

## 🛡️ 管理员 API

服务端管理 API。登录后，浏览器持有 `admin_token` Cookie。

### 🔐 鉴权

#### `POST` `/api/v1/admin/login`

以 `root` 或普通管理员身份登录。成功后设置 `admin_token`（以及 `admin_name`）Cookie。

- 🍪 **是否需要 Cookie：** 否
- 🔑 **所需权限：** —

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `username` | string | ✅ | 管理员用户名（root 管理员为 `root`） |
| `password` | string | ✅ | 明文密码（服务端使用 SHA-256 哈希） |
| `is_remember` | bool | ❌ | `true` = 30 天会话，`false` / 省略 = 1 小时会话 |

```json
{ "username": "root", "password": "secret", "is_remember": true }
```

**响应**

```json
{ "code": 200, "message": "login success", "token": "<token>", "admin": "root" }
```

> 凭据无效 → `400` `{"code":400,"message":"invalid credentials"}`。

#### `DELETE` `/api/v1/admin/logout`

使当前 `admin_token` 失效，并从数据库中删除。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** —

**请求体：** 无

**响应**

```json
{ "code": 200, "message": "Logout success" }
```

#### `POST` `/api/v1/info`

任何已登录管理员均可查看的仪表盘统计数据。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** —

**请求体：** 无

**响应**

```json
{
  "code": 200,
  "data": {
    "user_count": 10,
    "character_count": 25,
    "total_token_count": 7,
    "user_token_count": 5,
    "admin_token_count": 2,
    "admin_count": 3
  }
}
```

> 若已配置，`admin_count` 会包含 root 管理员。

#### `POST` `/api/v1/admin/check_token`

校验当前 `admin_token`。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** —

**请求体：** 无

**响应**

```json
{ "code": 200, "message": "admin_token valid" }
```

> 令牌缺失 / 无效 → `400`。

#### `POST` `/api/v1/admin/info`

返回当前管理员的身份与权限。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** —

**请求体：** 无

**响应（普通管理员）**

```json
{
  "code": 200,
  "message": "Admin Info",
  "data": {
    "username": "alice",
    "uu_hash": "<uu_hash>",
    "is_root": false,
    "permission": { "can_add_admin": true, "can_operate_user": true }
  }
}
```

**响应（root 管理员）** —— `is_root: true`，且所有权限均为 `true`。

### 👤 管理员管理

#### `POST` `/api/v1/admin/count`

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_get_admin`
- **请求体：** 无

```json
{ "code": 200, "data": { "admin_count": 3 } }
```

#### `POST` `/api/v1/admin/add`

创建一个新的（非 root）管理员。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_add_admin`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `username` | string | ✅ | 不能为 `root`，且必须唯一 |
| `password` | string | ✅ | 明文密码 |
| `permission` | object | ✅ | 权限映射（见[权限模型](#-权限模型)） |

```json
{
  "username": "alice",
  "password": "p@ssw0rd",
  "permission": { "can_add_admin": true, "can_operate_user": true }
}
```

**响应**

```json
{ "code": 200, "message": "Admin Added", "admin": { "ID": 2, "Username": "alice", "UUHash": "<uu_hash>", "Permission": { } } }
```

> 用户名为 `root` 会被拒绝；用户名重复或用户名 / 密码为空 → `400`。

#### `DELETE` `/api/v1/admin/delete`

根据 `uu_hash` 删除管理员。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_delete_admin`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `uu_hash` | string | ✅ | 目标管理员的 `uu_hash` |

```json
{ "uu_hash": "<target_uu_hash>" }
```

**响应**

```json
{ "code": 200, "message": "Delete Admin Success" }
```

> root 管理员不可删除，管理员也不能删除自己 → `400`。

#### `PUT` `/api/v1/admin/edit`

修改管理员的密码和 / 或权限。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_edit_admin`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `uu_hash` | string | ✅ | 目标管理员的 `uu_hash` |
| `admin_info.password` | string | ❌ | 新密码；为空表示保持不变 |
| `admin_info.permission` | object | ❌ | 新的权限映射 |

```json
{
  "uu_hash": "<target_uu_hash>",
  "admin_info": {
    "password": "new-password",
    "permission": { "can_operate_user": true }
  }
}
```

**响应**

```json
{ "code": 200, "message": "Success" }
```

> root 管理员不可修改，管理员也不能修改自己 → `400`。

#### `POST` `/api/v1/admin/get`

根据 `uu_hash` 获取单个管理员。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_get_admin`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `uu_hash` | string | ✅ | 目标管理员的 `uu_hash` |

**响应（普通管理员）**

```json
{ "code": 200, "admin": { "ID": 2, "Username": "alice", "UUHash": "<uu_hash>", "Permission": { } }, "is_root": false }
```

**响应（root 管理员）** —— `{ "code": 200, "admin": null, "is_root": true }`
未找到 → `404`。

#### `POST` `/api/v1/admin/range`

分页 / 过滤的管理员列表。**不包含 root 管理员。**

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_get_admin`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `begin_table_id` | int | ✅ | 页码，从 `1` 开始（必须 `>= 0`） |
| `length` | int | ✅ | 每页数量；`-1` 表示返回全部（必须 `>= -1`） |
| `search_condition` | object | ❌ | 过滤条件（见下，空字段忽略，使用 `AND` 组合） |

**`search_condition`**

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `username` | string | 模糊匹配（`LIKE %value%`） |
| `uu_hash` | string | 模糊匹配（`LIKE %value%`） |
| `can_add_admin` | bool | 为 `true` 时仅返回**拥有**该权限的管理员 |
| `can_delete_admin` | bool | 同上 |
| `can_edit_admin` | bool | 同上 |
| `can_get_admin` | bool | 同上 |
| `can_operate_user` | bool | 同上 |
| `can_operate_character` | bool | 同上 |

```json
{
  "begin_table_id": 1,
  "length": 20,
  "search_condition": { "username": "ali", "can_get_admin": true }
}
```

**响应**

```json
{ "code": 200, "message": "Success(not include root admin)", "data": [ { "ID": 2, "Username": "alice" } ] }
```

### 🧑 用户管理

#### `POST` `/api/v1/user/count`

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_user`
- **请求体：** 无

```json
{ "code": 200, "data": { "user_count": 10 } }
```

#### `POST` `/api/v1/user/add`

创建用户。会自动创建一个**同名角色**。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_user`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `username` | string | ✅ | 唯一 |
| `password` | string | ✅ | 明文密码 |
| `meta_data` | string | ❌ | 任意元数据（JSON / XML / YAML / 纯文本） |

**响应**

```json
{ "code": 200, "message": "User Added Successfully,And Same Name Character Added Successfully", "user": { "ID": 1, "Username": "bob", "UUHash": "<uu_hash>", "MetaData": "" } }
```

#### `DELETE` `/api/v1/user/delete`

删除用户**及其所属的所有角色**。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_user`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `uu_hash` | string | ✅ | 目标用户的 `uu_hash` |

**响应**

```json
{ "code": 200, "message": "User Deleted" }
```

#### `PUT` `/api/v1/user/edit`

修改用户。字段为空时保持原值。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_user`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `uu_hash` | string | ✅ | 目标用户的 `uu_hash` |
| `username` | string | ❌ | 新用户名（会同时重命名同名角色） |
| `password` | string | ❌ | 新密码 |
| `meta_data` | string | ❌ | 新元数据 |

**响应**

```json
{ "code": 200, "message": "Success" }
```

#### `POST` `/api/v1/user/get`

获取单个用户及其拥有的角色数量。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_user`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `uu_hash` | string | ✅ | 目标用户的 `uu_hash` |

**响应**

```json
{ "code": 200, "message": "Success", "user": { "ID": 1, "Username": "bob" }, "character_number": 2 }
```

#### `POST` `/api/v1/user/range`

分页 / 过滤的用户列表。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_user`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `begin_table_id` | int | ✅ | 页码，从 `1` 开始 |
| `length` | int | ✅ | 每页数量；`-1` 表示返回全部 |
| `search_condition` | object | ❌ | 过滤条件（见下） |

**`search_condition`**

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `username` | string | 模糊匹配（`LIKE %value%`） |
| `uu_hash` | string | 模糊匹配（`LIKE %value%`） |

**响应**

```json
{ "code": 200, "message": "Success", "data": [ { "ID": 1, "Username": "bob" } ] }
```

### 🎭 角色管理

> ⚠️ 当配置开关 **`allow_multi_character`** 为 `true` 时，这些端点（以及下文整个 `/api/v1/app/*` 分组）才可访问。否则 → `403` `Server Not Allow Multi Character`。

#### `POST` `/api/v1/character/count`

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_character`

```json
{ "code": 200, "data": { "character_count": 25 } }
```

#### `POST` `/api/v1/character/add`

在已存在的用户下新增角色。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_character`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `character_name` | string | ✅ | 全局唯一 |
| `password` | string | ✅ | 明文密码 |
| `user_uu_hash` | string | ✅ | 所属用户的 `uu_hash`（必须存在） |
| `meta_data` | string | ❌ | 任意元数据 |

**响应**

```json
{ "code": 200, "message": "Character Added Successfully", "character": { "ID": 5, "CharacterName": "bob_alt", "UserUUHash": "<user_uu_hash>", "UUHash": "<uu_hash>" } }
```

#### `DELETE` `/api/v1/character/delete`

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_character`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `user_uu_hash` | string | ✅ | 所属用户的 `uu_hash` |
| `character_uu_hash` | string | ✅ | 目标角色的 `uu_hash` |

**响应**

```json
{ "code": 200, "message": "Character Deleted" }
```

> 同名角色（`character_uu_hash == user_uu_hash`）不能在此删除；请改为删除该用户。

#### `PUT` `/api/v1/character/edit`

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_character`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `user_uu_hash` | string | ✅ | 所属用户的 `uu_hash` |
| `character_uu_hash` | string | ✅ | 目标角色的 `uu_hash` |
| `character_name` | string | ❌ | 新名称（不能冲突） |
| `password` | string | ❌ | 新密码 |
| `meta_data` | string | ❌ | 新元数据 |

**响应**

```json
{ "code": 200, "message": "Edit Character Success" }
```

> 同名角色必须通过用户相关端点进行修改。

#### `POST` `/api/v1/character/get`

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_character`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `user_uu_hash` | string | ✅ | 所属用户的 `uu_hash` |
| `character_uu_hash` | string | ✅ | 目标角色的 `uu_hash` |

**响应**

```json
{ "code": 200, "character": { "ID": 5, "CharacterName": "bob_alt", "UserUUHash": "<user_uu_hash>" } }
```

#### `POST` `/api/v1/character/range`

分页 / 过滤的角色列表。

- 🍪 **是否需要 Cookie：** 是 —— `admin_token`
- 🔑 **所需权限：** `can_operate_character`

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `begin_table_id` | int | ✅ | 页码，从 `1` 开始 |
| `length` | int | ✅ | 每页数量；`-1` 表示返回全部 |
| `search_condition` | object | ❌ | 过滤条件（见下） |

**`search_condition`**

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `character_name` | string | 模糊匹配（`LIKE %value%`） |
| `user_uu_hash` | string | 模糊匹配（`LIKE %value%`） |
| `uu_hash` | string | 模糊匹配（`LIKE %value%`） |

**响应**

```json
{ "code": 200, "message": "Success", "data": [ { "ID": 5, "CharacterName": "bob_alt" } ] }
```

---

## 📱 应用 API

位于 `/api/v1/app` 下的第三方应用 API。所有端点都要求 **`app_key`** Cookie —— 两个 `metadata` 端点除外，它们使用请求体中的登录 `token` 进行鉴权。

> ⚠️ 由于分组中间件继承机制，**所有** `/api/v1/app/*` 端点也会受 **`allow_multi_character`** 开关限制（包括用户相关端点）。

### 🔐 鉴权

#### `POST` `/api/v1/app/user/login`

以指定角色登录用户。返回 `token`（用于应用用户自身的操作以及 metadata 端点）。

- 🍪 **是否需要 Cookie：** 是 —— `app_key`
- 🔑 **所需权限：** —

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `username` | string | ✅ | 用户名 |
| `character_name` | string | ✅ | 属于该用户的角色 |
| `password` | string | ✅ | 用户密码（明文） |
| `is_remember` | bool | ❌ | `true` = 30 天令牌，否则 1 小时 |

```json
{ "username": "bob", "character_name": "bob", "password": "secret", "is_remember": true }
```

**响应**

```json
{ "code": 200, "message": "Login success", "token": "<token>" }
```

> 用户不存在 / 密码错误 / 角色名不存在 → `400`。

#### `DELETE` `/api/v1/app/user/logout`

使用户 `token` 失效。

- 🍪 **是否需要 Cookie：** 是 —— `app_key`
- 🔑 **所需权限：** —

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `token` | string | ✅ | 应用登录返回的令牌 |

**响应**

```json
{ "code": 200, "message": "Logout success" }
```

### 🧑 用户端点（应用）

与对应的管理员端点使用相同的处理函数，但改用 `app_key` 校验，而非 `admin_token` + 权限。

| 方法 | 路径 | 🍪 Cookie | 请求体 |
| --- | --- | --- | --- |
| `POST` | `/api/v1/app/user/add` | `app_key` | `username`、`password`、`meta_data` |
| `DELETE` | `/api/v1/app/user/delete` | `app_key` | `uu_hash` |
| `PUT` | `/api/v1/app/user/edit` | `app_key` | `uu_hash`、`username?`、`password?`、`meta_data?` |
| `POST` | `/api/v1/app/user/get` | `app_key` | `uu_hash` |
| `POST` | `/api/v1/app/user/range` | `app_key` | `begin_table_id`、`length`、`search_condition{ username, uu_hash }` |

### 🎭 角色端点（应用）

| 方法 | 路径 | 🍪 Cookie | 请求体 |
| --- | --- | --- | --- |
| `POST` | `/api/v1/app/character/add` | `app_key` | `character_name`、`password`、`user_uu_hash`、`meta_data` |
| `DELETE` | `/api/v1/app/character/delete` | `app_key` | `user_uu_hash`、`character_uu_hash` |
| `PUT` | `/api/v1/app/character/edit` | `app_key` | `user_uu_hash`、`character_uu_hash`、`character_name?`、`password?`、`meta_data?` |
| `POST` | `/api/v1/app/character/get` | `app_key` | `user_uu_hash`、`character_uu_hash` |
| `POST` | `/api/v1/app/character/range` | `app_key` | `begin_table_id`、`length`、`search_condition{ character_name, user_uu_hash, uu_hash }` |

### 📦 Metadata 端点

> 🍪 这两个端点**不**要求 `app_key`；它们使用请求体中的应用用户 `token` 进行鉴权。

#### `POST` `/api/v1/app/user/get/metadata`

返回 token 所属用户的元数据。

- 🍪 **是否需要 Cookie：** 否
- 🔑 **所需权限：** —

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `token` | string | ✅ | 应用用户令牌 |

**响应**

```json
{ "code": 200, "message": "Success", "meta_data": "<metadata>" }
```

#### `POST` `/api/v1/app/character/get/metadata`

返回 token 所属角色的元数据。

- 🍪 **是否需要 Cookie：** 否
- 🔑 **所需权限：** —

**请求体**

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `token` | string | ✅ | 应用用户令牌 |

**响应**

```json
{ "code": 200, "message": "success", "meta_data": "<metadata>" }
```
