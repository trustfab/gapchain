package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// GiaoDichContract cung cấp các chức năng quản lý giao dịch mua bán nông sản
type GiaoDichContract struct {
	contractapi.Contract
}

// GiaoDich đại diện cho một giao dịch mua bán lô hàng giữa HTX và NPP
type GiaoDich struct {
	MaGiaoDich  string  `json:"ma_giao_dich"`
	MaLo        string  `json:"ma_lo"`          // Liên kết lô hàng từ lohang_cc — BẮT BUỘC
	MaHTX       string  `json:"ma_htx"`
	MaNPP       string  `json:"ma_npp"`
	SanPham     string  `json:"san_pham"`       // Tên sản phẩm (bản sao từ lô hàng)
	SoLuong     float64 `json:"so_luong"`
	DonViTinh   string  `json:"don_vi_tinh"`    // kg | tan | thung | qua
	DonGia      float64 `json:"don_gia"`        // VND
	TongTien    float64 `json:"tong_tien"`      // = SoLuong * DonGia
	TyLeHoaHong float64 `json:"ty_le_hoa_hong"` // % hoa hồng Platform, thiết lập khi tạo GD
	TienHoaHong float64 `json:"tien_hoa_hong"`  // = TongTien * TyLeHoaHong / 100 (tính sẵn)
	TrangThai   string  `json:"trang_thai"`     // cho_duyet | da_duyet | dang_giao | da_giao | cho_thanh_toan | da_thanh_toan | huy_bo
	GhiChu      string  `json:"ghi_chu"`
	NgayTao     string  `json:"ngay_tao"`
	NgayCapNhat string  `json:"ngay_cap_nhat"`
}

// LichSuRecord đại diện cho 1 bản ghi trong lịch sử thay đổi
type LichSuRecord struct {
	TxID      string `json:"tx_id"`
	Timestamp string `json:"timestamp"`
	IsDeleted bool   `json:"is_deleted"`
	Data      string `json:"data"`
}

// TinhHoaHongKetQua chứa kết quả tính hoa hồng NPP
type TinhHoaHongKetQua struct {
	MaNPP           string  `json:"ma_npp"`
	TongGiaoDich    int     `json:"tong_giao_dich"`
	TongDoanhThu    float64 `json:"tong_doanh_thu"`
	TongTienHoaHong float64 `json:"tong_tien_hoa_hong"`
}

