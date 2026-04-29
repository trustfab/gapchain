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

// MockFabricContract giả lập FabricGateway
type MockFabricContract struct {
	SubmitTransactionFunc   func(name string, args ...string) ([]byte, error)
	EvaluateTransactionFunc func(name string, args ...string) ([]byte, error)
}

func (m *MockFabricContract) SubmitTransaction(name string, args ...string) ([]byte, error) {
	if m.SubmitTransactionFunc != nil {
		return m.SubmitTransactionFunc(name, args...)
	}
	return []byte(""), nil
}

func (m *MockFabricContract) EvaluateTransaction(name string, args ...string) ([]byte, error) {
	if m.EvaluateTransactionFunc != nil {
		return m.EvaluateTransactionFunc(name, args...)
	}
	return []byte(""), nil
}

var _ = Describe("Lohang Integration API", Label("integration"), func() {
	var (
		registry     *infra.GatewayRegistry
		mockContract *MockFabricContract
		router       *gin.Engine
		w            *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		
		// 1. Khởi tạo Mock duy nhất là ranh giới ngoài cùng (Fabric)
		mockContract = &MockFabricContract{}

		// 2. Setup GatewayRegistry thật (nhưng bỏ qua việc kết nối gRPC)
		registry = &infra.GatewayRegistry{
			FallbackMspID: "PlatformOrgMSP",
		}
		
		// Bơm mockContract vào HTXNongSanOrgMSP (cho API private)
		registry.SetOrgGatewayForTest("HTXNongSanOrgMSP", &infra.OrgGateway{
			LotHangContract: mockContract,
			NhatKyContract:  mockContract,
		})
		
		// Bơm mockContract vào PlatformOrgMSP (cho API public)
		registry.SetOrgGatewayForTest("PlatformOrgMSP", &infra.OrgGateway{
			LotHangContract: mockContract,
			NhatKyContract:  mockContract,
		})

		// 3. Khởi tạo Repository thật
		lohangRepo := repo.NewLohangRepo(registry)
		nhatkyRepo := repo.NewNhatkyRepo(registry)

		// 4. Khởi tạo Usecase thật
		lohangUC := usecase.NewLohangUsecase(lohangRepo, nhatkyRepo)

		// 5. Khởi tạo Handler thật
		lohangHandler := myhttp.NewLohangHandler(lohangUC)

		// 6. Gắn vào Gin router
		router = gin.New()
		
		// Giả lập Middleware lấy JWT -> mspID
		router.Use(func(c *gin.Context) {
			msp := c.GetHeader("X-MSP-ID")
			if msp == "" {
				msp = "HTXNongSanOrgMSP" // mặc định user login là HTX Nông Sản
			}
			c.Set("mspID", msp)
			c.Next()
		})

		router.POST("/lohang", lohangHandler.TaoLotHang)
		router.GET("/lohang/:ma_lo", lohangHandler.DocLotHang)
		router.GET("/consumer/:ma_lo", lohangHandler.LayThongTinTraCuu)

		w = httptest.NewRecorder()
	})

	It("thực hiện luồng TaoLotHang đi xuyên qua 3 layers", func() {
		// 1. HTTP payload đầu vào
		payload := model.TaoLotHangReq{
			MaLo:        "LO_INT_999",
			MaHTX:       "HTX_01",
			TenSanPham:  "Dưa Hấu",
			LoaiSanPham: "Trái Cây",
			SoLuong:     150,
			DonViTinh:   "kg",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/lohang", bytes.NewReader(body))

		// 2. Chặn ở cuối đường ống (MockFabricContract) để verify dữ liệu đã chui xuống đúng
		calledFabric := false
		mockContract.SubmitTransactionFunc = func(name string, args ...string) ([]byte, error) {
			calledFabric = true
			Expect(name).To(Equal("TaoLotHang"))
			// args từ Usecase truyền xuống Repo, rồi tới Fabric
			Expect(args[0]).To(Equal("LO_INT_999")) // MaLo
			Expect(args[2]).To(Equal("HTX_01"))     // MaHTX
			Expect(args[3]).To(Equal("Dưa Hấu"))    // TenSanPham
			Expect(args[5]).To(Equal("150"))        // SoLuong (dạng chuỗi do fmt.Sprintf("%g"))
			
			return []byte("SUCCESS_SUBMIT"), nil
		}

		// 3. Thực thi
		router.ServeHTTP(w, req)

		// 4. Assert đầu ra ở tầng HTTP
		Expect(calledFabric).To(BeTrue(), "Nên gọi tới Fabric Contract")
		Expect(w.Code).To(Equal(http.StatusCreated))
		Expect(w.Body.String()).To(ContainSubstring("LO_INT_999"))
		Expect(w.Body.String()).To(ContainSubstring("Tao lo hang thanh cong"))
	})

	It("thực hiện luồng Public API (LayThongTinTraCuu) kết hợp nhiều Repo gọi Evaluate", func() {
		req, _ := http.NewRequest("GET", "/consumer/LO_INT_999", nil)
		
		lohangCalled := false
		nhatkyCalled := false
		
		mockContract.EvaluateTransactionFunc = func(name string, args ...string) ([]byte, error) {
			Expect(args[0]).To(Equal("LO_INT_999"))
			if name == "LayThongTinTraCuu" { // Lời gọi từ LohangRepo
				lohangCalled = true
				return []byte(`{"ma_lo": "LO_INT_999", "ma_lo_me": ""}`), nil
			}
			if name == "DocNhatKyTheoLo" { // Lời gọi từ NhatkyRepo
				nhatkyCalled = true
				return []byte(`[{"id": "NK1", "hoat_dong": "Thu hoach"}]`), nil
			}
			return nil, nil
		}

		router.ServeHTTP(w, req)

		Expect(lohangCalled).To(BeTrue(), "Usecase nên gọi LohangRepo")
		Expect(nhatkyCalled).To(BeTrue(), "Usecase nên gọi NhatkyRepo để gộp vào trả về")
		Expect(w.Code).To(Equal(http.StatusOK))
		
		// Assert kết quả trả về là một object tổng hợp lo_hang và nhat_ky (logic của usecase)
		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		Expect(result).To(HaveKey("lo_hang"))
		Expect(result).To(HaveKey("nhat_ky"))
	})
})
