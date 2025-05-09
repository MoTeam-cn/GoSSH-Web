FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# 安装依赖并构建
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o gossh-web

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/gossh-web .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./gossh-web"] 