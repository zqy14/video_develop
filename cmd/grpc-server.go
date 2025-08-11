package main

import (
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"rider-management/internal/pb"
	"rider-management/internal/repository"
	"rider-management/internal/service"
	"rider-management/pkg/notification"
)

func main() {
	// 1. 初始化数据库连接
	dsn := "root:zqy123456@tcp(14.103.243.153:3306)/device?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 3. 初始化依赖组件
	riderRepo := repository.NewRiderRepository(db)               // 骑手仓库
	notifier := &notification.SMSNotifier{}                      // 通知服务
	riderService := service.NewRiderService(riderRepo, notifier) // 骑手服务

	// 4. 创建gRPC服务器
	grpcServer := grpc.NewServer()
	pb.RegisterRiderServiceServer(grpcServer, riderService)

	// 5. 启动gRPC5. 启动gRPC服务
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}

	log.Printf("gRPC服务启动，监听地址: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC服务启动失败: %v", err)
	}
}
