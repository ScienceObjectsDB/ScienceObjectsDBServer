package util

import (
	"os"

	"github.com/spf13/viper"
)

//InitTestEnv Initiates the test env
func InitTestEnv() error {

	if os.Getenv("MONGO_INITDB_ROOT_PASSWORD") == "" {
		os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "test123")
	}

	viper.Set("Config.Database.Mongo.URL", "localhost")
	viper.Set("Config.Database.Mongo.AuthSource", "admin")
	viper.Set("Config.Database.Mongo.Username", "root")

	if os.Getenv("MONGO_INITDB_URL") != "" {
		viper.Set("Config.Database.Mongo.URL", os.Getenv("MONGO_INITDB_URL"))
	}

	return nil
}
