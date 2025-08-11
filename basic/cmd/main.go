package main

import (
	"devicemanage/devicerpc/basic/config"
	_ "devicemanage/devicerpc/basic/init"
	__ "devicemanage/devicerpc/courier"
	"devicemanage/devicerpc/handler/server"
	"devicemanage/devicerpc/utils"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", config.GloBalConfig.Server.Port, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	__.RegisterRiderServiceServer(s, &service.Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// 获取当天日期
	today := time.Now().Format("2006-01-02")

	// a) 当天配送数统计
	deliveryCount := utils.CountDeliveries(today)
	fmt.Printf("当日配送总数: %d\n", deliveryCount)

	// b) 快递员分级统计
	courierStats := utils.CountCouriers()
	fmt.Println("\n快递员分级统计:")
	for level, percent := range courierStats {
		fmt.Printf("%s级: %.1f%%\n", level, percent)
	}

	// c) 单日开箱数统计
	openBoxCount := utils.CountDailyOpenBoxes(today)
	fmt.Printf("\n单日开箱数(去重): %d\n", openBoxCount)
}
