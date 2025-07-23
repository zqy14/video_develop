package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func main() {
	// 添加用户名和密码认证信息（示例值，请替换为实际凭证）
	opts := mqtt.NewClientOptions().AddBroker("tcp://14.103.243.153:1883").
		//opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").
		SetUsername("zqy").
		SetPassword("zqy123456")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := "2212a"
	if token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		//todo 将数据入库
	}); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}

	fmt.Println("Subscribed to topic:", topic)
	time.Sleep(5 * time.Second) // Keep the program running for a while
	select {}
}
