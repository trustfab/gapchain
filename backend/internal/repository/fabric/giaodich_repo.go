package fabric

import (
	"fmt"
	"log"

	infra "github.com/trustfab/gapchain/backend/internal/infrastructure/fabric"
)

type GiaodichRepo interface {
	Submit(mspID string, funcName string, args ...string) ([]byte, error)
	Evaluate(mspID string, funcName string, args ...string) ([]byte, error)
}

type giaodichRepo struct {
	registry *infra.GatewayRegistry
}

func NewGiaodichRepo(registry *infra.GatewayRegistry) GiaodichRepo {
	return &giaodichRepo{registry: registry}
}

func (r *giaodichRepo) Submit(mspID string, funcName string, args ...string) ([]byte, error) {
	orgGw, err := r.registry.GetOrgGateway(mspID)
	if err != nil {
		log.Printf("[GiaodichRepo Submit] Lỗi lấy Gateway cho mspID %s: %v", mspID, err)
		return nil, err
	}
	if orgGw.GiaoDichContract == nil {
		log.Printf("[GiaodichRepo Submit] Lỗi: org %s không tham gia giaodich-channel", mspID)
		return nil, fmt.Errorf("org %s khong tham gia giaodich-channel", mspID)
	}

	log.Printf("[GiaodichRepo Submit] Đang gọi '%s' bằng mspID '%s', args: %v", funcName, mspID, args)
	res, err := orgGw.GiaoDichContract.SubmitTransaction(funcName, args...)
	if err != nil {
		log.Printf("[GiaodichRepo Submit] ❌ LỖI FABRIC ('%s'): %v", funcName, err)
	} else {
		log.Printf("[GiaodichRepo Submit] ✅ THÀNH CÔNG ('%s')", funcName)
	}
	return res, err
}

func (r *giaodichRepo) Evaluate(mspID string, funcName string, args ...string) ([]byte, error) {
	orgGw, err := r.registry.GetOrgGateway(mspID)
	if err != nil {
		log.Printf("[GiaodichRepo Evaluate] Lỗi lấy Gateway cho mspID %s: %v", mspID, err)
		return nil, err
	}
	if orgGw.GiaoDichContract == nil {
		log.Printf("[GiaodichRepo Evaluate] Lỗi: org %s không tham gia giaodich-channel", mspID)
		return nil, fmt.Errorf("org %s khong tham gia giaodich-channel", mspID)
	}

	log.Printf("[GiaodichRepo Evaluate] Đang gọi '%s' bằng mspID '%s', args: %v", funcName, mspID, args)
	res, err := orgGw.GiaoDichContract.EvaluateTransaction(funcName, args...)
	if err != nil {
		log.Printf("[GiaodichRepo Evaluate] ❌ LỖI FABRIC ('%s'): %v", funcName, err)
	} else {
		log.Printf("[GiaodichRepo Evaluate] ✅ THÀNH CÔNG ('%s')", funcName)
	}
	return res, err
}
