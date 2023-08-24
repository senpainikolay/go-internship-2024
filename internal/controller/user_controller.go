package controller

import (
	"mime/multipart"
	"net/http"
	"senpainikolay/go-internship-smartdata/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	GetById(uint) (models.UserInfo, error)
	Register(*models.UserModel) error
	LogIn(*models.UserCredentials) error
	DeleteById(uint) error
	UpdateImage(uint, *multipart.File, string) error
	GetImage(uint) (*[]byte, error)
}

type UserController struct {
	userSvc IUserService
}

func NewUserController(userSvc IUserService) *UserController {
	return &UserController{
		userSvc: userSvc,
	}
}

func (ctrl *UserController) GetById(c *gin.Context) {
	id_param := c.Param("id")
	u64, err := strconv.ParseUint(id_param, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid ID parameter",
		})
		return
	}

	id := uint(u64)
	user, err := ctrl.userSvc.GetById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "User not Found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (ctrl *UserController) DeleteById(c *gin.Context) {
	id_param := c.Param("id")
	u64, err := strconv.ParseUint(id_param, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid ID parameter",
		})
		return
	}

	id := uint(u64)
	err = ctrl.userSvc.DeleteById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "User succefully deleted",
	})
}

func (ctrl *UserController) Register(c *gin.Context) {
	var user models.UserModel
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	err := ctrl.userSvc.Register(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "successfully create an user",
	})

}

func (ctrl *UserController) LogIn(c *gin.Context) {
	var userCredentials models.UserCredentials
	if err := c.BindJSON(&userCredentials); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	err := ctrl.userSvc.LogIn(&userCredentials)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "succesfully logged in",
	})
}

func (ctrl *UserController) UpdateImage(c *gin.Context) {
	// Shift magic. Foloseste constante.
	err := c.Request.ParseMultipartForm(32 << 20) // approx 32 MB is the maximum file size; 2^20 bytes

	if err != nil {
		// StatusUnprocessableEntity sau ContentTooLarge ? https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/413
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	id_param := c.Param("id")
	u64, err := strconv.ParseUint(id_param, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid ID parameter",
		})
		return
	}

	id := uint(u64)

	file, handler, err := c.Request.FormFile("avatar")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	defer file.Close()

	err = ctrl.userSvc.UpdateImage(id, &file, handler.Filename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "image updated succesful",
	})

}

func (ctrl *UserController) GetImage(c *gin.Context) {
	id_param := c.Param("id")
	u64, err := strconv.ParseUint(id_param, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid ID parameter",
		})
		return
	}

	id := uint(u64)
	binaryImgAddr, err := ctrl.userSvc.GetImage(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return

	}
	c.Data(http.StatusOK, "image/jpeg", *binaryImgAddr)

}
