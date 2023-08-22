package main

import (
	"fmt"
	"senpainikolay/go-internship-smartdata/internal/controller"
	"senpainikolay/go-internship-smartdata/internal/repository"
	"senpainikolay/go-internship-smartdata/internal/service"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "test_db"
)

func main() {
	router := gin.Default()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	userRepo := repository.NewUserRepository(psqlInfo)
	userService := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userService)

	userRouter := router.Group("/user")
	{
		userRouter.GET("/me/:id", userCtrl.GetById)
		userRouter.POST("/register", userCtrl.Register)
		userRouter.POST("/login", userCtrl.LogIn)
		userRouter.DELETE("/unregister/:id", userCtrl.DeleteById)

	}

	_ = router.Run(":8888") // listen and serve on 0.0.0.0:8888
}
