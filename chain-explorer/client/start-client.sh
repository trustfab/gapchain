#!/bin/bash
cd "$(dirname "$0")"

echo "========================================="
echo "KHỞI ĐỘNG HYPERLEDGER EXPLORER CLIENT"
echo "========================================="

echo "[*] Đang khởi động Client (npm start) ở chế độ chạy ngầm (Background)..."

# Chạy npm start ngầm, output lưu vào file client-run.log
nohup npm start > client-run.log 2>&1 &

echo "[*] App Client đã chạy (Hot Reload & Proxy) trên http://localhost:3000"
echo "    - Để xem log trực tiếp: tail -f client-run.log"
echo "    - Để tắt Client: ./stop-client.sh"
