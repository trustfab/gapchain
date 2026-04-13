# GAPChain - Setup Network

Hướng dẫn và hỗ trợ khởi tạo mạng Hyperledger Fabric cho GAPChain.

## Ngữ cảnh dự án

GAPChain là mạng Hyperledger Fabric v3.1.1 với 4 tổ chức:
- **PlatformOrg** - Quản lý nền tảng (platform.gapchain.vn)
- **HTXNongSanOrg** - Hợp tác xã nông sản (htxnongsan.gapchain.vn)
- **ChiCucBVTVOrg** - Chi cục Bảo vệ Thực vật (chicucbvtv.gapchain.vn)
- **NPPXanhOrg** - Nhà phân phối xanh (nppxanh.gapchain.vn)

2 Channel:
- `nhatky-htx-channel` - Nhật ký hoạt động (chaincode: nhatky_cc)
- `giaodich-channel` - Giao dịch (chaincode: giao_dich_cc)

## Nhiệm vụ

Người dùng muốn setup hoặc debug mạng GAPChain. Hãy:

1. Đọc file `chain-setup/README.md` để nắm rõ hướng dẫn hiện tại
2. Kiểm tra Docker đang chạy không: `docker ps`
3. Kiểm tra trạng thái các container: `docker ps -a --filter "name=gapchain" --filter "name=orderer" --filter "name=peer" --filter "name=couchdb"`
4. Dựa vào câu hỏi/vấn đề cụ thể của người dùng, đề xuất lệnh phù hợp từ:
   - `chain-setup/scripts/setup-network.sh` - Khởi tạo toàn bộ mạng
   - `chain-setup/scripts/cleanup-network.sh` - Dọn dẹp và reset
   - `chain-setup/scripts/verify-network.sh` - Kiểm tra mạng
   - `chain-setup/scripts/verify-anchor-peers.sh` - Kiểm tra anchor peers

5. Nếu có lỗi, đọc log container liên quan và đề xuất fix

## Lưu ý quan trọng

- Fabric binaries nằm tại `bin/` - cần thêm vào PATH: `export PATH=$PATH:$(pwd)/../bin`
- Docker images cần: `hyperledger/fabric-peer:3.1.1`, `hyperledger/fabric-orderer:3.1.1`, `hyperledger/fabric-ca:1.5.15`
- Working directory cho scripts là `chain-setup/`
- Certificates được tạo bởi cryptogen tại `chain-setup/organizations/`

$ARGUMENTS
