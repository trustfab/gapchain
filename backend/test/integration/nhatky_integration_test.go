package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	myhttp "github.com/trustfab/gapchain/backend/internal/handler/http"
	infra "github.com/trustfab/gapchain/backend/internal/infrastructure/fabric"
	"github.com/trustfab/gapchain/backend/internal/model"
	repo "github.com/trustfab/gapchain/backend/internal/repository/fabric"
	"github.com/trustfab/gapchain/backend/internal/usecase"
)

var _ = Describe("Nhatky Integration API", Label("integration"), func() {
	var (
		registry     *infra.GatewayRegistry
		mockContract *MockFabricContract // using the same mock struct defined in lohang_integration_test.go
		router       *gin.Engine
		w            *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		
		mockContract = &MockFabricContract{}
		registry = &infra.GatewayRegistry{
			FallbackMspID: "PlatformOrgMSP",
		}
		
		// Setup cho HTX Nông Sản
		registry.SetOrgGatewayForTest("HTXNongSanOrgMSP", &infra.OrgGateway{
			LotHangContract: mockContract,
			NhatKyContract:  mockContract,
		})

		// Setup cho Chi Cục BVTV
		registry.SetOrgGatewayForTest("ChiCucBVTVOrgMSP", &infra.OrgGateway{
			LotHangContract: mockContract,
			NhatKyContract:  mockContract,
		})

		nhatkyRepo := repo.NewNhatkyRepo(registry)
		lohangRepo := repo.NewLohangRepo(registry)
		nhatkyUC := usecase.NewNhatkyUsecase(nhatkyRepo, lohangRepo)
		nhatkyHandler := myhttp.NewNhatkyHandler(nhatkyUC)

		router = gin.New()
		
		// Giả lập Middleware lấy JWT -> mspID. 
		router.Use(func(c *gin.Context) {
			msp := c.GetHeader("X-MSP-ID")
			if msp == "" {
				msp = "HTXNongSanOrgMSP" // mặc định user login là HTX
			}
			c.Set("mspID", msp)
			c.Next()
		})

		router.POST("/nhatky", nhatkyHandler.GhiNhatKy)
		router.PUT("/nhatky/:id/duyet", nhatkyHandler.DuyetNhatKy)
		router.GET("/nhatky/lo/:ma_lo", nhatkyHandler.DocNhatKyTheoLo)

		w = httptest.NewRecorder()
	})

	It("thực hiện luồng GhiNhatKy đi xuyên qua 3 layers", func() {
		payload := model.GhiNhatKyReq{
			MaNhatKy:      "NK_001",
			MaLo:          "LO_INT_999",
			MaHTX:         "HTX_01",
			LoaiHoatDong:  "bon_phan",
			ChiTiet:       "Bón phân NPK",
			ViTri:         "Khu A",
			NguoiThucHien: "Nguyễn Văn A",
			NgayGhi:       "2026-04-29",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/nhatky", bytes.NewReader(body))
		// Không truyền header X-MSP-ID => dùng msp mặc định là HTXNongSanOrgMSP

		calledFabric := false
		mockContract.EvaluateTransactionFunc = func(name string, args ...string) ([]byte, error) {
			return []byte(`{"ma_lo": "LO_INT_999", "ma_htx": "HTX_01", "trang_thai": "dang_trong"}`), nil
		}

		mockContract.SubmitTransactionFunc = func(name string, args ...string) ([]byte, error) {
			calledFabric = true
			Expect(name).To(Equal("GhiNhatKy"))
			Expect(args[0]).To(Equal("NK_001"))
			Expect(args[1]).To(Equal("LO_INT_999"))
			Expect(args[3]).To(Equal("bon_phan"))
			
			return []byte("SUCCESS"), nil
		}

		router.ServeHTTP(w, req)

		Expect(calledFabric).To(BeTrue(), "Error: "+w.Body.String())
		Expect(w.Code).To(Equal(http.StatusCreated))
		Expect(w.Body.String()).To(ContainSubstring("NK_001"))
	})

	It("thực hiện luồng Duyệt Nhật Ký bởi Chi Cục BVTV đi xuyên qua 3 layers", func() {
		payload := model.DuyetNhatKyReq{
			QuyetDinh:  "duyet",
			LyDoTuChoi: "",
			NguoiDuyet: "Cán Bộ B",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", "/nhatky/NK_001/duyet", bytes.NewReader(body))
		
		// Set header giả lập login bằng user của Chi Cục BVTV
		req.Header.Set("X-MSP-ID", "ChiCucBVTVOrgMSP")

		calledFabric := false
		
		// Setup mock EvaluateTransaction to simulate that the Lot Hang exists
		mockContract.EvaluateTransactionFunc = func(name string, args ...string) ([]byte, error) {
			return []byte(`{"ma_lo": "LO_INT_999", "ma_htx": "HTX_01", "trang_thai": "Đang trồng"}`), nil
		}
		mockContract.SubmitTransactionFunc = func(name string, args ...string) ([]byte, error) {
			calledFabric = true
			Expect(name).To(Equal("DuyetNhatKy"))
			Expect(args[0]).To(Equal("NK_001"))
			Expect(args[1]).To(Equal("Cán Bộ B"))
			Expect(args[2]).To(Equal("duyet"))
			
			return []byte("SUCCESS"), nil
		}

		router.ServeHTTP(w, req)

		Expect(calledFabric).To(BeTrue())
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(ContainSubstring("Duyet nhat ky thanh cong"))
	})
})
