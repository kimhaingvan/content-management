package database

import (
	"cms-project/cmd/content-server/config"
	"fmt"
	//"github.com/kiem-toan/cmd/audit-server/config"
	//"github.com/kiem-toan/infrastructure/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	DB *gorm.DB
}

func New(d config.Config) *Database {
	c := d.Databases.Postgres
	connString := fmt.Sprintf("dbname=%v user=%v password=%v host=%v port=%v sslmode=%v", c.Database, c.Username, c.Password, c.Host, c.Port, c.SSLMode)
	db, err := gorm.Open("postgres", connString)
	//db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
	//	SkipDefaultTransaction: true,
	//	//Logger:logger.New(),
	//})
	if err != nil {
		panic(err)
	}
	return &Database{DB: db}
}
