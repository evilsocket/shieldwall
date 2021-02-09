package database

import (
	"fmt"
	"github.com/evilsocket/islazy/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db = (*gorm.DB)(nil)
)

func Setup(config Config) (err error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Hostname,
		config.Port,
		config.User,
		config.Password,
		config.Name)
	db, err = gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err == nil {
		log.Info("connected to database")
		if err = db.Debug().AutoMigrate(&Agent{}); err != nil {
			return
		} else if err = db.Debug().AutoMigrate(&User{}); err != nil {
			return
		}
	}
	return
}