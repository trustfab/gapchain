package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// LotHangContract cung cấp các chức năng quản lý lô hàng nông sản
type LotHangContract struct {
	contractapi.Contract
}

// LotHang đại diện cho một lô hàng nông sản
type LotHang struct {
	MaLo          string      `json:"ma_lo"`          // VD: LH-HTX001-2025-001
	MaLoMe        string      `json:"ma_lo_me"`       // Mã lô mẹ (nếu có, dùng để Traceability)
	MaHTX         string      `json:"ma_htx"`
	TenSanPham    string      `json:"ten_san_pham"`   // VD: Gạo ST25, Cà phê Lâm Đồng
	LoaiSanPham   string      `json:"loai_san_pham"`  // lua | ca_phe | rau | qua | khac
	SoLuong       float64     `json:"so_luong"`
	SoLuongConLai float64     `json:"so_luong_con_lai"`
	DonViTinh     string      `json:"don_vi_tinh"`    // kg | tan | thung | qua
	VuMua         string      `json:"vu_mua"`          // VD: Dong Xuan 2025
	DiaDiem       string      `json:"dia_diem"`        // Tên vùng trồng / địa chỉ
	ChungNhan     []ChungNhan `json:"chung_nhan"`      // Danh sách chứng nhận đính kèm
	TrangThai     string      `json:"trang_thai"`      // dang_trong | da_thu_hoach | cho_chung_nhan | san_sang_ban | het_hang | dinh_chi
	NgayTao       string      `json:"ngay_tao"`
	NgayCapNhat   string      `json:"ngay_cap_nhat"`
}

// ChungNhan đại diện cho một chứng nhận được cấp bởi cơ quan có thẩm quyền
type ChungNhan struct {
	LoaiChungNhan string `json:"loai_chung_nhan"` // VietGAP | GlobalGAP | Organic | TCVN
	MaChungNhan   string `json:"ma_chung_nhan"`
	CoQuanCap     string `json:"co_quan_cap"`     // VD: Chi Cuc BVTV Can Tho
	NgayCap       string `json:"ngay_cap"`
	NgayHetHan    string `json:"ngay_het_han"`
	GhiChu        string `json:"ghi_chu"`
}

// LichSuRecord đại diện cho 1 bản ghi trong lịch sử thay đổi
type LichSuRecord struct {
	TxID      string `json:"tx_id"`
	Timestamp string `json:"timestamp"`
	IsDeleted bool   `json:"is_deleted"`
	Data      string `json:"data"`
}

// Vòng đời trạng thái lô hàng:
// dang_trong → da_thu_hoach → cho_chung_nhan → san_sang_ban → het_hang
// Bất kỳ trạng thái nào (trừ het_hang) → dinh_chi (Platform hoặc BVTV)
// dinh_chi → trạng thái trước đó (chỉ Platform)
var chuyenTrangThaiHopLe = map[string][]string{
	"dang_trong":     {"da_thu_hoach", "dinh_chi"},
	"da_thu_hoach":   {"cho_chung_nhan", "dinh_chi"},
	"cho_chung_nhan": {"san_sang_ban", "dinh_chi"},
	"san_sang_ban":   {"het_hang", "dinh_chi"},
	"het_hang":       {},                                                                    // trạng thái kết thúc
	"dinh_chi":       {"dang_trong", "da_thu_hoach", "cho_chung_nhan", "san_sang_ban"}, // phục hồi
}

// trangThaiHopLe dùng cho validate query (DocLotHangTheoTrangThai)
var trangThaiHopLe = map[string]bool{
	"dang_trong":     true,
	"da_thu_hoach":   true,
	"cho_chung_nhan": true,
	"san_sang_ban":   true,
	"het_hang":       true,
	"dinh_chi":       true,
}

// loaiChungNhanHopLe là các loại chứng nhận hợp lệ
var loaiChungNhanHopLe = map[string]bool{
	"VietGAP":  true,
	"GlobalGAP": true,
	"Organic":  true,
	"TCVN":     true,
	"Khac":     true,
}

