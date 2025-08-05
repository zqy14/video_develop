package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 消息处理函数，打印接收到的消息
var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("收到消息 - 主题: %s, 内容: %s\n", msg.Topic(), msg.Payload())
}

// 连接成功回调函数
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("成功连接到MQTT broker")
	// 连接成功后订阅oncard主题
	token := client.Subscribe("oncard", 1, nil)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("订阅主题失败: %v\n", token.Error())
	} else {
		fmt.Println("已订阅主题: oncard")
	}
}

// 连接丢失回调函数
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("连接丢失: %v，正在尝试重连...\n", err)
}

func main() {
	// MQTT broker地址和端口
	var broker = "broker.emqx.io"
	var port = 1883
	// 客户端ID，使用时间戳确保唯一性
	clientID := fmt.Sprintf("mqtt-consumer-%d", time.Now().UnixNano())

	// 创建连接选项
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)
	opts.SetDefaultPublishHandler(messageHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)

	// 创建客户端
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("连接失败: %v", token.Error()))
	}

	fmt.Println("消费者服务已启动，等待接收消息... (按Ctrl+C退出)")

	// 保持程序运行
	select {}
}
