package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// MQTT broker地址和端口
	broker := "broker.emqx.io"
	port := 1883
	// 要发布的主题
	topic := "oncard"
	// 客户端ID，使用时间戳确保唯一性
	clientID := fmt.Sprintf("mqtt-producer-%d", time.Now().UnixNano())

	// 创建连接选项
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)
	opts.SetConnectTimeout(5 * time.Second)

	// 创建客户端
	client := mqtt.NewClient(opts)
	token := client.Connect()
	
	// 等待连接完成
	if token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("无法连接到MQTT broker: %v", token.Error()))
	}
	defer client.Disconnect(250)

	fmt.Println("生产者已连接到MQTT broker，开始发送消息...")

	// 发送5条测试消息
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("这是来自客户端的第%d条消息，时间: %s", i, time.Now().Format("15:04:05"))
		// 发布消息到oncard主题，QoS为1
		token := client.Publish(topic, 1, false, message)
		token.Wait()

		if token.Error() != nil {
			fmt.Printf("发送消息失败: %v\n", token.Error())
		} else {
			fmt.Printf("成功发送消息: %s\n", message)
		}
		// 每条消息间隔1秒
		time.Sleep(1 * time.Second)
	}

	fmt.Println("所有消息发送完毕，程序退出")
}