package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LLIEPJIOK/forum/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DBInterface interface {
	AddUser(user *database.User) error
	GetUserByID(id uint) (*database.User, error)
	GetAllUsers() ([]*database.User, error)
	UpdateUser(user *database.User) (*database.User, error)
	DeleteUser(id uint) error

	AddPost(post *database.Post) error
	GetPost(id uint) (*database.Post, error)
	GetAllPosts() ([]*database.Post, error)
	UpdatePost(post *database.Post) (*database.Post, error)
	DeletePost(id uint) error

	AddMessage(message *database.Message) error
	GetMessage(id uint) (*database.Message, error)
	GetAllMessages() ([]*database.Message, error)
	UpdateMessage(message *database.Message) (*database.Message, error)
	DeleteMessage(id uint) error

	AddChat(chat *database.Chat) error
	GetChat(id uint) (*database.Chat, error)
	GetAllChats() ([]*database.Chat, error)
	UpdateChat(chat *database.Chat) (*database.Chat, error)
	DeleteChat(id uint) error
}

type Controller struct {
	db     DBInterface
	logger *slog.Logger
}

func New(db DBInterface, logger *slog.Logger) *Controller {
	return &Controller{
		db:     db,
		logger: logger,
	}
}

func (ctrl *Controller) AddUser(c *gin.Context) {
	var user database.User
	if err := c.BindJSON(&user); err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid user json: %s", err), "method", "ctrl.AddUser")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "json is invalid"})
		c.Abort()
		return
	}

	if err := ctrl.db.AddUser(&user); err != nil {
		if errors.Is(err, database.ErrUniqueConstraint) {
			c.IndentedJSON(
				http.StatusBadRequest,
				gin.H{"error": "user with this email already registered"},
			)
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.AddUser(%#v): %s", &user, err),
			"method",
			"ctrl.AddUser",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (ctrl *Controller) GetUser(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid user id: %s", err), "method", "ctrl.GetUser")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	user, err := ctrl.db.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no user with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Info(
			fmt.Sprintf("ctrl.db.GetUserByID(uint(%d)): %s", id, err),
			"method",
			"ctrl.GetUser",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (ctrl *Controller) GetAllUsers(c *gin.Context) {
	users, err := ctrl.db.GetAllUsers()
	if err != nil {
		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.GetAllUsers(): %s", err),
			"method",
			"ctrl.GetAllUsers",
		)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func (ctrl *Controller) UpdateUser(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid user id: %s", err), "method", "ctrl.GetUser")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	var user database.User
	if err := c.BindJSON(&user); err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid user json: %s", err), "method", "ctrl.UpdateUser")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "json is invalid"})
		c.Abort()
		return
	}

	user.ID = uint(id)
	updatedUser, err := ctrl.db.UpdateUser(&user)
	if err != nil {
		if errors.Is(err, database.ErrUniqueConstraint) {
			c.IndentedJSON(
				http.StatusBadRequest,
				gin.H{"error": "user with this email already registered"},
			)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no user with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.UpdateUser(%#v): %s", &user, err),
			"method",
			"ctrl.UpdateUser",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, updatedUser)
}

func (ctrl *Controller) DeleteUser(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid user id: %s", err), "method", "ctrl.DeleteUser")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	if err := ctrl.db.DeleteUser(uint(id)); err != nil {
		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.DeleteUser(%d): %s", id, err),
			"method",
			"ctrl.DeleteUser",
		)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

func (ctrl *Controller) AddPost(c *gin.Context) {
	var post database.Post
	if err := c.BindJSON(&post); err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid post json: %s", err), "method", "ctrl.AddPost")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "json is invalid"})
		c.Abort()
		return
	}

	if err := ctrl.db.AddPost(&post); err != nil {
		if errors.Is(err, database.ErrForeignKeyConstraint) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no such author with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.AddPost(%#v): %s", &post, err),
			"method",
			"ctrl.AddPost",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, post)
}

func (ctrl *Controller) GetPost(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid post id: %s", err), "method", "ctrl.GetPost")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		c.Abort()
		return
	}

	post, err := ctrl.db.GetPost(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no post with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Info(
			fmt.Sprintf("ctrl.db.GetPost(uint(%d)): %s", id, err),
			"method",
			"ctrl.GetPost",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, post)
}

func (ctrl *Controller) GetAllPosts(c *gin.Context) {
	posts, err := ctrl.db.GetAllPosts()
	if err != nil {
		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.GetAllUsers(): %s", err),
			"method",
			"ctrl.GetAllPosts",
		)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func (ctrl *Controller) UpdatePost(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid post id: %s", err), "method", "ctrl.UpdatePost")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		c.Abort()
		return
	}

	var post database.Post
	if err := c.BindJSON(&post); err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid post json: %s", err), "method", "ctrl.UpdatePost")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "json is invalid"})
		c.Abort()
		return
	}

	post.ID = uint(id)
	updatedPost, err := ctrl.db.UpdatePost(&post)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no post with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.UpdatePost(%#v): %s", &post, err),
			"method",
			"ctrl.UpdatePost",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, updatedPost)
}

func (ctrl *Controller) DeletePost(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid post id: %s", err), "method", "ctrl.DeletePost")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		c.Abort()
		return
	}

	if err := ctrl.db.DeletePost(uint(id)); err != nil {
		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.DeletePost(%d): %s", id, err),
			"method",
			"ctrl.DeletePost",
		)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

func (ctrl *Controller) AddMessage(c *gin.Context) {
	var message database.Message
	if err := c.BindJSON(&message); err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid message json: %s", err), "method", "ctrl.AddMessage")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "json is invalid"})
		c.Abort()
		return
	}

	if err := ctrl.db.AddMessage(&message); err != nil {
		if errors.Is(err, database.ErrForeignKeyConstraint) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no such creator with this id or chat with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.AddMessage(%#v): %s", &message, err),
			"method",
			"ctrl.AddMessage",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, message)
}

func (ctrl *Controller) GetMessage(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid message id: %s", err), "method", "ctrl.GetMessage")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid message id"})
		c.Abort()
		return
	}

	message, err := ctrl.db.GetMessage(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no message with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Info(
			fmt.Sprintf("ctrl.db.GetMessage(uint(%d)): %s", id, err),
			"method",
			"ctrl.GetMessage",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, message)
}

func (ctrl *Controller) GetAllMessages(c *gin.Context) {
	messages, err := ctrl.db.GetAllMessages()
	if err != nil {
		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.GetAllMessages(): %s", err),
			"method",
			"ctrl.GetAllMessages",
		)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, messages)
}

func (ctrl *Controller) UpdateMessage(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid message id: %s", err), "method", "ctrl.UpdateMessage")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid message id"})
		c.Abort()
		return
	}

	var message database.Message
	if err := c.BindJSON(&message); err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid message json: %s", err), "method", "ctrl.UpdateMessage")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "json is invalid"})
		c.Abort()
		return
	}

	message.ID = uint(id)
	updatedMessage, err := ctrl.db.UpdateMessage(&message)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no message with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.UpdateMessage(%#v): %s", &message, err),
			"method",
			"ctrl.UpdateMessage",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, updatedMessage)
}

