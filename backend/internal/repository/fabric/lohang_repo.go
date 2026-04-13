package fabric

import (
	"fmt"
	"log"

	infra "github.com/trustfab/gapchain/backend/internal/infrastructure/fabric"
)

type LohangRepo interface {
	Submit(mspID string, funcName string, args ...string) ([]byte, error)
	Evaluate(mspID string, funcName string, args ...string) ([]byte, error)
}

type lohangRepo struct {
	registry *infra.GatewayRegistry
}

func NewLohangRepo(registry *infra.GatewayRegistry) LohangRepo {
	return &lohangRepo{registry: registry}
}

func (r *lohangRepo) Submit(mspID string, funcName string, args ...string) ([]byte, error) {
	orgGw, err := r.registry.GetOrgGateway(mspID)
	if err != nil {
		log.Printf("[LohangRepo Submit] Lỗi lấy Gateway cho mspID %s: %v", mspID, err)
		return nil, err
	}
	if orgGw.LotHangContract == nil {
		log.Printf("[LohangRepo Submit] Lỗi: org %s không tham gia nhatky-htx-channel", mspID)
		return nil, fmt.Errorf("org %s khong tham gia nhatky-htx-channel", mspID)
	}

	log.Printf("[LohangRepo Submit] Đang gọi '%s' bằng mspID '%s', args: %v", funcName, mspID, args)
	res, err := orgGw.LotHangContract.SubmitTransaction(funcName, args...)
	if err != nil {
		log.Printf("[LohangRepo Submit] ❌ LỖI FABRIC ('%s'): %v", funcName, err)
	} else {
		log.Printf("[LohangRepo Submit] ✅ THÀNH CÔNG ('%s')", funcName)
	}
	return res, err
}

func (r *lohangRepo) Evaluate(mspID string, funcName string, args ...string) ([]byte, error) {
	orgGw, err := r.registry.GetOrgGateway(mspID)
	if err != nil {
		log.Printf("[LohangRepo Evaluate] Lỗi lấy Gateway cho mspID %s: %v", mspID, err)
		return nil, err
	}
	if orgGw.LotHangContract == nil {
		log.Printf("[LohangRepo Evaluate] Lỗi: org %s không tham gia nhatky-htx-channel", mspID)
		return nil, fmt.Errorf("org %s khong tham gia nhatky-htx-channel", mspID)
	}

	log.Printf("[LohangRepo Evaluate] Đang gọi '%s' bằng mspID '%s', args: %v", funcName, mspID, args)
	res, err := orgGw.LotHangContract.EvaluateTransaction(funcName, args...)
	if err != nil {
		log.Printf("[LohangRepo Evaluate] ❌ LỖI FABRIC ('%s'): %v", funcName, err)
	} else {
		log.Printf("[LohangRepo Evaluate] ✅ THÀNH CÔNG ('%s')", funcName)
	}
	return res, err
}
