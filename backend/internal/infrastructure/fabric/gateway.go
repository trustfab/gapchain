package fabric

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/trustfab/gapchain/backend/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// FabricContract định nghĩa các hàm tương tác với chaincode, giúp mock dễ dàng trong Unit Test
type FabricContract interface {
	SubmitTransaction(name string, args ...string) ([]byte, error)
	EvaluateTransaction(name string, args ...string) ([]byte, error)
}

// OrgGateway chứa Gateway và Contracts cho một org cụ thể
type OrgGateway struct {
	Gateway          *client.Gateway
	LotHangContract  FabricContract // nil nếu org không tham gia nhatky-htx-channel
	NhatKyContract   FabricContract // nil nếu org không tham gia nhatky-htx-channel
	GiaoDichContract FabricContract // nil nếu org không tham gia giaodich-channel
}

// GatewayRegistry quản lý Gateway cho tất cả orgs
type GatewayRegistry struct {
	mu            sync.RWMutex
	gateways      map[string]*OrgGateway // key = MspID
	FallbackMspID string
}

func NewGatewayRegistry(cfg *config.Config) (*GatewayRegistry, error) {
	registry := &GatewayRegistry{
		gateways:      make(map[string]*OrgGateway),
		FallbackMspID: "PlatformOrgMSP",
	}

	for mspID, orgCfg := range cfg.Orgs {
		orgGw, err := connectOrg(mspID, orgCfg)
		if err != nil {
			return nil, fmt.Errorf("khong the ket noi org %s: %w", mspID, err)
		}
		log.Printf("Connected gateway cho org: %s (channels: %v)", mspID, orgCfg.Channels)
		registry.gateways[mspID] = orgGw
	}

	return registry, nil
}

// GetOrgGateway trả về OrgGateway theo mspID
func (r *GatewayRegistry) GetOrgGateway(mspID string) (*OrgGateway, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	gw, ok := r.gateways[mspID]
	if !ok {
		return nil, fmt.Errorf("khong tim thay gateway cho org: %s", mspID)
	}
	return gw, nil
}

// GetFallbackGateway trả về gateway mặc định (Platform) cho public endpoints
func (r *GatewayRegistry) GetFallbackGateway() (*OrgGateway, error) {
	return r.GetOrgGateway(r.FallbackMspID)
}

// SetOrgGatewayForTest hỗ trợ inject mock Gateway cho Unit Test
func (r *GatewayRegistry) SetOrgGatewayForTest(mspID string, gw *OrgGateway) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.gateways == nil {
		r.gateways = make(map[string]*OrgGateway)
	}
	r.gateways[mspID] = gw
}

// CloseAll đóng tất cả gateway connections
func (r *GatewayRegistry) CloseAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for mspID, gw := range r.gateways {
		if gw.Gateway != nil {
			gw.Gateway.Close()
			log.Printf("Closed gateway cho org: %s", mspID)
		}
	}
}

func connectOrg(mspID string, orgCfg *config.OrgConfig) (*OrgGateway, error) {
	id, err := newIdentity(mspID, orgCfg.CertPath)
	if err != nil {
		return nil, fmt.Errorf("identity %s: %w", mspID, err)
	}
	sign, err := newSigner(orgCfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("signer %s: %w", mspID, err)
	}

	tlsCert, err := os.ReadFile(orgCfg.TLSCertPath)
	if err != nil {
		return nil, fmt.Errorf("tls cert %s: %w", mspID, err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(tlsCert)
	creds := credentials.NewClientTLSFromCert(certPool, orgCfg.GatewayPeer)

	conn, err := grpc.NewClient(orgCfg.PeerEndpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("grpc %s: %w", mspID, err)
	}

	gw, err := client.Connect(id, client.WithSign(sign), client.WithClientConnection(conn))
	if err != nil {
		return nil, err
	}

	orgGw := &OrgGateway{Gateway: gw}

	channelSet := make(map[string]bool)
	for _, ch := range orgCfg.Channels {
		channelSet[ch] = true
	}

	if channelSet["nhatky-htx-channel"] {
		network := gw.GetNetwork("nhatky-htx-channel")
		orgGw.LotHangContract = network.GetContract("lohang_cc")
		orgGw.NhatKyContract = network.GetContract("nhatky_cc")
	}

	if channelSet["giaodich-channel"] {
		network := gw.GetNetwork("giaodich-channel")
		orgGw.GiaoDichContract = network.GetContract("giao_dich_cc")
	}

	return orgGw, nil
}

func newIdentity(mspID, certPath string) (*identity.X509Identity, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	cert, err := identity.CertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}
	return identity.NewX509Identity(mspID, cert)
}

func newSigner(keyDir string) (identity.Sign, error) {
	entries, err := os.ReadDir(keyDir)
	if err != nil || len(entries) == 0 {
		return nil, fmt.Errorf("khong tim thay private key trong %s", keyDir)
	}
	keyPEM, err := os.ReadFile(keyDir + "/" + entries[0].Name())
	if err != nil {
		return nil, err
	}
	key, err := identity.PrivateKeyFromPEM(keyPEM)
	if err != nil {
		return nil, err
	}
	return identity.NewPrivateKeySign(key)
}
