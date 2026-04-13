package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// NhatKyContract cung cấp các chức năng quản lý nhật ký canh tác HTX
type NhatKyContract struct {
	contractapi.Contract
}

// NhatKyHTX đại diện cho một bản ghi nhật ký canh tác
type NhatKyHTX struct {
	MaNhatKy      string `json:"ma_nhat_ky"`
	MaLo          string `json:"ma_lo"`           // Liên kết lô hàng trong lohang_cc — BẮT BUỘC
	MaHTX         string `json:"ma_htx"`
	LoaiHoatDong  string `json:"loai_hoat_dong"`  // Danh mục động — chaincode chỉ validate != "", backend validate theo danh mục
	ChiTiet       string `json:"chi_tiet"`
	ViTri         string `json:"vi_tri"`
	NguoiThucHien string `json:"nguoi_thuc_hien"` // ID kỹ thuật viên thực hiện
	MinhChungHash string `json:"minh_chung_hash"` // SHA256 của ảnh/file bằng chứng (upload ở app layer)
	TrangThai     string `json:"trang_thai"`      // cho_nong_dan_xac_nhan | cho_duyet | da_duyet | tu_choi
	LyDoTuChoi    string `json:"ly_do_tu_choi"`   // Bắt buộc điền khi trang_thai = tu_choi
	NguoiDuyet    string `json:"nguoi_duyet"`      // ID người duyệt
	NgayGhi       string `json:"ngay_ghi"`         // Ngày thực hiện hoạt động (do người dùng nhập)
	NgayTao       string `json:"ngay_tao"`
	NgayCapNhat   string `json:"ngay_cap_nhat"`
}

// LichSuRecord đại diện cho 1 bản ghi trong lịch sử thay đổi
type LichSuRecord struct {
	TxID      string `json:"tx_id"`
	Timestamp string `json:"timestamp"`
	IsDeleted bool   `json:"is_deleted"`
	Data      string `json:"data"`
}

// loaiHoatDong: KHÔNG hardcode enum trong chaincode.
// Danh mục loại hoạt động do backend/frontend quản lý (có thể thay đổi theo loại cây trồng).
// Chaincode chỉ validate loai_hoat_dong không rỗng.

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

// GhiNhatKy tạo một bản ghi nhật ký canh tác mới, gắn với lô hàng
// MSP được phép: HTXNongSanOrgMSP, PlatformOrgMSP
func (c *NhatKyContract) GhiNhatKy(
	ctx contractapi.TransactionContextInterface,
	maNhatKy string,
	maLo string,
	maHTX string,
	loaiHoatDong string,
	chiTiet string,
	viTri string,
	nguoiThucHien string,
	ngayGhi string,
) error {
	if err := kiemTraMSP(ctx, "HTXNongSanOrgMSP", "PlatformOrgMSP"); err != nil {
		return err
	}

	if maNhatKy == "" || maLo == "" || maHTX == "" {
		return fmt.Errorf("ma_nhat_ky, ma_lo va ma_htx khong duoc de trong")
	}
	if loaiHoatDong == "" {
		return fmt.Errorf("loai_hoat_dong khong duoc de trong")
	}

	// Kiểm tra nhật ký chưa tồn tại
	existing, err := ctx.GetStub().GetState(maNhatKy)
	if err != nil {
		return fmt.Errorf("loi khi kiem tra trang thai: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("nhat ky %s da ton tai", maNhatKy)
	}

	now, err := txNow(ctx)
	if err != nil {
		return err
	}

	nhatKy := NhatKyHTX{
		MaNhatKy:      maNhatKy,
		MaLo:          maLo,
		MaHTX:         maHTX,
		LoaiHoatDong:  loaiHoatDong,
		ChiTiet:       chiTiet,
		ViTri:         viTri,
		NguoiThucHien: nguoiThucHien,
		MinhChungHash: "",
		TrangThai:     "cho_nong_dan_xac_nhan",
		LyDoTuChoi:    "",
		NguoiDuyet:    "",
		NgayGhi:       ngayGhi,
		NgayTao:       now,
		NgayCapNhat:   now,
	}

	nhatKyJSON, err := json.Marshal(nhatKy)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maNhatKy, nhatKyJSON)
}

// XacNhanNhatKy xác nhận nhật ký thực địa bằng Hash OTP/Ảnh (Farmer Bridge)
// MSP được phép: HTXNongSanOrgMSP, PlatformOrgMSP
func (c *NhatKyContract) XacNhanNhatKy(
	ctx contractapi.TransactionContextInterface,
	maNhatKy string,
	minhChungHash string,
) error {
	if err := kiemTraMSP(ctx, "HTXNongSanOrgMSP", "PlatformOrgMSP"); err != nil {
		return err
	}

	if maNhatKy == "" || minhChungHash == "" {
		return fmt.Errorf("ma_nhat_ky va minh_chung_hash khong duoc de trong")
	}

	nhatKy, err := c.layNhatKy(ctx, maNhatKy)
	if err != nil {
		return err
	}

	if nhatKy.TrangThai != "cho_nong_dan_xac_nhan" {
		return fmt.Errorf("nhat ky %s khong o trang thai cho_nong_dan_xac_nhan (hien tai: %s)", maNhatKy, nhatKy.TrangThai)
	}

	now, err := txNow(ctx)
	if err != nil {
		return err
	}

	nhatKy.TrangThai = "cho_duyet"
	nhatKy.MinhChungHash = minhChungHash
	nhatKy.NgayCapNhat = now

	nhatKyJSON, err := json.Marshal(nhatKy)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maNhatKy, nhatKyJSON)
}

