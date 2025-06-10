package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"regexp"
	"time"
	"video_develop-main/video-api/api/request"
	__ "video_develop-main/video-api/basic/proto"
	"video_develop-main/video-api/sdk"
)

func Send(c *gin.Context) {
	var req request.SendUser
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "参数验证",
			"data": err.Error(),
		})
		return
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	C := __.NewUserClient(conn)

	code, err := C.SendSmsCode(c, &__.SendSmsCodeRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "发送失败",
			"data": nil,
		})
		return
	}
	// 验证手机号格式
	if !validatePhone(req.Mobile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号格式不正确"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "发送成功",
		"data": code,
	})

}

// 验证手机号格式
func validatePhone(mobile string) bool {
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	return matched
}

func Login(c *gin.Context) {
	var req request.LoginUser
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "参数验证",
			"data": err.Error(),
		})
		return
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	C := __.NewUserClient(conn)

	login, err := C.PhoneLogin(c, &__.PhoneLoginRequest{
		Mobile: req.Mobile,
		Code:   req.Code,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "登录失败",
			"data": nil,
		})
		return
	}

	token, err := sdk.NewJWT("2211a").CreateToken(sdk.CustomClaims{
		ID: uint(login.Id),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 86400,
		},
	})
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "登录成功",
		"data": token,
	})

}
