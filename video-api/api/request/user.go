package request

// 发验证码
type SendUser struct {
	Mobile string `form:"mobile" binding:"required"` // 手机号
}

// 登录
type LoginUser struct {
	Mobile string `form:"mobile" binding:"required"` // 手机号
	Code   string `form:"code" binding:"required"`
}

// 作品添加
type VideoWorks struct {
	Title          string ` form:"title" binding:"required"`           // 标题
	Desc           string ` form:"desc" binding:"required"`            // 描述
	MusicId        uint64 ` form:"music_id" binding:"required"`        // 选择音乐
	WorkType       string ` form:"work_type" binding:"required"`       // 作品类型
	CheckStatus    string ` form:"check_status" binding:"required"`    // 审核状态
	CheckUser      uint64 ` form:"check_user" binding:"required"`      // 审核人
	IpAddress      string ` form:"ip_address" binding:"required"`      // IP地址
	WorkPermission string ` form:"work_permission" binding:"required"` // 作品权限
	LikeCount      uint64 ` form:"like_count"`                         // 喜欢数量
	CommentCount   uint64 ` form:"comment_count" binding:"required"`   // 评论数
	ShareCount     uint64 ` form:"share_count" binding:"required"`     // 分享数
	CollectCount   uint64 ` form:"collect_count" binding:"required"`   // 收藏数
	BrowseCount    uint64 ` form:"browse_count" binding:"required"`    // 浏览量
}

// 作品展示列表
type ListVideoWorks struct {
	Page uint64 ` form:"page" binding:"required"`
	Size uint64 ` form:"size" binding:"required"`
}

// 评论修改
type UpdateVideoWorkComment struct {
	Id      uint64 `form:"id" binding:"required"`
	WorkId  uint64 `form:"work_id" binding:"required"` // 作品ID
	UserId  uint64 `form:"user_id" binding:"required"` // 用户ID
	Content string `form:"content" binding:"required"` // 评论内容
	Tag     uint64 `form:"tag" binding:"required"`     // 评论标签表
	Pid     uint64 `form:"pid" binding:"required"`     // 父级ID
}

// 评论发布
type AddVideoWorkComment struct {
	Id      uint64 `form:"id" binding:"required"`
	WorkId  uint64 `form:"work_id" binding:"required"` // 作品ID
	UserId  uint64 `form:"user_id" binding:"required"` // 用户ID
	Content string `form:"content" binding:"required"` // 评论内容
	Tag     uint64 `form:"tag" binding:"required"`     // 评论标签表
	Pid     uint64 `form:"pid" binding:"required"`     // 父级ID
}

// 评论删除
type DeleteVideoWorkComment struct {
	Id uint64 `form:"id" binding:"required"`
}

// 评论展示
type ListVideoWorkComment struct {
}

type AnthorCheck struct {
	UserId      uint64 `form:"user_id" binding:"required"`      // 用户
	CheckStatus string `form:"check_status" binding:"required"` // 审核状态
	Remark      string `form:"remark" binding:"required"`       // 备注
}
