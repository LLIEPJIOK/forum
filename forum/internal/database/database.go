package database

import (
	"fmt"

	"gorm.io/gorm"
)

type Database struct {
	gormDB *gorm.DB
}

func New(gormDB *gorm.DB) *Database {
	return &Database{
		gormDB: gormDB,
	}
}

func (db *Database) CreateTables() error {
	err := db.gormDB.AutoMigrate(User{}, Post{}, Message{}, Chat{})
	if err != nil {
		return fmt.Errorf("cannot create tables: %w", err)
	}

	return nil
}

func (db *Database) AddUser(user *User) error {
	result := db.gormDB.Create(user)
	if result.Error != nil {
		return fmt.Errorf("cannot add user %#v to db: %w", user, result.Error)
	}

	return nil
}

func (db *Database) GetUserByEmail(email string) (*User, error) {
	var user *User
	result := db.gormDB.Where("email = ?", email).First(user)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get user by email = %q: %w", email, result.Error)
	}

	return user, nil
}

func (db *Database) GetUserByID(id uint) (*User, error) {
	var user *User
	result := db.gormDB.Where("id = ?", id).First(user)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get user by id = %d: %w", id, result.Error)
	}

	return user, nil
}

func (db *Database) UpdateUser(user *User) error {
	result := db.gormDB.Save(user)

	if result.Error != nil {
		return fmt.Errorf("cannot update user %#v: %w", user, result.Error)
	}

	return nil
}

func (db *Database) DeleteUser(id uint) error {
	result := db.gormDB.Delete(&User{}, id)
	if result.Error != nil {
		return fmt.Errorf("cannot delete user with id = %d: %w", id, result.Error)
	}

	return nil
}

func (db *Database) AddPost(post *Post) error {
	result := db.gormDB.Create(post)
	if result.Error != nil {
		return fmt.Errorf("cannot add post %#v to db: %w", post, result.Error)
	}

	return nil
}

func (db *Database) GetPost(id uint) (*Post, error) {
	var post *Post
	result := db.gormDB.Where("id = ?", id).First(post)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get post by id = %d: %w", id, result.Error)
	}

	return post, nil
}

func (db *Database) UpdatePost(post *Post) error {
	result := db.gormDB.Save(post)

	if result.Error != nil {
		return fmt.Errorf("cannot update post %#v: %w", post, result.Error)
	}

	return nil
}

func (db *Database) DeletePost(id uint) error {
	result := db.gormDB.Delete(&Post{}, id)
	if result.Error != nil {
		return fmt.Errorf("cannot delete post with id = %d: %w", id, result.Error)
	}

	return nil
}

func (db *Database) AddMessage(message *Message) error {
	result := db.gormDB.Create(message)
	if result.Error != nil {
		return fmt.Errorf("cannot add message %#v to db: %w", message, result.Error)
	}

	return nil
}

func (db *Database) GetMessage(id uint) (*Message, error) {
	var message *Message
	result := db.gormDB.Where("id = ?", id).First(message)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get message by id = %d: %w", id, result.Error)
	}

	return message, nil
}

func (db *Database) UpdateMessage(message *Message) error {
	result := db.gormDB.Save(message)

	if result.Error != nil {
		return fmt.Errorf("cannot update message %#v: %w", message, result.Error)
	}

	return nil
}

func (db *Database) DeleteMessage(id uint) error {
	result := db.gormDB.Delete(&Message{}, id)
	if result.Error != nil {
		return fmt.Errorf("cannot delete message with id = %d: %w", id, result.Error)
	}

	return nil
}

func (db *Database) AddChat(chat *Chat) error {
	result := db.gormDB.Create(chat)
	if result.Error != nil {
		return fmt.Errorf("cannot add chat %#v to db: %w", chat, result.Error)
	}

	return nil
}

func (db *Database) GetChat(id uint) (*Chat, error) {
	var chat *Chat
	result := db.gormDB.Where("id = ?", id).First(chat)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get chat by id = %d: %w", id, result.Error)
	}

	return chat, nil
}

func (db *Database) UpdateChat(chat *Chat) error {
	result := db.gormDB.Save(chat)

	if result.Error != nil {
		return fmt.Errorf("cannot update chat %#v: %w", chat, result.Error)
	}

	return nil
}

func (db *Database) DeleteChat(id uint) error {
	result := db.gormDB.Delete(&Chat{}, id)
	if result.Error != nil {
		return fmt.Errorf("cannot delete chat with id = %d: %w", id, result.Error)
	}

	return nil
}

func (db *Database) AddUserToChat(user *User, chat *Chat) error {
	err := db.gormDB.Model(chat).Association("Members").Append(user)
	if err != nil {
		return fmt.Errorf("cannot add user %#v to chat %#v: %w", user, chat, err)
	}

	return nil
}
