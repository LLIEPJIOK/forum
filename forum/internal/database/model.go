package database

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primarykey; autoIncrement"`
	Nickname     string    `gorm:"not null;"`
	Email        string    `gorm:"not null; unique;"`
	HashPassword string    `gorm:"not null;" json:"Password"`
	RegisteredAt time.Time `gorm:"autoCreateTime"`
	Posts        []Post    `gorm:"foreignKey:AuthorID;" json:"-"`
	Messages     []Message `gorm:"foreignKey:SenderID;" json:"-"`
	Chats        []Chat    `gorm:"many2many:user_x_chat;" json:"-"`
}

type Post struct {
	ID        uint   `gorm:"primarykey; autoIncrement"`
	Content   string `gorm:"not null;"`
	AuthorID  uint
	CreatedAt time.Time
}

type Message struct {
	ID       uint   `gorm:"primarykey; autoIncrement"`
	Content  string `gorm:"not null;"`
	SenderID uint
	ChatID   uint
	SendedAt time.Time `gorm:"autoCreateTime"`
}

type Chat struct {
	ID        uint   `gorm:"primarykey; autoIncrement"`
	Name      string `gorm:"not null;"`
	CreatedAt time.Time
	Members   []User    `gorm:"many2many:user_x_chat;" json:"-"`
	Messages  []Message `gorm:"foreignKey:ChatID;" json:"-"`
}
