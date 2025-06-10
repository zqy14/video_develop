package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        uint64         `gorm:"column:id;type:bigint(20) UNSIGNED;primaryKey;not null;" json:"id"`
	Mobile    string         `gorm:"column:mobile;type:char(11);comment:手机号;not null;" json:"mobile"` // 手机号
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);default:NULL;" json:"deleted_at"`
}

func (u *User) TableName() string {
	return "user"
}

type VideoWorks struct {
	Id             uint64 `gorm:"column:id;type:int(11);primaryKey;not null;" json:"id"`
	Title          string `gorm:"column:title;type:varchar(100);comment:标题;default:NULL;" json:"title"`                      // 标题
	Desc           string `gorm:"column:desc;type:varchar(255);comment:描述;default:NULL;" json:"desc"`                        // 描述
	MusicId        uint64 `gorm:"column:music_id;type:int(11);comment:选择音乐;default:NULL;" json:"music_id"`                   // 选择音乐
	WorkType       string `gorm:"column:work_type;type:varchar(20);comment:作品类型;default:NULL;" json:"work_type"`             // 作品类型
	CheckStatus    string `gorm:"column:check_status;type:varchar(10);comment:审核状态;default:1;" json:"check_status"`          // 审核状态
	CheckUser      uint64 `gorm:"column:check_user;type:int(11);comment:审核人;default:NULL;" json:"check_user"`                // 审核人
	IpAddress      string `gorm:"column:ip_address;type:varchar(20);comment:IP地址;default:NULL;" json:"ip_address"`           // IP地址
	WorkPermission string `gorm:"column:work_permission;type:varchar(20);comment:作品权限;default:NULL;" json:"work_permission"` // 作品权限
	LikeCount      uint64 `gorm:"column:like_count;type:int(11);comment:喜欢数量;default:NULL;" json:"like_count"`               // 喜欢数量
	CommentCount   uint64 `gorm:"column:comment_count;type:int(11);comment:评论数;default:NULL;" json:"comment_count"`          // 评论数
	ShareCount     uint64 `gorm:"column:share_count;type:int(11);comment:分享数;default:NULL;" json:"share_count"`              // 分享数
	CollectCount   uint64 `gorm:"column:collect_count;type:int(11);comment:收藏数;default:NULL;" json:"collect_count"`          // 收藏数
	BrowseCount    uint64 `gorm:"column:browse_count;type:int(11);comment:浏览量;default:NULL;" json:"browse_count"`            // 浏览量
}

func (v *VideoWorks) TableName() string {
	return "video_works"
}

// 评论表
type VideoWorkComment struct {
	Id      uint64 `gorm:"column:id;type:int(11);primaryKey;not null;" json:"id"`
	WorkId  uint64 `gorm:"column:work_id;type:int(11);comment:作品ID;default:NULL;" json:"work_id"`      // 作品ID
	UserId  uint64 `gorm:"column:user_id;type:int(11);comment:用户ID;default:NULL;" json:"user_id"`      // 用户ID
	Content string `gorm:"column:content;type:varchar(100);comment:评论内容;default:NULL;" json:"content"` // 评论内容
	Tag     uint64 `gorm:"column:tag;type:int(11);comment:评论标签表;default:NULL;" json:"tag"`             // 评论标签表
	Pid     uint64 `gorm:"column:pid;type:int(11);comment:父级ID;default:0;" json:"pid"`                 // 父级ID
}

func (v *VideoWorkComment) TableName() string {
	return "video_work_comment"
}

type VideoAnthorCheck struct {
	Id          uint64 `gorm:"column:id;type:int(11);primaryKey;not null;" json:"id"`
	UserId      uint64 `gorm:"column:user_id;type:int(11);comment:用户;primaryKey;not null;" json:"user_id"`          // 用户
	CheckStatus string `gorm:"column:check_status;type:varchar(20);comment:审核状态;default:NULL;" json:"check_status"` // 审核状态
	Remark      string `gorm:"column:remark;type:varchar(255);comment:备注;default:NULL;" json:"remark"`              // 备注
}

func (v *VideoAnthorCheck) TableName() string {
	return "video_anthor_check"
}
