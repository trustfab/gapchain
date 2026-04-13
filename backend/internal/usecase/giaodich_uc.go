package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/trustfab/gapchain/backend/internal/infrastructure"
	"github.com/trustfab/gapchain/backend/internal/model"
	repo "github.com/trustfab/gapchain/backend/internal/repository/fabric"
)

type GiaodichUsecase interface {
	TaoGiaoDich(mspID string, req *model.TaoGiaoDichReq) (string, error)
	DuyetGiaoDich(mspID string, id string) error
	CapNhatTrangThai(mspID string, id string, trangThai string) error
	DocGiaoDich(mspID string, id string) ([]byte, error)
	LichSuGiaoDich(mspID string, id string) ([]byte, error)
	DocCongNoNPP(mspID string, maNPP string) ([]byte, error)
	TinhHoaHongNPP(mspID string, maNPP string) ([]byte, error)
	DocGiaoDichTheoHTX(mspID string, maHTX string) ([]byte, error)
	DocGiaoDichTheoNPP(mspID string, maNPP string) ([]byte, error)
}

type giaodichUsecase struct {
	repo       repo.GiaodichRepo
	lohangRepo repo.LohangRepo
}

func NewGiaodichUsecase(r repo.GiaodichRepo, lr repo.LohangRepo) GiaodichUsecase {
	return &giaodichUsecase{repo: r, lohangRepo: lr}
}

func (uc *giaodichUsecase) TaoGiaoDich(mspID string, req *model.TaoGiaoDichReq) (string, error) {
	// 1. Kiểm tra Lô hàng
	lotHangBytes, err := uc.lohangRepo.Evaluate(mspID, "DocLotHang", req.MaLo)
	if err != nil {
		return "", fmt.Errorf("không thể truy xuất lô hàng từ channel: %v", err)
	}

	var loHang model.LotHangDTO
	if err := json.Unmarshal(lotHangBytes, &loHang); err != nil {
		return "", fmt.Errorf("lỗi parse dữ liệu lô hàng: %v", err)
	}

	// 2. Validate inventory (Theo rule R4, R5 trong SKILL.md)
	if loHang.TrangThai != "san_sang_ban" {
		return "", fmt.Errorf("lô hàng chưa sẵn sàng bán (trạng thái: %s)", loHang.TrangThai)
	}
	if loHang.SoLuongConLai < req.SoLuong {
		return "", fmt.Errorf("số lượng giao dịch (%g) vượt số lượng còn lại của lô (%g)", req.SoLuong, loHang.SoLuongConLai)
	}

	// 3. Thực thi tạo giao dịch ở channel giao dịch
	_, err = uc.repo.Submit(
		mspID, "TaoGiaoDich",
		req.MaGiaoDich, req.MaLo, req.MaHTX, req.MaNPP, req.SanPham,
		fmt.Sprintf("%g", req.SoLuong), req.DonViTinh, fmt.Sprintf("%g", req.DonGia), fmt.Sprintf("%g", req.TyLeHoaHong), req.GhiChu,
	)
	if err != nil {
		return "", err
	}

	// 4. Saga/Retry: Cập nhật giảm trừ số lượng vào lô hàng
	platformMsp := "PlatformOrgMSP" // Việc update inventory cross-channel thường được do Platform điều phối
	infrastructure.EnqueueRetryTask(
		fmt.Sprintf("DeductInventory-%s", req.MaGiaoDich),
		5, // Max Retries
		func() error {
			_, e := uc.lohangRepo.Submit(platformMsp, "CapNhatInventory", req.MaLo, fmt.Sprintf("%g", -req.SoLuong))
			return e
		},
	)

	return req.MaGiaoDich, nil
}

