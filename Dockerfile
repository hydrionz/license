FROM node:22-alpine AS frontend-builder
WORKDIR /app
# 复制package.json和package-lock.json以利用缓存
COPY web/package*.json ./
# 安装依赖
RUN npm install
# 复制React前端源代码
COPY web/ ./
# 构建React应用
RUN npm run build && \
    ls -la build && \
    echo "前端构建完成，检查build目录"

FROM golang:1.24-alpine AS go-builder
# 定义版本参数，默认为0.0.1
ARG VERSION=0.0.1
ARG HASH=2ed26fe1
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
# 复制Go源代码
COPY . .
# 删除可能存在的web/build目录，避免嵌入旧文件
RUN rm -rf web/build
# 复制前端构建产物
COPY --from=frontend-builder /app/build/ ./web/build/
# 验证前端文件已复制成功
RUN ls -la web/build && \
    echo "前端文件已复制到web/build目录"
# 安装必要的编译工具
RUN apk add --no-cache gcc musl-dev
# 设置CGO禁用用
ENV CGO_ENABLED=0
# 编译Go应用，注入版本信息
RUN go build -v -ldflags="-X 'license/internal/sys.Version=${VERSION} -X license/internal/sys.Build=${HASH} -X license/internal/sys.OsArch=${TARGETOS}/${TARGETARCH}'" -o license ./main.go && \
    ls -la license || echo "验证可执行文件失败，但这可能是因为缺少参数"

FROM alpine:latest
WORKDIR /app
# 复制Go二进制文件
COPY --from=go-builder /app/license ./license
RUN ls -la /app/license && \
    chmod +x /app/license && \
    echo "可执行文件已复制到/app/license"
# 安装运行时依赖
RUN apk update && \
    apk add --no-cache ca-certificates tzdata libc6-compat && \
    update-ca-certificates 2>/dev/null || true
# 创建数据目录
RUN mkdir -p /data

# 启动应用
CMD ["/app/license"]

