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
		ctrl.logger.Info("invalid user json", "method", "ctrl.AddUser")

		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "json is invalid"})
		c.Abort()
		return
	}

	if err := ctrl.db.AddUser(&user); err != nil {
		if errors.Is(err, database.ErrUniqueConstraint) {
			ctrl.logger.Info("user with this email already registered", "method", "ctrl.AddUser")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user with this email already registered"})
		} else {
			ctrl.logger.Error(fmt.Sprintf("ctrl.db.AddUser(%#v): %s", user, err), "method", "ctrl.AddUser")
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (ctrl *Controller) GetUser(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid user id = %d", id), "method", "ctrl.GetUser")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	user, err := ctrl.db.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctrl.logger.Info(
				fmt.Sprintf("no user with id = %d: %s", id, err),
				"method",
				"ctrl.GetUser",
			)
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no user with this id"})
		} else {
			ctrl.logger.Info(
				fmt.Sprintf("ctrl.db.GetUserByID(uint(%d)): %s", id, err),
				"method",
				"ctrl.GetUser",
			)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (ctrl *Controller) UpdateUser(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid user id = %d", id), "method", "ctrl.GetUser")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	var user database.User
	if err := c.BindJSON(&user); err != nil {
		ctrl.logger.Info("invalid user json", "method", "ctrl.UpdateUser")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "json is invalid"})
		c.Abort()
		return
	}

	user.ID = uint(id)
	updatedUser, err := ctrl.db.UpdateUser(&user)
	if err != nil {
		if errors.Is(err, database.ErrUniqueConstraint) {
			ctrl.logger.Info("user with this email already registered", "method", "ctrl.UpdateUser")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user with this email already registered"})
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			ctrl.logger.Info(
				fmt.Sprintf("no user with id = %d: %s", id, err),
				"method",
				"ctrl.UpdateUser",
			)
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no user with this id"})
		} else {
			ctrl.logger.Error(fmt.Sprintf("ctrl.db.AddUser(%#v): %s", user, err), "method", "ctrl.UpdateUser")
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "server is unavailable now"})
		}

		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, updatedUser)
}

func (ctrl *Controller) DeleteUser(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Info(fmt.Sprintf("invalid user id = %d", id), "method", "ctrl.DeleteUser")
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