func (uc *giaodichUsecase) DuyetGiaoDich(mspID string, id string) error {
	gdBytes, err := uc.repo.Evaluate(mspID, "DocGiaoDich", id)
	if err != nil {
		return fmt.Errorf("không thể truy xuất thông tin giao dịch: %v", err)
	}

	var gd model.GiaoDichDTO
	if err := json.Unmarshal(gdBytes, &gd); err != nil {
		return fmt.Errorf("lỗi parse dữ liệu giao dịch: %v", err)
	}

	lotHangBytes, err := uc.lohangRepo.Evaluate(mspID, "DocLotHang", gd.MaLo)
	if err != nil {
		return fmt.Errorf("không thể truy xuất lô hàng từ channel: %v", err)
	}

	var loHang model.LotHangDTO
	if err := json.Unmarshal(lotHangBytes, &loHang); err != nil {
		return fmt.Errorf("lỗi parse dữ liệu lô hàng: %v", err)
	}

	if loHang.TrangThai == "dinh_chi" {
		return fmt.Errorf("Action blocked: Giao dịch bị khoá do lô hàng đang trong trạng thái đình chỉ thu hồi")
	}

	_, err = uc.repo.Submit(mspID, "DuyetGiaoDich", id)
	return err
}

func (uc *giaodichUsecase) CapNhatTrangThai(mspID string, id string, trangThai string) error {
	var gd model.GiaoDichDTO

	gdBytes, err := uc.repo.Evaluate(mspID, "DocGiaoDich", id)
	if err != nil {
		return fmt.Errorf("không thể truy xuất thông tin giao dịch để tiến hành cập nhật: %v", err)
	}
	if err := json.Unmarshal(gdBytes, &gd); err != nil {
		return fmt.Errorf("lỗi parse dữ liệu giao dịch nội bộ: %v", err)
	}

	// Cross-check lock: Nếu lô hàng bị đình chỉ thì chỉ cho phép Platform Hủy bỏ giao dịch
	lotHangBytes, err := uc.lohangRepo.Evaluate(mspID, "DocLotHang", gd.MaLo)
	if err == nil {
		var loHang model.LotHangDTO
		if err := json.Unmarshal(lotHangBytes, &loHang); err == nil {
			if loHang.TrangThai == "dinh_chi" && trangThai != "huy_bo" {
				return fmt.Errorf("Action blocked: Giao dịch bị khoá do lô hàng đang trong trạng thái đình chỉ thu hồi")
			}
		}
	}

	_, err = uc.repo.Submit(mspID, "CapNhatTrangThai", id, trangThai, "")
	if err != nil {
		return err
	}

	// Rollback Inventory (Rule R6 SKILL.md)
	if trangThai == "huy_bo" {
		platformMsp := "PlatformOrgMSP"
		infrastructure.EnqueueRetryTask(
			fmt.Sprintf("RollbackInventory-%s", id),
			5,
			func() error {
				_, e := uc.lohangRepo.Submit(platformMsp, "CapNhatInventory", gd.MaLo, fmt.Sprintf("%g", gd.SoLuong))
				return e
			},
		)
	}

	return nil
}

func (uc *giaodichUsecase) DocGiaoDich(mspID string, id string) ([]byte, error) {
	return uc.repo.Evaluate(mspID, "DocGiaoDich", id)
}

func (uc *giaodichUsecase) LichSuGiaoDich(mspID string, id string) ([]byte, error) {
	return uc.repo.Evaluate(mspID, "LichSuGiaoDich", id)
}

func (uc *giaodichUsecase) DocCongNoNPP(mspID string, maNPP string) ([]byte, error) {
	return uc.repo.Evaluate(mspID, "DocCongNoNPP", maNPP)
}

func (uc *giaodichUsecase) TinhHoaHongNPP(mspID string, maNPP string) ([]byte, error) {
	// Ghi đè mspID thành PlatformOrgMSP vì chaincode TinhHoaHongNPP chỉ cấp quyền cho Platform
	return uc.repo.Evaluate("PlatformOrgMSP", "TinhHoaHongNPP", maNPP)
}

func (uc *giaodichUsecase) DocGiaoDichTheoHTX(mspID string, maHTX string) ([]byte, error) {
	return uc.repo.Evaluate(mspID, "DocGiaoDichTheoHTX", maHTX)
}

func (uc *giaodichUsecase) DocGiaoDichTheoNPP(mspID string, maNPP string) ([]byte, error) {
	return uc.repo.Evaluate(mspID, "DocGiaoDichTheoNPP", maNPP)
}
