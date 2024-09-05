package database

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primarykey; autoIncrement" json:"id"`
	Nickname     string    `gorm:"not null;" json:"nickname"`
	Email        string    `gorm:"not null; unique;" json:"email"`
	HashPassword string    `gorm:"not null;" json:"password"`
	RegisteredAt time.Time `gorm:"autoCreateTime" json:"registered_at"`
	Posts        []Post    `gorm:"foreignKey:AuthorID;" json:"-"`
	Messages     []Message `gorm:"foreignKey:SenderID;" json:"-"`
	Chats        []Chat    `gorm:"many2many:user_x_chat;" json:"-"`
}

type Post struct {
	ID        uint      `gorm:"primarykey; autoIncrement" json:"id"`
	Content   string    `gorm:"not null;" json:"content"`
	AuthorID  uint      `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID       uint      `gorm:"primarykey; autoIncrement" json:"id"`
	Content  string    `gorm:"not null;" json:"content"`
	SenderID uint      `json:"sender_id"`
	ChatID   uint      `json:"chat_id"`
	SendedAt time.Time `gorm:"autoCreateTime" json:"sended_at"`
}

type Chat struct {
	ID        uint      `gorm:"primarykey; autoIncrement" json:"id"`
	Name      string    `gorm:"not null;" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Members   []User    `gorm:"many2many:user_x_chat;" json:"-"`
	Messages  []Message `gorm:"foreignKey:ChatID;" json:"-"`
}
