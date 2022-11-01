package userschema

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Content    string         `json:"content"`
	SenderID   uint           `json:"senderID"`
	ReceiverID uint           `json:"receiverID"`
}
