package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/trustfab/gapchain/backend/internal/model"
	repo "github.com/trustfab/gapchain/backend/internal/repository/fabric"
)

type LohangUsecase interface {
	TaoLotHang(mspID string, req *model.TaoLotHangReq) (string, error)
	CapNhatTrangThai(mspID string, maLo string, trangThai string) error
	ThemChungNhan(mspID string, maLo string, req *model.ThemChungNhanReq) error
	DocLotHang(mspID string, maLo string) ([]byte, error)
	LichSuLotHang(mspID string, maLo string) ([]byte, error)
	DocLotHangTheoHTX(mspID string, maHTX string) ([]byte, error)
	LayThongTinTraCuu(maLo string) ([]byte, error)
	CapNhatInventory(mspID string, maLo string, thayDoi float64) error
	TachLo(mspID string, maLoMe string, req *model.TachLoReq) error
}

type lohangUsecase struct {
	lohangRepo repo.LohangRepo
	nhatkyRepo repo.NhatkyRepo
}

func NewLohangUsecase(lr repo.LohangRepo, nr repo.NhatkyRepo) LohangUsecase {
	return &lohangUsecase{
		lohangRepo: lr,
		nhatkyRepo: nr,
	}
}

func (uc *lohangUsecase) TaoLotHang(mspID string, req *model.TaoLotHangReq) (string, error) {
	_, err := uc.lohangRepo.Submit(
		mspID, "TaoLotHang",
		req.MaLo, req.MaLoMe, req.MaHTX, req.TenSanPham, req.LoaiSanPham,
		fmt.Sprintf("%g", req.SoLuong), req.DonViTinh, req.VuMua, req.DiaDiem,
	)
	if err != nil {
		return "", err
	}
	return req.MaLo, nil
}

func (uc *lohangUsecase) TachLo(mspID string, maLoMe string, req *model.TachLoReq) error {
	_, err := uc.lohangRepo.Submit(
		mspID, "TachLo",
		maLoMe, req.MaLoMoi, fmt.Sprintf("%g", req.SoLuongTach),
	)
	return err
}

func (uc *lohangUsecase) CapNhatTrangThai(mspID string, maLo string, trangThai string) error {
	_, err := uc.lohangRepo.Submit(mspID, "CapNhatTrangThaiLo", maLo, trangThai)
	return err
}

func (uc *lohangUsecase) ThemChungNhan(mspID string, maLo string, req *model.ThemChungNhanReq) error {
	_, err := uc.lohangRepo.Submit(
		mspID, "ThemChungNhan",
		maLo, req.LoaiChungNhan, req.MaChungNhan, req.CoQuanCap, req.NgayCap, req.NgayHetHan, req.GhiChu,
	)
	return err
}

func (uc *lohangUsecase) DocLotHang(mspID string, maLo string) ([]byte, error) {
	return uc.lohangRepo.Evaluate(mspID, "DocLotHang", maLo)
}

func (uc *lohangUsecase) LichSuLotHang(mspID string, maLo string) ([]byte, error) {
	return uc.lohangRepo.Evaluate(mspID, "LichSuLotHang", maLo)
}

func (uc *lohangUsecase) DocLotHangTheoHTX(mspID string, maHTX string) ([]byte, error) {
	return uc.lohangRepo.Evaluate(mspID, "DocLotHangTheoHTX", maHTX)
}

func (uc *lohangUsecase) LayThongTinTraCuu(maLo string) ([]byte, error) {
	fallbackMsp := "PlatformOrgMSP"

	lotHangBytes, err := uc.lohangRepo.Evaluate(fallbackMsp, "LayThongTinTraCuu", maLo)
	if err != nil {
		return nil, err
	}

	var loHang struct {
		MaLoMe string `json:"ma_lo_me"`
	}
	json.Unmarshal(lotHangBytes, &loHang)

	nhatKyBytes, _ := uc.nhatkyRepo.Evaluate(fallbackMsp, "DocNhatKyTheoLo", maLo)
	var nkJSON json.RawMessage
	if len(nhatKyBytes) > 0 {
		nkJSON = json.RawMessage(nhatKyBytes)
	} else {
		nkJSON = json.RawMessage(`[]`)
	}

	combined := map[string]interface{}{
		"lo_hang": json.RawMessage(lotHangBytes),
		"nhat_ky": nkJSON,
	}

	// Traceability logic: If part of a split batch, dynamically query the parent's info
	if loHang.MaLoMe != "" {
		parentLotHangBytes, _ := uc.lohangRepo.Evaluate(fallbackMsp, "LayThongTinTraCuu", loHang.MaLoMe)
		if len(parentLotHangBytes) > 0 {
			combined["lo_hang_me"] = json.RawMessage(parentLotHangBytes)
		}

		parentNhatKyBytes, _ := uc.nhatkyRepo.Evaluate(fallbackMsp, "DocNhatKyTheoLo", loHang.MaLoMe)
		var pnkJSON json.RawMessage
		if len(parentNhatKyBytes) > 0 {
			pnkJSON = json.RawMessage(parentNhatKyBytes)
		} else {
			pnkJSON = json.RawMessage(`[]`)
		}
		combined["nhat_ky_me"] = pnkJSON
	}

	return json.Marshal(combined)
}

func (uc *lohangUsecase) CapNhatInventory(mspID string, maLo string, thayDoi float64) error {
	_, err := uc.lohangRepo.Submit(mspID, "CapNhatInventory", maLo, fmt.Sprintf("%g", thayDoi))
	return err
}
