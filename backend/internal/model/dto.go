package model

type TaoLotHangReq struct {
	MaLo        string  `json:"ma_lo" binding:"required"`
	MaLoMe      string  `json:"ma_lo_me"`
	MaHTX       string  `json:"ma_htx" binding:"required"`
	TenSanPham  string  `json:"ten_san_pham" binding:"required"`
	LoaiSanPham string  `json:"loai_san_pham" binding:"required"`
	SoLuong     float64 `json:"so_luong" binding:"required,gt=0"`
	DonViTinh   string  `json:"don_vi_tinh" binding:"required"`
	VuMua       string  `json:"vu_mua"`
	DiaDiem     string  `json:"dia_diem"`
}

type CapNhatTrangThaiLoReq struct {
	TrangThai string `json:"trang_thai" binding:"required"`
}

type TachLoReq struct {
	MaLoMoi     string  `json:"ma_lo_moi" binding:"required"`
	SoLuongTach float64 `json:"so_luong_tach" binding:"required,gt=0"`
}

type ThemChungNhanReq struct {
	LoaiChungNhan string `json:"loai_chung_nhan" binding:"required"`
	MaChungNhan   string `json:"ma_chung_nhan" binding:"required"`
	CoQuanCap     string `json:"co_quan_cap" binding:"required"`
	NgayCap       string `json:"ngay_cap" binding:"required"`
	NgayHetHan    string `json:"ngay_het_han" binding:"required"`
	GhiChu        string `json:"ghi_chu"`
}

type GhiNhatKyReq struct {
	MaNhatKy      string `json:"ma_nhat_ky" binding:"required"`
	MaLo          string `json:"ma_lo" binding:"required"`
	MaHTX         string `json:"ma_htx" binding:"required"`
	LoaiHoatDong  string `json:"loai_hoat_dong" binding:"required"`
	ChiTiet       string `json:"chi_tiet" binding:"required"`
	ViTri         string `json:"vi_tri" binding:"required"`
	NguoiThucHien string `json:"nguoi_thuc_hien" binding:"required"`
	NgayGhi       string `json:"ngay_ghi" binding:"required"`
}

type XacNhanNhatKyReq struct {
	MinhChungHash string `json:"minh_chung_hash" binding:"required"`
}

type DuyetNhatKyReq struct {
	QuyetDinh  string `json:"quyet_dinh" binding:"required"` // duyet | tu_choi
	LyDoTuChoi string `json:"ly_do_tu_choi"`                 // bat buoc khi huy
	NguoiDuyet string `json:"nguoi_duyet" binding:"required"`
}

type TaoGiaoDichReq struct {
	MaGiaoDich  string  `json:"ma_giao_dich" binding:"required"`
	MaLo        string  `json:"ma_lo" binding:"required"`
	MaHTX       string  `json:"ma_htx" binding:"required"`
	MaNPP       string  `json:"ma_npp" binding:"required"`
	SanPham     string  `json:"san_pham" binding:"required"`
	SoLuong     float64 `json:"so_luong" binding:"required,gt=0"`
	DonViTinh   string  `json:"don_vi_tinh" binding:"required"`
	DonGia      float64 `json:"don_gia" binding:"required,gt=0"`
	TyLeHoaHong float64 `json:"ty_le_hoa_hong" binding:"required,gte=0"`
	GhiChu      string  `json:"ghi_chu"`
}

type DuyetGiaoDichReq struct {
	TrangThai string `json:"trang_thai" binding:"required"` // duyet / tu_choi
}

type CapNhatTrangThaiGiaoDichReq struct {
	TrangThai string `json:"trang_thai" binding:"required"`
}

// LotHangDTO dùng để parse dữ liệu từ chaincode trả về ở các layer Usecase
type LotHangDTO struct {
	MaLo          string  `json:"ma_lo"`
	MaLoMe        string  `json:"ma_lo_me"`
	MaHTX         string  `json:"ma_htx"`
	TrangThai     string  `json:"trang_thai"`
	SoLuong       float64 `json:"so_luong"`
	SoLuongConLai float64 `json:"so_luong_con_lai"`
}

// GiaoDichDTO dùng để parse dữ liệu giao dịch từ chaincode
type GiaoDichDTO struct {
	MaLo     string  `json:"ma_lo"`
	SoLuong  float64 `json:"so_luong"`
}
