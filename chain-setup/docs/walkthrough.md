# GAPChain Installation & Setup Walkthrough

Dưới đây là tổng hợp các bước và các lỗi đã được xử lý để tự động hoá việc cấu hình, khởi tạo, triển khai chaincode và xác minh mạng blockchain GAPChain sử dụng Hyperledger Fabric 3.1.1:

## 1. Cập nhật Images & Package Chaincode
- Sửa đổi [scripts/update-images.sh](file:///datame/trustfab-workspace/trustfab-projects/gapchain/chain-setup/scripts/update-images.sh) để sử dụng các bản image mới nhất cho `fabric-tools` thay vì version 3.1.1 (do không tồn tại tag này trên DockerHub tương ứng công cụ tool).
- Cập nhật quá trình đóng gói Chaincode trong [scripts/deploy-chaincode.sh](file:///datame/trustfab-workspace/trustfab-projects/gapchain/chain-setup/scripts/deploy-chaincode.sh) bằng việc chạy `go mod tidy` trong volume `chaincode_pkg_vol` với `golang:1.21-alpine` thay vì script cũ gây lỗi permission host user trên Linux.

## 2. Kết nối Peer tới Orderer & Endpoints
- **Lỗi `no endpoints currently defined`**: Cấu hình thêm `OrdererEndpoints` theo list `Addresses` cho Orderer của mọi Org trong file [configtx.yaml](file:///datame/trustfab-workspace/trustfab-projects/gapchain/chain-setup/configtx.yaml), đồng thời update Profile Genesis để orderer client mapping chính xác hostname `orderer.gapchain.vn:7050`.
- Thêm `CORE_PEER_TLS_ENABLED` và truyền Root CA TLS Cert của Orderer vào các lệnh tạo và gia nhập kênh (channel join) của peer để bypass proxy.

## 3. Xác thực Orderer & NodeOUs (Lỗi `FORBIDDEN`)
- **Lỗi `received bad status FORBIDDEN from orderer`**: Orderer từ chối Peer thực hiện hành động do policy ImplicitMeta `Readers` / `Writers` không nhận diện được loại identity (client/peer/admin).
- Đã khắc phục bằng việc set `EnableNodeOUs: true` cho tất cả các tổ chức trong [crypto-config.yaml](file:///datame/trustfab-workspace/trustfab-projects/gapchain/chain-setup/crypto-config.yaml) giúp Certificate Authority gắn đuôi định danh chính xác vào chứng chỉ phát hành.

## 4. Tắt Block Dissemination qua Gossip (Theo Fabric v3.1.x Release Notes)
- Dựa vào [thông tin release notes của Fabric 3.1.4](https://github.com/hyperledger/fabric/releases/tag/v3.1.4), các version 3.x đã deprecated Block dissemination bằng Gossip. Để khắc phục lỗi *timeout waiting for txid on all peers*, tôi đã disable Gossip block dissemination và cấu hình peer nhận block trực tiếp từ Orderer bằng cách set:
```yaml
- CORE_PEER_GOSSIP_ORGLEADER=true
- CORE_PEER_GOSSIP_USELEADERELECTION=false
- CORE_PEER_GOSSIP_STATE_ENABLED=false
- CORE_PEER_DELIVERYCLIENT_BLOCKGOSSIPENABLED=false
```
Cho toàn bộ Peer Containers trong [docker-compose.yaml](file:///datame/trustfab-workspace/trustfab-projects/gapchain/chain-setup/docker-compose.yaml).

## 5. Peer Chaincode Container Network
- **Lỗi khởi tạo Docker cho Chaincode Runtime (`API error 404`)**: Mạng của docker compose chạy mặc định là `gapchain-network` (định nghĩa ở cuối [docker-compose.yaml](file:///datame/trustfab-workspace/trustfab-projects/gapchain/chain-setup/docker-compose.yaml)), nhưng biến môi trường của peer lại trỏ là `gapchain_gapchain`.
- Tôi đã sửa lại `CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=gapchain-network` để Container chứa mã nguồn (Chaincode docker image) có thể giao tiếp với Peer.

## 6. Xác minh Mạng (Verify Network)
- Tôi đã viết thêm script [scripts/verify-network.sh](file:///datame/trustfab-workspace/trustfab-projects/gapchain/chain-setup/scripts/verify-network.sh) giúp tự động Invoke hàm `InitLedger` để khởi tạo dữ liệu ban đầu cho cả sổ cái nhật ký (`nhatky-htx-channel`) lẫn giao dịch (`giaodich-channel`), và gọi chaincode `query` để in ra danh sách bản ghi mới trong Terminal. Xác minh mã nguồn Chaincode được cài đặt thành công và sổ cái hoạt động logic trơn tru.
