package request

type PostComment struct {
	WorkId  int64  `form:"work_id" binding:"required"`
	UserId  int64  `form:"user_id" binding:"required"`
	Content string `form:"content"`
	Tag     int64  `form:"tag" binding:"required"`
	Pid     int64  `form:"pid" binding:"required"`
}
