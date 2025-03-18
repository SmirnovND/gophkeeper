package config

import (
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Db       `yaml:"db"`
	App      `yaml:"app"`
	RabbitMQ `yaml:"rabbitmq"`
}

type Db struct {
	Dsn string `yaml:"dsn"`
}

type App struct {
	JwtSecret string `yaml:"jwt_secret"`
	RunAddr   string `yaml:"run_addr"`
}

type RabbitMQ struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	VHost    string `yaml:"vhost"`
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

func (c *Config) GetRabbitMQURI() string {
	return "amqp://" + c.RabbitMQ.User + ":" + c.RabbitMQ.Password + "@" +
		c.RabbitMQ.Host + ":" + c.RabbitMQ.Port + "/" + c.RabbitMQ.VHost
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