// kiemTraMSP kiểm tra quyền truy cập dựa trên MSP ID của caller
func kiemTraMSP(ctx contractapi.TransactionContextInterface, mspChoPhep ...string) error {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("khong the lay MSP ID: %v", err)
	}
	for _, m := range mspChoPhep {
		if mspID == m {
			return nil
		}
	}
	return fmt.Errorf("MSP %s khong co quyen thuc hien thao tac nay", mspID)
}

// txNow lấy thời điểm transaction (deterministic — dùng thay time.Now() trong chaincode)
func txNow(ctx contractapi.TransactionContextInterface) (string, error) {
	ts, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("khong the lay tx timestamp: %v", err)
	}
	return time.Unix(ts.Seconds, int64(ts.Nanos)).UTC().Format(time.RFC3339), nil
}

// TaoLotHang tạo một lô hàng mới trên blockchain
// MSP được phép: HTXNongSanOrgMSP, PlatformOrgMSP
func (c *LotHangContract) TaoLotHang(
	ctx contractapi.TransactionContextInterface,
	maLo string,
	maLoMe string,
	maHTX string,
	tenSanPham string,
	loaiSanPham string,
	soLuong float64,
	donViTinh string,
	vuMua string,
	diaDiem string,
) error {
	if err := kiemTraMSP(ctx, "HTXNongSanOrgMSP", "PlatformOrgMSP"); err != nil {
		return err
	}

	if maLo == "" || maHTX == "" || tenSanPham == "" {
		return fmt.Errorf("ma_lo, ma_htx va ten_san_pham khong duoc de trong")
	}
	if soLuong <= 0 {
		return fmt.Errorf("so_luong phai lon hon 0")
	}

	// Kiểm tra lô hàng chưa tồn tại
	existing, err := ctx.GetStub().GetState(maLo)
	if err != nil {
		return fmt.Errorf("loi khi kiem tra trang thai: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("lo hang %s da ton tai", maLo)
	}

	now, err := txNow(ctx)
	if err != nil {
		return err
	}
	lotHang := LotHang{
		MaLo:          maLo,
		MaLoMe:        maLoMe,
		MaHTX:         maHTX,
		TenSanPham:    tenSanPham,
		LoaiSanPham:   loaiSanPham,
		SoLuong:       soLuong,
		SoLuongConLai: soLuong, // Khởi tạo inventory bằng soLuong
		DonViTinh:     donViTinh,
		VuMua:         vuMua,
		DiaDiem:       diaDiem,
		ChungNhan:     []ChungNhan{},
		TrangThai:     "dang_trong",
		NgayTao:       now,
		NgayCapNhat:   now,
	}

	lotHangJSON, err := json.Marshal(lotHang)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maLo, lotHangJSON)
}

