package main

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// DeviceStatus 设备状态结构体
type DeviceStatus struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Docker string `json:"docker"`
	Count  string `json:"count"`
	Status string `json:"status"` // 新增状态字段
	Time   string `json:"time"`   // 新增时间戳字段
}

func main() {
	// 1. 创建MQTT客户端并连接
	opts := mqtt.NewClientOptions().AddBroker("tcp://14.103.243.153:1883")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("成功连接到MQTT服务器")

	// 2. 初始化设备信息
	device := DeviceStatus{
		ID:     "1",
		Name:   "zqy",
		Docker: "林医生",
		Count:  "300",
	}

	// 3. 主循环 - 持续更新状态
	for {
		// 更新状态和时间戳
		if time.Now().Second()%2 == 0 {
			device.Status = "正常"
		} else {
			device.Status = "异常"
		}
		device.Time = time.Now().Format("2006-01-02 15:04:05")

		// 序列化为JSON
		payload, err := json.Marshal(device)
		if err != nil {
			fmt.Println("JSON序列化错误:", err)
			continue
		}

		// 发布到主题"bird"
		token := client.Publish("bird", 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			fmt.Println("发布失败:", token.Error())
		} else {
			fmt.Printf("发布成功: %s\n", string(payload))
		}

		// 每秒更新一次
		time.Sleep(1 * time.Second)
	}
}
