package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"rider-management/pkg/emqx"
	"time"
)

func main() {
	// 1. 初始化EMQX客户端 (模拟快递柜设备)
	emqxClient, err := emqx.New(
		"tcp://emqx-server:1883", // EMQX服务器地址
		"locker-device-001",      // 设备ID
		"admin",                  // 用户名
		"public",                 // 密码
		"bird",                   // 主题
	)
	if err != nil {
		log.Fatalf("EMQX客户端初始化失败: %v", err)
	}
	defer emqxClient.Close()

	// 2. 订阅开锁指令
	err = emqxClient.Subscribe(context.Background(), func(payload []byte) {
		var cmd map[string]interface{}
		if err := json.Unmarshal(payload, &cmd); err != nil {
			log.Printf("解析指令失败: %v", err)
			return
		}

		// 处理开锁指令
		if command, ok := cmd["command"].(string); ok && command == "open" {
			log.Printf("收到开锁指令: %v", cmd)
			// 这里可以添加实际控制快递柜开锁的代码
		}
	})
	if err != nil {
		log.Fatalf("订阅主题失败: %v", err)
	}

	// 3. 定时上报设备状态
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 模拟设备状态
		status := map[string]interface{}{
			"locker_id":   "device-001",       // 设备ID
			"is_online":   true,               // 在线状态
			"is_open":     rand.Intn(100) < 5, // 5%概率模拟门被打开
			"temperature": 20 + rand.Intn(15), // 20-35度随机温度
			"timestamp":   time.Now().Unix(),  // 时间戳
		}

		// 发布状态到EMQX
		if err := emqxClient.Publish(context.Background(), status); err != nil {
			log.Printf("上报状态失败: %v", err)
			continue
		}

		log.Printf("状态上报成功: %+v", status)
	}
}
