package user

import (
	"time"
)

type UserID string

type Source string

const (
	WechatIOS     Source = "wechat_ios"
	WechatAndroid Source = "wechat_android"
	IOS           Source = "ios"
	Android       Source = "android"
	Web           Source = "web"
)

type User struct {
	ID           UserID
	Email        string
	PhoneNumber  string
	AvatarURL    string
	Source       Source
	OnboardedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastActiveAt *time.Time
	DeletedAt    *time.Time
}
