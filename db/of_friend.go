package db

import (
	"time"
)

// 好友配置表
type OFFriendInfo struct {
	UserId uint32 `gorm:"primary_key;not null"` // 用户id
}

// 好友申请表
type OFFriendRequest struct {
	SenderUserId  uint32 `gorm:"primary_key;not null;index:request"` // 申请者
	RequestUserId uint32 `gorm:"primary_key;not null;index:request"` // 被申请者
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// 好友关系表
type OFFriend struct {
	UserId    uint32 `gorm:"primary_key;not null;index:friend"` // 用户id
	FriendId  uint32 `gorm:"primary_key;not null;index:friend"` // 好友id
	CreatedAt time.Time
	UpdatedAt time.Time
}
