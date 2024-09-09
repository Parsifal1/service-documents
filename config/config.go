package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
)

type Configs struct {
	AdminToken string `split_words:"true"`
}

var Config Configs

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	Config = Configs{}
	if err := envconfig.Process("DOCUMENT", &Config); err != nil {
		fmt.Printf("ошибка конфигурации: %v", err)
		os.Exit(1)
	}
}
