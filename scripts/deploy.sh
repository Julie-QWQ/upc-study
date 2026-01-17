#!/bin/bash
# 自动部署脚本
# 用途: 拉取最新代码、构建前后端、重启服务

set -e  # 遇到错误立即退出

echo "========================================="
echo "Study-UPC 自动部署脚本"
echo "========================================="

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

success() {
  echo -e "${GREEN}✅ $1${NC}"
}

error() {
  echo -e "${RED}❌ $1${NC}"
  exit 1
}

# ==================== 1. 拉取最新代码 ====================
echo ""
echo "=== 1. 拉取最新代码 ==="
cd /root/upc-study
git pull

success "代码拉取完成"

# ==================== 2. 安装前端依赖 ====================
echo ""
echo "=== 2. 安装前端依赖 ==="
cd /root/upc-study/frontend

# 使用国内 npm 镜像加速
npm config set registry https://registry.npmmirror.com
npm ci || error "npm ci 失败"

success "前端依赖安装完成"

# ==================== 3. 构建前端 ====================
echo ""
echo "=== 3. 构建前端 ==="
npm run build || error "前端构建失败"

if [ ! -d "dist" ]; then
  error "dist 目录不存在"
fi

success "前端构建完成"

# ==================== 4. 部署前端文件 ====================
echo ""
echo "=== 4. 部署前端文件 ==="
rm -rf /var/www/upc-study
cp -r dist /var/www/upc-study

success "前端文件部署完成"

# ==================== 5. 重启 Nginx ====================
echo ""
echo "=== 5. 重启 Nginx ==="
sudo systemctl restart nginx || error "Nginx 重启失败"
sudo systemctl status nginx --no-pager

success "Nginx 重启完成"

# ==================== 6. 下载 Go 依赖 ====================
echo ""
echo "=== 6. 下载 Go 依赖 ==="
cd /root/upc-study/backend

# 使用国内 Go 代理加速
export GOPROXY=https://goproxy.cn,direct
go mod download || error "Go 依赖下载失败"

success "Go 依赖下载完成"

# ==================== 7. 构建后端 ====================
echo ""
echo "=== 7. 构建后端 ==="
# 检测服务器架构
ARCH=$(uname -m)
echo "服务器架构: $ARCH"

if [ "$ARCH" = "aarch64" ]; then
  echo "检测到 ARM 架构，使用交叉编译"
  GOOS=linux GOARCH=arm64 go build -o upc-study-server cmd/server/main.go || error "后端构建失败"
else
  echo "检测到 x86 架构，使用原生编译"
  go build -o upc-study-server cmd/server/main.go || error "后端构建失败"
fi

if [ ! -f "upc-study-server" ]; then
  error "后端可执行文件不存在"
fi

success "后端构建完成"

# ==================== 8. 部署后端服务 ====================
echo ""
echo "=== 8. 部署后端服务 ==="
sudo chmod +x upc-study-server

success "后端服务部署完成"

# ==================== 9. 重启后端服务 ====================
echo ""
echo "=== 9. 重启后端服务 ==="
sudo systemctl restart upc-study || error "后端服务重启失败"
sudo systemctl status upc-study --no-pager

success "后端服务重启完成"

# ==================== 完成 ====================
echo ""
echo "========================================="
success "部署完成！"
echo "========================================="
