package config

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp"
)

const (
	ServerPort = "8011"
)

var (
	appConfig = Config{}
)

func GetAppConfig() Config {
	return appConfig
}

func SetAppConfig(cfg Config) {
	appConfig = cfg
	appConfig.assignEnv()
	appConfig.printlnConfig()
}

// Config ...
type Config struct {
	ApplicationName string
	Databases       DBConfig    `json:"databases"`
	Minio           MinioConfig `json:"minio"`
	Log             Log         `json:"log"`
	LogStash        LogStash    `json:"logstash"`
	Zipkin          Zipkin      `json:"zipkin"`
	Consul          Consul
	ServerPort      string
}

type OracleConfig struct {
	UsernameOracleOTP         string `json:"username_oracle_otp"`
	PasswordOracleOTP         string `json:"password_oracle_otp"`
	ConnectionStringOracleOTP string `json:"connection_string_oracle_otp"`
	MaxOpenConnnsOracle       string `json:"max_open_connns_oracle"`
	MaxIdleConnsOracle        string `json:"max_idle_conns_oracle"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
	Timeout  int    `json:"timeout"`

	MaxOpenConns    int `json:"max_open_conns"`
	MaxIdleConns    int `json:"max_idle_conns"`
	MaxConnLifetime int `json:"max_conn_lifetime"`

	GoogleAuthFile string `json:"google_auth_file"`
}

type DBConfig struct {
	PostgresConfig PostgresConfig `json:"postgres_db"`
	OracleConfig   OracleConfig   `json:"oracle_db"`
}

// DefaultPostgres ...
func DefaultPostgres() PostgresConfig {
	return PostgresConfig{
		Host:           "localhost",
		Port:           5432,
		Username:       "postgres",
		Password:       "postgres",
		Database:       "postgres",
		SSLMode:        "disable",
		Timeout:        15,
		GoogleAuthFile: "",
	}
}

func DefaultConfig() *Config {
	return &Config{
		ApplicationName: "",
		Databases: DBConfig{
			PostgresConfig: PostgresConfig{
				Host:            "",
				Port:            0,
				Username:        "",
				Password:        "",
				Database:        "",
				SSLMode:         "",
				Timeout:         0,
				MaxOpenConns:    0,
				MaxIdleConns:    0,
				MaxConnLifetime: 0,
				GoogleAuthFile:  "",
			},
		},
		Minio: MinioConfig{
			Endpoint:        "",
			AccessKey:       "",
			SecretAccessKey: "",
			BucketName:      "",
		},
		Log: Log{
			Level: "",
		},
		LogStash: LogStash{
			Port: "",
			IP:   "",
		},
		Zipkin: Zipkin{
			URL: "",
		},
		Consul: Consul{
			ACLToken: "",
			IP:       "",
			Port:     "",
		},
		ServerPort: "",
	}
}

type MinioConfig struct {
	Endpoint        string `json:"endpoint"    `
	AccessKey       string `json:"access_key"   `
	SecretAccessKey string `json:"secret_access_key"   `
	BucketName      string `json:"bucket_name"    `
}

type Log struct {
	Level string `json:"level"`
}

type LogStash struct {
	Port string `json:"port"`
	IP   string `json:"ip"`
}

type Zipkin struct {
	URL string `json:"url"`
}

type Consul struct {
	ACLToken string
	IP       string
	Port     string
}

func (c *Config) assignEnv() {
	if os.Getenv("APPLICATION_NAME") != "" {
		c.ApplicationName = os.Getenv("APPLICATION_NAME")
	}
	if os.Getenv("CONSUL_IP") != "" {
		c.Consul.IP = os.Getenv("CONSUL_IP")
	}
	if os.Getenv("CONSUL_PORT") != "" {
		c.Consul.Port = os.Getenv("CONSUL_PORT")
	}
	if os.Getenv("CONSUL_ACL_TOKEN") != "" {
		c.Consul.ACLToken = os.Getenv("CONSUL_ACL_TOKEN")
	}
	if os.Getenv("LOGSTASH_IP") != "" {
		c.LogStash.IP = os.Getenv("LOGSTASH_IP")
	}
	if os.Getenv("LOGSTASH_PORT") != "" {
		c.LogStash.Port = os.Getenv("LOGSTASH_PORT")
	}
	if os.Getenv("ZIPKIN_URL") != "" {
		c.Zipkin.URL = os.Getenv("ZIPKIN_URL")
	}
	c.ServerPort = ServerPort
}

func (c *Config) printlnConfig() {
	fmt.Println("Thông số biến môi trường:")
	pp.Println(c)
}
