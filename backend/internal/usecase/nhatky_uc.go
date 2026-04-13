package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/trustfab/gapchain/backend/internal/model"
	repo "github.com/trustfab/gapchain/backend/internal/repository/fabric"
)

type NhatkyUsecase interface {
	GhiNhatKy(mspID string, req *model.GhiNhatKyReq) (string, error)
	DuyetNhatKy(mspID string, id string, req *model.DuyetNhatKyReq) error
	DocNhatKyTheoLo(mspID string, maLo string) ([]byte, error)
	DocNhatKyTheoHTX(mspID string, maHTX string) ([]byte, error)
	XacNhanNhatKy(mspID string, id string, req *model.XacNhanNhatKyReq) error
	LichSuNhatKy(mspID string, id string) ([]byte, error)
	ThongKeNhatKy(mspID string, maHTX string) ([]byte, error)
}

type nhatkyUsecase struct {
	repo       repo.NhatkyRepo
	lohangRepo repo.LohangRepo
}

func NewNhatkyUsecase(r repo.NhatkyRepo, lr repo.LohangRepo) NhatkyUsecase {
	return &nhatkyUsecase{repo: r, lohangRepo: lr}
}

func (uc *nhatkyUsecase) GhiNhatKy(mspID string, req *model.GhiNhatKyReq) (string, error) {
	lotHangBytes, err := uc.lohangRepo.Evaluate(mspID, "DocLotHang", req.MaLo)
	if err != nil {
		return "", fmt.Errorf("không tìm thấy dữ liệu lô hàng: %v", err)
	}

	var loHang model.LotHangDTO
	if err := json.Unmarshal(lotHangBytes, &loHang); err != nil {
		return "", fmt.Errorf("lỗi parse lô hàng: %v", err)
	}

	hopLe := false
	tt := loHang.TrangThai
	switch req.LoaiHoatDong {
	case "gieo_hat", "bon_phan", "tuoi_nuoc", "phun_thuoc":
		if tt == "dang_trong" {
			hopLe = true
		}
	case "thu_hoach":
		if tt == "dang_trong" || tt == "da_thu_hoach" {
			hopLe = true
		}
	case "kiem_tra":
		if tt == "dang_trong" || tt == "da_thu_hoach" {
			hopLe = true
		}
	case "dong_goi":
		if tt == "da_thu_hoach" {
			hopLe = true
		}
	case "van_chuyen":
		if tt == "da_thu_hoach" || tt == "san_sang_ban" {
			hopLe = true
		}
	case "khac":
		hopLe = true
	}

	if !hopLe {
		return "", fmt.Errorf("hành động '%s' không hợp lệ khi lô hàng đang ở trạng thái '%s'", req.LoaiHoatDong, tt)
	}

	_, err = uc.repo.Submit(
		mspID, "GhiNhatKy",
		req.MaNhatKy, req.MaLo, req.MaHTX, req.LoaiHoatDong, req.ChiTiet, req.ViTri, req.NguoiThucHien, req.NgayGhi,
	)
	if err != nil {
		return "", err
	}
	return req.MaNhatKy, nil
}

func (uc *nhatkyUsecase) DuyetNhatKy(mspID string, id string, req *model.DuyetNhatKyReq) error {
	_, err := uc.repo.Submit(mspID, "DuyetNhatKy", id, req.NguoiDuyet, req.QuyetDinh, req.LyDoTuChoi)
	return err
}

func (uc *nhatkyUsecase) DocNhatKyTheoLo(mspID string, maLo string) ([]byte, error) {
	return uc.repo.Evaluate(mspID, "DocNhatKyTheoLo", maLo)
}

func (uc *nhatkyUsecase) DocNhatKyTheoHTX(mspID string, maHTX string) ([]byte, error) {
	return uc.repo.Evaluate(mspID, "DocNhatKyTheoHTX", maHTX)
}

func (uc *nhatkyUsecase) XacNhanNhatKy(mspID string, id string, req *model.XacNhanNhatKyReq) error {
	_, err := uc.repo.Submit(mspID, "XacNhanNhatKy", id, req.MinhChungHash)
	return err
}

func (uc *nhatkyUsecase) LichSuNhatKy(mspID string, id string) ([]byte, error) {
	return uc.repo.Evaluate(mspID, "LichSuNhatKy", id)
}

func (uc *nhatkyUsecase) ThongKeNhatKy(mspID string, maHTX string) ([]byte, error) {
	if maHTX == "" {
		maHTX = "ALL"
	}
	return uc.repo.Evaluate(mspID, "ThongKeNhatKy", maHTX)
}
