package authservice

import (
	"io/ioutil"
	"net/http"
	"senpainikolay/go-internship-smartdata/gateway/pkg/models"
	tokenvalservice "senpainikolay/go-internship-smartdata/gateway/pkg/token-val-service"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	NUM_FILE_BYTES      = 32
	SHIFTED_BYTES       = 20
	wsSocketBufferSize  = 1024
	wsMessageBufferSize = 256
)

type UserController struct {
	userAuthSvcCllient UserServiceClient
	TokenValSvcClient  tokenvalservice.TokenValidationServiceClient
}

func NewUserController(userAuthSvcCllient UserServiceClient, tokenSvcClient tokenvalservice.TokenValidationServiceClient) *UserController {
	return &UserController{
		userAuthSvcCllient: userAuthSvcCllient,
		TokenValSvcClient:  tokenSvcClient,
	}
}

func AttachUserAuthRoutesToRouter(r *gin.Engine, c *UserController) {

	userRouter := r.Group("/user")
	{

		userRouter.GET("/me", RequireAuth(c), c.GetById)
		userRouter.POST("/register", c.Register)
		userRouter.POST("/login", c.LogIn)
		userRouter.DELETE("/unregister/:id", RequireAuth(c), c.DeleteById)
		userRouter.PUT("/img", RequireAuth(c), c.UpdateImage)
		userRouter.GET("/:id/img", RequireAuth(c), c.GetImage)
		userRouter.POST("/refreshToken", RequireAuth(c), c.RefreshToken)

	}
}

func (ctrl *UserController) GetById(c *gin.Context) {

	val, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "required authentication",
		})
		return
	}

	usr, ok := val.(models.UserJWTInfo)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "wrong info extrated from token",
		})
		return
	}

	user, err := ctrl.userAuthSvcCllient.GetById(usr.ID)
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
	err = ctrl.userAuthSvcCllient.DeleteById(id)
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
	err := ctrl.userAuthSvcCllient.Register(&user)
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

	tokens, err := ctrl.userAuthSvcCllient.LogIn(&userCredentials)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	seconds_to_expire := 3600 * 8
	c.SetCookie("Authorization", tokens.RefreshToken, seconds_to_expire, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token":   tokens.AccesToken,
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

	file, _, err := c.Request.FormFile("avatar")
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

	err = ctrl.userAuthSvcCllient.UploadUsrImage(uint64(usr.ID), &fileContents)
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
	binaryImgAddr, err := ctrl.userAuthSvcCllient.GetUsrImage(u64)
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

	usrJwtInfo, err := ctrl.TokenValSvcClient.ValidateUsrRefreshToken(refresh_token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	newTokens, err := ctrl.TokenValSvcClient.GenerateUsrTokenPair(uint64(usrJwtInfo.ID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return

	}

	c.SetSameSite(http.SameSiteLaxMode)
	seconds_to_expire := 3600 * 8
	c.SetCookie("Authorization", newTokens.RefreshToken, seconds_to_expire, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token":   newTokens.AccesToken,
		"message": "succesfully logged in",
	})

}
