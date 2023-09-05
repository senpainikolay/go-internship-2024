package models

type Config struct {
	Host       string `yaml:"HOST"`
	Port       string `yaml:"PORT"`
	Enviroment string `yaml:"ENVIRONMENT"`
}
