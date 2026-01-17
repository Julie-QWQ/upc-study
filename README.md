# DocHub

> 学院学习资料托管平台 - 面向学院的学习资料共享与管理平台

## 📘 项目简介

DocHub 是面向学院的学习资料管理平台,提供资料上传下载、审核与检索服务。项目采用前后端分离的微服务架构,支持多角色权限管理(学生、学委、管理员),确保内容质量与合规性。平台目前支持约 1,200 用户。

### ✅ 核心功能

- **用户认证**: 注册登录、基于角色的权限控制(RBAC)
- **资料管理**: 资料上传、下载、分类、标签
- **审核流程**: 学委身份申请、资料内容审核、违规举报处理
- **检索与推荐**: 全文搜索、多条件筛选、热门资料推荐
- **通知系统**: 审核结果通知、系统公告、站内消息
- **数据统计**: 上传下载统计、用户活跃度分析

## 🛠 技术栈

| 层次 | 技术选型 | 说明 |
|------|----------|------|
| **前端** | Vue 3 + TypeScript + Vite | Composition API + `<script setup>` |
| **UI 框架** | Element Plus | Vue 3 组件库 |
| **状态管理** | Pinia | 官方状态管理方案 |
| **路由** | Vue Router 4 | 客户端路由 |
| **HTTP** | Axios | Promise based HTTP client |
| **验证** | VeeValidate + Yup | 表单验证 |
| **后端** | Go 1.21 + Gin | 高性能 HTTP 框架 |
| **ORM** | GORM | Go ORM 库 |
| **数据库** | PostgreSQL 15 | 关系型数据库,支持全文搜索 |
| **缓存** | Redis 7 | 会话管理、JWT 黑名单、热点数据缓存 |
| **认证** | JWT (golang-jwt/jwt/v5) | 无状态认证 |
| **配置** | Viper (YAML) | 配置管理 |
| **日志** | Zap + Lumberjack | 结构化日志与日志轮转 |
| **存储** | MinIO / 阿里云 OSS | 对象存储,支持预签名 URL |
| **部署** | Docker / Nginx + systemd | 容器化部署或传统部署 |

## 📁 项目结构

```
DocHub/
├── backend/                # Go 后端服务
│   ├── cmd/server/         # 程序入口
│   ├── configs/            # YAML 配置文件
│   ├── internal/           # 私有代码
│   │   ├── handler/        # HTTP 请求处理器
│   │   ├── service/        # 业务逻辑层
│   │   ├── repository/     # 数据访问层
│   │   ├── model/          # 数据模型
│   │   └── pkg/            # 工具包
│   ├── logs/               # 日志目录
│   ├── migrations/         # 数据库迁移
│   ├── go.mod
│   └── go.sum
├── frontend/               # Vue 3 前端应用
│   ├── src/
│   │   ├── api/            # API 接口封装
│   │   ├── assets/         # 静态资源
│   │   ├── components/     # 公共组件
│   │   ├── views/          # 页面组件
│   │   ├── router/         # 路由配置
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── composables/    # 组合式函数
│   │   └── utils/          # 工具函数
│   ├── dist/               # 构建产物
│   ├── package.json
│   └── vite.config.ts
├── docker/                 # Docker 编排文件
│   └── docker-compose.yml
├── scripts/                # 运维/辅助脚本
│   ├── dev.sh / dev.bat    # Docker 开发环境启动
│   ├── build.sh            # 构建脚本
│   ├── install.sh          # 依赖安装
│   ├── docker-up.bat       # Windows Docker 启动
│   └── docker-down.bat     # Windows Docker 停止
└── README.md
```

### 后端分层架构

严格遵循三层架构原则:

```
Handler (HTTP 层) → Service (业务逻辑层) → Repository (数据访问层)
```

- **Handler 层**: 处理 HTTP 请求/响应,参数验证,调用 Service
- **Service 层**: 实现所有业务逻辑,事务管理
- **Repository 层**: 数据库 CRUD 操作,不包含业务逻辑

## 🚀 快速启动

### Windows 下 Docker 开发环境

使用 scripts/docker-up.bat 一键启动所有服务（PostgreSQL、Redis、MinIO）

### 本地开发

#### 1. 环境准备

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- MinIO

#### 2. 后端配置

```bash
cd backend
cp configs/config.dev.yaml configs/config.dev.local.yaml
# 编辑 configs/config.yaml 配置数据库、Redis、OSS 等
go run cmd/server/main.go
```

后端服务将运行在 `http://localhost:8080`

#### 3. 前端启动

```bash
cd frontend
npm install          # 安装依赖
npm run dev          # 开发服务器
npm run lint         # ESLint 检查
npm run format       # Prettier 格式化
```

前端开发服务器将运行在 `http://localhost:3000`

## 🧩 生产部署

环境准备:

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- 阿里云 OSS 服务

### 1. 构建应用

```bash
# 构建后端
cd backend
go build -o upc-study-server cmd/server/main.go

# 构建前端
cd frontend
npm ci
npm run build
```

### 2. 配置后端服务

创建 systemd 服务文件 `/etc/systemd/system/upc-study.service`:

```ini
[Unit]
Description=UPC-STUDY Backend Service
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=~/upc-study/backend
ExecStart=~/upc-study/backend/upc-study-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务:

```bash
sudo systemctl daemon-reload
sudo systemctl enable upc-study
sudo systemctl start upc-study
sudo systemctl status upc-study
```

健康检查: `http://127.0.0.1:8080/health`