// TachLo tách một lô hàng (partial batch) từ lô gốc
// MSP được phép: HTXNongSanOrgMSP, PlatformOrgMSP
func (c *LotHangContract) TachLo(
	ctx contractapi.TransactionContextInterface,
	maLoMe string,
	maLoMoi string,
	soLuongTach float64,
) error {
	if err := kiemTraMSP(ctx, "HTXNongSanOrgMSP", "PlatformOrgMSP"); err != nil {
		return err
	}
	if maLoMe == "" || maLoMoi == "" {
		return fmt.Errorf("ma_lo_me va ma_lo_moi khong duoc de trong")
	}
	if soLuongTach <= 0 {
		return fmt.Errorf("so_luong_tach phai lon hon 0")
	}

	parent, err := c.layLotHang(ctx, maLoMe)
	if err != nil {
		return fmt.Errorf("loi lay lo me: %v", err)
	}

	if parent.TrangThai != "san_sang_ban" {
		return fmt.Errorf("chi duoc tach lo khi tren trang thai san_sang_ban (hien tai: %s)", parent.TrangThai)
	}

	if soLuongTach > parent.SoLuongConLai {
		return fmt.Errorf("so luong tach (%g) vuot qua so luong con lai cua lo me (%g)", soLuongTach, parent.SoLuongConLai)
	}

	// Kiểm tra lô mới chưa tồn tại
	existing, err := ctx.GetStub().GetState(maLoMoi)
	if err != nil {
		return fmt.Errorf("loi khi kiem tra trang thai lo moi: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("lo con %s da ton tai", maLoMoi)
	}

	now, err := txNow(ctx)
	if err != nil {
		return err
	}

	// Cập nhật lô mẹ
	parent.SoLuongConLai -= soLuongTach
	if parent.SoLuongConLai == 0 {
		parent.TrangThai = "het_hang"
	}
	parent.NgayCapNhat = now

	parentJSON, err := json.Marshal(parent)
	if err != nil {
		return fmt.Errorf("loi ma hoa lo me: %v", err)
	}

	// Tạo lô con bằng cách copy thuộc tính của lô mẹ
	loCon := LotHang{
		MaLo:          maLoMoi,
		MaLoMe:        maLoMe,
		MaHTX:         parent.MaHTX,
		TenSanPham:    parent.TenSanPham,
		LoaiSanPham:   parent.LoaiSanPham,
		SoLuong:       soLuongTach,
		SoLuongConLai: soLuongTach, // Inventory khởi tạo bằng soLuongTach
		DonViTinh:     parent.DonViTinh,
		VuMua:         parent.VuMua,
		DiaDiem:       parent.DiaDiem,
		// Không copy list chứng nhận theo yêu cầu người dùng, 
		// Traceability Lookup ở Web/API sẽ đệ quy truy vết Lô Mẹ.
		ChungNhan:     []ChungNhan{},
		TrangThai:     "san_sang_ban",
		NgayTao:       now,
		NgayCapNhat:   now,
	}

	loConJSON, err := json.Marshal(loCon)
	if err != nil {
		return fmt.Errorf("loi ma hoa lo con: %v", err)
	}

	// Lưu State cho Lô Mẹ và Lô Con đồng thời trong cùng 1 transaction atomic
	if err := ctx.GetStub().PutState(maLoMe, parentJSON); err != nil {
		return err
	}
	return ctx.GetStub().PutState(maLoMoi, loConJSON)
}

// CapNhatTrangThaiLo cập nhật trạng thái lô hàng theo vòng đời
// Quy tắc MSP theo trạng thái đích:
//
//	da_thu_hoach, cho_chung_nhan, san_sang_ban, het_hang → HTXNongSanOrgMSP, PlatformOrgMSP
//	dinh_chi → PlatformOrgMSP, ChiCucBVTVOrgMSP
//	phục hồi từ dinh_chi → PlatformOrgMSP (chỉ Platform mới phục hồi)
func (c *LotHangContract) CapNhatTrangThaiLo(
	ctx contractapi.TransactionContextInterface,
	maLo string,
	trangThaiMoi string,
) error {
	lotHang, err := c.layLotHang(ctx, maLo)
	if err != nil {
		return err
	}

	// Kiểm tra quyền theo trạng thái đích
	switch {
	case trangThaiMoi == "dinh_chi":
		// Đình chỉ: Platform hoặc BVTV
		if err := kiemTraMSP(ctx, "PlatformOrgMSP", "ChiCucBVTVOrgMSP"); err != nil {
			return err
		}
	case lotHang.TrangThai == "dinh_chi":
		// Phục hồi từ đình chỉ: chỉ Platform
		if err := kiemTraMSP(ctx, "PlatformOrgMSP"); err != nil {
			return err
		}
	case trangThaiMoi == "san_sang_ban":
		// Theo BA: HTX, Platform (hoặc BVTV khi cấp chứng nhận) đều được phép chuyển tới san_sang_ban.
		// Điệu kiện tiên quyết ở Backend (chaincode enforce để strict): Phải có >= 1 chứng nhận
		if err := kiemTraMSP(ctx, "HTXNongSanOrgMSP", "PlatformOrgMSP", "ChiCucBVTVOrgMSP"); err != nil {
			return err
		}
		if len(lotHang.ChungNhan) == 0 {
			return fmt.Errorf("dieu kien tien quyet: phai co it nhat 1 chung nhan de chuyen sang san_sang_ban")
		}
	default:
		// Chuyển trạng thái bình thường: HTX hoặc Platform
		if err := kiemTraMSP(ctx, "HTXNongSanOrgMSP", "PlatformOrgMSP"); err != nil {
			return err
		}
	}

	// Kiểm tra vòng đời hợp lệ
	choPhep := chuyenTrangThaiHopLe[lotHang.TrangThai]
	hopLe := false
	for _, tt := range choPhep {
		if tt == trangThaiMoi {
			hopLe = true
			break
		}
	}
	if !hopLe {
		return fmt.Errorf("khong the chuyen trang thai tu '%s' sang '%s'", lotHang.TrangThai, trangThaiMoi)
	}

	now, err := txNow(ctx)
	if err != nil {
		return err
	}
	lotHang.TrangThai = trangThaiMoi
	lotHang.NgayCapNhat = now

	lotHangJSON, err := json.Marshal(lotHang)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maLo, lotHangJSON)
}

