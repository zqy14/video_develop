package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"video_develop-main/video-rpc/configs"

	//_ "video_develop-main/video-rpc/basic/init"
	"video_develop-main/video-rpc/handler/server"
	__ "video_develop-main/video-rpc/proto"
	"video_develop-main/video-rpc/utils"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	configs.Viper()
	utils.GlobalMysql()
	utils.ExampleClient()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	__.RegisterUserServer(s, &server.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
