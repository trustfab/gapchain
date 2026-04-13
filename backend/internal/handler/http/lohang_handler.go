package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustfab/gapchain/backend/internal/model"
	"github.com/trustfab/gapchain/backend/internal/usecase"
)

type LohangHandler struct {
	uc usecase.LohangUsecase
}

func NewLohangHandler(uc usecase.LohangUsecase) *LohangHandler {
	return &LohangHandler{uc: uc}
}

func (h *LohangHandler) TaoLotHang(c *gin.Context) {
	mspID := getMspID(c)
	var req model.TaoLotHangReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	maLo, err := h.uc.TaoLotHang(mspID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ma_lo": maLo, "message": "Tao lo hang thanh cong"})
}

func (h *LohangHandler) TachLo(c *gin.Context) {
	mspID := getMspID(c)
	maLoMe := c.Param("ma_lo")
	var req model.TachLoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.TachLo(mspID, maLoMe, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ma_lo_moi": req.MaLoMoi, "message": "Tach lo hang thanh cong"})
}

func (h *LohangHandler) CapNhatTrangThaiLo(c *gin.Context) {
	mspID := getMspID(c)
	maLo := c.Param("ma_lo")
	var req model.CapNhatTrangThaiLoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.CapNhatTrangThai(mspID, maLo, req.TrangThai)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cap nhat trang thai lo hang thanh cong"})
}

func (h *LohangHandler) ThemChungNhan(c *gin.Context) {
	mspID := getMspID(c)
	maLo := c.Param("ma_lo")
	var req model.ThemChungNhanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.ThemChungNhan(mspID, maLo, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Them chung nhan thanh cong"})
}

func (h *LohangHandler) DocLotHang(c *gin.Context) {
	mspID := getMspID(c)
	maLo := c.Param("ma_lo")
	result, err := h.uc.DocLotHang(mspID, maLo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

func (h *LohangHandler) LichSuLotHang(c *gin.Context) {
	mspID := getMspID(c)
	maLo := c.Param("ma_lo")
	result, err := h.uc.LichSuLotHang(mspID, maLo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

func (h *LohangHandler) DocLotHangTheoHTX(c *gin.Context) {
	mspID := getMspID(c)
	maHTX := c.Param("ma_htx")
	result, err := h.uc.DocLotHangTheoHTX(mspID, maHTX)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

// LayThongTinTraCuu - Public endpoint, không cần JWT → dùng Platform fallback
func (h *LohangHandler) LayThongTinTraCuu(c *gin.Context) {
	maLo := c.Param("ma_lo")
	resultBytes, err := h.uc.LayThongTinTraCuu(maLo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lo hang khong ton tai"})
		return
	}

	c.JSON(http.StatusOK, json.RawMessage(resultBytes))
}
