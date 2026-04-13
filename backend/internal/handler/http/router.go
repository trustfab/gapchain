package http

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/trustfab/gapchain/backend/internal/config"
	"github.com/trustfab/gapchain/backend/internal/infrastructure/fabric"
	"github.com/trustfab/gapchain/backend/internal/middleware"
	"go.uber.org/fx"
)

// getMspID lấy mspID từ gin context (đã set bởi JWT middleware)
func getMspID(c *gin.Context) string {
	mspID, _ := c.Get("mspID")
	if mspID == nil {
		return ""
	}
	return mspID.(string)
}

func SetupRouter(
	lc fx.Lifecycle,
	cfg *config.Config,
	registry *fabric.GatewayRegistry,
	lohangH *LohangHandler,
	nhatkyH *NhatkyHandler,
	giaodichH *GiaodichHandler,
	authH *AuthHandler,
) *gin.Engine {

	gin.SetMode(cfg.GinMode)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public Routes
	r.POST("/api/v1/auth/login", authH.Login)
	r.GET("/api/v1/consumer/:ma_lo", lohangH.LayThongTinTraCuu)

	// Private Routes
	api := r.Group("/api/v1", middleware.JWTAuth())
	{
		// Lô hàng
		api.POST("/lohang", lohangH.TaoLotHang)
		api.POST("/lohang/:ma_lo/tach", lohangH.TachLo)
		api.GET("/lohang/:ma_lo", lohangH.DocLotHang)
		api.GET("/lohang/:ma_lo/lichsu", lohangH.LichSuLotHang)
		api.PUT("/lohang/:ma_lo/trangthai", lohangH.CapNhatTrangThaiLo)
		api.POST("/lohang/:ma_lo/chungnhan", lohangH.ThemChungNhan)
		api.GET("/lohang/htx/:ma_htx", lohangH.DocLotHangTheoHTX)

		// Nhật ký
		api.POST("/nhatky", nhatkyH.GhiNhatKy)
		api.PUT("/nhatky/:id/duyet", nhatkyH.DuyetNhatKy)
		api.PUT("/nhatky/:id/xacnhan", nhatkyH.XacNhanNhatKy)
		api.GET("/nhatky/lo/:ma_lo", nhatkyH.DocNhatKyTheoLo)
		api.GET("/nhatky/htx/:ma_htx", nhatkyH.DocNhatKyTheoHTX)
		api.GET("/nhatky/:id/lichsu", nhatkyH.LichSuNhatKy)
		api.GET("/nhatky/thongke", nhatkyH.ThongKeNhatKy)

		// Giao dịch
		api.POST("/giaodich", giaodichH.TaoGiaoDich)
		api.PUT("/giaodich/:id/duyet", giaodichH.DuyetGiaoDich)
		api.PUT("/giaodich/:id/trangthai", giaodichH.CapNhatTrangThai)
		api.GET("/giaodich/:id", giaodichH.DocGiaoDich)
		api.GET("/giaodich/:id/lichsu", giaodichH.LichSuGiaoDich)
		api.GET("/giaodich/htx/:ma_htx", giaodichH.DocGiaoDichTheoHTX)
		api.GET("/giaodich/npp/:ma_npp", giaodichH.DocGiaoDichTheoNPP)
		api.GET("/giaodich/npp/:ma_npp/congno", giaodichH.DocCongNoNPP)
		api.GET("/giaodich/npp/:ma_npp/hoahong", giaodichH.TinhHoaHongNPP)
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Printf("Starting Gin server on port %s", cfg.Port)
				if err := r.Run(":" + cfg.Port); err != nil {
					log.Printf("Server failed: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			log.Println("Stopping Server...")
			registry.CloseAll()
			return nil
		},
	})

	return r
}
