package main

import (
	"os"
	authservice "senpainikolay/go-internship-smartdata/gateway/pkg/auth-service"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Host        string `yaml:"HOST"`
	Port        string `yaml:"PORT"`
	Enviroment  string `yaml:"ENVIRONMENT"`
	UserSvcPort string `yaml:"SVC_PORT1"`
}

func main() {
	config := configInit()

	router := gin.Default()
	router.Use(gin.Recovery())

	usrSvcClient := authservice.NewUserServiceClient(config.UserSvcPort)
	usrController := authservice.NewUserController(*usrSvcClient)
	authservice.AttachUserAuthRoutesToRouter(router, usrController)

	_ = router.Run(config.Port)
}

func configInit() Config {
	yamlFile := "config/config.yaml"

	data, err := os.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	return config

}
