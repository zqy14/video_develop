package server

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"time"
	"video_develop-main/video-rpc/global"
	"video_develop-main/video-rpc/handler/model"
	__ "video_develop-main/video-rpc/proto"
	"video_develop-main/video-rpc/sdk"
)

type Server struct {
	__.UnimplementedUserServer
}

func (s *Server) Send(_ context.Context, in *__.SendSmsCodeRequest) (*__.SendSmsCodeResponse, error) {

	if !validatePhone(in.Mobile) {
		return &__.SendSmsCodeResponse{
			Success: false,
			Message: "手机号格式不正确",
		}, nil
	}

	return &__.SendSmsCodeResponse{
		Success: true,
		Message: "验证码已发送",
	}, nil
}

// validatePhone 验证手机号格式（国内）
func validatePhone(mobile string) bool {
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	return matched
}

// PhoneLogin 手机号登录/注册
func (s *Server) PhoneLogin(_ context.Context, in *__.PhoneLoginRequest) (*__.PhoneLoginResponse, error) {

	user := model.User{
		Mobile: in.Mobile,
	}

	err := global.DB.Create(&user).Error
	if err != nil {
		return nil, fmt.Errorf("用户未注册")
	}

	code := rand.Intn(9000) + 1000

	global.Red.Set(context.Background(), fmt.Sprintf("%s", in.Mobile), code, time.Minute*1)

	return &__.PhoneLoginResponse{
		Id:      user.Id,
		Message: "登录成功",
	}, nil
}

// 作品添加
func (s *Server) VideoWorks(_ context.Context, in *__.VideoWorksRequest) (*__.VideoWorksResponse, error) {

	works := model.VideoWorks{
		Title:          in.Title,
		Desc:           in.Desc,
		MusicId:        in.MusicId,
		WorkType:       in.WorkType,
		CheckStatus:    in.CheckStatus,
		CheckUser:      in.CheckUser,
		IpAddress:      in.IpAddress,
		WorkPermission: in.WorkPermission,
		LikeCount:      in.LikeCount,
		CommentCount:   in.CommentCount,
		ShareCount:     in.ShareCount,
		CollectCount:   in.CollectCount,
		BrowseCount:    in.BrowseCount,
	}

	err := global.DB.Create(&works).Error
	if err != nil {
		return nil, fmt.Errorf("作品发布失败")
	}

	return &__.VideoWorksResponse{
		Id:      works.Id,
		Message: "作品发布成功",
	}, nil
}

// 作品展示
func (s *Server) VideoList(_ context.Context, in *__.VideoListRequest) (*__.VideoListResponse, error) {
	var find []model.VideoWorks
	err := global.DB.Find(&find).Error
	if err != nil {
		return nil, fmt.Errorf("作品展示失败")
	}
	var finds []*__.ListWorks
	for _, v := range find {
		i := &__.ListWorks{
			Title:          v.Title,
			Desc:           v.Desc,
			MusicId:        v.MusicId,
			WorkType:       v.WorkType,
			CheckStatus:    v.CheckStatus,
			CheckUser:      v.CheckUser,
			IpAddress:      v.IpAddress,
			WorkPermission: v.WorkPermission,
			LikeCount:      v.LikeCount,
			CommentCount:   v.CommentCount,
			ShareCount:     v.ShareCount,
			CollectCount:   v.CollectCount,
			BrowseCount:    v.BrowseCount,
		}
		finds = append(finds, i)
	}

	return &__.VideoListResponse{
		List: finds,
	}, nil
}

// 评论更新
func (s *Server) UpdateVideoWorkComment(_ context.Context, in *__.UpdateVideoWorkCommentRequest) (*__.UpdateVideoWorkCommentResponse, error) {
	comment := model.VideoWorkComment{
		Id:      in.Id,
		WorkId:  in.WorkId,
		UserId:  in.UserId,
		Content: in.Content,
		Tag:     in.Tag,
		Pid:     in.Pid,
	}

	err := global.DB.Updates(&comment).Error
	if err != nil {
		return nil, fmt.Errorf("评论更新失败")
	}

	return &__.UpdateVideoWorkCommentResponse{
		Success: true,
		Message: "评论更新成功",
	}, nil
}

// 评论删除
func (s *Server) DeleteVideoWorkComment(_ context.Context, in *__.DeleteVideoWorkCommentRequest) (*__.DeleteVideoWorkCommentResponse, error) {
	comment := model.VideoWorkComment{
		Id: in.Id,
	}

	err := global.DB.Delete(&comment).Error
	if err != nil {
		return nil, fmt.Errorf("评论删除失败")
	}

	return &__.DeleteVideoWorkCommentResponse{
		Success: true,
		Message: "评论删除成功",
	}, nil
}

// 评论列表
func (s *Server) ListVideoWorkComment(_ context.Context, in *__.ListVideoWorkCommentRequest) (*__.ListVideoWorkCommentResponse, error) {
	var find []model.VideoWorkComment
	err := global.DB.Find(&find).Error
	if err != nil {
		return nil, fmt.Errorf("评论展示失败")
	}

	var finds []*__.ListComment
	for _, v := range find {
		i := &__.ListComment{
			WorkId:  v.WorkId,
			UserId:  v.UserId,
			Content: v.Content,
			Tag:     v.Tag,
			Pid:     v.Pid,
		}
		finds = append(finds, i)
	}

	return &__.ListVideoWorkCommentResponse{
		List: finds,
	}, nil
}

// 发布评论
func (s *Server) AddVideoWorkComment(_ context.Context, in *__.AddVideoWorkCommentRequest) (*__.AddVideoWorkCommentResponse, error) {

	comment := model.VideoWorkComment{
		Id:      in.WorkId,
		WorkId:  in.WorkId,
		UserId:  in.UserId,
		Content: in.Content,
		Tag:     in.Tag,
		Pid:     in.Pid,
	}

	err := global.DB.Create(&comment).Error
	if err != nil {
		return nil, fmt.Errorf("发布评论失败")
	}

	if !sdk.Content(in.Content) {
		return &__.AddVideoWorkCommentResponse{
			Success: false,
			Message: "内容含有敏感词",
		}, nil
	}

	return &__.AddVideoWorkCommentResponse{
		Success: true,
		Message: "发布评论成功",
	}, nil

}
