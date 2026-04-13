package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type OrgConfig struct {
	MspID        string
	CertPath     string
	KeyPath      string
	PeerEndpoint string
	TLSCertPath  string
	GatewayPeer  string
	Channels     []string
}

type Config struct {
	Port    string
	GinMode string
	Orgs    map[string]*OrgConfig
}

func LoadConfig() *Config {
	err := godotenv.Overload()
	if err != nil {
		log.Println("WARNING: Không tìm thấy file .env, sử dụng biến môi trường hệ thống")
	}

	cfg := &Config{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),
		Orgs:    make(map[string]*OrgConfig),
	}

	orgList := getEnv("FABRIC_ORGS", "PlatformOrgMSP,HTXNongSanOrgMSP,ChiCucBVTVOrgMSP,NPPXanhOrgMSP,NPPTieuChuanOrgMSP")
	for _, mspID := range strings.Split(orgList, ",") {
		mspID = strings.TrimSpace(mspID)
		if mspID == "" {
			continue
		}

		channels := getEnv("FABRIC_"+mspID+"_CHANNELS", "")
		var channelList []string
		if channels != "" {
			for _, ch := range strings.Split(channels, ",") {
				ch = strings.TrimSpace(ch)
				if ch != "" {
					channelList = append(channelList, ch)
				}
			}
		}

		cfg.Orgs[mspID] = &OrgConfig{
			MspID:        mspID,
			CertPath:     getEnv("FABRIC_"+mspID+"_CERT_PATH", ""),
			KeyPath:      getEnv("FABRIC_"+mspID+"_KEY_PATH", ""),
			PeerEndpoint: getEnv("FABRIC_"+mspID+"_PEER_ENDPOINT", ""),
			TLSCertPath:  getEnv("FABRIC_"+mspID+"_TLS_CERT", ""),
			GatewayPeer:  getEnv("FABRIC_"+mspID+"_GATEWAY_PEER", ""),
			Channels:     channelList,
		}
	}

	return cfg
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
