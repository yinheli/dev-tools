package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yinheli/go-toolbox/logger/log"
	"github.com/yinheli/go-toolbox/orm/gormext"
	"go.uber.org/zap"
)

var (
	DB          *gorm.DB
	currentAddr string
)

func InitDB(addr string) error {
	if DB != nil {
		if addr == currentAddr {
			return nil
		}

		DB.Close()
	}
	var err error
	log.Info("open...", zap.String("addr", addr))
	DB, err = gorm.Open("mysql", addr)
	if err != nil {
		return err
	}

	err = DB.DB().Ping()
	if err != nil {
		return err
	}

	currentAddr = addr

	DB.DB().SetMaxIdleConns(0)
	DB.SingularTable(true)
	DB.SetLogger(gormext.NewDBLogger(log.WithOptions(zap.AddCallerSkip(6))))

	DB.LogMode(true)

	return nil
}
