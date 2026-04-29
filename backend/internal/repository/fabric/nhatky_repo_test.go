package fabric

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	infra "github.com/trustfab/gapchain/backend/internal/infrastructure/fabric"
)

// MockFabricContract implements infra.FabricContract
type MockFabricContract struct {
	SubmitTransactionFunc   func(name string, args ...string) ([]byte, error)
	EvaluateTransactionFunc func(name string, args ...string) ([]byte, error)
}

func (m *MockFabricContract) SubmitTransaction(name string, args ...string) ([]byte, error) {
	if m.SubmitTransactionFunc != nil {
		return m.SubmitTransactionFunc(name, args...)
	}
	return nil, nil
}

func (m *MockFabricContract) EvaluateTransaction(name string, args ...string) ([]byte, error) {
	if m.EvaluateTransactionFunc != nil {
		return m.EvaluateTransactionFunc(name, args...)
	}
	return nil, nil
}

var _ = Describe("NhatkyRepo", func() {
	var (
		registry     *infra.GatewayRegistry
		mockContract *MockFabricContract
		repo         NhatkyRepo
	)

	BeforeEach(func() {
		// Tạo registry rỗng để test
		registry = &infra.GatewayRegistry{}
		mockContract = &MockFabricContract{}
		repo = NewNhatkyRepo(registry)
	})

	Describe("Submit", func() {
		Context("khi mspID không tồn tại trong registry", func() {
			It("trả về lỗi không tìm thấy gateway", func() {
				_, err := repo.Submit("UnknownMSP", "GhiNhatKy", "arg1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("khong tim thay gateway"))
			})
		})

		Context("khi org tồn tại nhưng NhatKyContract là nil (không tham gia channel)", func() {
			It("trả về lỗi không tham gia channel", func() {
				// Inject OrgGateway không có NhatKyContract
				registry.SetOrgGatewayForTest("NPPXanhOrgMSP", &infra.OrgGateway{})

				_, err := repo.Submit("NPPXanhOrgMSP", "GhiNhatKy", "arg1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("khong tham gia nhatky-htx-channel"))
			})
		})

		Context("khi org hợp lệ và có quyền truy cập", func() {
			It("gọi SubmitTransaction của contract và trả về kết quả", func() {
				registry.SetOrgGatewayForTest("HTXNongSanOrgMSP", &infra.OrgGateway{
					NhatKyContract: mockContract,
				})

				mockContract.SubmitTransactionFunc = func(name string, args ...string) ([]byte, error) {
					Expect(name).To(Equal("GhiNhatKy"))
					Expect(args[0]).To(Equal("arg1"))
					return []byte("SUCCESS"), nil
				}

				res, err := repo.Submit("HTXNongSanOrgMSP", "GhiNhatKy", "arg1")
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal([]byte("SUCCESS")))
			})

			It("trả về lỗi nếu SubmitTransaction thất bại", func() {
				registry.SetOrgGatewayForTest("HTXNongSanOrgMSP", &infra.OrgGateway{
					NhatKyContract: mockContract,
				})

				mockContract.SubmitTransactionFunc = func(name string, args ...string) ([]byte, error) {
					return nil, errors.New("chaincode panic")
				}

				_, err := repo.Submit("HTXNongSanOrgMSP", "GhiNhatKy", "arg1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("chaincode panic"))
			})
		})
	})

	Describe("Evaluate", func() {
		Context("khi org hợp lệ và có quyền truy cập", func() {
			It("gọi EvaluateTransaction của contract và trả về kết quả", func() {
				registry.SetOrgGatewayForTest("PlatformOrgMSP", &infra.OrgGateway{
					NhatKyContract: mockContract,
				})

				mockContract.EvaluateTransactionFunc = func(name string, args ...string) ([]byte, error) {
					Expect(name).To(Equal("DocNhatKyTheoLo"))
					Expect(args[0]).To(Equal("LO_123"))
					return []byte(`[{"id":"1"}]`), nil
				}

				res, err := repo.Evaluate("PlatformOrgMSP", "DocNhatKyTheoLo", "LO_123")
				Expect(err).NotTo(HaveOccurred())
				Expect(string(res)).To(Equal(`[{"id":"1"}]`))
			})
		})
	})
})
