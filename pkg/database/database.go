package database

import (
	"content-management/core/config"
	"content-management/pkg/log"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.elastic.co/apm/module/apmgorm"
)

type Database struct {
	GormDB *gorm.DB
}

func New(c config.PostgresConfig) *Database {
	connString := fmt.Sprintf("dbname=%v user=%v password=%v host=%v port=%v sslmode=%v", c.Database, c.Username, c.Password, c.Host, c.Port, c.SSLMode)
	db, err := apmgorm.Open("postgres", connString)
	if err != nil {
		log.Fatal(err, nil, nil)
	}
	db.LogMode(true)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetConnMaxLifetime(time.Minute * 5)

	if err = db.DB().Ping(); err != nil {
		log.Fatal(err, nil, nil)
	}
	return &Database{GormDB: db}
}
