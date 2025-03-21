package config

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Db  `yaml:"db"`
	App `yaml:"app"`
	S3  `yaml:"s3"`
}

type Db struct {
	Dsn string `yaml:"dsn"`
}

type S3 struct {
	BucketName string `yaml:"bucket_name"`
	Region     string `yaml:"region"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
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

func (c *Config) GetS3BucketName() string {
	return c.S3.BucketName
}

func (c *Config) GetS3Region() string {
	return c.S3.Region
}

func (c *Config) GetS3AccessKey() string {
	return c.S3.AccessKey
}

func (c *Config) GetS3SecretKey() string {
	return c.S3.SecretKey
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