// Vòng đời trạng thái hợp lệ:
// cho_duyet → da_duyet → dang_giao → da_giao → cho_thanh_toan → da_thanh_toan
// Bất kỳ trạng thái nào (trừ da_thanh_toan) → huy_bo (chỉ Platform)
var chuyenTrangThaiHopLe = map[string][]string{
	"cho_duyet":      {"da_duyet", "huy_bo"},
	"da_duyet":       {"dang_giao", "huy_bo"},
	"dang_giao":      {"da_giao", "huy_bo"},
	"da_giao":        {"cho_thanh_toan", "huy_bo"},
	"cho_thanh_toan": {"da_thanh_toan", "huy_bo"},
	"da_thanh_toan":  {}, // trạng thái kết thúc
	"huy_bo":         {}, // trạng thái kết thúc
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

// TaoGiaoDich tạo một giao dịch mua bán lô hàng mới
// MSP được phép: HTXNongSanOrgMSP, PlatformOrgMSP
func (c *GiaoDichContract) TaoGiaoDich(
	ctx contractapi.TransactionContextInterface,
	maGiaoDich string,
	maLo string,
	maHTX string,
	maNPP string,
	sanPham string,
	soLuong float64,
	donViTinh string,
	donGia float64,
	tyLeHoaHong float64,
	ghiChu string,
) error {
	if err := kiemTraMSP(ctx, "HTXNongSanOrgMSP", "PlatformOrgMSP"); err != nil {
		return err
	}

	if maGiaoDich == "" || maLo == "" || maHTX == "" || maNPP == "" {
		return fmt.Errorf("ma_giao_dich, ma_lo, ma_htx va ma_npp khong duoc de trong")
	}
	if soLuong <= 0 {
		return fmt.Errorf("so_luong phai lon hon 0")
	}
	if donGia <= 0 {
		return fmt.Errorf("don_gia phai lon hon 0")
	}
	if tyLeHoaHong < 0 || tyLeHoaHong > 100 {
		return fmt.Errorf("ty_le_hoa_hong phai tu 0 den 100")
	}

	// Kiểm tra giao dịch chưa tồn tại
	existing, err := ctx.GetStub().GetState(maGiaoDich)
	if err != nil {
		return fmt.Errorf("loi khi kiem tra trang thai: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("giao dich %s da ton tai", maGiaoDich)
	}

	tongTien := soLuong * donGia
	tienHoaHong := tongTien * tyLeHoaHong / 100

	now, err := txNow(ctx)
	if err != nil {
		return err
	}

	giaoDich := GiaoDich{
		MaGiaoDich:  maGiaoDich,
		MaLo:        maLo,
		MaHTX:       maHTX,
		MaNPP:       maNPP,
		SanPham:     sanPham,
		SoLuong:     soLuong,
		DonViTinh:   donViTinh,
		DonGia:      donGia,
		TongTien:    tongTien,
		TyLeHoaHong: tyLeHoaHong,
		TienHoaHong: tienHoaHong,
		TrangThai:   "cho_duyet",
		GhiChu:      ghiChu,
		NgayTao:     now,
		NgayCapNhat: now,
	}

	giaoDichJSON, err := json.Marshal(giaoDich)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maGiaoDich, giaoDichJSON)
}

// DuyetGiaoDich phê duyệt giao dịch, chuyển từ cho_duyet → da_duyet
// MSP được phép: PlatformOrgMSP
func (c *GiaoDichContract) DuyetGiaoDich(
	ctx contractapi.TransactionContextInterface,
	maGiaoDich string,
) error {
	if err := kiemTraMSP(ctx, "PlatformOrgMSP"); err != nil {
		return err
	}

	giaoDich, err := c.layGiaoDich(ctx, maGiaoDich)
	if err != nil {
		return err
	}

	if giaoDich.TrangThai != "cho_duyet" {
		return fmt.Errorf("giao dich %s khong o trang thai cho_duyet (hien tai: %s)", maGiaoDich, giaoDich.TrangThai)
	}

	now, err := txNow(ctx)
	if err != nil {
		return err
	}

	giaoDich.TrangThai = "da_duyet"
	giaoDich.NgayCapNhat = now

	giaoDichJSON, err := json.Marshal(giaoDich)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maGiaoDich, giaoDichJSON)
}

// CapNhatTrangThai chuyển trạng thái giao dịch theo vòng đời
// Quy tắc MSP theo trạng thái đích:
//   da_duyet       → PlatformOrgMSP (đã xử lý riêng trong DuyetGiaoDich)
//   dang_giao      → PlatformOrgMSP
//   da_giao        → NPPXanhOrgMSP, PlatformOrgMSP
//   cho_thanh_toan → PlatformOrgMSP
//   da_thanh_toan  → NPPXanhOrgMSP, PlatformOrgMSP
//   huy_bo         → PlatformOrgMSP (chỉ Platform mới hủy được)
func (c *GiaoDichContract) CapNhatTrangThai(
	ctx contractapi.TransactionContextInterface,
	maGiaoDich string,
	trangThaiMoi string,
	ghiChu string,
) error {
	// Kiểm tra quyền theo trạng thái đích
	switch trangThaiMoi {
	case "dang_giao", "cho_thanh_toan", "huy_bo":
		if err := kiemTraMSP(ctx, "PlatformOrgMSP"); err != nil {
			return err
		}
	case "da_giao", "da_thanh_toan":
		if err := kiemTraMSP(ctx, "NPPXanhOrgMSP", "PlatformOrgMSP"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("trang_thai_moi khong hop le: %s", trangThaiMoi)
	}

	giaoDich, err := c.layGiaoDich(ctx, maGiaoDich)
	if err != nil {
		return err
	}

	// Kiểm tra vòng đời hợp lệ
	choPhep := chuyenTrangThaiHopLe[giaoDich.TrangThai]
	hopLe := false
	for _, tt := range choPhep {
		if tt == trangThaiMoi {
			hopLe = true
			break
		}
	}
	if !hopLe {
		return fmt.Errorf("khong the chuyen trang thai tu '%s' sang '%s'", giaoDich.TrangThai, trangThaiMoi)
	}

	now, err := txNow(ctx)
	if err != nil {
		return err
	}

	giaoDich.TrangThai = trangThaiMoi
	if ghiChu != "" {
		giaoDich.GhiChu = ghiChu
	}
	giaoDich.NgayCapNhat = now

	giaoDichJSON, err := json.Marshal(giaoDich)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maGiaoDich, giaoDichJSON)
}