// DuyetNhatKy phê duyệt hoặc từ chối nhật ký canh tác
// MSP được phép: ChiCucBVTVOrgMSP, PlatformOrgMSP
func (c *NhatKyContract) DuyetNhatKy(
	ctx contractapi.TransactionContextInterface,
	maNhatKy string,
	nguoiDuyet string,
	trangThaiMoi string,
	lyDoTuChoi string,
) error {
	if err := kiemTraMSP(ctx, "ChiCucBVTVOrgMSP", "PlatformOrgMSP"); err != nil {
		return err
	}

	if trangThaiMoi != "da_duyet" && trangThaiMoi != "tu_choi" {
		return fmt.Errorf("trang_thai_moi phai la 'da_duyet' hoac 'tu_choi'")
	}
	if trangThaiMoi == "tu_choi" && lyDoTuChoi == "" {
		return fmt.Errorf("ly_do_tu_choi bat buoc khi tu choi nhat ky")
	}

	nhatKy, err := c.layNhatKy(ctx, maNhatKy)
	if err != nil {
		return err
	}

	if nhatKy.TrangThai != "cho_duyet" {
		return fmt.Errorf("nhat ky %s khong o trang thai cho_duyet (hien tai: %s)", maNhatKy, nhatKy.TrangThai)
	}

	now, err := txNow(ctx)
	if err != nil {
		return err
	}

	nhatKy.TrangThai = trangThaiMoi
	nhatKy.NguoiDuyet = nguoiDuyet
	nhatKy.LyDoTuChoi = lyDoTuChoi
	nhatKy.NgayCapNhat = now

	nhatKyJSON, err := json.Marshal(nhatKy)
	if err != nil {
		return fmt.Errorf("loi khi ma hoa du lieu: %v", err)
	}

	return ctx.GetStub().PutState(maNhatKy, nhatKyJSON)
}

// DocNhatKy đọc thông tin một nhật ký theo mã nhật ký
// MSP được phép: Tất cả
func (c *NhatKyContract) DocNhatKy(
	ctx contractapi.TransactionContextInterface,
	maNhatKy string,
) (*NhatKyHTX, error) {
	return c.layNhatKy(ctx, maNhatKy)
}

// DocNhatKyTheoLo trả về tất cả nhật ký của một lô hàng (CouchDB rich query)
// MSP được phép: Tất cả
func (c *NhatKyContract) DocNhatKyTheoLo(
	ctx contractapi.TransactionContextInterface,
	maLo string,
) ([]*NhatKyHTX, error) {
	if maLo == "" {
		return nil, fmt.Errorf("ma_lo khong duoc de trong")
	}
	queryString := fmt.Sprintf(`{"selector":{"ma_lo":"%s","ngay_ghi":{"$gt":null}},"sort":[{"ngay_ghi":"asc"}]}`, maLo)
	return c.queryNhatKy(ctx, queryString)
}

// DocNhatKyTheoHTX trả về tất cả nhật ký của một HTX (CouchDB rich query)
// MSP được phép: Tất cả
func (c *NhatKyContract) DocNhatKyTheoHTX(
	ctx contractapi.TransactionContextInterface,
	maHTX string,
) ([]*NhatKyHTX, error) {
	if maHTX == "" {
		return nil, fmt.Errorf("ma_htx khong duoc de trong")
	}
	queryString := fmt.Sprintf(`{"selector":{"ma_htx":"%s","ngay_tao":{"$gt":null}},"sort":[{"ngay_tao":"desc"}]}`, maHTX)
	return c.queryNhatKy(ctx, queryString)
}

// DocNhatKyTheoTrangThai trả về nhật ký theo trạng thái (CouchDB rich query)
// MSP được phép: Tất cả
func (c *NhatKyContract) DocNhatKyTheoTrangThai(
	ctx contractapi.TransactionContextInterface,
	trangThai string,
) ([]*NhatKyHTX, error) {
	validTrangThai := map[string]bool{
		"cho_nong_dan_xac_nhan": true,
		"cho_duyet":             true,
		"da_duyet":              true,
		"tu_choi":               true,
	}
	if !validTrangThai[trangThai] {
		return nil, fmt.Errorf("trang_thai khong hop le: %s. Cac gia tri hop le: cho_nong_dan_xac_nhan, cho_duyet, da_duyet, tu_choi", trangThai)
	}
	queryString := fmt.Sprintf(`{"selector":{"trang_thai":"%s","ngay_tao":{"$gt":null}},"sort":[{"ngay_tao":"desc"}]}`, trangThai)
	return c.queryNhatKy(ctx, queryString)
}

