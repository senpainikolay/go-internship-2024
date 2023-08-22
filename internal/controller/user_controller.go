package controller

import (
	"net/http"
	"senpainikolay/go-internship-smartdata/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	GetById(int) string
	Register(string, string) error
	LogIn(string, string) string
	DeleteById(int) error
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
	id_param := c.Query("id")
	id, _ := strconv.Atoi(id_param)
	// aici trebuie sa avem procesam eroarea de la convertire. si daca nu e int returnam bad request
	// vad ca ai facut prin query parameter, mai bine parameter pentru ca asta este ceva concret ce tine de un utilizator
	// nu e ceva adaugator. Cand noi folosim ca URL parameter, gin singur valideaza acest URL si daca nu exista acest id returneaza 404
	
	user := ctrl.userSvc.GetById(id)
	c.JSON(http.StatusOK, gin.H{
		"msg": user,
	})
}

func (ctrl *UserController) DeleteById(c *gin.Context) {
	id_param := c.Query("id")
	id, _ := strconv.Atoi(id_param)
	ctrl.userSvc.DeleteById(id)
	c.JSON(http.StatusOK, gin.H{
		"msg": "deleted request sent.",
	})
}

func (ctrl *UserController) Register(c *gin.Context) {
	var requestData models.UserCredentials

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.userSvc.Register(requestData.Email, requestData.Password)

	c.JSON(http.StatusOK, gin.H{"message": "Registered"})
}

func (ctrl *UserController) LogIn(c *gin.Context) {
	var requestData models.UserCredentials

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := ctrl.userSvc.LogIn(requestData.Email, requestData.Password)

	c.JSON(http.StatusOK, gin.H{"message": msg})
}
