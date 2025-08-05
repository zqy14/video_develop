package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"sharecar/internal/model"
	"sharecar/internal/mqtt"
	"sharecar/internal/service"

	"github.com/spf13/cobra"
)

var (
	brokerURL string
	clientID  string
	topic     string
	qos       int
)

var rootCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consume messages from EMQX broker",
	Long:  "A command-line tool to consume messages from EMQX broker for shared bike system",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("------>连接到MQTT服务器 %s\n", brokerURL)
		// 创建单车服务
		bikeService := setupBikeService()

		// 等待中断信号
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		// 打印所有骑行记录
		rides := bikeService.GetAllRides()
		fmt.Printf("\nTotal rides: %d\n", len(rides))
		for i, ride := range rides {
			fmt.Printf("Ride %d: User %s, Bike %s, Duration: %.2f minutes, Fee: %.2f yuan\n",
				i+1, ride.UserID, ride.DeviceID, ride.Duration.Minutes(), ride.Fee)
		}
	},
}

// 新增client子命令
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Interactive MQTT client for sending and receiving messages",
	Long:  "An interactive MQTT client that can send and receive messages on specified topics",
	Run: func(cmd *cobra.Command, args []string) {
		// 创建MQTT客户端
		mqttClient := mqtt.NewMQTTClient(
			brokerURL,
			clientID,
			topic,
			byte(qos),
			func(event model.BikeEvent) {
				// 检查是否有状态信息
				statusMsg := ""
				if event.Status != "" {
					statusMsg = fmt.Sprintf("状态: %s", event.Status)
				}
				fmt.Printf("收到消息: UserID=%s, DeviceID=%s, EventType=%s, 位置=(%.6f, %.6f) %s\n",
					event.UserID, event.DeviceID, event.EventType, event.X, event.Y, statusMsg)
			},
		)

		// 连接MQTT服务器
		err := mqttClient.Connect()
		if err != nil {
			log.Fatalf("连接MQTT服务器失败: %v\n", err)
		}

		fmt.Printf("已连接到MQTT服务器 %s 并订阅主题 %s\n", brokerURL, topic)
		fmt.Println("输入消息格式: <用户ID> <设备ID> <事件类型> <X坐标> <Y坐标>")
		fmt.Println("例如: user123 bike456 unlock 121.456789 31.123456")
		fmt.Println("输入 'exit' 退出")

		// 创建一个扫描器读取用户输入
		scanner := bufio.NewScanner(os.Stdin)

		// 设置中断信号处理
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// 使用通道来处理用户输入和信号
		go func() {
			for scanner.Scan() {
				input := scanner.Text()
				if input == "exit" {
					os.Exit(0)
				}

				// 解析用户输入
				parts := strings.Fields(input)
				if len(parts) != 5 {
					fmt.Println("输入格式错误，请按照格式输入: <用户ID> <设备ID> <事件类型> <X坐标> <Y坐标>")
					continue
				}

				// 创建BikeEvent对象
				event := model.BikeEvent{
					UserID:    parts[0],
					DeviceID:  parts[1],
					EventType: parts[2],
				}

				// 解析坐标
				fmt.Sscanf(parts[3], "%f", &event.X)
				fmt.Sscanf(parts[4], "%f", &event.Y)

				// 发布消息
				err := mqttClient.Publish(topic, byte(qos), event)
				if err != nil {
					fmt.Printf("发送消息失败: %v\n", err)
				} else {
					fmt.Printf("消息已发送到主题 %s\n", topic)
				}
			}
		}()

		// 等待中断信号
		<-sigChan
		fmt.Println("\n正在退出...")
	},
}

func setupBikeService() *service.BikeService {
	// 创建MQTT客户端
	mqttClient := mqtt.NewMQTTClient(
		brokerURL,
		clientID,
		topic,
		byte(qos),
		nil, // 暂时不设置回调
	)

	// 创建单车服务
	bikeService := service.NewBikeService(mqttClient)

	// 设置消息处理回调
	mqttClient.OnMessage = func(event model.BikeEvent) {
		bikeService.HandleBikeEvent(event)
	}

	// 连接MQTT服务器
	err := mqttClient.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v\n", err)
	}

	return bikeService
}

func init() {
	rootCmd.PersistentFlags().StringVar(&brokerURL, "broker", "tcp://localhost:1883", "MQTT broker URL")
	rootCmd.PersistentFlags().StringVar(&clientID, "client-id", "sharecar-consumer", "MQTT client ID")
	rootCmd.PersistentFlags().StringVar(&topic, "topic", "shanghai", "MQTT topic to subscribe")
	rootCmd.PersistentFlags().IntVar(&qos, "qos", 1, "MQTT QoS level")

	// 添加client子命令
	rootCmd.AddCommand(clientCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
