package main

import (
	"encoding/json"
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
	// 回复主题
	replyTopic := "oncard/response"

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

	// ====== 修复1：在消息发送前订阅，并检查订阅结果 ======
	fmt.Println("订阅回复主题:", replyTopic)
	token = client.Subscribe(replyTopic, 1, responseHandler)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("订阅失败: %v\n", token.Error())
		return
	}
	fmt.Println("订阅成功，可以发送消息")

	// 发送消息（使用JSON格式包含clientID）
	var clientMsg struct {
		ClientID string `json:"clientID"`
		Content  string `json:"content"`
	}
	clientMsg.ClientID = clientID
	clientMsg.Content = fmt.Sprintf("这是来自客户端的消息，时间: %s", time.Now().Format("15:04:05"))

	jsonMessage, _ := json.Marshal(clientMsg)

	// 发布消息到oncard主题，QoS为1
	token = client.Publish(topic, 1, false, jsonMessage)
	token.Wait()

	if token.Error() != nil {
		fmt.Printf("发送消息失败: %v\n", token.Error())
	} else {
		fmt.Printf("成功发送消息: %s\n", jsonMessage)
	}

	// ====== 修复2：添加明确的等待机制，确保能收到回复 ======
	fmt.Println("等待服务端回复... (10秒超时)")
	// 使用带超时的等待
	timer := time.NewTimer(10 * time.Second)
	select {
	case <-timer.C:
		fmt.Println("等待超时，未收到服务端回复")
	}

	fmt.Println("程序退出")
}

// 回复消息处理函数 - 解析JSON格式开锁指令
var responseHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("收到服务端回复 - 主题: %s\n", msg.Topic(), msg.Payload())

	// 解析JSON格式的开锁指令
	var unlockCmd struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}

	if err := json.Unmarshal(msg.Payload(), &unlockCmd); err != nil {
		fmt.Printf("解析开锁指令失败: %v\n", err)
		return
	}

	// 验证开锁状态并执行操作
	if unlockCmd.Status == "success" {
		fmt.Printf("[开闸操作] %s (时间: %s)\n", unlockCmd.Message, unlockCmd.Timestamp)
		fmt.Println("[系统提示] 闸机已开启，允许通行")
	} else {
		fmt.Printf("[开闸失败] %s\n", unlockCmd.Message)
	}

	// 通知主程序可以退出
	//wg.Done()
}