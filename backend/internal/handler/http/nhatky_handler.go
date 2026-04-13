package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustfab/gapchain/backend/internal/model"
	"github.com/trustfab/gapchain/backend/internal/usecase"
)

type NhatkyHandler struct {
	uc usecase.NhatkyUsecase
}

func NewNhatkyHandler(uc usecase.NhatkyUsecase) *NhatkyHandler {
	return &NhatkyHandler{uc: uc}
}

func (h *NhatkyHandler) GhiNhatKy(c *gin.Context) {
	mspID := getMspID(c)
	var req model.GhiNhatKyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	maNhatKy, err := h.uc.GhiNhatKy(mspID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ma_nhat_ky": maNhatKy, "message": "Ghi nhat ky thanh cong"})
}

func (h *NhatkyHandler) XacNhanNhatKy(c *gin.Context) {
	mspID := getMspID(c)
	id := c.Param("id")
	var req model.XacNhanNhatKyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.XacNhanNhatKy(mspID, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xac nhan nhat ky thanh cong"})
}

func (h *NhatkyHandler) DuyetNhatKy(c *gin.Context) {
	mspID := getMspID(c)
	id := c.Param("id")
	var req model.DuyetNhatKyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.DuyetNhatKy(mspID, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Duyet nhat ky thanh cong"})
}

func (h *NhatkyHandler) DocNhatKyTheoLo(c *gin.Context) {
	mspID := getMspID(c)
	maLo := c.Param("ma_lo")
	result, err := h.uc.DocNhatKyTheoLo(mspID, maLo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

func (h *NhatkyHandler) DocNhatKyTheoHTX(c *gin.Context) {
	mspID := getMspID(c)
	maHTX := c.Param("ma_htx")
	result, err := h.uc.DocNhatKyTheoHTX(mspID, maHTX)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

func (h *NhatkyHandler) LichSuNhatKy(c *gin.Context) {
	mspID := getMspID(c)
	id := c.Param("id")
	result, err := h.uc.LichSuNhatKy(mspID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}

func (h *NhatkyHandler) ThongKeNhatKy(c *gin.Context) {
	mspID := getMspID(c)
	maHTX := c.Query("ma_htx")
	result, err := h.uc.ThongKeNhatKy(mspID, maHTX)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", result)
}
