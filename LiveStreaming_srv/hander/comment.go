package hander

import (
	"LiveStreaming_srv/basic/global"
	"LiveStreaming_srv/model"
	__ "LiveStreaming_srv/proto"
	"context"
	"gorm.io/gorm"
)

func (c *UserServer) PostComment(_ context.Context, in *__.PostCommentRequest) (*__.PostCommentResponse, error) {

	var user model.VideoUser
	if err := global.DB.Where("id = ?", in.UserId).First(&user).Error; err != nil {
		return &__.PostCommentResponse{
			Success: false,
			Message: "用户不存在",
		}, nil
	}

	var work model.VideoWorks
	if err := global.DB.Where("id = ?", in.WorkId).First(&work).Error; err != nil {
		return &__.PostCommentResponse{
			Success: false,
			Message: "作品不存在",
		}, nil
	}

	if in.Content == "" {
		return &__.PostCommentResponse{
			Success: false,
			Message: "评论内容不能为空",
		}, nil
	}

	works := model.VideoWorkComment{
		WorkId:  int(in.WorkId),
		UserId:  int(in.UserId),
		Content: in.Content,
		Tag:     int(in.Tag),
		Pid:     int(in.Pid),
	}

	if err := global.DB.Model(&work).Where("id = ?", in.WorkId).
		//UpdateColumn直接执行 SQL 更新，不触发 GORM 的钩子函数（如 BeforeUpdate），适合纯计数更新
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
		return &__.PostCommentResponse{
			Success: false,
			Message: "更新评论数失败: " + err.Error(),
		}, nil
	}

	if err := global.DB.Create(&works).Error; err != nil {
		return &__.PostCommentResponse{
			Success: false,
			Message: "发布作品失败",
		}, nil
	}

	return &__.PostCommentResponse{
		Success: true,
		Message: "发布评论成功",
	}, nil

}

func (c *UserServer) LikeWork(_ context.Context, in *__.LikeWorkRequest) (*__.LikeWorkResponse, error) {
	var user model.VideoUser
	if err := global.DB.Where("id = ?", in.UserId).First(&user).Error; err != nil {
		return &__.LikeWorkResponse{
			Code:    500,
			Message: "用户不存在",
		}, nil
	}

	var work model.VideoWorks
	if err := global.DB.Where("id = ?", in.WorkId).First(&work).Error; err != nil {
		return &__.LikeWorkResponse{
			Code:    500,
			Message: "作品不存在",
		}, nil
	}

	if err := global.DB.Model(&work).Where("id = ?", in.WorkId).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		return &__.LikeWorkResponse{
			Code:    500,
			Message: "点赞失败: " + err.Error(),
		}, nil
	}

	return &__.LikeWorkResponse{
		Code:    200,
		Message: "作品点赞成功",
	}, nil
}
