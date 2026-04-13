package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustfab/gapchain/backend/internal/model"
	"github.com/trustfab/gapchain/backend/internal/usecase"
)

type GiaodichHandler struct {
	uc usecase.GiaodichUsecase
}

func NewGiaodichHandler(uc usecase.GiaodichUsecase) *GiaodichHandler {
	return &GiaodichHandler{uc: uc}
}

func (h *GiaodichHandler) TaoGiaoDich(c *gin.Context) {
	mspID := getMspID(c)
	var req model.TaoGiaoDichReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	maGD, err := h.uc.TaoGiaoDich(mspID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ma_giao_dich": maGD, "message": "Tao giao dich thanh cong"})
}

func (h *GiaodichHandler) DuyetGiaoDich(c *gin.Context) {
	mspID := getMspID(c)
	id := c.Param("id")
	err := h.uc.DuyetGiaoDich(mspID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Duyet giao dich thanh cong"})
}

func (h *GiaodichHandler) CapNhatTrangThai(c *gin.Context) {
	mspID := getMspID(c)
	id := c.Param("id")
	var req model.CapNhatTrangThaiGiaoDichReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.uc.CapNhatTrangThai(mspID, id, req.TrangThai)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cap nhat trang thai giao dich thanh cong"})
}

func (h *GiaodichHandler) DocGiaoDich(c *gin.Context) {
	mspID := getMspID(c)
	id := c.Param("id")
	result, err := h.uc.DocGiaoDich(mspID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

func (h *GiaodichHandler) LichSuGiaoDich(c *gin.Context) {
	mspID := getMspID(c)
	id := c.Param("id")
	result, err := h.uc.LichSuGiaoDich(mspID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

func (h *GiaodichHandler) DocCongNoNPP(c *gin.Context) {
	mspID := getMspID(c)
	maNPP := c.Param("ma_npp")
	result, err := h.uc.DocCongNoNPP(mspID, maNPP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

func (h *GiaodichHandler) TinhHoaHongNPP(c *gin.Context) {
	mspID := getMspID(c)
	maNPP := c.Param("ma_npp")
	result, err := h.uc.TinhHoaHongNPP(mspID, maNPP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

func (h *GiaodichHandler) DocGiaoDichTheoHTX(c *gin.Context) {
	mspID := getMspID(c)
	maHTX := c.Param("ma_htx")
	result, err := h.uc.DocGiaoDichTheoHTX(mspID, maHTX)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

func (h *GiaodichHandler) DocGiaoDichTheoNPP(c *gin.Context) {
	mspID := getMspID(c)
	maNPP := c.Param("ma_npp")
	result, err := h.uc.DocGiaoDichTheoNPP(mspID, maNPP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}
