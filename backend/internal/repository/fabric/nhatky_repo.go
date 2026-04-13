package fabric

import (
	"fmt"
	"log"

	infra "github.com/trustfab/gapchain/backend/internal/infrastructure/fabric"
)

type NhatkyRepo interface {
	Submit(mspID string, funcName string, args ...string) ([]byte, error)
	Evaluate(mspID string, funcName string, args ...string) ([]byte, error)
}

type nhatkyRepo struct {
	registry *infra.GatewayRegistry
}

func NewNhatkyRepo(registry *infra.GatewayRegistry) NhatkyRepo {
	return &nhatkyRepo{registry: registry}
}

func (r *nhatkyRepo) Submit(mspID string, funcName string, args ...string) ([]byte, error) {
	orgGw, err := r.registry.GetOrgGateway(mspID)
	if err != nil {
		log.Printf("[NhatkyRepo Submit] Lỗi lấy Gateway cho mspID %s: %v", mspID, err)
		return nil, err
	}
	if orgGw.NhatKyContract == nil {
		log.Printf("[NhatkyRepo Submit] Lỗi: org %s không tham gia nhatky-htx-channel", mspID)
		return nil, fmt.Errorf("org %s khong tham gia nhatky-htx-channel", mspID)
	}

	log.Printf("[NhatkyRepo Submit] Đang gọi '%s' bằng mspID '%s', args: %v", funcName, mspID, args)
	res, err := orgGw.NhatKyContract.SubmitTransaction(funcName, args...)
	if err != nil {
		log.Printf("[NhatkyRepo Submit] ❌ LỖI FABRIC ('%s'): %v", funcName, err)
	} else {
		log.Printf("[NhatkyRepo Submit] ✅ THÀNH CÔNG ('%s')", funcName)
	}
	return res, err
}

func (r *nhatkyRepo) Evaluate(mspID string, funcName string, args ...string) ([]byte, error) {
	orgGw, err := r.registry.GetOrgGateway(mspID)
	if err != nil {
		log.Printf("[NhatkyRepo Evaluate] Lỗi lấy Gateway cho mspID %s: %v", mspID, err)
		return nil, err
	}
	if orgGw.NhatKyContract == nil {
		log.Printf("[NhatkyRepo Evaluate] Lỗi: org %s không tham gia nhatky-htx-channel", mspID)
		return nil, fmt.Errorf("org %s khong tham gia nhatky-htx-channel", mspID)
	}

	log.Printf("[NhatkyRepo Evaluate] Đang gọi '%s' bằng mspID '%s', args: %v", funcName, mspID, args)
	res, err := orgGw.NhatKyContract.EvaluateTransaction(funcName, args...)
	if err != nil {
		log.Printf("[NhatkyRepo Evaluate] ❌ LỖI FABRIC ('%s'): %v", funcName, err)
	} else {
		log.Printf("[NhatkyRepo Evaluate] ✅ THÀNH CÔNG ('%s')", funcName)
	}
	return res, err
}
