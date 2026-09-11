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

## 🌐Language
[English](README.md) | [简体中文](README_zh-cn.md)

## 📖 简介
Go Account Hub（简称 GAH）是一个用 Go 语言编写的用户中心，可以存储用户的元数据，并支持多个「子用户」（在本项目中称为「角色」），同时提供一个 Vue 前端管理页面。

本项目需要一个 PostgreSQL 数据库来存储数据。

## 🛠️ 安装与启动
**目前只能通过编译安装。我们计划在未来添加 docker-compose 文件。**

- 📦 安装  
        1. 克隆本仓库。  
        2. 安装 PostgreSQL 数据库和 Node.js。  
        3. 在 `./GoAccountHub` 目录下执行 `go mod tidy` 以安装 Go 依赖。  
        4. 在 `./GoAccountHub/GAHFrontend` 目录下执行 `npm install` 以安装 Vue 依赖。  
        5. 在 `./GoAccountHub` 目录下编译 Go 项目，执行 `go build` 进行编译。  
- ▶️ 启动
    **⚠️注意：在启动 API 服务之前，你**必须**先启动 API 服务。同时请确保前端根目录下存在 `.env` 文件。因为 API 服务启动时会根据配置修改该 `.env` 文件。**

    + 🖥️ Api 服务
        1. 执行以下命令生成配置文件（默认配置名为 `config.json`）
        ```bash
        gah generate
        ```
        2. 执行以下命令设置 root 管理员密码：
        ```bash
        gah password <your_password>
        ```
        **__⚠️警告：GAH 的 `root` 管理员是最高权限管理员，可以操作系统中的任何内容。__**
        3. 打开配置文件，设置配置项：
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
        配置项说明：
        - `port`：服务的端口号。
        - `database_name`：数据库名称。
        - `database_host`：数据库主机地址。
        - `database_port`：数据库端口号。
        - `database_user`：数据库用户名。
        - `database_password`：数据库密码。
        - `frontend_port`：前端端口号。API 服务启动时会写入 `.env` 文件。

        - `root_admin_password_hash`：root 管理员的密码哈希。由 CLI 生成，无需手动设置。
        - `root_admin_uu_hash`：root 管理员的 uu 哈希。由 CLI 生成，无需手动设置。

        - `switch_config`：系统的开关配置。
        - `switch_config.allow_multi_character`：是否允许多角色。
        - `switch_config.allow_admin_logout`：是否允许管理员登出（开发有误，将会移除）。
        4. 执行以下命令启动 API 服务：
        ```bash
        gah start
        ```
        现在你可以通过 `http://localhost:<your_port>` 访问 GAH API。它将返回：
        ```json
        {"code":200,"message":"Welcome to GoAccountHub"}
        ```
    + 🎨 前端服务
        1. 执行以下命令启动前端服务：
        ```bash
        npm run dev
        ```
        现在你可以通过 `http://localhost:<your_frontend_port>` 访问 GAH 前端。

**__更多信息，请参阅 [`docs/api_zh-cn.md`](docs/api_zh-cn.md)。__**
