# GoAccountHub CLI 文档

GoAccountHub 使用 [urfave/cli v2](https://github.com/urfave/cli) 构建命令行界面。以下内容根据 `cliAction/` 目录下的代码生成。

- 应用名：`account_hub`
- 应用说明：`account hub example`

> 下文示例以编译出的二进制文件名 `GoAccountHub` 为例（可自行重命名或设置别名，`README` 中写作 `gah`）。

## 命令总览

| 命令 | 说明 |
| --- | --- |
| `start` | 启动 GoAccountHub 服务 |
| `password` | 修改 Root Admin 密码 |
| `generate` | 生成配置文件 |
| `generate-test` | 生成测试用配置文件 |
| `add-key` | 新增 Application Key |
| `del-key` | 删除 Application Key |

查看全局帮助与版本信息：

```bash
GoAccountHub --help
GoAccountHub <command> --help
```

---

## start

启动 GoAccountHub 服务。

**用法**

```bash
GoAccountHub start [选项]
```

**选项**

| 选项 | 别名 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--config` | `-c` | `config.json` | 配置文件路径 |

**行为**

1. 读取指定配置文件中配置。
2. 将前端相关配置写入 `./GAHFrontend/.env`：
   ```
   PORT=<frontend_port>
   VITE_API_PROXY_TARGET=http://127.0.0.1:<port>
   ```
3. 启动 API 服务。

> ⚠️ 注意：启动服务前需要确保前端根目录存在 `.env` 文件，服务启动时会按配置覆写该文件。

**示例**

```bash
GoAccountHub start
GoAccountHub start -c ./config.json
```

---

## password

修改 Root Admin 密码，密码以 SHA256 形式写入配置文件。

**用法**

```bash
GoAccountHub password <password> [选项]
```

**参数**

| 参数 | 是否必填 | 说明 |
| --- | --- | --- |
| `<password>` | 是 | 新的 Root Admin 密码；缺省时会提示 `⚠️ Password is Empty` 且不做修改 |

**选项**

| 选项 | 别名 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | 配置文件路径 |

**行为**

1. 计算 `root_admin_password_hash = SHA256(password)`。
2. 计算 `root_admin_uu_hash = SHA256("root" + passwordHash + Unix 时间戳)`。
3. 将以上两个字段写回配置文件。

**示例**

```bash
GoAccountHub password your_password
GoAccountHub password your_password -c ./config.json
```

---

## generate

生成一份空白配置文件。

**用法**

```bash
GoAccountHub generate [选项]
```

**选项**

| 选项 | 别名 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | 配置文件路径 |

**行为**

将空的 `Config` 结构序列化写入目标文件，所有字段均为零值，需手动填写。

**示例**

```bash
GoAccountHub generate
GoAccountHub generate -c ./config.json
```

---

## generate-test

生成一份预置好测试数据的配置文件，默认 Root Admin 密码为 `123`。

**用法**

```bash
GoAccountHub generate-test [选项]
```

**选项**

| 选项 | 别名 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | 配置文件路径 |

**行为**

写入如下预置配置：

| 字段 | 值 |
| --- | --- |
| `port` | `8080` |
| `database_name` | `account_hub` |
| `database_host` | `127.0.0.1` |
| `database_port` | `3306` |
| `database_user` | `postgres` |
| `database_password` | `postgres` |
| `root_admin_password_hash` | `SHA256("123")` |
| `root_admin_uu_hash` | `SHA256("root" + SHA256("123") + Unix 时间戳)` |
| `switch_config.allow_multi_character` | `false` |

执行后输出 `✅ Config File is Generated, Default Password is 123`。

**示例**

```bash
GoAccountHub generate-test
```

---

## add-key

新增一个 Application Key 并写入数据库。

**用法**

```bash
GoAccountHub add-key [选项]
```

**选项**

| 选项 | 别名 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | 配置文件路径 |
| `--site` | `-s` | `default` | 该 Key 对应的应用站点 |

**行为**

1. 读取配置并连接数据库（执行完成后自动关闭连接）。
2. 构造 `ApplicationKey`：
   - `application_using_site = <site>`
   - `key = SHA256(site + time.Now().String())`
3. 写入数据库，输出 `✅ Add Key Success`。

> ℹ️ 说明：生成的 Key 仅存入数据库，不会打印到终端，如需使用请从数据库中查询。

**示例**

```bash
GoAccountHub add-key
GoAccountHub add-key -s mysite
```

---

## del-key

删除数据库中的 Application Key。

**用法**

```bash
GoAccountHub del-key [选项]
```

**选项**

| 选项 | 别名 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | 配置文件路径 |
| `--key` | `-k` | `default` | 要删除的 Application Key |

**行为**

1. 读取配置并连接数据库（执行完成后自动关闭连接）。
2. 按 key 删除记录，输出 `✅ Delete Key Success`。

**示例**

```bash
GoAccountHub del-key -k <your_key>
```

---

## 配置文件字段说明

配置文件为 JSON 格式，结构定义见 `config/config.go`。

```json
{
    "port": "8081",
    "database_name": "GoAccountHub",
    "database_host": "127.0.0.1",
    "database_port": "5432",
    "database_user": "postgres",
    "database_password": "postgres",
    "frontend_port": "8082",
    "root_admin_password_hash": "a665a45920422f9d417......",
    "root_admin_uu_hash": "0fbaf45ee863c0......",
    "switch_config": {
        "allow_multi_character": true,
        "allow_admin_logout": false
    }
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `port` | string | API 服务端口 |
| `database_name` | string | 数据库名 |
| `database_host` | string | 数据库地址 |
| `database_port` | string | 数据库端口 |
| `database_user` | string | 数据库用户名 |
| `database_password` | string | 数据库密码 |
| `frontend_port` | string | 前端端口；启动服务时写入 `GAHFrontend/.env` |
| `root_admin_password_hash` | string | Root Admin 密码哈希，由 `password` / `generate-test` 生成，无需手填 |
| `root_admin_uu_hash` | string | Root Admin 唯一标识哈希，由 CLI 生成，无需手填 |
| `switch_config` | object | 系统开关配置 |
| `switch_config.allow_multi_character` | bool | 是否允许一个账号拥有多个角色 |
| `switch_config.allow_admin_logout` | bool | 是否允许管理员登出（开发有误，后续将移除） |

---

## 典型使用流程

```bash
# 1. 生成配置文件
GoAccountHub generate

# 2. 设置 Root Admin 密码
GoAccountHub password <your_password>

# 3. 编辑配置文件，填写数据库等参数（见上文配置字段说明）

# 4. 启动服务
GoAccountHub start
```

服务启动后访问 `http://localhost:<port>`，正常返回：

```json
{"code":200,"message":"Welcome to GoAccountHub"}
```
