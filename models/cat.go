package models

import "time"

type Cat struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Username    string    `gorm:"unique" json:"username"`
	Description string    `json:"description"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Likes       []Cat     `gorm:"many2many:cat_likes;" json:"likes,omitempty"`
	Dislikes    []Cat     `gorm:"many2many:cat_dislikes;" json:"dislikes,omitempty"`
	Avatar      string    `json:"avatar"`
}
