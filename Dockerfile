FROM golang:1.23.12-alpine AS builder

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译 API 服务
RUN cd service/api && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/media-api media.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件和配置
COPY --from=builder /app/media-api .
COPY --from=builder /app/service/api/etc ./etc

# 设置时区
ENV TZ=Asia/Shanghai

EXPOSE 8888

CMD ["./media-api", "-f", "etc/media-api.yaml"]
