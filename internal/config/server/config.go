package config

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Db    `yaml:"db"`
	App   `yaml:"app"`
	Minio `yaml:"minio"`
}

type Db struct {
	Dsn string `yaml:"dsn"`
}

type Minio struct {
	BucketName string `yaml:"bucket_name"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	Host       string `yaml:"host"`
}

type App struct {
	JwtSecret string `yaml:"jwt_secret"`
	RunAddr   string `yaml:"run_addr"`
}

func (c *Config) GetDBDsn() string {
	return c.Db.Dsn
}

func (c *Config) GetJwtSecret() string {
	return c.App.JwtSecret
}

func (c *Config) GetRunAddr() string {
	return c.App.RunAddr
}

func (c *Config) GetMinioBucketName() string {
	return c.Minio.BucketName
}

func (c *Config) GetMinioAccessKey() string {
	return c.Minio.AccessKey
}

func (c *Config) GetMinioSecretKey() string {
	return c.Minio.SecretKey
}

func (c *Config) GetMinioHost() string {
	return c.Minio.Host
}

func NewConfig() interfaces.ConfigServer {
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
