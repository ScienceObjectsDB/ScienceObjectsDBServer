package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

//ParseConfig Config parser
func ParseConfig() error {
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func GetMainDatabaseName() {
	viper.GetString("Config.Store.DatasetDatabaseName")
}
