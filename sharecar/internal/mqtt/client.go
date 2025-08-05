package mqtt

import (
	"encoding/json"
	"log"
	"time"

	"sharecar/internal/model"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient MQTT客户端
type MQTTClient struct {
	client    paho.Client
	BrokerURL string
	ClientID  string
	Topic     string
	Qos       byte
	OnMessage func(event model.BikeEvent)
}

// NewMQTTClient 创建新的MQTT客户端
func NewMQTTClient(brokerURL, clientID, topic string, qos byte, onMessage func(event model.BikeEvent)) *MQTTClient {
	return &MQTTClient{
		BrokerURL: brokerURL,
		ClientID:  clientID,
		Topic:     topic,
		Qos:       qos,
		OnMessage: onMessage,
	}
}

// Connect 连接到MQTT服务器
func (m *MQTTClient) Connect() error {
	opts := paho.NewClientOptions()
	opts.AddBroker(m.BrokerURL)
	opts.SetClientID(m.ClientID)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(func(client paho.Client) {
		log.Printf("Connected to MQTT broker: %s\n", m.BrokerURL)
		if token := client.Subscribe(m.Topic, m.Qos, m.messageHandler); token.Wait() && token.Error() != nil {
			log.Printf("Subscribe error: %s\n", token.Error())
		}
	})

	m.client = paho.NewClient(opts)
	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// messageHandler 处理接收到的消息
func (m *MQTTClient) messageHandler(client paho.Client, msg paho.Message) {
	var event model.BikeEvent
	err := json.Unmarshal(msg.Payload(), &event)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v\n", err)
		return
	}

	log.Printf("Received message: %+v\n", event)
	if m.OnMessage != nil {
		m.OnMessage(event)
	}
}

// Publish 发布消息到指定主题
func (m *MQTTClient) Publish(topic string, qos byte, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	token := m.client.Publish(topic, qos, true, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// SendUnlockSignal 发送开锁信号
func (m *MQTTClient) SendUnlockSignal(deviceID string) error {
	//topic := fmt.Sprintf("bike/%s/command", deviceID)
	topic := "shanghai"
	command := map[string]string{
		"action":   "unlock",
		"deviceId": deviceID,
		"time":     time.Now().Format("20060102150405"),
		"status":   "设备已开启",
	}
	log.Printf("Sending unlock signal to topic %s: %+v\n", topic, command)
	return m.Publish(topic, m.Qos, command)
}

// SendLockSignal 发送关锁信号
func (m *MQTTClient) SendLockSignal(deviceID string) error {
	topic := "shanghai"
	command := model.BikeEvent{
		UserID:    "system",
		DeviceID:  deviceID,
		EventType: "lock",
		Time:      time.Now().Format("20060102150405"),
		Status:    "设备已关闭111",
		X:         0,
		Y:         0,
	}

	return m.Publish(topic, m.Qos, command)
}

// Disconnect 断开连接
func (m *MQTTClient) Disconnect() {
	if m.client != nil && m.client.IsConnected() {
		m.client.Disconnect(250)
	}
}
