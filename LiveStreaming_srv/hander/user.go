package hander

import (
	"LiveStreaming_srv/basic/global"
	"LiveStreaming_srv/model"
	__ "LiveStreaming_srv/proto"
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

type UserServer struct {
	__.UnimplementedUserServer
}

func isValidMobile(mobile string) bool {
	pattern := `^1[3-9]\d{9}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(mobile)
}

func (c *UserServer) SendSms(_ context.Context, in *__.SendSmsRequest) (*__.SendSmsResponse, error) {
	if !isValidMobile(in.Mobile) {
		return nil, fmt.Errorf("手机号格式不正确")
	}

	var userCount model.VideoUser
	global.DB.Where("mobile = ?", in.Mobile).Find(&userCount)

	key := "SendSms" + in.Mobile
	existingCode := global.Rdb.Get(context.Background(), key)
	if existingCode.Err() == nil {
		return nil, fmt.Errorf("验证码已发送，请稍后再试")
	}

	code := rand.Intn(9000) + 1000

	sendCountKey := "SendSmsCount_" + in.Mobile
	count, err := global.Rdb.Incr(context.Background(), sendCountKey).Result()
	if err == nil && count > 5 {
		return nil, fmt.Errorf("验证码发送过于频繁，请明天再试")
	}
	global.Rdb.Expire(context.Background(), sendCountKey, time.Minute*5)

	global.Rdb.Set(context.Background(), key, code, time.Minute)

	return &__.SendSmsResponse{}, nil
}

func (c *UserServer) Login(_ context.Context, in *__.LoginRequest) (*__.LoginResponse, error) {

	if !isValidMobile(in.Mobile) {
		return nil, fmt.Errorf("手机号格式不正确")
	}

	if in.SendSmsCode == "" {
		return nil, fmt.Errorf("验证码不能为空")
	}

	var user model.VideoUser
	global.DB.Where("mobile=?", in.Mobile).Find(&user)

	if user.Id == 0 {
		var count int64
		global.DB.Where("name = ?", "用户"+in.Mobile).Count(&count)
		if count > 0 {
			return nil, fmt.Errorf("用户名已存在")
		}

		newUser := model.VideoUser{
			Name:   "用户" + in.Mobile,
			Status: strconv.Itoa(1),
			Mobile: in.Mobile,
		}

		result := global.DB.Create(&newUser)
		if result.Error != nil {
			return nil, fmt.Errorf("注册失败")
		}
		user = newUser
	}

	key := "SendSms" + in.Mobile
	get := global.Rdb.Get(context.Background(), key)
	if get.Err() != nil {
		return nil, fmt.Errorf("验证码已过期")
	}

	if get.Val() != in.SendSmsCode {
		return nil, fmt.Errorf("验证码错误")
	}

	global.Rdb.Del(context.Background(), key)

	global.Rdb.Del(context.Background(), "LoginError_"+in.Mobile)

	return &__.LoginResponse{
		Id: int64(user.Id),
	}, nil
}

func (c *UserServer) Personal(_ context.Context, in *__.PersonalRequest) (*__.PersonalResponse, error) {

	var user model.VideoUser
	result := global.DB.First(&user, in.Id)

	if result.Error != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	resp := &__.PersonalResponse{
		Name:          user.Name,
		NickName:      user.NickName,
		UserCode:      user.UserCode,
		Signature:     user.Signature,
		Sex:           user.Sex,
		IpAddress:     "",
		Constellation: user.Constellation,
		AttendCount:   user.AttendCount,
		FansCount:     float32(user.FansCount),
		ZanCount:      float32(user.ZanCount),
		AvatorFileId:  int64(user.AvatorFileId),
		AuthriryInfo:  user.AuthriryInfo,
		Mobile:        "",
		RealNameAuth:  user.RealNameAuth,
		Age:           0,
		OnlineStatus:  strconv.Itoa(0),
		AuthrityType:  user.AuthrityType,
		Level:         int64(user.Level),
		Balance:       0,
	}

	if user.Id == int(in.Id) {
		resp.IpAddress = user.IpAddress
		resp.Mobile = user.Mobile
		resp.Age = user.Age
		resp.OnlineStatus = user.OnlineStatus
		resp.Balance = int64(user.Balance)
	}

	return resp, nil

}

func (c *UserServer) UpdatePersonal(_ context.Context, in *__.UpdatePersonalRequest) (*__.UpdatePersonalResponse, error) {

	if in.Id <= 0 {
		return &__.UpdatePersonalResponse{
			Code:    500,
			Message: "无效的用户ID",
		}, nil
	}

	var user model.VideoUser
	if err := global.DB.First(&user, in.Id).Error; err != nil {
		return &__.UpdatePersonalResponse{
			Code:    500,
			Message: "用户不存在",
		}, nil
	}
	users := model.VideoUser{
		Id:            int(in.Id),
		Name:          in.Name,
		NickName:      in.NickName,
		Signature:     in.Signature,
		Sex:           in.Sex,
		Constellation: in.Constellation,
		AvatorFileId:  int(in.AvatorFileId),
		Mobile:        in.Mobile,
		Age:           in.Age,
		OnlineStatus:  in.OnlineStatus,
	}
	if err := global.DB.Updates(&users).Error; err != nil {
		return &__.UpdatePersonalResponse{
			Code:    500,
			Message: "编辑个人信息失败: " + err.Error(),
		}, nil
	}

	return &__.UpdatePersonalResponse{
		Code:    200,
		Message: "编辑个人信息成功",
	}, nil

}
