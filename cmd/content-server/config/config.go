package config

import (
	"content-management/pkg/redis"
	"fmt"
	"os"

	consulAPI "github.com/hashicorp/consul/api"

	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Databases DBConfig `yaml:",inline"`
	Env       string   `yaml:"env"`
	Port      string   `yaml:"port"`
	Redis     Redis    `yaml:"redis"`
	S3        S3       `yaml:"s3"`
}

type ConfigPostgres struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
	Timeout  int    `yaml:"timeout"`

	MaxOpenConns    int `yaml:"max_open_conns"`
	MaxIdleConns    int `yaml:"max_idle_conns"`
	MaxConnLifetime int `yaml:"max_conn_lifetime"`

	GoogleAuthFile string `yaml:"google_auth_file"`
}

type DBConfig struct {
	Postgres ConfigPostgres `yaml:"postgres"`
}

// DefaultPostgres ...
func DefaultPostgres() ConfigPostgres {
	return ConfigPostgres{
		Protocol:       "",
		Host:           "localhost",
		Port:           5432,
		Username:       "postgres",
		Password:       "postgres",
		Database:       "cms_postgres",
		SSLMode:        "disable",
		Timeout:        15,
		GoogleAuthFile: "",
	}
}

type Redis = redis.Redis

// DefaultRedis ...
func DefaultRedis() Redis {
	return Redis{
		Host:     "redis",
		Port:     "6379",
		Username: "",
		Password: "",
	}
}

type S3 struct {
	AwsS3Region        string `yaml:"aws_s3_region"    valid:"required"`
	AwsS3Bucket        string `yaml:"aws_s3_bucket"    valid:"required"`
	AwsAccessKey       string `yaml:"aws_access_key"    valid:"required"`
	AwsSecretAccessKey string `yaml:"aws_secret_access_key"    valid:"required"`
	AwsSessionToken    string `yaml:"aws_session_token"`
}

// Default ...
func Default() Config {
	cfg := Config{
		Databases: DBConfig{
			Postgres: DefaultPostgres(),
		},
		Env:   "dev",
		Port:  "8080",
		Redis: DefaultRedis(),
		S3:    S3{},
	}
	return cfg
}

// Load loads config from file
func LoadCfgFromConsul(addr, port string) (Config, error) {
	consulCfg := consulAPI.DefaultConfig()
	consulCfg.Address = fmt.Sprintf("%v:%v", addr, port)
	consulClient, err := consulAPI.NewClient(consulCfg)
	kv := consulClient.KV()
	pair, _, err := kv.Get(os.Getenv("CONSUL_CONFIG_KEY_VALUE"), nil)
	if err != nil {
		panic(err)
	}
	var cfg Config
	err = yaml.Unmarshal(pair.Value, &cfg)
	return cfg, err
}