// ThemChungNhan thêm một chứng nhận vào lô hàng
// MSP được phép: ChiCucBVTVOrgMSP, PlatformOrgMSP
func (c *LotHangContract) ThemChungNhan(
	ctx contractapi.TransactionContextInterface,
	maLo string,
	loaiChungNhan string,
	maChungNhan string,
	coQuanCap string,
	ngayCap string,
	ngayHetHan string,
	ghiChu string,
) error {
	// Chỉ Chi Cục BVTV và Platform mới được cấp chứng nhận
	if err := kiemTraMSP(ctx, "ChiCucBVTVOrgMSP", "PlatformOrgMSP"); err != nil {
		return err
	}

	if maChungNhan == "" || coQuanCap == "" || ngayCap == "" {
		return fmt.Errorf("ma_chung_nhan, co_quan_cap va ngay_cap khong duoc de trong")
	}

	if !loaiChungNhanHopLe[loaiChungNhan] {
		return fmt.Errorf("loai_chung_nhan khong hop le: %s. Cac gia tri hop le: VietGAP, GlobalGAP, Organic, TCVN, Khac", loaiChungNhan)
	}

	lotHang, err := c.layLotHang(ctx, maLo)
	if err != nil {
		return err
	}

	// Kiểm tra chứng nhận cùng loại đã tồn tại chưa
	for _, cn := range lotHang.ChungNhan {
		if cn.MaChungNhan == maChungNhan {
			return fmt.Errorf("chung nhan %s da ton tai trong lo hang nay", maChungNhan)
		}
	}

	chungNhan := ChungNhan{
		LoaiChungNhan: loaiChungNhan,
		MaChungNhan:   maChungNhan,
		CoQuanCap:     coQuanCap,
		NgayCap:       ngayCap,
		NgayHetHan:    ngayHetHan,
		GhiChu:        ghiChu,
	}

	now, err := txNow(ctx)
	if err != nil {
		return err
	}
	lotHang.ChungNhan = append(lotHang.ChungNhan, chungNhan)
	lotHang.NgayCapNhat = now

	lotHangJSON, err := json.Marshal(lotHang)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maLo, lotHangJSON)
}

// DocLotHang đọc thông tin của một lô hàng theo mã lô
// MSP được phép: Tất cả
func (c *LotHangContract) DocLotHang(
	ctx contractapi.TransactionContextInterface,
	maLo string,
) (*LotHang, error) {
	return c.layLotHang(ctx, maLo)
}

// DocLotHangTheoHTX trả về danh sách lô hàng của một HTX (CouchDB rich query)
// MSP được phép: Tất cả
func (c *LotHangContract) DocLotHangTheoHTX(
	ctx contractapi.TransactionContextInterface,
	maHTX string,
) ([]*LotHang, error) {
	if maHTX == "" {
		return nil, fmt.Errorf("ma_htx khong duoc de trong")
	}
	queryString := fmt.Sprintf(`{"selector":{"ma_htx":"%s","ngay_tao":{"$gt":null}},"sort":[{"ngay_tao":"desc"}]}`, maHTX)
	return c.queryLotHang(ctx, queryString)
}

// DocLotHangTheoTrangThai trả về danh sách lô hàng theo trạng thái (CouchDB rich query)
// MSP được phép: Tất cả
func (c *LotHangContract) DocLotHangTheoTrangThai(
	ctx contractapi.TransactionContextInterface,
	trangThai string,
) ([]*LotHang, error) {
	if !trangThaiHopLe[trangThai] {
		return nil, fmt.Errorf("trang_thai khong hop le: %s", trangThai)
	}
	queryString := fmt.Sprintf(`{"selector":{"trang_thai":"%s","ngay_tao":{"$gt":null}},"sort":[{"ngay_tao":"desc"}]}`, trangThai)
	return c.queryLotHang(ctx, queryString)
}

