package config

import (
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Db  `yaml:"db"`
	App `yaml:"app"`
}

type Db struct {
	Dsn string `yaml:"dsn"`
}

type App struct {
	JwtSecret string `yaml:"jwt_secret"`
}

func (c *Config) GetDBDsn() string {
	return c.Db.Dsn
}
func (c *Config) GetJwtSecret() string {
	return c.App.JwtSecret
}

func NewConfig() (cf interfaces.Config) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	cPath := os.Args[1]
	cf := &Config{}
	cf.LoadConfig(cPath)

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
