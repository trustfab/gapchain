package usecase_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/trustfab/gapchain/backend/internal/model"
	"github.com/trustfab/gapchain/backend/internal/usecase"
)

// MockLohangRepo implements repo.LohangRepo
type MockLohangRepo struct {
	SubmitFunc   func(mspID string, funcName string, args ...string) ([]byte, error)
	EvaluateFunc func(mspID string, funcName string, args ...string) ([]byte, error)
}

func (m *MockLohangRepo) Submit(mspID string, funcName string, args ...string) ([]byte, error) {
	if m.SubmitFunc != nil {
		return m.SubmitFunc(mspID, funcName, args...)
	}
	return nil, nil
}
func (m *MockLohangRepo) Evaluate(mspID string, funcName string, args ...string) ([]byte, error) {
	if m.EvaluateFunc != nil {
		return m.EvaluateFunc(mspID, funcName, args...)
	}
	return nil, nil
}

// MockNhatkyRepo implements repo.NhatkyRepo
type MockNhatkyRepo struct {
	SubmitFunc   func(mspID string, funcName string, args ...string) ([]byte, error)
	EvaluateFunc func(mspID string, funcName string, args ...string) ([]byte, error)
}

func (m *MockNhatkyRepo) Submit(mspID string, funcName string, args ...string) ([]byte, error) {
	if m.SubmitFunc != nil {
		return m.SubmitFunc(mspID, funcName, args...)
	}
	return nil, nil
}
func (m *MockNhatkyRepo) Evaluate(mspID string, funcName string, args ...string) ([]byte, error) {
	if m.EvaluateFunc != nil {
		return m.EvaluateFunc(mspID, funcName, args...)
	}
	return nil, nil
}

var _ = Describe("LohangUsecase", func() {
	var (
		mockLR *MockLohangRepo
		mockNR *MockNhatkyRepo
		uc     usecase.LohangUsecase
	)

	BeforeEach(func() {
		mockLR = &MockLohangRepo{}
		mockNR = &MockNhatkyRepo{}
		uc = usecase.NewLohangUsecase(mockLR, mockNR)
	})

	Describe("TaoLotHang", func() {
		Context("khi dữ liệu hợp lệ", func() {
			It("gọi repo Submit và trả về ma_lo", func() {
				req := &model.TaoLotHangReq{
					MaLo:        "LO_123",
					MaLoMe:      "",
					MaHTX:       "HTX_01",
					TenSanPham:  "Xoài",
					LoaiSanPham: "Trái cây",
					SoLuong:     100,
					DonViTinh:   "kg",
					VuMua:       "Mùa Hè 2026",
					DiaDiem:     "Đồng Tháp",
				}

				mockLR.SubmitFunc = func(mspID string, funcName string, args ...string) ([]byte, error) {
					Expect(mspID).To(Equal("HTXNongSanOrgMSP"))
					Expect(funcName).To(Equal("TaoLotHang"))
					Expect(args[0]).To(Equal("LO_123"))
					Expect(args[1]).To(Equal("")) // MaLoMe
					Expect(args[5]).To(Equal("100")) // SoLuong formatted as %g
					return []byte("OK"), nil
				}

				maLo, err := uc.TaoLotHang("HTXNongSanOrgMSP", req)
				Expect(err).NotTo(HaveOccurred())
				Expect(maLo).To(Equal("LO_123"))
			})
		})

		Context("khi repo Submit trả về lỗi", func() {
			It("trả về lỗi tương ứng", func() {
				mockLR.SubmitFunc = func(mspID string, funcName string, args ...string) ([]byte, error) {
					return nil, errors.New("lỗi chaincode")
				}

				_, err := uc.TaoLotHang("HTXNongSanOrgMSP", &model.TaoLotHangReq{MaLo: "LO_123"})
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("lỗi chaincode"))
			})
		})
	})

	Describe("LayThongTinTraCuu", func() {
		Context("khi lô hàng là lô gốc (không có MaLoMe)", func() {
			It("trả về thông tin lô hàng và nhật ký", func() {
				maLo := "LO_123"
				
				// Mock Evaluate của LohangRepo
				mockLR.EvaluateFunc = func(mspID string, funcName string, args ...string) ([]byte, error) {
					Expect(mspID).To(Equal("PlatformOrgMSP")) // fallback msp
					if funcName == "LayThongTinTraCuu" {
						Expect(args[0]).To(Equal(maLo))
						return []byte(`{"ma_lo": "LO_123", "ten_san_pham": "Xoài"}`), nil
					}
					return nil, nil
				}

				// Mock Evaluate của NhatkyRepo
				mockNR.EvaluateFunc = func(mspID string, funcName string, args ...string) ([]byte, error) {
					Expect(mspID).To(Equal("PlatformOrgMSP"))
					if funcName == "DocNhatKyTheoLo" {
						Expect(args[0]).To(Equal(maLo))
						return []byte(`[{"hoat_dong": "Thu hoạch"}]`), nil
					}
					return nil, nil
				}

				resultBytes, err := uc.LayThongTinTraCuu(maLo)
				Expect(err).NotTo(HaveOccurred())

				var result map[string]interface{}
				err = json.Unmarshal(resultBytes, &result)
				Expect(err).NotTo(HaveOccurred())
				
				Expect(result).To(HaveKey("lo_hang"))
				Expect(result).To(HaveKey("nhat_ky"))
				Expect(result).NotTo(HaveKey("lo_hang_me"))
			})
		})

		Context("khi lô hàng là lô tách (có MaLoMe)", func() {
			It("trả về thêm thông tin lo_hang_me và nhat_ky_me", func() {
				maLo := "LO_CHILD"
				maLoMe := "LO_PARENT"

				mockLR.EvaluateFunc = func(mspID string, funcName string, args ...string) ([]byte, error) {
					if funcName == "LayThongTinTraCuu" {
						if args[0] == maLo {
							return []byte(`{"ma_lo": "LO_CHILD", "ma_lo_me": "LO_PARENT"}`), nil
						}
						if args[0] == maLoMe {
							return []byte(`{"ma_lo": "LO_PARENT", "ten_san_pham": "Xoài"}`), nil
						}
					}
					return nil, nil
				}

				mockNR.EvaluateFunc = func(mspID string, funcName string, args ...string) ([]byte, error) {
					if funcName == "DocNhatKyTheoLo" {
						if args[0] == maLo {
							return []byte(`[{"hoat_dong": "Đóng gói"}]`), nil
						}
						if args[0] == maLoMe {
							return []byte(`[{"hoat_dong": "Thu hoạch"}]`), nil
						}
					}
					return nil, nil
				}

				resultBytes, err := uc.LayThongTinTraCuu(maLo)
				Expect(err).NotTo(HaveOccurred())

				var result map[string]interface{}
				json.Unmarshal(resultBytes, &result)
				
				Expect(result).To(HaveKey("lo_hang"))
				Expect(result).To(HaveKey("nhat_ky"))
				Expect(result).To(HaveKey("lo_hang_me"))
				Expect(result).To(HaveKey("nhat_ky_me"))
			})
		})

		Context("khi không tìm thấy lô hàng", func() {
			It("trả về lỗi", func() {
				mockLR.EvaluateFunc = func(mspID string, funcName string, args ...string) ([]byte, error) {
					return nil, errors.New("không tìm thấy")
				}

				_, err := uc.LayThongTinTraCuu("LO_999")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("không tìm thấy"))
			})
		})
	})
})
