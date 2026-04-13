# GAPChain - Deploy Chaincode

Hỗ trợ đóng gói, cài đặt và commit chaincode lên mạng GAPChain.

## Chaincode trong dự án (MVP — 3 chaincode)

| Chaincode | Channel | Thành viên | Chức năng |
|-----------|---------|-----------|-----------|
| `lohang_cc` | nhatky-htx-channel | Platform, HTXNongSan, ChiCucBVTV | Lô hàng, chứng nhận VietGAP/GlobalGAP, QR truy xuất |
| `nhatky_cc` | nhatky-htx-channel | Platform, HTXNongSan, ChiCucBVTV | Nhật ký canh tác gắn với lô hàng |
| `giao_dich_cc` | giaodich-channel | Platform, HTXNongSan, NPPXanh | Giao dịch mua bán, công nợ NPP, hoa hồng |

> **Thứ tự deploy bắt buộc**: `lohang_cc` → `nhatky_cc` → `giao_dich_cc`
> (nhatky_cc và giao_dich_cc đều tham chiếu `ma_lo` từ lohang_cc)

Thư mục chaincode: `chain-setup/chaincode/`

## Quy trình deploy (Fabric v3.1.1 lifecycle)

1. **Package** - Đóng gói chaincode
2. **Install** - Cài đặt lên từng peer
3. **Approve** - Mỗi tổ chức approve definition
4. **Commit** - Commit lên channel

## Nhiệm vụ

Người dùng muốn deploy hoặc update chaincode. Hãy:

1. Hỏi rõ: deploy `lohang_cc`, `nhatky_cc`, `giao_dich_cc`, hay cả ba?
2. Đọc script `chain-setup/scripts/deploy-chaincode.sh` để hiểu flow hiện tại
3. Kiểm tra containers đang chạy: `docker ps --format "table {{.Names}}\t{{.Status}}"`
4. Nếu người dùng muốn chạy deploy, hướng dẫn:

```bash
cd chain-setup
# Deploy cả 3 chaincode (theo thứ tự: lohang → nhatky → giaodich)
./scripts/deploy-chaincode.sh

# Kiểm tra từng chaincode đã committed chưa
export PATH=$PATH:$(pwd)/../bin
export FABRIC_CFG_PATH=$(pwd)/config/platform

peer lifecycle chaincode querycommitted \
  --channelID nhatky-htx-channel \
  --name lohang_cc \
  --output json

peer lifecycle chaincode querycommitted \
  --channelID nhatky-htx-channel \
  --name nhatky_cc \
  --output json

peer lifecycle chaincode querycommitted \
  --channelID giaodich-channel \
  --name giao_dich_cc \
  --output json
```

5. Nếu có lỗi trong quá trình deploy, đọc log và đề xuất fix
6. Sau khi deploy, hướng dẫn test bằng `invoke` hoặc `query`

## Test chaincode sau deploy

```bash
# Biến môi trường dùng chung
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem
export HTX_TLS_CERT=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt
export BVTV_TLS_CERT=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/chicucbvtv.gapchain.vn/peers/peer0.chicucbvtv.gapchain.vn/tls/ca.crt

# Test lohang_cc — Tạo lô hàng mới
peer chaincode invoke \
  -o orderer.gapchain.vn:7050 \
  --ordererTLSHostnameOverride orderer.gapchain.vn \
  --tls --cafile $ORDERER_CA \
  -C nhatky-htx-channel -n lohang_cc \
  --peerAddresses peer0.htxnongsan.gapchain.vn:7051 --tlsRootCertFiles $HTX_TLS_CERT \
  -c '{"function":"TaoLotHang","Args":["LH-HTX001-2025-001","HTX001","Gao ST25","lua","1000","kg","Dong Xuan 2025","Can Tho"]}'

# Test lohang_cc — Đọc lô hàng
peer chaincode query \
  -C nhatky-htx-channel -n lohang_cc \
  -c '{"function":"DocLotHang","Args":["LH-HTX001-2025-001"]}'

# Test nhatky_cc — Ghi nhật ký canh tác
peer chaincode invoke \
  -o orderer.gapchain.vn:7050 \
  --ordererTLSHostnameOverride orderer.gapchain.vn \
  --tls --cafile $ORDERER_CA \
  -C nhatky-htx-channel -n nhatky_cc \
  --peerAddresses peer0.htxnongsan.gapchain.vn:7051 --tlsRootCertFiles $HTX_TLS_CERT \
  -c '{"function":"GhiNhatKy","Args":["NK-001","LH-HTX001-2025-001","HTX001","san_xuat","Giao trong lua","Canh dong A","kt_vien_001","2025-03-01",""]}'

# Test nhatky_cc — Xem nhật ký theo lô hàng
peer chaincode query \
  -C nhatky-htx-channel -n nhatky_cc \
  -c '{"function":"DocNhatKyTheoLo","Args":["LH-HTX001-2025-001"]}'

# Test giao_dich_cc — Tạo giao dịch
peer chaincode invoke \
  -o orderer.gapchain.vn:7050 \
  --ordererTLSHostnameOverride orderer.gapchain.vn \
  --tls --cafile $ORDERER_CA \
  -C giaodich-channel -n giao_dich_cc \
  --peerAddresses peer0.htxnongsan.gapchain.vn:7051 --tlsRootCertFiles $HTX_TLS_CERT \
  -c '{"function":"TaoGiaoDich","Args":["GD-001","LH-HTX001-2025-001","HTX001","NPP001","Gao ST25","500","kg","15000","2",""]}'

# Test giao_dich_cc — Xem công nợ NPP
peer chaincode query \
  -C giaodich-channel -n giao_dich_cc \
  -c '{"function":"DocCongNoNPP","Args":["NPP001"]}'
```

$ARGUMENTS
