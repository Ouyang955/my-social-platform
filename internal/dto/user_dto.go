package dto

type UserDTO struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Bio         string `json:"bio"`          // 个性签名
	FollowCount int    `json:"follow_count"` // 关注数
	FansCount   int    `json:"fans_count"`   // 粉丝数
	LikeCount   int    `json:"like_count"`   // 获赞数
}
