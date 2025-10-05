package utils

import (
	// "github.com/zohaibsoomro/users_api_golang/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDb() *gorm.DB {
	dsn := "host=40.81.227.40 user=zhs password=123 dbname=user_db port=5432 sslmode=require"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to datbase:" + err.Error())
	}
	db.Config.Logger = logger.Default.LogMode(logger.Silent)

	// if err := db.AutoMigrate(&models.User{}); err != nil {
	// 	panic("failed to migrate schema: " + err.Error())
	// }
	return db
}
