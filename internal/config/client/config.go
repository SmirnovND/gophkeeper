package config

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	App `yaml:"app"`
}

type App struct {
	ServerAddr string `yaml:"server_addr"`
}

func (c *Config) GetServerAddr() string {
	return c.App.ServerAddr
}

func NewConfig() interfaces.ConfigClient {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	cf := &Config{}
	// Проверяем, есть ли аргументы командной строки
	if len(os.Args) > 1 {
		cPath := os.Args[1]
		cf.LoadConfig(cPath)
		fmt.Println(cf)
	} else {
		// Если аргументов нет, используем значения по умолчанию или ищем config.yaml в текущей директории
		defaultPath := "cmd/client/config.yaml"
		if _, err := os.Stat(defaultPath); err == nil {
			cf.LoadConfig(defaultPath)
			fmt.Println(cf)
		} else {
			// Если файл не найден, используем значения по умолчанию
			cf.App.ServerAddr = "127.0.0.1:8085"
			fmt.Println("Конфигурационный файл не найден. Используются значения по умолчанию.")
			fmt.Println(cf)
		}
	}

	return cf
}

func (c *Config) LoadConfig(patch string) {
	file, err := os.Open(patch)
	if err != nil {
		log.Fatal("ReadConfigFile: ", err)
	}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatal("DecodeConfigFile: ", err)
	}
}