// DocNhatKyTheoLoVaTrangThai trả về nhật ký theo lô hàng và trạng thái (CouchDB rich query)
// MSP được phép: Tất cả
func (c *NhatKyContract) DocNhatKyTheoLoVaTrangThai(
	ctx contractapi.TransactionContextInterface,
	maLo string,
	trangThai string,
) ([]*NhatKyHTX, error) {
	if maLo == "" {
		return nil, fmt.Errorf("ma_lo khong duoc de trong")
	}
	queryString := fmt.Sprintf(
		`{"selector":{"ma_lo":"%s","trang_thai":"%s","ngay_ghi":{"$gt":null}},"sort":[{"ngay_ghi":"asc"}]}`,
		maLo, trangThai,
	)
	return c.queryNhatKy(ctx, queryString)
}

// LichSuNhatKy trả về toàn bộ lịch sử thay đổi của một nhật ký (audit trail thực sự)
// MSP được phép: Tất cả
func (c *NhatKyContract) LichSuNhatKy(
	ctx contractapi.TransactionContextInterface,
	maNhatKy string,
) ([]*LichSuRecord, error) {
	// Kiểm tra nhật ký tồn tại
	if _, err := c.layNhatKy(ctx, maNhatKy); err != nil {
		return nil, err
	}

	iterator, err := ctx.GetStub().GetHistoryForKey(maNhatKy)
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

// ThongKeNhatKy thống kê nhật ký theo HTX
// MSP được phép: Tất cả
func (c *NhatKyContract) ThongKeNhatKy(
	ctx contractapi.TransactionContextInterface,
	maHTX string,
) (string, error) {
	var queryString string
	if maHTX == "ALL" || maHTX == "" {
		queryString = `{"selector":{"ma_nhat_ky":{"$gt":null}}}`
	} else {
		queryString = fmt.Sprintf(`{"selector":{"ma_htx":"%s"}}`, maHTX)
	}
	danhSach, err := c.queryNhatKy(ctx, queryString)
	if err != nil {
		return "", err
	}

	thongKe := map[string]interface{}{
		"ma_htx":     maHTX,
		"tong_so":    len(danhSach),
		"theo_trang_thai": map[string]int{
			"cho_nong_dan_xac_nhan": 0,
			"cho_duyet":             0,
			"da_duyet":              0,
			"tu_choi":               0,
		},
		"theo_loai_hoat_dong": map[string]int{},
	}

	ttMap := thongKe["theo_trang_thai"].(map[string]int)
	loaiMap := thongKe["theo_loai_hoat_dong"].(map[string]int)

	for _, nk := range danhSach {
		if _, ok := ttMap[nk.TrangThai]; ok {
			ttMap[nk.TrangThai]++
		}
		loaiMap[nk.LoaiHoatDong]++
	}

	result, err := json.Marshal(thongKe)
	if err != nil {
		return "", fmt.Errorf("loi khi ma hoa thong ke: %v", err)
	}
	return string(result), nil
}

// --- Helper methods ---

// layNhatKy lấy và deserialize dữ liệu nhật ký từ state
func (c *NhatKyContract) layNhatKy(
	ctx contractapi.TransactionContextInterface,
	maNhatKy string,
) (*NhatKyHTX, error) {
	if maNhatKy == "" {
		return nil, fmt.Errorf("ma_nhat_ky khong duoc de trong")
	}

	nhatKyJSON, err := ctx.GetStub().GetState(maNhatKy)
	if err != nil {
		return nil, fmt.Errorf("loi khi doc trang thai: %v", err)
	}
	if nhatKyJSON == nil {
		return nil, fmt.Errorf("nhat ky %s khong ton tai", maNhatKy)
	}

	var nhatKy NhatKyHTX
	if err := json.Unmarshal(nhatKyJSON, &nhatKy); err != nil {
		return nil, fmt.Errorf("loi khi giai ma du lieu: %v", err)
	}
	return &nhatKy, nil
}

// queryNhatKy thực hiện CouchDB rich query và trả về danh sách nhật ký
func (c *NhatKyContract) queryNhatKy(
	ctx contractapi.TransactionContextInterface,
	queryString string,
) ([]*NhatKyHTX, error) {
	iterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("loi khi thuc hien query: %v", err)
	}
	defer iterator.Close()

	var danhSach []*NhatKyHTX
	for iterator.HasNext() {
		item, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("loi khi doc ket qua query: %v", err)
		}

		var nhatKy NhatKyHTX
		if err := json.Unmarshal(item.Value, &nhatKy); err != nil {
			return nil, fmt.Errorf("loi khi giai ma du lieu: %v", err)
		}
		danhSach = append(danhSach, &nhatKy)
	}

	return danhSach, nil
}

func main() {
	cc, err := contractapi.NewChaincode(&NhatKyContract{})
	if err != nil {
		log.Panicf("Loi khi tao chaincode nhatky_cc: %v", err)
	}
	if err := cc.Start(); err != nil {
		log.Panicf("Loi khi khoi dong chaincode nhatky_cc: %v", err)
	}
}
