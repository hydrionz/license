#!/bin/bash

cd web
npm install
npm run build
cd ../

# VERSION - 获取最新的tag并去除开头的'v'
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "0.0.1")
# 如果版本号以'v'开头，则去除这个'v'
VERSION=$(echo $VERSION | sed 's/^v//')

# 创建 build 目录，如果不存在的话
mkdir -p build

# 获取Git提交的哈希值
HASH=$(git rev-parse --short=12 HEAD)

# Linux amd64
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=linux/amd64" -o build/license-linux-amd64 .
# Linux arm64
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=linux/arm64" -o build/license-linux-arm64 .
# FreeBSD amd64
GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=freebsd/amd64" -o build/license-freebsd-amd64 .
# Windows amd64
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=windows/amd64" -o build/license-windows-amd64.exe .
# macOS arm
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=macos/arm64" -o build/license-macos-arm64 .
