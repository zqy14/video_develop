package model

import (
	"time"
)

// BikeEvent 表示单车事件数据
type BikeEvent struct {
	UserID    string  `json:"userid"`
	Time      string  `json:"time"`
	DeviceID  string  `json:"deviceId"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	EventType string  `json:"eventType"` // "unlock" 或 "lock"
	Status    string  `json:"status"`    // 设备状态信息
}

// RideRecord 表示一次骑行记录
type RideRecord struct {
	UserID      string
	DeviceID    string
	StartTime   time.Time
	EndTime     time.Time
	StartX      float64
	StartY      float64
	EndX        float64
	EndY        float64
	Duration    time.Duration
	Fee         float64
	IsCompleted bool
}