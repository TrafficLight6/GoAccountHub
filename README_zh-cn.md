<h1 align="center">GoAccountHub</h1>
<!-- <p align="center">
  <img src="./docs/assets/GAH.png" alt="Logo" width="100" height="100">
</p> -->

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.27.0-00ADD8?logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/PostgreSQL-17-4169E1?logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Vue-3.5.41-4FC08D?logo=vuedotjs&logoColor=white" alt="Vue">
  <img src="https://img.shields.io/badge/Node.js-24.19.0-5FA04E?logo=nodedotjs&logoColor=white" alt="Node.js">
  <img src="https://img.shields.io/badge/npm-12.0.2-CB3837?logo=npm&logoColor=white" alt="npm">
</p>

## 🌐 Language
[English](README.md) | [简体中文](README_zh-cn.md)

> ⚠️ **分支说明：** 当前是 `dev-integration` 分支，用于把 Web 界面与 API 服务合并为一体（构建好的前端会被嵌入 Go 二进制）。它仍在开发中，需要稳定版本请使用 `main` 分支。

## 📖 项目简介

GoAccountHub（简称 GAH）是一个用 Go 编写的用户中心。它可以存储用户元数据，并支持一个用户拥有多个子用户（本项目称为**角色**），同时提供 Vue 3 管理页面和面向第三方应用的 API。数据存储使用 PostgreSQL。

主要能力：

- **单文件部署** —— 构建好的 Web 界面通过 `go:embed` 嵌入二进制，`gah start` 即可在同一端口同时提供管理页面与 API，磁盘上不需要任何额外的静态文件。
- **管理员管理** —— 一个 root 管理员（保存在配置文件中，绕过所有权限校验）加上若干普通管理员，权限细粒度可配；支持新增、修改、删除、列表查询。
- **用户与角色管理** —— 用户可携带任意元数据（JSON / XML / YAML / 纯文本）；开启开关后一个用户可拥有多个角色。
- **应用密钥** —— 第三方应用以 `app_key` Cookie 使用的密钥。密钥用户名只能是可打印 ASCII 字符且创建后不可修改；完整的 key **只返回一次**，之后的响应中一律打码。
- **仪表盘** —— 管理员、用户、角色、密钥、登录令牌的数量统计，以及当前管理员的权限一览。

## 🧱 技术栈

| 层次 | 技术 |
| --- | --- |
| 后端 | Go 1.27、Gin、GORM、PostgreSQL |
| 前端 | Vue 3、Vite、Element Plus、vue-router、monaco-editor |
| 命令行 | urfave/cli v2 |

## ✅ 环境要求

- Go 1.27 或更高版本
- PostgreSQL（推荐 17）
- Node.js ≥ 20.19（22 / 24 LTS 均可，CI 使用 22）与 npm —— 仅在构建前端时需要

## 🚀 编译

**__目前只能通过编译安装，计划之后提供 docker-compose 文件。__**

最简单的方式是使用一键脚本：它会先构建前端，再把前端嵌入后端一起编译：

```bash
./build.cmd          # Windows（cmd 或双击运行）
./build.sh           # Linux / macOS / Git Bash
```

脚本会打印 `GAHversion` 中的版本号与最终二进制大小；任一步失败都会给出明确报错并中断。

手动等价操作 —— **顺序很重要**：

```bash
cd GAHFrontend
npm ci               # 或者 npm install
npm run build        # 产出 GAHFrontend/dist
cd ..
go build -o gah .
```

> ⚠️ `go build` 必须在 `npm run build` **之后**执行。`GAHFrontend/dist` 是在编译期嵌入的，顺序颠倒只会把占位文件 `GAHFrontend/dist/.gitkeep` 嵌进去（它被纳入版本控制，用于保证刚克隆下来也能编译）。此时服务不会返回界面，而是提示「前端未嵌入」。

## ▶️ 启动

1. 生成配置文件（默认文件名 `config.json`）：

   ```bash
   gah generate
   ```

   想要一份预填测试数据的配置，可以用 `gah generate-test`，它的 root 管理员默认密码是 `123`。

2. 设置 root 管理员密码：

   ```bash
   gah password <your_password>
   ```

   **__⚠️ `root` 管理员是系统最高权限账号，可以做任何操作。__**

