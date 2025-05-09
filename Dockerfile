FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gossh-web

# 使用更小的基础镜像
FROM alpine:3.19

# 安装基本工具和SSL证书
RUN apk --no-cache add ca-certificates tzdata

# 创建非root用户
RUN adduser -D -h /app appuser
WORKDIR /app
USER appuser

# 从构建阶段复制二进制文件和资源
COPY --from=builder --chown=appuser:appuser /app/gossh-web .
COPY --from=builder --chown=appuser:appuser /app/templates ./templates
COPY --from=builder --chown=appuser:appuser /app/static ./static

# 设置默认环境变量
ENV PORT=8080 \
    HOST=0.0.0.0 \
    LOG_LEVEL=info

# 元数据
LABEL org.opencontainers.image.title="GoSSH-Web"
LABEL org.opencontainers.image.description="强大而现代的Web终端解决方案"
LABEL org.opencontainers.image.source="https://github.com/MoTeam-cn/GoSSH-Web"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.vendor="MoTeam-cn"

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/ || exit 1

EXPOSE ${PORT}

# 启动命令
CMD ["./gossh-web"] 