package database

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nickname     string    `gorm:"not null;"`
	Email        string    `gorm:"not null; unique;"`
	HashPassword string    `gorm:"not null;"`
	Posts        []Post    `gorm:"foreignKey:AuthorID;"`
	Messages     []Message `gorm:"foreignKey:SenderID;"`
	Chats        []Chat    `gorm:"many2many:user_x_chat;"`
}

type Post struct {
	gorm.Model
	Content  string `gorm:"not null;"`
	AuthorID uint
}

type Message struct {
	gorm.Model
	Content  string `gorm:"not null;"`
	SenderID uint
	ChatID   uint
}

type Chat struct {
	gorm.Model
	Name     string    `gorm:"not null;"`
	Members  []User    `gorm:"many2many:user_x_chat;"`
	Messages []Message `gorm:"foreignKey:ChatID;"`
}