3. 打开配置文件，填写数据库相关配置：

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

   | 字段 | 说明 |
   | --- | --- |
   | `port` | API 服务端口；同时也用于提供内嵌的 Web 界面 |
   | `database_name` | 数据库名称 |
   | `database_host` | 数据库地址 |
   | `database_port` | 数据库端口 |
   | `database_user` | 数据库用户名 |
   | `database_password` | 数据库密码 |
   | `frontend_port` | Vite 开发服务器端口；服务启动时会写入 `GAHFrontend/.env` |
   | `root_admin_password_hash` | root 管理员密码哈希，由命令行工具生成，无需手动填写 |
   | `root_admin_uu_hash` | root 管理员唯一哈希，由命令行工具生成，无需手动填写 |
   | `switch_config.allow_multi_character` | 是否允许账号拥有多个角色 |
   | `switch_config.allow_admin_logout` | 是否允许管理员登出（当前实现有问题，计划移除） |

4. 启动服务：

   ```bash
   gah start
   gah start -c ./config.json   # 自定义配置文件路径
   ```

5. 访问：

   - Web 界面：`http://localhost:<your_port>/`
   - API：`http://localhost:<your_port>/api/v1/...`

   使用第 2 步设置的密码以 `root` 账号登录。所有非 API 路由的路径都会回退到 Web 界面，因此前端路由可以直接打开或刷新。

> ℹ️ 启动时若表不存在会自动创建（`AutoMigrate`）；每次启动都会按配置创建/覆写 `GAHFrontend/.env`。

## 🎨 开发调试

进行前端开发时，保持 API 服务运行，再启动 Vite 开发服务器即可获得热更新，并把 `/api` 代理到后端：

```bash
cd GAHFrontend
npm run dev
```

- 开发服务器端口取自 `GAHFrontend/.env` 中的 `PORT`，该文件由 `gah start` 按 `frontend_port` 写入。
- 单文件部署不需要它：Go 进程自己托管前端。

## 🔑 权限模型

普通管理员带有一个 `permission` 对象，权限缺失或为 `false` 时返回 HTTP `403`：

| 权限 | 可访问 |
| --- | --- |
| `can_add_admin` | 新增管理员 |
| `can_delete_admin` | 删除管理员 |
| `can_edit_admin` | 修改管理员 |
| `can_get_admin` | 管理员数量 / 查询 / 列表 |
| `can_operate_user` | 全部 `/api/v1/user/*` 管理端接口 |
| `can_operate_character` | 全部 `/api/v1/character/*` 管理端接口 |
| `can_operate_app_key` | 全部 `/api/v1/key/*` 管理端接口 |

- root 管理员绕过所有权限校验，不能通过 API 删除或修改自己，也不会出现在管理员列表中。
- 管理员会话使用 Cookie（`admin_token`），勾选「记住我」时为 30 天，否则 1 小时。
- `/api/v1/app/*` 使用 `app_key` Cookie 鉴权，而不是管理员会话。

## 🗂️ 项目结构

```
GoAccountHub/
├── main.go, frontend.go        # 入口，以及对 GAHFrontend/dist 的 go:embed
├── adminControllor/            # 管理端接口（管理员、用户、角色、密钥、仪表盘）
├── appControllor/              # 第三方应用接口
├── checkControllor/            # 管理员 token 校验
├── middleware/                 # 配置/数据库注入、管理员鉴权、权限、app_key 校验
├── router/                     # 路由注册与内嵌前端处理
├── sql/                        # 数据库连接与查询
├── sqlTable/                   # GORM 模型
├── cliAction/                  # 命令行命令（start / password / generate / generate-test）
├── GAHFrontend/                # Vue 3 + Vite 管理页面
├── build.cmd, build.sh         # 一键编译脚本
└── GAHversion                  # 版本文件
```

## 📚 文档

- API 文档 —— [English](./docs/api.md) | [简体中文](./docs/api_zh-cn.md)
- 命令行文档 —— [English](./docs/cli.md) | [简体中文](./docs/cil_zh-cn.md)

## 🏗️ 构建与发布

GitHub Actions（`.github/workflows/go.yml`）在向 `main` 推送或提交 PR 时运行：

- **build** —— `go build ./...` 与 `go test ./...`
- **frontend-build** —— `npm ci` + `npm run build`，产物以 `frontend-dist` 上传
- **cross-build** —— 下载上述产物到 `GAHFrontend/dist`，再交叉编译 windows / linux / darwin 的 amd64 与 arm64；每个二进制以 `GoAccountHub-<GAHversion>-<os>-<arch>` 上传

## 📄 许可证

[MIT](./LICENSE)
