package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUniqueConstraint     = errors.New("duplicate primary key value violates uniqueness constraint")
	ErrForeignKeyConstraint = errors.New("missing key in external table violates foreign key constraint")
)

type Database struct {
	gormDB *gorm.DB
}

func New(gormDB *gorm.DB) *Database {
	return &Database{
		gormDB: gormDB,
	}
}

func (db *Database) Migrate() error {
	err := db.gormDB.AutoMigrate(User{}, Post{}, Message{}, Chat{})
	if err != nil {
		return fmt.Errorf("cannot create tables: %w", err)
	}

	return nil
}

func (db *Database) AddUser(user *User) error {
	_, err := db.GetUserByEmail(*user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result := db.gormDB.Create(user)
			if result.Error != nil {
				return fmt.Errorf("cannot add user %#v to db: %w", user, result.Error)
			}

			return nil
		} else {
			return fmt.Errorf("cannot check email = %#v existence: %w", user.Email, err)
		}
	}

	return fmt.Errorf("cannot add user %#v to db: %w", user, ErrUniqueConstraint)
}

func (db *Database) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	result := db.gormDB.Where("email = ?", email).First(user)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get user by email = %q: %w", email, result.Error)
	}

	return user, nil
}

func (db *Database) GetUserByID(id uint) (*User, error) {
	user := &User{}
	result := db.gormDB.Where("id = ?", id).First(user)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get user by id = %d: %w", id, result.Error)
	}
	if user.RemovedAt.Valid {
		return nil, fmt.Errorf("cannot get user by id = %d: %w", id, gorm.ErrRecordNotFound)
	}

	return user, nil
}

func (db *Database) GetAllUsers() ([]*User, error) {
	var users []*User
	result := db.gormDB.Where("removed_at IS NULL").Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get all users: %w", result.Error)
	}

	return users, nil
}

func (db *Database) UpdateUser(user *User) (*User, error) {
	if *user.Email != "" {
		if err := db.gormDB.Where("email = ?", user.Email).First(&User{}).Error; err == nil {
			return nil, fmt.Errorf("cannot update user %#v: %w", user, ErrUniqueConstraint)
		}
	}

	_, err := db.GetUserByID(user.ID)
	if err == nil {
		return nil, fmt.Errorf("cannot update user %#v: %w", user, gorm.ErrRecordNotFound)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("cannot update user %#v: %w", user, err)
	}

	result := db.gormDB.Model(&User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot update user %#v: %w", user, result.Error)
	}

	updatedUser, err := db.GetUserByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("db.GetUserByID(%d): %w", user.ID, err)
	}

	return updatedUser, nil
}

func (db *Database) DeleteUser(id uint) error {
	result := db.gormDB.Model(&User{}).Select("nickname", "email", "removed_at").Where("id = ?", id).Updates(&User{
		Nickname: "Deleted user",
		Email:    nil,
		RemovedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if result.Error != nil {
		return fmt.Errorf("cannot delete user with id = %d: %w", id, result.Error)
	}

	result = db.gormDB.Table("user_x_chat").Where("user_id = ?", id).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("cannot delete user with id = %d from chats: %w", id, result.Error)
	}

	return nil
}

func (db *Database) AddPost(post *Post) error {
	_, err := db.GetUserByID(post.AuthorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cannot add post %#v to db: %w", post, ErrForeignKeyConstraint)
		} else {
			return fmt.Errorf("db.GetUserByID(%d): %w", post.AuthorID, err)
		}
	}

	result := db.gormDB.Create(post)
	if result.Error != nil {
		return fmt.Errorf("cannot add post %#v to db: %w", post, result.Error)
	}

	return nil
}

func (db *Database) GetAllPosts() ([]*Post, error) {
	var posts []*Post
	result := db.gormDB.Find(&posts)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get all posts: %w", result.Error)
	}

	return posts, nil
}

func (db *Database) GetPost(id uint) (*Post, error) {
	post := &Post{}
	result := db.gormDB.Where("id = ?", id).First(post)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get post by id = %d: %w", id, result.Error)
	}

	return post, nil
}

func (db *Database) UpdatePost(post *Post) (*Post, error) {
	result := db.gormDB.Model(&Post{}).Select("content").Where("id = ?", post.ID).Updates(post)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot update post %#v: %w", post, result.Error)
	}

	updatedPost, err := db.GetPost(post.ID)
	if err != nil {
		return nil, fmt.Errorf("db.GetPost(%d): %w", post.ID, err)
	}

	return updatedPost, nil
}

func (db *Database) DeletePost(id uint) error {
	result := db.gormDB.Delete(&Post{}, id)
	if result.Error != nil {
		return fmt.Errorf("cannot delete post with id = %d: %w", id, result.Error)
	}

	return nil
}

func (db *Database) AddMessage(message *Message) error {
	if _, err := db.GetUserByID(message.SenderID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cannot add message %#v to db: %w", message, ErrForeignKeyConstraint)
		} else {
			return fmt.Errorf("db.GetUserByID(%d): %w", message.SenderID, err)
		}
	}

	if _, err := db.GetChat(message.ChatID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cannot add message %#v to db: %w", message, ErrForeignKeyConstraint)
		} else {
			return fmt.Errorf("db.GetChat(%d): %w", message.ChatID, err)
		}
	}

	result := db.gormDB.Create(message)
	if result.Error != nil {
		return fmt.Errorf("cannot add message %#v to db: %w", message, result.Error)
	}

	return nil
}

func (db *Database) GetMessage(id uint) (*Message, error) {
	message := &Message{}
	result := db.gormDB.Where("id = ?", id).First(message)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get message by id = %d: %w", id, result.Error)
	}

	return message, nil
}

func (db *Database) GetAllMessages() ([]*Message, error) {
	var messages []*Message
	result := db.gormDB.Find(&messages)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get all messages: %w", result.Error)
	}

	return messages, nil
}

func (db *Database) UpdateMessage(message *Message) (*Message, error) {
	result := db.gormDB.Model(&Message{}).Select("content").Where("id = ?", message.ID).Updates(message)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot update message %#v: %w", message, result.Error)
	}

	updatedMessage, err := db.GetMessage(message.ID)
	if err != nil {
		return nil, fmt.Errorf("db.GetMessage(%d): %w", message.ID, err)
	}

	return updatedMessage, nil
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
	chat := &Chat{}
	result := db.gormDB.Where("id = ?", id).First(chat)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get chat by id = %d: %w", id, result.Error)
	}

	return chat, nil
}

func (db *Database) GetAllChats() ([]*Chat, error) {
	var chats []*Chat
	result := db.gormDB.Find(&chats)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get all chats: %w", result.Error)
	}

	return chats, nil
}

func (db *Database) UpdateChat(chat *Chat) (*Chat, error) {
	result := db.gormDB.Model(&Chat{}).Select("name").Where("id = ?", chat.ID).Updates(chat)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot update chat %#v: %w", chat, result.Error)
	}

	updatedChat, err := db.GetChat(chat.ID)
	if err != nil {
		return nil, fmt.Errorf("db.GetChat(%d): %w", chat.ID, err)
	}

	return updatedChat, nil
}

func (db *Database) DeleteChat(id uint) error {
	result := db.gormDB.Table("user_x_chat").Where("chat_id = ?", id).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("cannot delete users from chat with id = %d: %w", id, result.Error)
	}

	result = db.gormDB.Model(&Message{}).Where("chat_id = ?", id).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("cannot delete users from chat with id = %d: %w", id, result.Error)
	}

	result = db.gormDB.Delete(&Chat{}, id)
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