// DocGiaoDich đọc thông tin một giao dịch theo mã giao dịch
// MSP được phép: Tất cả
func (c *GiaoDichContract) DocGiaoDich(
	ctx contractapi.TransactionContextInterface,
	maGiaoDich string,
) (*GiaoDich, error) {
	return c.layGiaoDich(ctx, maGiaoDich)
}

// DocGiaoDichTheoHTX trả về danh sách giao dịch của một HTX (CouchDB rich query)
// MSP được phép: Tất cả
func (c *GiaoDichContract) DocGiaoDichTheoHTX(
	ctx contractapi.TransactionContextInterface,
	maHTX string,
) ([]*GiaoDich, error) {
	if maHTX == "" {
		return nil, fmt.Errorf("ma_htx khong duoc de trong")
	}
	queryString := fmt.Sprintf(`{"selector":{"ma_htx":"%s","ngay_tao":{"$gt":null}},"sort":[{"ngay_tao":"desc"}]}`, maHTX)
	return c.queryGiaoDich(ctx, queryString)
}

// DocGiaoDichTheoNPP trả về danh sách giao dịch của một NPP (CouchDB rich query)
// MSP được phép: Tất cả
func (c *GiaoDichContract) DocGiaoDichTheoNPP(
	ctx contractapi.TransactionContextInterface,
	maNPP string,
) ([]*GiaoDich, error) {
	if maNPP == "" {
		return nil, fmt.Errorf("ma_npp khong duoc de trong")
	}
	queryString := fmt.Sprintf(`{"selector":{"ma_npp":"%s","ngay_tao":{"$gt":null}},"sort":[{"ngay_tao":"desc"}]}`, maNPP)
	return c.queryGiaoDich(ctx, queryString)
}

// DocCongNoNPP trả về các giao dịch đang ở trạng thái "cho_thanh_toan" của một NPP
// Dùng để NPP xem danh sách công nợ cần thanh toán
// MSP được phép: NPPXanhOrgMSP, PlatformOrgMSP
func (c *GiaoDichContract) DocCongNoNPP(
	ctx contractapi.TransactionContextInterface,
	maNPP string,
) ([]*GiaoDich, error) {
	if err := kiemTraMSP(ctx, "NPPXanhOrgMSP", "PlatformOrgMSP"); err != nil {
		return nil, err
	}
	if maNPP == "" {
		return nil, fmt.Errorf("ma_npp khong duoc de trong")
	}
	queryString := fmt.Sprintf(
		`{"selector":{"ma_npp":"%s","trang_thai":"cho_thanh_toan","ngay_tao":{"$gt":null}},"sort":[{"ngay_tao":"asc"}]}`,
		maNPP,
	)
	return c.queryGiaoDich(ctx, queryString)
}

// TinhHoaHongNPP tính tổng hoa hồng của một NPP từ các giao dịch đã thanh toán
// MSP được phép: PlatformOrgMSP
func (c *GiaoDichContract) TinhHoaHongNPP(
	ctx contractapi.TransactionContextInterface,
	maNPP string,
) (string, error) {
	if err := kiemTraMSP(ctx, "PlatformOrgMSP"); err != nil {
		return "", err
	}
	if maNPP == "" {
		return "", fmt.Errorf("ma_npp khong duoc de trong")
	}

	queryString := fmt.Sprintf(
		`{"selector":{"ma_npp":"%s","trang_thai":"da_thanh_toan"}}`,
		maNPP,
	)
	danhSach, err := c.queryGiaoDich(ctx, queryString)
	if err != nil {
		return "", err
	}

	ketQua := TinhHoaHongKetQua{
		MaNPP: maNPP,
	}
	for _, gd := range danhSach {
		ketQua.TongGiaoDich++
		ketQua.TongDoanhThu += gd.TongTien
		ketQua.TongTienHoaHong += gd.TienHoaHong
	}

	result, err := json.Marshal(ketQua)
	if err != nil {
		return "", fmt.Errorf("loi khi ma hoa ket qua: %v", err)
	}
	return string(result), nil
}

