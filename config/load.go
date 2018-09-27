package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/Eric-GreenComb/one-sendeth/bean"
)

// Ethereum Ethereum Config
var Ethereum bean.EthereumConfig

// Server Server Config
var Server bean.ServerConfig

// MariaDB 数据库相关配置
var MariaDB bean.DBConfig

// Nsq NsqConfig
var Nsq bean.NsqConfig

func init() {
	readConfig()
	initConfig()
}

func readConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
}

func initConfig() {
	Server.Port = strings.Split(viper.GetString("server.port"), ",")
	Server.Mode = viper.GetString("server.mode")
	Server.GormLogMode = viper.GetString("server.gorm.LogMode")
	Server.CallBack = viper.GetString("server.callback")

	Ethereum.ChainID = viper.GetInt64("ethereum.chainID")
	Ethereum.Host = viper.GetString("ethereum.host")
	Ethereum.Address = viper.GetString("ethereum.address")
	Ethereum.Passphrase = viper.GetString("ethereum.passphrase")

	MariaDB.Dialect = viper.GetString("database.dialect")
	MariaDB.Database = viper.GetString("database.database")
	MariaDB.User = viper.GetString("database.user")
	MariaDB.Password = viper.GetString("database.password")
	MariaDB.Host = viper.GetString("database.host")
	MariaDB.Port = viper.GetInt("database.port")
	MariaDB.Charset = viper.GetString("database.charset")
	MariaDB.MaxIdleConns = viper.GetInt("database.maxIdleConns")
	MariaDB.MaxOpenConns = viper.GetInt("database.maxOpenConns")
	MariaDB.URL = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		MariaDB.User, MariaDB.Password, MariaDB.Host, MariaDB.Port, MariaDB.Database, MariaDB.Charset)

	Nsq.Host = viper.GetString("nsq.host")
	Nsq.Topic = viper.GetString("nsq.topic")
}
