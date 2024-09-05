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
	UpdateUser(user *database.User) (*database.User, error)
	DeleteUser(id uint) error

	AddPost(post *database.Post) error
	GetPost(id uint) (*database.Post, error)
	GetAllPosts() ([]*database.Post, error)
	UpdatePost(post *database.Post) (*database.Post, error)
	DeletePost(id uint) error
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
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user with this email already registered"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Error(fmt.Sprintf("ctrl.db.AddUser(%#v): %s", &user, err), "method", "ctrl.AddUser")
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
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user with this email already registered"})
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no user with this id"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		ctrl.logger.Error(fmt.Sprintf("ctrl.db.UpdateUser(%#v): %s", &user, err), "method", "ctrl.UpdateUser")
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
		ctrl.logger.Error(fmt.Sprintf("ctrl.db.DeleteUser(%d): %s", id, err), "method", "ctrl.DeleteUser")
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

		ctrl.logger.Error(fmt.Sprintf("ctrl.db.AddPost(%#v): %s", &post, err), "method", "ctrl.AddPost")
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
		ctrl.logger.Error(fmt.Sprintf("ctrl.db.GetAllUsers(): %s", err), "method", "ctrl.GetAllPosts")
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

		ctrl.logger.Error(fmt.Sprintf("ctrl.db.UpdatePost(%#v): %s", &post, err), "method", "ctrl.UpdatePost")
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
		ctrl.logger.Error(fmt.Sprintf("ctrl.db.DeletePost(%d): %s", id, err), "method", "ctrl.DeletePost")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}
