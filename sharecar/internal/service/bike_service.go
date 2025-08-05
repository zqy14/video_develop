package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"sharecar/internal/model"
	"sharecar/internal/mqtt"
)

// BikeService 单车服务
type BikeService struct {
	mqttClient *mqtt.MQTTClient
	rides      map[string]*model.RideRecord // 使用 "userID:deviceID" 作为键
	mutex      sync.RWMutex
}

// NewBikeService 创建新的单车服务
func NewBikeService(mqttClient *mqtt.MQTTClient) *BikeService {
	return &BikeService{
		mqttClient: mqttClient,
		rides:      make(map[string]*model.RideRecord),
		mutex:      sync.RWMutex{},
	}
}

// HandleBikeEvent 处理单车事件
func (s *BikeService) HandleBikeEvent(event model.BikeEvent) {
	// 解析时间
	eventTime := time.Now()
	if event.Time != "" {
		parsedTime, err := time.Parse("20060102150405", event.Time)
		if err == nil {
			eventTime = parsedTime
		}
	}

	// 生成键
	key := fmt.Sprintf("%s:%s", event.UserID, event.DeviceID)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 根据事件类型处理
	switch event.EventType {
	case "unlock":
		// 发送开锁信号并返回设备已开启
		err := s.mqttClient.SendUnlockSignal(event.DeviceID)
		if err != nil {
			log.Printf("Failed to send unlock signal: %v\n", err)
			return
		}

		// 创建新的骑行记录
		s.rides[key] = &model.RideRecord{
			UserID:      event.UserID,
			DeviceID:    event.DeviceID,
			StartTime:   eventTime,
			StartX:      event.X,
			StartY:      event.Y,
			IsCompleted: false,
		}

		log.Printf("Started new ride for user %s on bike %s\n", event.UserID, event.DeviceID)

	case "lock":
		// 检查是否存在进行中的骑行
		existingRide, exists := s.rides[key]
		if !exists || existingRide.IsCompleted {
			log.Printf("No active ride found for user %s on bike %s\n", event.UserID, event.DeviceID)
			return
		}

		// 结束现有骑行（关锁）
		existingRide.EndTime = eventTime
		existingRide.EndX = event.X
		existingRide.EndY = event.Y
		existingRide.Duration = eventTime.Sub(existingRide.StartTime)
		existingRide.IsCompleted = true

		// 计算费用（示例：每分钟1元）
		minutes := float64(existingRide.Duration.Minutes())
		existingRide.Fee = minutes * 1.0

		// 发送关锁信号
		err := s.mqttClient.SendLockSignal(event.DeviceID)
		if err != nil {
			log.Printf("Failed to send lock signal: %v\n", err)
		}

		log.Printf("Completed ride for user %s on bike %s\n", event.UserID, event.DeviceID)
		log.Printf("Ride duration: %.2f minutes, Fee: %.2f yuan\n", minutes, existingRide.Fee)

	default:
		log.Printf("Unknown event type: %s\n", event.EventType)
	}
}

// GetAllRides 获取所有骑行记录
func (s *BikeService) GetAllRides() []*model.RideRecord {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	rides := make([]*model.RideRecord, 0, len(s.rides))
	for _, ride := range s.rides {
		rides = append(rides, ride)
	}

	return rides
}