// DocLotHangTheoHTXVaTrangThai trả về lô hàng theo HTX và trạng thái (CouchDB rich query)
// MSP được phép: Tất cả
func (c *LotHangContract) DocLotHangTheoHTXVaTrangThai(
	ctx contractapi.TransactionContextInterface,
	maHTX string,
	trangThai string,
) ([]*LotHang, error) {
	if maHTX == "" {
		return nil, fmt.Errorf("ma_htx khong duoc de trong")
	}
	if !trangThaiHopLe[trangThai] {
		return nil, fmt.Errorf("trang_thai khong hop le: %s", trangThai)
	}
	queryString := fmt.Sprintf(
		`{"selector":{"ma_htx":"%s","trang_thai":"%s","ngay_tao":{"$gt":null}},"sort":[{"ngay_tao":"desc"}]}`,
		maHTX, trangThai,
	)
	return c.queryLotHang(ctx, queryString)
}

// LichSuLotHang trả về toàn bộ lịch sử thay đổi của một lô hàng (audit trail)
// MSP được phép: Tất cả
func (c *LotHangContract) LichSuLotHang(
	ctx contractapi.TransactionContextInterface,
	maLo string,
) ([]*LichSuRecord, error) {
	// Kiểm tra lô hàng tồn tại
	if _, err := c.layLotHang(ctx, maLo); err != nil {
		return nil, err
	}

	iterator, err := ctx.GetStub().GetHistoryForKey(maLo)
	if err != nil {
		return nil, fmt.Errorf("loi khi lay lich su: %v", err)
	}
	defer iterator.Close()

	var lichSu []*LichSuRecord
	for iterator.HasNext() {
		record, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("loi khi doc lich su: %v", err)
		}
		lichSu = append(lichSu, &LichSuRecord{
			TxID:      record.TxId,
			Timestamp: record.Timestamp.String(),
			IsDeleted: record.IsDelete,
			Data:      string(record.Value),
		})
	}

	return lichSu, nil
}

// LayThongTinTraCuu trả về thông tin đầy đủ của lô hàng để hiển thị trên QR truy xuất
// Trả về LotHang bao gồm toàn bộ chứng nhận. Backend API sẽ tổng hợp thêm nhật ký canh tác.
// MSP được phép: Tất cả (bao gồm public — consumer quét QR)
func (c *LotHangContract) LayThongTinTraCuu(
	ctx contractapi.TransactionContextInterface,
	maLo string,
) (*LotHang, error) {
	return c.layLotHang(ctx, maLo)
}

// CapNhatSoLuong cập nhật số lượng lô hàng (VD: sau khi thu hoạch thực tế)
// MSP được phép: HTXNongSanOrgMSP, PlatformOrgMSP
func (c *LotHangContract) CapNhatSoLuong(
	ctx contractapi.TransactionContextInterface,
	maLo string,
	soLuongMoi float64,
	donViTinhMoi string,
) error {
	if err := kiemTraMSP(ctx, "HTXNongSanOrgMSP", "PlatformOrgMSP"); err != nil {
		return err
	}
	if soLuongMoi <= 0 {
		return fmt.Errorf("so_luong phai lon hon 0")
	}

	lotHang, err := c.layLotHang(ctx, maLo)
	if err != nil {
		return err
	}

	lotHang.SoLuong = soLuongMoi
	now, err := txNow(ctx)
	if err != nil {
		return err
	}
	if donViTinhMoi != "" {
		lotHang.DonViTinh = donViTinhMoi
	}
	lotHang.NgayCapNhat = now

	lotHangJSON, err := json.Marshal(lotHang)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maLo, lotHangJSON)
}

