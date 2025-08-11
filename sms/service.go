package main

import (
	"context"
	"rider-management/internal/model"
	"rider-management/pkg/backoff"
	"sync"
	"time"
)

// 短信发送状态
const (
	StatusSuccess  = "success"
	StatusFailed   = "failed"
	StatusRetrying = "retrying"
)

// Provider 短信供应商接口
type Provider interface {
	Send(ctx context.Context, phone, content string) error
	Name() string // 供应商名称
}

// SMSService 短信服务
type SMSService struct {
	providers []Provider
	repo      SMSRepository
	backoff   *backoff.ExponentialBackoff
}

// SMSRepository 短信记录仓库接口
type SMSRepository interface {
	Create(ctx context.Context, record *model.SMSRecord) error
	Update(ctx context.Context, record *model.SMSRecord) error
	GetFailedRecords(ctx context.Context) ([]*model.SMSRecord, error)
}

func NewSMSService(providers []Provider, repo SMSRepository) *SMSService {
	return &SMSService{
		providers: providers,
		repo:      repo,
		backoff:   backoff.NewExponentialBackoff(3, time.Second), // 最多重试3次，基础延迟1秒
	}
}

// SendSMS 发送短信 (带重试和供应商切换)
func (s *SMSService) SendSMS(ctx context.Context, riderID int64, phone, content string) error {
	// 创建初始记录
	record := &model.SMSRecord{
		RiderID: riderID,
		Phone:   phone,
		Content: content,
		Status:  StatusRetrying,
	}

	if err := s.repo.Create(ctx, record); err != nil {
		return err
	}

	// 使用第一个供应商尝试发送
	var _ error
	for _, provider := range s.providers {
		record.Provider = provider.Name()

		err := provider.Send(ctx, phone, content)
		if err == nil {
			record.Status = StatusSuccess
			return s.repo.Update(ctx, record)
		}

		_ = err
	}

	// 所有供应商都失败，设置重试
	record.Status = StatusFailed
	record.RetryCount++
	record.NextRetryAt = time.Now().Add(s.backoff.CalculateDelay(record.RetryCount))
	return s.repo.Update(ctx, record)
}

// RetryFailedSMS 重试失败的短信
func (s *SMSService) RetryFailedSMS(ctx context.Context) error {
	records, err := s.repo.GetFailedRecords(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, record := range records {
		if record.NextRetryAt.After(time.Now()) {
			continue // 未到重试时间
		}

		if !s.backoff.ShouldRetry(record.RetryCount) {
			continue // 超过最大重试次数
		}

		wg.Add(1)
		go func(r *model.SMSRecord) {
			defer wg.Done()
			_ = s.SendSMS(ctx, r.RiderID, r.Phone, r.Content)
		}(record)
	}

	wg.Wait()
	return nil
}
