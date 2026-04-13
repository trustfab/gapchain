package http

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Xử lý logic Multi-Tenant (Mock dựa trên tiền tố username)
	// Thực tế sẽ dùng Database để check username & mapping.
	var role, mspID, tenantID string

	switch {
	case req.Username == "platform":
		role = "platform"
		mspID = "PlatformOrgMSP"
		tenantID = "PLATFORM"
	case req.Username == "bvtv":
		role = "bvtv"
		mspID = "ChiCucBVTVOrgMSP"
		tenantID = "BVTV_CT"
	case len(req.Username) >= 3 && req.Username[:3] == "npp":
		role = "npp"
		mspID = "NPPXanhOrgMSP"
		tenantID = "NPP001" // Thực tế map theo DB
	case len(req.Username) >= 4 && req.Username[:4] == "nptc":
		role = "nptc"
		mspID = "NPPTieuChuanOrgMSP"
		tenantID = "NPTC001" // Thực tế map theo DB
	case len(req.Username) >= 3 && req.Username[:3] == "htx":
		role = "htx"
		mspID = "HTXNongSanOrgMSP"
		tenantID = "HTX001" // VD: Nếu username = htx001, tenant = HTX001
	default:
		// Default fallback cho MVP
		role = "htx"
		mspID = "HTXNongSanOrgMSP"
		tenantID = "HTX001"
	}

	// Hardcode password check
	if req.Password != "123456" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai tài khoản hoặc mật khẩu (Pass chung: 123456)"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-here"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   req.Username,
		"role":      role,
		"msp_id":    mspID,
		"tenant_id": tenantID,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"username":  req.Username,
			"role":      role,
			"msp_id":    mspID,
			"tenant_id": tenantID,
		},
	})
}