// LichSuGiaoDich trả về toàn bộ lịch sử thay đổi trạng thái của một giao dịch (audit trail)
// MSP được phép: Tất cả
func (c *GiaoDichContract) LichSuGiaoDich(
	ctx contractapi.TransactionContextInterface,
	maGiaoDich string,
) ([]*LichSuRecord, error) {
	// Kiểm tra giao dịch tồn tại
	if _, err := c.layGiaoDich(ctx, maGiaoDich); err != nil {
		return nil, err
	}

	iterator, err := ctx.GetStub().GetHistoryForKey(maGiaoDich)
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

// ThongKeGiaoDich thống kê giao dịch theo HTX hoặc NPP
// Truyền maHTX hoặc maNPP, truyền chuỗi rỗng nếu không lọc theo trường đó
// MSP được phép: Tất cả
func (c *GiaoDichContract) ThongKeGiaoDich(
	ctx contractapi.TransactionContextInterface,
	maHTX string,
	maNPP string,
) (string, error) {
	var queryString string
	switch {
	case maHTX != "" && maNPP != "":
		queryString = fmt.Sprintf(`{"selector":{"ma_htx":"%s","ma_npp":"%s"}}`, maHTX, maNPP)
	case maHTX != "":
		queryString = fmt.Sprintf(`{"selector":{"ma_htx":"%s"}}`, maHTX)
	case maNPP != "":
		queryString = fmt.Sprintf(`{"selector":{"ma_npp":"%s"}}`, maNPP)
	default:
		return "", fmt.Errorf("phai truyen it nhat ma_htx hoac ma_npp")
	}

	danhSach, err := c.queryGiaoDich(ctx, queryString)
	if err != nil {
		return "", err
	}

	thongKe := map[string]interface{}{
		"ma_htx":  maHTX,
		"ma_npp":  maNPP,
		"tong_so": len(danhSach),
		"theo_trang_thai": map[string]int{
			"cho_duyet":      0,
			"da_duyet":       0,
			"dang_giao":      0,
			"da_giao":        0,
			"cho_thanh_toan": 0,
			"da_thanh_toan":  0,
			"huy_bo":         0,
		},
		"tong_tien":         0.0,
		"tong_tien_hoa_hong": 0.0,
	}

	ttMap := thongKe["theo_trang_thai"].(map[string]int)

	for _, gd := range danhSach {
		if _, ok := ttMap[gd.TrangThai]; ok {
			ttMap[gd.TrangThai]++
		}
		thongKe["tong_tien"] = thongKe["tong_tien"].(float64) + gd.TongTien
		thongKe["tong_tien_hoa_hong"] = thongKe["tong_tien_hoa_hong"].(float64) + gd.TienHoaHong
	}

	result, err := json.Marshal(thongKe)
	if err != nil {
		return "", fmt.Errorf("loi khi ma hoa thong ke: %v", err)
	}
	return string(result), nil
}

// --- Helper methods ---

// layGiaoDich lấy và deserialize dữ liệu giao dịch từ state
func (c *GiaoDichContract) layGiaoDich(
	ctx contractapi.TransactionContextInterface,
	maGiaoDich string,
) (*GiaoDich, error) {
	if maGiaoDich == "" {
		return nil, fmt.Errorf("ma_giao_dich khong duoc de trong")
	}

	giaoDichJSON, err := ctx.GetStub().GetState(maGiaoDich)
	if err != nil {
		return nil, fmt.Errorf("loi khi doc trang thai: %v", err)
	}
	if giaoDichJSON == nil {
		return nil, fmt.Errorf("giao dich %s khong ton tai", maGiaoDich)
	}

	var giaoDich GiaoDich
	if err := json.Unmarshal(giaoDichJSON, &giaoDich); err != nil {
		return nil, fmt.Errorf("loi khi giai ma du lieu: %v", err)
	}
	return &giaoDich, nil
}

// queryGiaoDich thực hiện CouchDB rich query và trả về danh sách giao dịch
func (c *GiaoDichContract) queryGiaoDich(
	ctx contractapi.TransactionContextInterface,
	queryString string,
) ([]*GiaoDich, error) {
	iterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("loi khi thuc hien query: %v", err)
	}
	defer iterator.Close()

	var danhSach []*GiaoDich
	for iterator.HasNext() {
		item, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("loi khi doc ket qua query: %v", err)
		}

		var giaoDich GiaoDich
		if err := json.Unmarshal(item.Value, &giaoDich); err != nil {
			return nil, fmt.Errorf("loi khi giai ma du lieu: %v", err)
		}
		danhSach = append(danhSach, &giaoDich)
	}

	return danhSach, nil
}

func main() {
	cc, err := contractapi.NewChaincode(&GiaoDichContract{})
	if err != nil {
		log.Panicf("Loi khi tao chaincode giao_dich_cc: %v", err)
	}
	if err := cc.Start(); err != nil {
		log.Panicf("Loi khi khoi dong chaincode giao_dich_cc: %v", err)
	}
}
