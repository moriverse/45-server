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
	Name         string
	PhoneNumber  string
	AvatarURL    string
	Source       Source
	OnboardedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastActiveAt *time.Time
	DeletedAt    *time.Time
}

// UpdateProfileParams contains the set of fields that can be updated on a user's profile.
type UpdateProfileParams struct {
	Name        *string
	PhoneNumber *string
	AvatarURL   *string
}

// Update modifies the user's profile based on the provided parameters.
func (u *User) Update(params UpdateProfileParams) {
	if params.Name != nil {
		u.Name = *params.Name
	}
	if params.PhoneNumber != nil {
		u.PhoneNumber = *params.PhoneNumber
	}
	if params.AvatarURL != nil {
		u.AvatarURL = *params.AvatarURL
	}
	u.UpdatedAt = time.Now()
}
