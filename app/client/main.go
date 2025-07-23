package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://14.103.243.153:1883")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := "2212a"
	payload := []byte("{'id':'1234567890','name':'zqy','docker':'林医生','count':'300'}")
	if token := client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}

	fmt.Println("Message published")
	time.Sleep(1 * time.Second) // Wait before exiting
}
