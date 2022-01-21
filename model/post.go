package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	UserId     uint      `json:"user_id" gorm:"not null"`
	CategoryId uint      `json:"category_id" gorm:"not null"`
	Category   *Category `gorm:"foreignKey:CategoryId"` //外键
	Title      string    `json:"title" gorm:"type:varchar(50);not null"'`
	HeadImg    string    `json:"head_img"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	CreatedAt  Time      `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt  Time      `json:"updated_at" gorm:"type:timestamp"`
}

func (p *Post) BeforeCreate(db *gorm.DB) error {
	p.ID = uuid.NewV4()
	return nil
}
