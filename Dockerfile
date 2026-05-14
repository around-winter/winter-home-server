
# 构建阶段
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# 运行阶段
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./server"]

