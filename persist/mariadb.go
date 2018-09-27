package persist

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/Eric-GreenComb/one-sendeth/bean"
	"github.com/Eric-GreenComb/one-sendeth/config"
)

// ConnectDb connect Db
func ConnectDb() (*gorm.DB, error) {
	db, err := gorm.Open(config.MariaDB.Dialect, config.MariaDB.URL)

	if config.Server.GormLogMode == "false" {
		db.LogMode(false)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

// InitDatabase Init Db
func InitDatabase() {
	db, err := gorm.Open(config.MariaDB.Dialect, config.MariaDB.URL)
	defer db.Close()

	if config.Server.GormLogMode == "false" {
		db.LogMode(false)
	}

	if err != nil {
		panic(err)
	}

	if !db.HasTable(&bean.Order{}) {
		db.CreateTable(&bean.Order{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&bean.Order{})
	}

	fmt.Println(">>>>>>>>> db is ok")

	return
}
