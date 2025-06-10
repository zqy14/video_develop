package handler

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"video_develop-main/video-api/api/request"
	__ "video_develop-main/video-api/basic/proto"
)

// 作品添加
func VideoAdd(c *gin.Context) {
	var req request.VideoWorks
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

	works, err := C.VideoWorks(c, &__.VideoWorksRequest{
		Title:          req.Title,
		Desc:           req.Desc,
		MusicId:        req.MusicId,
		WorkType:       req.WorkType,
		CheckStatus:    req.CheckStatus,
		CheckUser:      req.CheckUser,
		IpAddress:      req.IpAddress,
		WorkPermission: req.WorkPermission,
		LikeCount:      req.LikeCount,
		CommentCount:   req.CommentCount,
		ShareCount:     req.ShareCount,
		CollectCount:   req.CollectCount,
		BrowseCount:    req.BrowseCount,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "作品添加失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "作品添加成功",
		"data": works,
	})

}

// 作品展示及分页
func VideoList(c *gin.Context) {
	var req request.ListVideoWorks
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

	list, err := C.VideoList(c, &__.VideoListRequest{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "作品展示失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "作品展示成功",
		"data": list,
	})
}

// 评论更新
func CommentUpdate(c *gin.Context) {
	var req request.UpdateVideoWorkComment
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

	comment, err := C.UpdateVideoWorkComment(c, &__.UpdateVideoWorkCommentRequest{
		Id:      req.Id,
		WorkId:  req.WorkId,
		UserId:  req.UserId,
		Content: req.Content,
		Tag:     req.Tag,
		Pid:     req.Pid,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "评论更新失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "评论更新成功",
		"data": comment,
	})
}

// 评论删除
func CommentDelete(c *gin.Context) {
	var req request.DeleteVideoWorkComment
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

	comment, err := C.DeleteVideoWorkComment(c, &__.DeleteVideoWorkCommentRequest{
		Id: req.Id,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "评论删除失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "评论删除成功",
		"data": comment,
	})
}

// 评论列表
func CommentList(c *gin.Context) {
	var req request.ListVideoWorkComment
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

	comment, err := C.ListVideoWorkComment(c, &__.ListVideoWorkCommentRequest{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "评论展示失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "评论展示成功",
		"data": comment,
	})
}

func CommentAdd(c *gin.Context) {
	var req request.AddVideoWorkComment
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

	comment, err := C.AddVideoWorkComment(c, &__.AddVideoWorkCommentRequest{
		WorkId:  req.WorkId,
		UserId:  req.UserId,
		Content: req.Content,
		Tag:     req.Tag,
		Pid:     req.Pid,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "发布评论失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "发布评论成功",
		"data": comment,
	})
}

func AnthorCheck(c *gin.Context) {
	var req request.AnthorCheck
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

	check, err := C.AnthorCheck(c, &__.AnthorCheckRequest{
		UserId:      c.GetUint64("id"),
		CheckStatus: req.CheckStatus,
		Remark:      req.Remark,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "审核失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "审核成功",
		"data": check,
	})

}
