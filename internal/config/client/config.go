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

	cPath := os.Args[1]
	cf := &Config{}
	cf.LoadConfig(cPath)
	fmt.Println(cf)
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
