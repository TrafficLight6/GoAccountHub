# GoAccountHub 命令行文档

GoAccountHub 使用 [urfave/cli v2](https://github.com/urfave/cli) 构建命令行界面，以下内容根据 `cliAction/` 目录中的代码生成。

## 🌐Language
[English](cli.md) | [简体中文](cil_zh-cn.md)

- 应用名称：`account_hub`
- 应用用途：`account hub example`

> 下面的示例使用编译后的二进制名 `GoAccountHub`（你可以重命名或设置别名，`README` 中称其为 `gah`）。

## 命令总览

| 命令 | 说明 |
| --- | --- |
| `start` | 启动 GoAccountHub |
| `password` | 修改 root 管理员密码 |
| `generate` | 生成配置文件 |
| `generate-test` | 生成测试配置文件 |

查看全局帮助与版本信息：

```bash
GoAccountHub --help
GoAccountHub <command> --help
```

---

## start

启动 GoAccountHub 服务。同一个进程既提供 REST API，也提供编译时嵌入二进制的前端界面。

**用法**

```bash
GoAccountHub start [选项]
```

**选项**

| 选项 | 别名 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--config` | `-c` | `config.json` | 配置文件路径 |

**行为**

1. 读取指定配置文件中的配置。
2. 将前端相关配置写入 `./GAHFrontend/.env`（文件不存在会自动创建）：
   ```
   PORT=<frontend_port>
   VITE_API_PROXY_TARGET=http://127.0.0.1:<port>
   ```
3. 启动 API 服务，并自动创建缺失的数据库表（`AutoMigrate`）。
4. 在同一端口上托管编译进二进制的 Web 界面（`GAHFrontend/dist`），因此 `http://localhost:<port>` 即可打开管理前端。

> ℹ️ `.env` 仅供 Vite 开发服务器（`npm run dev`）使用。生产部署并不需要它：Go 进程自己托管前端。

> ⚠️ 请在 `go build` **之前**先在 `GAHFrontend` 下执行 `npm run build`，否则二进制里只有占位文件（仓库跟踪了 `GAHFrontend/public/.gitkeep`，保证刚克隆下来也能编译，每次前端构建都会把它复制进 `GAHFrontend/dist`）。这种情况下服务只会返回「前端未嵌入」的提示，而不是界面。

**示例**

```bash
GoAccountHub start
GoAccountHub start -c ./config.json
```

---

## password

修改 root 管理员密码，密码会以 SHA256 哈希写入配置文件。

**用法**

```bash
GoAccountHub password <密码> [选项]
```

**参数**

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `<密码>` | 是 | 新的 root 管理员密码；不传时会输出 `⚠️ Password is Empty` 且不做任何修改 |

**选项**

| 选项 | 别名 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | 配置文件路径 |

**行为**

1. 计算 `root_admin_password_hash = SHA256(password)`。
2. 计算 `root_admin_uu_hash = SHA256("root" + passwordHash + Unix 时间戳)`。
3. 把这两个字段写回配置文件，并输出 `✅ Password is Updated`。

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

将空的 `Config` 结构体序列化后写入目标文件，所有字段均为零值，需要手动填写。

执行后输出 `✅ Config File is Generated`。

**示例**

```bash
GoAccountHub generate
GoAccountHub generate -c ./config.json
```

---

## generate-test

生成一份预填测试数据的配置文件，默认 root 管理员密码为 `123`。

**用法**

```bash
GoAccountHub generate-test [选项]
```

**选项**

| 选项 | 别名 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | 配置文件路径 |

**行为**

写入如下预设配置：

| 字段 | 值 |
| --- | --- |
| `port` | `8080` |
| `database_name` | `account_hub` |
| `database_host` | `127.0.0.1` |
| `database_port` | `3306` |
| `database_user` | `postgres` |
| `database_password` | `postgres` |
| `frontend_port` | *（空值，需手动填写）* |
| `root_admin_password_hash` | `SHA256("123")` |
| `root_admin_uu_hash` | `SHA256("root" + SHA256("123") + Unix 时间戳)` |
| `switch_config.allow_multi_character` | `false` |
| `switch_config.allow_admin_logout` | `false` |

执行后输出 `✅ Config File is Generated, Default Password is 123`。

**示例**

```bash
GoAccountHub generate-test
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
| `port` | string | API 服务端口；同时也用于提供内嵌的 Web 界面 |
| `database_name` | string | 数据库名称 |
| `database_host` | string | 数据库地址 |
| `database_port` | string | 数据库端口 |
| `database_user` | string | 数据库用户名 |
| `database_password` | string | 数据库密码 |
| `frontend_port` | string | Vite 开发服务器端口；服务启动时会写入 `GAHFrontend/.env` |
| `root_admin_password_hash` | string | root 管理员密码哈希，由 `password` / `generate-test` 生成，无需手动填写 |
| `root_admin_uu_hash` | string | root 管理员唯一哈希，由命令行工具生成，无需手动填写 |
| `switch_config` | object | 系统开关配置 |
| `switch_config.allow_multi_character` | bool | 是否允许账号拥有多个角色 |
| `switch_config.allow_admin_logout` | bool | 是否允许管理员登出（当前实现有问题，计划移除） |

---

## 典型流程

```bash
# 1. 生成配置文件
GoAccountHub generate

# 2. 设置 root 管理员密码
GoAccountHub password <your_password>

# 3. 打开配置文件，填写数据库等参数（见上方「配置文件字段说明」）

# 4. 先构建前端，再编译服务端（或直接执行 `build.cmd` / `build.sh`）
cd GAHFrontend && npm ci && npm run build && cd ..
go build -o GoAccountHub .

# 5. 启动服务
GoAccountHub start
```

服务启动后，用浏览器打开 `http://localhost:<port>` 即可使用管理前端，API 位于 `http://localhost:<port>/api/v1/...`。使用第 2 步设置的密码登录 root 管理员。