func (ctrl *Controller) DeleteMessage(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid message id: %s", err), "method", "ctrl.DeleteMessage")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid message id"})
		c.Abort()
		return
	}

	if err := ctrl.db.DeleteMessage(uint(id)); err != nil {
		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.DeleteMessage(%d): %s", id, err),
			"method",
			"ctrl.DeleteMessage",
		)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

func (ctrl *Controller) AddChat(c *gin.Context) {
	var chat database.Chat
	if err := c.BindJSON(&chat); err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid chat json: %s", err), "method", "ctrl.AddChat")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "json is invalid"})
		c.Abort()
		return
	}

	if err := ctrl.db.AddChat(&chat); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.AddChat(%#v): %s", &chat, err),
			"method",
			"ctrl.AddChat",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, chat)
}

func (ctrl *Controller) GetChat(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid chat id: %s", err), "method", "ctrl.GetChat")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		c.Abort()
		return
	}

	chat, err := ctrl.db.GetChat(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no chat with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Info(
			fmt.Sprintf("ctrl.db.GetChat(uint(%d)): %s", id, err),
			"method",
			"ctrl.GetChat",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, chat)
}

func (ctrl *Controller) GetAllChats(c *gin.Context) {
	chats, err := ctrl.db.GetAllChats()
	if err != nil {
		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.GetAllChats(): %s", err),
			"method",
			"ctrl.GetAllChats",
		)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, chats)
}

func (ctrl *Controller) UpdateChat(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid chat id: %s", err), "method", "ctrl.UpdateChat")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		c.Abort()
		return
	}

	var chat database.Chat
	if err := c.BindJSON(&chat); err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid chat json: %s", err), "method", "ctrl.UpdateChat")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "json is invalid"})
		c.Abort()
		return
	}

	chat.ID = uint(id)
	updatedChat, err := ctrl.db.UpdateChat(&chat)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no chat with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.UpdateChat(%#v): %s", &chat, err),
			"method",
			"ctrl.UpdateChat",
		)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, updatedChat)
}

func (ctrl *Controller) DeleteChat(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid chat id: %s", err), "method", "ctrl.DeleteChat")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		c.Abort()
		return
	}

	if err := ctrl.db.DeleteChat(uint(id)); err != nil {
		ctrl.logger.Error(
			fmt.Sprintf("ctrl.db.DeleteChat(%d): %s", id, err),
			"method",
			"ctrl.DeleteChat",
		)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}
