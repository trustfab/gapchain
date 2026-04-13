# GAPChain - Cài đặt Hyperledger Fabric v3.x

Hỗ trợ cài đặt toàn bộ môi trường Hyperledger Fabric v3.x cho GAPChain trên máy mới.

---

## Phiên bản yêu cầu của dự án này

| Thành phần | Phiên bản |
|------------|-----------|
| fabric-peer / fabric-orderer | **3.1.1** |
| fabric-ca | **1.5.15** |
| fabric-tools (CLI) | **3.1.1** |
| Go (chaincode & client) | **1.21+** |
| Docker Engine | **24.0+** |
| Docker Compose | **v2.x** (plugin) |

---

## Nhiệm vụ

Người dùng muốn cài đặt môi trường Fabric. Hãy:

1. Hỏi hệ điều hành: **macOS (Apple Silicon / Intel)** hay **Linux (Ubuntu/Debian)**?
2. Kiểm tra từng điều kiện tiên quyết bằng các lệnh ở phần "Kiểm tra prerequisites" bên dưới.
3. Hướng dẫn theo đúng OS, **từng bước một**, chỉ chuyển bước tiếp theo sau khi bước trước thành công.
4. Sau mỗi bước, yêu cầu người dùng paste output để xác nhận trước khi tiếp tục.

---

## Bước 0 — Kiểm tra prerequisites

```bash
# Docker Engine
docker --version          # cần >= 24.0
docker compose version    # cần Compose v2 (plugin, không phải docker-compose v1)

# Go
go version                # cần >= 1.21

# Git
git --version

# curl
curl --version

# (macOS) Xcode Command Line Tools - cần cho một số build tool
xcode-select -p
```

Nếu thiếu bất kỳ thứ gì, hướng dẫn cài trước khi tiếp tục.

---

## Bước 1 — Cài Docker (nếu chưa có)

### macOS
```bash
# Cài Docker Desktop (bao gồm Docker Engine + Compose v2)
brew install --cask docker
# Sau đó mở Docker Desktop từ Applications và chờ nó start
open -a Docker
```

### Ubuntu / Debian
```bash
# Gỡ bản cũ nếu có
sudo apt-get remove docker docker-engine docker.io containerd runc

# Cài dependencies
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg lsb-release

# Thêm Docker GPG key và repo
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
  https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
  | sudo tee /etc/apt/sources.list.d/docker.list

# Cài Docker Engine + Compose plugin
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Thêm user vào group docker (không cần sudo mỗi lần)
sudo usermod -aG docker $USER
newgrp docker
```

**Kiểm tra:**
```bash
docker run --rm hello-world
docker compose version
```

---

## Bước 2 — Cài Go (nếu chưa có)

### macOS
```bash
brew install go
# Hoặc tải trực tiếp từ https://go.dev/dl/ (chọn darwin/arm64 cho Apple Silicon)
```

### Linux
```bash
# Tải Go 1.22 (thay bằng phiên bản mới nhất nếu cần)
curl -LO https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz

# Thêm vào PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
```

**Kiểm tra:**
```bash
go version   # phải hiện >= go1.21
```

---

## Bước 3 — Tải Docker Images Fabric

Pull các images mà project GAPChain cần (tổng ~1.5 GB):

```bash
# Fabric core images
docker pull hyperledger/fabric-peer:3.1.1
docker pull hyperledger/fabric-orderer:3.1.1
docker pull hyperledger/fabric-ca:1.5.15

# Tools image (dùng để package chaincode)
docker pull hyperledger/fabric-tools:3.1.1

# CouchDB (state database)
docker pull couchdb:3.3.3
```

**Kiểm tra đã có đủ images:**
```bash
docker images | grep -E "fabric|couchdb"
```

Output mong đợi:
```
hyperledger/fabric-peer      3.1.1   ...
hyperledger/fabric-orderer   3.1.1   ...
hyperledger/fabric-ca        1.5.15  ...
hyperledger/fabric-tools     3.1.1   ...
couchdb                      3.3.3   ...
```

