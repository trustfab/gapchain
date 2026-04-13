#!/bin/bash
echo "========================================="
echo "TẮT HYPERLEDGER EXPLORER CLIENT"
echo "========================================="

echo "[*] Đang tìm và tắt tiến trình đang giữ cổng 3000..."

# Tìm mã PID của tiến trình
PID=$(lsof -ti:3000)

if [ -z "$PID" ]; then
    echo "[!] Không có tiến trình nào đang chạy trên cổng 3000."
else
    kill -9 $PID
    echo "[*] Đã tắt thành công tiến trình Client (PID: $PID)."
fi