// CapNhatInventory thay đổi số lượng còn lại của lô (khi có giao dịch hoặc hủy GD)
// MSP được phép: PlatformOrgMSP (vì giao dịch được Platform quản lý/trigger)
func (c *LotHangContract) CapNhatInventory(
	ctx contractapi.TransactionContextInterface,
	maLo string,
	thayDoi float64,
) error {
	if err := kiemTraMSP(ctx, "PlatformOrgMSP"); err != nil {
		return err
	}

	lotHang, err := c.layLotHang(ctx, maLo)
	if err != nil {
		return err
	}

	lotHang.SoLuongConLai += thayDoi
	if lotHang.SoLuongConLai < 0 {
		return fmt.Errorf("tong so luong con lai khong the nho hon 0 (con lai truoc thuc thi: %g)", lotHang.SoLuongConLai-thayDoi)
	}

	// Tự động chuyển sang trạng thái het_hang nếu hết số lượng
	if lotHang.SoLuongConLai == 0 {
		lotHang.TrangThai = "het_hang"
	}

	now, err := txNow(ctx)
	if err != nil {
		return err
	}
	lotHang.NgayCapNhat = now

	lotHangJSON, err := json.Marshal(lotHang)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maLo, lotHangJSON)
}

// ThongKeLotHang trả về thống kê lô hàng theo HTX
// MSP được phép: Tất cả
func (c *LotHangContract) ThongKeLotHang(
	ctx contractapi.TransactionContextInterface,
	maHTX string,
) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"ma_htx":"%s"}}`, maHTX)
	danhSach, err := c.queryLotHang(ctx, queryString)
	if err != nil {
		return "", err
	}

	thongKe := map[string]interface{}{
		"ma_htx":     maHTX,
		"tong_so_lo": len(danhSach),
		"theo_trang_thai": map[string]int{
			"dang_trong":     0,
			"da_thu_hoach":   0,
			"cho_chung_nhan": 0,
			"san_sang_ban":   0,
			"het_hang":       0,
			"dinh_chi":       0,
		},
		"theo_loai_san_pham": map[string]int{},
	}

	ttMap := thongKe["theo_trang_thai"].(map[string]int)
	loaiMap := thongKe["theo_loai_san_pham"].(map[string]int)

	for _, lo := range danhSach {
		if _, ok := ttMap[lo.TrangThai]; ok {
			ttMap[lo.TrangThai]++
		}
		loaiMap[lo.LoaiSanPham]++
	}

	result, err := json.Marshal(thongKe)
	if err != nil {
		return "", fmt.Errorf("loi khi ma hoa thong ke: %v", err)
	}
	return string(result), nil
}

// --- Helper methods ---

// layLotHang lấy và deserialize dữ liệu lô hàng từ state
func (c *LotHangContract) layLotHang(
	ctx contractapi.TransactionContextInterface,
	maLo string,
) (*LotHang, error) {
	if maLo == "" {
		return nil, fmt.Errorf("ma_lo khong duoc de trong")
	}

	lotHangJSON, err := ctx.GetStub().GetState(maLo)
	if err != nil {
		return nil, fmt.Errorf("loi khi doc trang thai: %v", err)
	}
	if lotHangJSON == nil {
		return nil, fmt.Errorf("lo hang %s khong ton tai", maLo)
	}

	var lotHang LotHang
	if err := json.Unmarshal(lotHangJSON, &lotHang); err != nil {
		return nil, fmt.Errorf("loi khi giai ma du lieu: %v", err)
	}
	return &lotHang, nil
}

// queryLotHang thực hiện CouchDB rich query và trả về danh sách lô hàng
func (c *LotHangContract) queryLotHang(
	ctx contractapi.TransactionContextInterface,
	queryString string,
) ([]*LotHang, error) {
	iterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("loi khi thuc hien query: %v", err)
	}
	defer iterator.Close()

	var danhSach []*LotHang
	for iterator.HasNext() {
		item, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("loi khi doc ket qua query: %v", err)
		}

		var lotHang LotHang
		if err := json.Unmarshal(item.Value, &lotHang); err != nil {
			return nil, fmt.Errorf("loi khi giai ma du lieu: %v", err)
		}
		danhSach = append(danhSach, &lotHang)
	}

	return danhSach, nil
}

func main() {
	contract := new(LotHangContract)
	cc, err := contractapi.NewChaincode(contract)
	if err != nil {
		log.Panicf("Loi khi tao chaincode lohang_cc: %v", err)
	}
	if err := cc.Start(); err != nil {
		log.Panicf("Loi khi khoi dong chaincode lohang_cc: %v", err)
	}
}