---

## Bước 4 — Tải Fabric Binaries (CLI Tools)

Fabric binaries (`cryptogen`, `configtxgen`, `peer`, `orderer`, v.v.) dùng để tạo certs và quản lý mạng qua CLI.

### Cách A — Script tự động của Hyperledger (khuyến nghị)

```bash
# Tải vào thư mục hiện tại (sẽ tạo thư mục bin/ và config/)
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh \
  | bash -s -- --fabric-version 3.1.1 --ca-version 1.5.15 binary
```

Script này tạo thư mục `bin/` với đầy đủ binaries.

### Cách B — Tải thủ công (nếu script bị chặn)

```bash
# Xác định OS và arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')   # linux hoặc darwin
ARCH=$(uname -m)                               # x86_64 hoặc arm64

# Tải
curl -LO "https://github.com/hyperledger/fabric/releases/download/v3.1.1/hyperledger-fabric-${OS}-${ARCH}-3.1.1.tar.gz"
tar -xzf hyperledger-fabric-${OS}-${ARCH}-3.1.1.tar.gz
# Giải nén tạo thư mục bin/ ngay trong thư mục hiện tại
```

### Cài vào đúng vị trí cho GAPChain

Binaries phải nằm tại `gapchain/bin/`:
```bash
# Giả sử đang ở thư mục gapchain/
ls bin/peer bin/cryptogen bin/configtxgen   # phải tồn tại
```

**Thêm vào PATH khi làm việc:**
```bash
# Chạy lệnh này mỗi khi mở terminal mới, hoặc thêm vào ~/.zshrc / ~/.bashrc
export PATH=$PATH:/đường/dẫn/tới/gapchain/bin
```

**Kiểm tra binaries hoạt động:**
```bash
peer version
# Output: peer:
#   Version: 3.1.1
#   ...

cryptogen version
configtxgen --version
```

---

## Bước 5 — Kiểm tra tổng thể trước khi chạy mạng

```bash
# 1. Docker daemon đang chạy
docker ps

# 2. Tất cả images đã pull
docker images | grep -E "fabric-peer|fabric-orderer|fabric-ca|fabric-tools|couchdb"

# 3. Binaries hoạt động
peer version | head -3
cryptogen version | head -2
configtxgen --version | head -2

# 4. Go version hợp lệ
go version

# 5. Docker Compose v2
docker compose version
```

Nếu tất cả đều OK → chạy `/setup-network` để khởi động mạng GAPChain.

---

## Lỗi thường gặp khi cài đặt

| Lỗi | Nguyên nhân | Fix |
|-----|-------------|-----|
| `Cannot connect to Docker daemon` | Docker Desktop chưa mở | Mở Docker Desktop và chờ start |
| `docker-compose: command not found` | Dùng Compose v1 cũ | Dùng `docker compose` (không có dấu `-`) |
| `exec format error` khi chạy peer | Tải sai arch (x86 vs arm64) | Tải lại đúng ARCH cho máy |
| `permission denied` khi chạy docker | User chưa thuộc group docker | `sudo usermod -aG docker $USER` rồi logout/login |
| Pull image bị timeout | Mạng chậm hoặc bị chặn | Dùng Docker mirror hoặc VPN |
| `fabric-ca:1.5.15` không start | Cần compatibility mode | Kiểm tra env `FABRIC_CA_SERVER_COMPATIBILITY_MODE_V1_3=true` trong docker-compose-ca.yaml |

---

## Lưu ý đặc biệt cho Apple Silicon (M1/M2/M3)

```bash
# Kiểm tra arch
uname -m   # phải hiện arm64

# Fabric 3.1.1 hỗ trợ arm64 native
# Nếu tải binary thủ công, dùng file: hyperledger-fabric-darwin-arm64-3.1.1.tar.gz

# Docker Desktop trên Apple Silicon chạy images x86 qua Rosetta
# Một số images cũ (< Fabric 3.x) có thể chậm hơn trên arm64
```

$ARGUMENTS
