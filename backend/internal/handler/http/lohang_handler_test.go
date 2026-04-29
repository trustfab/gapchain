package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	myhttp "github.com/trustfab/gapchain/backend/internal/handler/http"
	"github.com/trustfab/gapchain/backend/internal/model"
)

// MockLohangUsecase implements usecase.LohangUsecase for testing
type MockLohangUsecase struct {
	TaoLotHangFunc        func(mspID string, req *model.TaoLotHangReq) (string, error)
	TachLoFunc            func(mspID string, maLoMe string, req *model.TachLoReq) error
	CapNhatTrangThaiFunc  func(mspID string, maLo string, trangThai string) error
	ThemChungNhanFunc     func(mspID string, maLo string, req *model.ThemChungNhanReq) error
	DocLotHangFunc        func(mspID string, maLo string) ([]byte, error)
	LichSuLotHangFunc     func(mspID string, maLo string) ([]byte, error)
	DocLotHangTheoHTXFunc func(mspID string, maHTX string) ([]byte, error)
	LayThongTinTraCuuFunc func(maLo string) ([]byte, error)
	CapNhatInventoryFunc  func(mspID string, maLo string, thayDoi float64) error
}

func (m *MockLohangUsecase) TaoLotHang(mspID string, req *model.TaoLotHangReq) (string, error) {
	if m.TaoLotHangFunc != nil {
		return m.TaoLotHangFunc(mspID, req)
	}
	return "", nil
}
func (m *MockLohangUsecase) TachLo(mspID string, maLoMe string, req *model.TachLoReq) error {
	if m.TachLoFunc != nil {
		return m.TachLoFunc(mspID, maLoMe, req)
	}
	return nil
}
func (m *MockLohangUsecase) CapNhatTrangThai(mspID string, maLo string, trangThai string) error {
	if m.CapNhatTrangThaiFunc != nil {
		return m.CapNhatTrangThaiFunc(mspID, maLo, trangThai)
	}
	return nil
}
func (m *MockLohangUsecase) ThemChungNhan(mspID string, maLo string, req *model.ThemChungNhanReq) error {
	if m.ThemChungNhanFunc != nil {
		return m.ThemChungNhanFunc(mspID, maLo, req)
	}
	return nil
}
func (m *MockLohangUsecase) DocLotHang(mspID string, maLo string) ([]byte, error) {
	if m.DocLotHangFunc != nil {
		return m.DocLotHangFunc(mspID, maLo)
	}
	return nil, nil
}
func (m *MockLohangUsecase) LichSuLotHang(mspID string, maLo string) ([]byte, error) {
	if m.LichSuLotHangFunc != nil {
		return m.LichSuLotHangFunc(mspID, maLo)
	}
	return nil, nil
}
func (m *MockLohangUsecase) DocLotHangTheoHTX(mspID string, maHTX string) ([]byte, error) {
	if m.DocLotHangTheoHTXFunc != nil {
		return m.DocLotHangTheoHTXFunc(mspID, maHTX)
	}
	return nil, nil
}
func (m *MockLohangUsecase) LayThongTinTraCuu(maLo string) ([]byte, error) {
	if m.LayThongTinTraCuuFunc != nil {
		return m.LayThongTinTraCuuFunc(maLo)
	}
	return nil, nil
}
func (m *MockLohangUsecase) CapNhatInventory(mspID string, maLo string, thayDoi float64) error {
	if m.CapNhatInventoryFunc != nil {
		return m.CapNhatInventoryFunc(mspID, maLo, thayDoi)
	}
	return nil
}

