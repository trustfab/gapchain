flowchart TD
    subgraph Users ["👥 Người dùng cuối (Ứng dụng Web / Mobile)"]
        A1["👨‍💼 KT viên (HTX)"]
        A2["👨‍💼 QL HTX"]
        A3["👨‍💼 QL Mua hàng (Phân phối)"]
        A4["🧑‍🌾 Nông dân"]
    end

    subgraph AppLayer ["🔧 Platform (Lớp Ứng dụng)"]
        B1["🌐 API Gateway / Web App"]
        B2["⚙️ Chaincode SDK (Fabric SDK Go / Node.js)"]
    end
    
    subgraph FabricNetwork ["🔗 Mạng lưới Blockchain GAPChain"]
        direction TB
        
        subgraph PlatformOrg ["💻 Tổ chức Platform"]
            P1["💻 PlatformOrg.Peer"]
            PCA["🔑 CA: PlatformOrg"]
        end
        
        subgraph Channel1 ["📋 Kênh: nhatky-htx-channel"]
            N1["🏢 HTXNongSanOrg.Peer"]
            N2["🏛️ ChiCucBVTVOrg.Peer"]
            CC1["📝 Chaincode: nhatky_cc"]
        end

        subgraph Channel2 ["💰 Kênh: giaodich-channel"]
            N3["🏢 HTXNongSanOrg.Peer"]
            N4["🏪 NPPXanhOrg.Peer"]
            CC2["💱 Chaincode: giao_dich_cc"]
        end

        subgraph OrdererSystem ["🔄 OrdererOrg (Ordering Service)"]
            O1["⚡ Orderer Node (Raft)"]
        end
    end

    %% Kết nối lớp người dùng với app
    A1 --> B1
    A2 --> B1
    A3 --> B1
    A4 --> B1

    %% Ứng dụng gọi SDK
    B1 --> B2

    %% SDK tương tác với peer của chính nó
    B2 -->|"Kết nối Gateway"| P1

    %% Peer của PlatformOrg tương tác với các Peer khác
    P1 -->|"📝 Endorse"| N1
    P1 -->|"📝 Endorse"| N2
    P1 -->|"💰 Endorse"| N3
    P1 -->|"💰 Endorse"| N4
    
    %% Chaincode chạy trên peer
    CC1 --> N1
    CC1 --> N2
    CC2 --> N3
    CC2 --> N4

    %% Các peer kết nối orderer
    P1 --> O1
    N1 --> O1
    N2 --> O1
    N3 --> O1
    N4 --> O1

    %% Định danh người dùng
    A1 -.->|"🔐 xác thực"| PCA
    A2 -.->|"🔐 xác thực"| PCA
    A3 -.->|"🔐 xác thực"| PCA
    A4 -.->|"🔐 xác thực"| PCA

    %% Styling
    classDef userClass fill:#e1f5fe,stroke:#0277bd,stroke-width:2px
    classDef appClass fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    classDef peerClass fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef chaincodeClass fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    classDef ordererClass fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    classDef caClass fill:#f1f8e9,stroke:#558b2f,stroke-width:2px

    class A1,A2,A3,A4 userClass
    class B1,B2 appClass
    class P1,N1,N2,N3,N4 peerClass
    class CC1,CC2 chaincodeClass
    class O1 ordererClass
    class PCA,CA1,CA2,CA3 caClass