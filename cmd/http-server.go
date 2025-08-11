package main

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	__ "rider-device/initenal/proto"
)

// 添加骑手请求结构体
type AddRiderRequest struct {
	Name  string `json:"name" binding:"required"`  // 姓名，必填
	Phone string `json:"phone" binding:"required"` // 手机号，必填
}

func main() {
	// 3. 初始化Gin路由
	router := gin.Default()

	// 1. 连接gRPC服务
	conn, err := grpc.Dial("localhost:50091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接gRPC服务失败: %v", err)
	}
	defer conn.Close()

	// 2. 创建gRPC客户端
	riderClient := __.NewRiderServiceClient(conn)

	// 4. 添加骑手路由
	router.POST("/api/riders", func(c *gin.Context) {
		// 4.1 绑定请求参数
		var req AddRiderRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "参数错误: " + err.Error()})
			return
		}

		// 4.2 构造gRPC请求
		grpcReq := &__.AddRiderRequest{
			Name:  req.Name,
			Phone: req.Phone,
		}

		// 4.3 调用gRPC服务
		resp, err := riderClient.AddRider(c, grpcReq)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		// 4.4 返回响应
		c.JSON(http.StatusOK, gin.H{
			"success":  resp.Success,
			"message":  resp.Message,
			"rider_id": resp.RiderId,
		})
	})

	// 5. 启动HTTP服务
	if err := router.Run(":8084"); err != nil {
		log.Fatalf("HTTP服务启动失败: %v", err)
	}
}
