package repository

import (
	"fmt"

	"github.com/tonybobo/go-chat/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _db *gorm.DB

func init() {
	username := config.GetConfig().MySql.User
	password := config.GetConfig().MySql.Password
	hostNport := config.GetConfig().MySql.HostnPort
	name := config.GetConfig().MySql.Name

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, hostNport, name, "10s")

	var err error

	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Fail to connect to DB err:" + err.Error())
	}

	sqlDB, err := _db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(100)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
}

func GetDB() *gorm.DB {
	return _db
}