var _ = Describe("LohangHandler", func() {
	var (
		mockUC  *MockLohangUsecase
		handler *myhttp.LohangHandler
		router  *gin.Engine
		w       *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockUC = &MockLohangUsecase{}
		handler = myhttp.NewLohangHandler(mockUC)
		router = gin.New()
		
		// Giả lập middleware lấy mspID
		router.Use(func(c *gin.Context) {
			msp := c.GetHeader("X-MSP-ID")
			if msp == "" {
				msp = "HTXNongSanOrgMSP" // default for private routes
			}
			c.Set("mspID", msp)
			c.Next()
		})
		
		router.POST("/lohang", handler.TaoLotHang)
		router.POST("/lohang/:ma_lo/tach", handler.TachLo)
		router.PUT("/lohang/:ma_lo/trangthai", handler.CapNhatTrangThaiLo)
		router.POST("/lohang/:ma_lo/chungnhan", handler.ThemChungNhan)
		router.GET("/lohang/:ma_lo", handler.DocLotHang)
		router.GET("/lohang/:ma_lo/lichsu", handler.LichSuLotHang)
		router.GET("/lohang/htx/:ma_htx", handler.DocLotHangTheoHTX)
		router.GET("/consumer/:ma_lo", handler.LayThongTinTraCuu)
		
		w = httptest.NewRecorder()
	})

	Describe("POST /lohang (TaoLotHang)", func() {
		Context("Happy Path", func() {
			It("trả về 201 và ma_lo khi payload hợp lệ", func() {
				mockUC.TaoLotHangFunc = func(mspID string, req *model.TaoLotHangReq) (string, error) {
					Expect(mspID).To(Equal("HTXNongSanOrgMSP"))
					Expect(req.MaLo).To(Equal("LO_123"))
					return "LO_123", nil
				}

				payload := model.TaoLotHangReq{
					MaLo: "LO_123",
					MaHTX: "HTX_01",
					TenSanPham: "Xoài Cát Chu",
					LoaiSanPham: "Trái cây",
					SoLuong: 100,
					DonViTinh: "kg",
				}
				body, _ := json.Marshal(payload)
				req, _ := http.NewRequest("POST", "/lohang", bytes.NewReader(body))
				
				router.ServeHTTP(w, req)
				
				Expect(w.Code).To(Equal(http.StatusCreated))
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				Expect(resp["ma_lo"]).To(Equal("LO_123"))
				Expect(resp["message"]).To(Equal("Tao lo hang thanh cong"))
			})
		})

		Context("Validation Error", func() {
			It("trả về 400 khi payload sai định dạng", func() {
				req, _ := http.NewRequest("POST", "/lohang", bytes.NewReader([]byte(`{bad json}`)))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("Usecase Error", func() {
			It("trả về 500 khi usecase (Fabric) trả về lỗi", func() {
				mockUC.TaoLotHangFunc = func(mspID string, req *model.TaoLotHangReq) (string, error) {
					return "", errors.New("chaincode error: unauthorized")
				}

				payload := model.TaoLotHangReq{
					MaLo: "LO_123",
					MaHTX: "HTX_01",
					TenSanPham: "Xoài Cát Chu",
					LoaiSanPham: "Trái cây",
					SoLuong: 100,
					DonViTinh: "kg",
				}
				body, _ := json.Marshal(payload)
				req, _ := http.NewRequest("POST", "/lohang", bytes.NewReader(body))
				
				router.ServeHTTP(w, req)
				
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(ContainSubstring("chaincode error: unauthorized"))
			})
		})
	})
    
    Describe("GET /lohang/:ma_lo (DocLotHang)", func() {
		Context("khi lô hàng tồn tại", func() {
			It("trả về 200 và dữ liệu lô hàng", func() {
				mockUC.DocLotHangFunc = func(mspID string, maLo string) ([]byte, error) {
					Expect(maLo).To(Equal("LO_123"))
					return []byte(`{"ma_lo": "LO_123", "ten_san_pham": "Xoài"}`), nil
				}

				req, _ := http.NewRequest("GET", "/lohang/LO_123", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("Xoài"))
			})
		})
        Context("khi lô hàng không tồn tại", func() {
            It("trả về 404 Not Found", func() {
                mockUC.DocLotHangFunc = func(mspID string, maLo string) ([]byte, error) {
					return nil, errors.New("lo hang khong ton tai")
				}

				req, _ := http.NewRequest("GET", "/lohang/LO_999", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
            })
        })
	})
    
    Describe("GET /consumer/:ma_lo (LayThongTinTraCuu - Public API)", func() {
		It("trả về thông tin tra cứu thành công", func() {
			mockUC.LayThongTinTraCuuFunc = func(maLo string) ([]byte, error) {
                Expect(maLo).To(Equal("LO_123"))
				return []byte(`{"lo_hang": {"ma_lo": "LO_123"}}`), nil
			}

			req, _ := http.NewRequest("GET", "/consumer/LO_123", nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("LO_123"))
		})
	})
})
