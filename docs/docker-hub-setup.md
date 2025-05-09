# Docker Hub 配置指南

本文档介绍如何配置 Docker Hub 凭据，以便 GitHub Actions 工作流能够构建并推送 Docker 镜像。

## 步骤 1: 创建 Docker Hub 访问令牌

1. 登录到您的 [Docker Hub 账户](https://hub.docker.com/)
2. 点击右上角您的头像，然后选择 "Account Settings"（账户设置）
3. 在左侧菜单中，点击 "Security"（安全）
4. 点击 "New Access Token"（新建访问令牌）
5. 给令牌起一个名称（例如 "GoSSH-Web GitHub Actions"）
6. 选择适当的权限（至少需要 "Read & Write" 权限）
7. 点击 "Generate"（生成）
8. **重要**：复制生成的令牌，因为它只会显示一次！

## 步骤 2: 将 Docker Hub 凭据添加到 GitHub Secrets

1. 转到您的 GitHub 仓库（MoTeam-cn/GoSSH-Web）
2. 点击 "Settings"（设置）标签
3. 在左侧菜单中，点击 "Secrets and variables" > "Actions"
4. 点击 "New repository secret"（新建仓库密钥）
5. 添加以下两个 Secrets：

   a. 第一个 Secret：
   - 名称：`DOCKERHUB_USERNAME`
   - 值：您的 Docker Hub 用户名
   
   b. 第二个 Secret：
   - 名称：`DOCKERHUB_TOKEN`
   - 值：您刚才生成的 Docker Hub 访问令牌

## 步骤 3: 验证配置

配置完成后，您可以通过以下方式验证：

1. 在 GitHub 仓库中，提交一个包含 `[build]` 前缀的提交消息
2. 查看 Actions 标签页中的工作流运行状态
3. 确认工作流能够成功登录到 Docker Hub 并推送镜像

## Docker 镜像使用

配置完成后，您的 Docker 镜像将发布到 Docker Hub，可以通过以下命令使用：

```bash
# 拉取最新版本
docker pull MoTeam-cn/gossh-web:latest

# 拉取特定版本
docker pull MoTeam-cn/gossh-web:v1.0.0

# 运行容器
docker run -d -p 8080:8080 --name gossh-web MoTeam-cn/gossh-web:latest
```

使用 Docker Compose：

```bash
# 使用环境变量指定镜像
export DOCKER_IMAGE=MoTeam-cn/gossh-web:latest

# 启动服务
docker-compose up -d
``` 