### 3. 配置 Nginx 反向代理

创建站点配置 `/etc/nginx/conf.d/upc-study.conf`:

```nginx
server {
    listen 80;
    server_name your.domain.com;

    # 前端静态文件
    root ~/upc-study/frontend/dist;
    index index.html;

    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 后端 API 代理
    location /api/v1/ {
        proxy_pass http://127.0.0.1:8080/api/v1/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

重载 Nginx:

```bash
sudo nginx -t && sudo systemctl reload nginx
```

## 🔐 核心设计

### 角色权限系统(RBAC)

- **Student**(学生): 浏览、下载、收藏资料
- **Committee**(学委): 上传资料、查看审核状态
- **Admin**(管理员): 审核资料、处理申请、管理用户

### 文件上传策略

- 使用预签名 URL 直接上传到 OSS
- 文件不经过后端服务器,避免性能瓶颈
- 后端在文件上传完成后记录元数据

### 缓存策略

- **会话管理**: 用户 Session 存储
- **JWT 黑名单**: 用户登出时将 Token 加入黑名单
- **热门资料**: TTL 5-10 分钟
- **用户信息**: TTL 10 分钟
- **搜索结果**: 减少数据库查询压力

### 审核流程

1. 学委身份需要申请审核
2. 资料上传后需要管理员审核
3. 举报内容由管理员处理
4. 所有审核操作记录审计日志

## 📊 数据库设计

核心表结构:

- **users**: 用户表,包含角色字段
- **materials**: 资料表,支持全文搜索(tsvector)
- **committee_applications**: 学委申请表
- **review_records**: 审核记录表
- **notifications**: 通知表
- **favorites**: 收藏表
- **download_records**: 下载记录表
- **reports**: 举报表

所有表都有 `created_at` 和 `updated_at` 字段。

详细设计见: [数据库设计文档](docs/02-数据库设计.md)

## 🌐 API 设计

RESTful API 设计,统一响应格式:

**成功响应**:

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

**错误响应**:

```json
{
  "code": 10001,
  "message": "参数错误",
  "data": null
}
```

详细 API 文档见: [API 接口设计](docs/03-API接口设计.md)

## 🔒 安全措施

1. **认证**: JWT token + Redis 黑名单
2. **授权**: 基于角色的访问控制中间件
3. **文件安全**: 类型限制、大小限制、文件名清理
4. **限流**: 登录尝试限制、下载频率限制
5. **输入验证**: Go validator 服务端验证 + VeeValidate 客户端验证

## 📈 项目进度

当前处于 **Phase 1: 基础框架搭建** 阶段。

核心模块完成度:

- 基础框架: 30%
- 用户认证: 0%
- 资料管理: 0%
- 审核流程: 0%
- 检索推荐: 0%
- 通知系统: 0%
- 管理后台: 0%

详细开发计划见: [项目开发计划](docs/项目开发计划.md)

## 🧰 脚本说明

### 构建脚本

- [`scripts/build.sh`](scripts/build.sh) - 构建前后端
- [`scripts/install.sh`](scripts/install.sh) - 安装前端依赖

### Docker 脚本

- [`scripts/dev.sh`](scripts/dev.sh) / [`scripts/dev.bat`](scripts/dev.bat) - Docker 开发环境启动
- [`scripts/docker-up.bat`](scripts/docker-up.bat) - Windows Docker 服务启动
- [`scripts/docker-down.bat`](scripts/docker-down.bat) - Windows Docker 服务停止

### 系统配置脚本

- `scripts/setup-postgres.sh` - 初始化 PostgreSQL 用户与数据库
- `scripts/setup-redis.sh` - 设置 Redis 密码并限制本机访问
- `scripts/setup-nginx.sh` - 生成 Nginx 站点配置并热加载

## 📚 文档

详细设计文档位于 [`docs/`](docs/) 目录:

- [总体架构设计](docs/01-总体架构设计.md)
- [数据库设计](docs/02-数据库设计.md)
- [API 接口设计](docs/03-API接口设计.md)
- [Go 后端详细设计](docs/04-Go后端详细设计.md)
- [前端详细设计](docs/05-前端详细设计.md)
- [部署与运维](docs/06-部署与运维.md)

## 🤝 贡献指南

欢迎提交 Issue 或 Pull Request。

### 开发规范

**Go 后端**:

- 遵循 Go 官方代码风格
- 使用 `gofmt` 格式化
- 函数必须添加注释
- 错误处理要完善,不能忽略错误

**Vue 前端**:

- 使用 Composition API 和 `<script setup>` 语法
- 组件命名使用 PascalCase(如 `LoginView`)
- 优先使用 TypeScript 类型定义
- 遵循 Vue 3 风格指南

**Git 提交**:

- 分支命名: `feature/xxx`, `fix/xxx`
- 提交信息格式: `feat: xxx`, `fix: xxx`, `docs: xxx`
- 提交前确保代码通过格式检查

### 贡献流程

1. Fork 仓库
2. 创建分支(`git checkout -b feature/xxx`)
3. 提交修改(`git commit -m 'feat: xxx'`)
4. 推送分支(`git push origin feature/xxx`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 [MIT](LICENSE) 许可证。

## 📞 联系方式

如有问题或建议,欢迎提交 Issue。
