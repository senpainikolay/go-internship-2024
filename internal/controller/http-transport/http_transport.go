package htttptransport

import (
	"io/ioutil"
	"net/http"
	"senpainikolay/go-internship-smartdata/internal/auth"
	"senpainikolay/go-internship-smartdata/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	NUM_FILE_BYTES = 32
	SHIFTED_BYTES  = 20
)

type IUserService interface {
	GetById(uint) (models.UserInfo, error)
	Register(*models.UserModel) error
	LogIn(*models.UserCredentials) (map[string]string, error)
	DeleteById(uint) error
	UpdateImage(uint, *[]byte, string) error
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

	val, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "user info not taken from context ",
		})
		return
	}

	usr, ok := val.(models.UserJWTInfo)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "wrong user info saved in context",
		})
		return
	}

	user, err := ctrl.userSvc.GetById(usr.ID)
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

	tokensMap, err := ctrl.userSvc.LogIn(&userCredentials)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	seconds_to_expire := 3600 * 8
	c.SetCookie("Authorization", tokensMap["refresh_token"], seconds_to_expire, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token":   tokensMap["access_token"],
		"message": "succesfully logged in",
	})
}

func (ctrl *UserController) UpdateImage(c *gin.Context) {
	err := c.Request.ParseMultipartForm(NUM_FILE_BYTES << SHIFTED_BYTES) // approx 32 MB is the maximum file size; 2^20 bytes

	if err != nil {
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	val, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "user info not taken from context ",
		})
		return
	}

	usr, ok := val.(models.UserJWTInfo)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "wrong user info saved in context",
		})
		return
	}

	file, handler, err := c.Request.FormFile("avatar")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Unable to read file contents",
		})
		return
	}

	err = ctrl.userSvc.UpdateImage(usr.ID, &fileContents, handler.Filename)
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

func (ctrl *UserController) RefreshToken(c *gin.Context) {
	refresh_token, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	usrId, err := auth.ValidateRefreshToken(refresh_token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	newTokens, err := auth.GenerateTokenPair(usrId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return

	}

	c.SetSameSite(http.SameSiteLaxMode)
	seconds_to_expire := 3600 * 8
	c.SetCookie("Authorization", newTokens["refresh_token"], seconds_to_expire, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token":   newTokens["access_token"],
		"message": "succesfully logged in",
	})

}
