package main

import (
	"log"

	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/config"
	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/server"
	"github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

var opts struct {
	ConfigFile string `short:"c" long:"configfile" description:"File of the config file" default:".config/dev/config.yaml"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalln(err.Error())
	}

	viper.SetConfigFile(opts.ConfigFile)
	err = config.ParseConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}

	config.ParseConfig()

	port := viper.GetInt64("Config.API.Endpointport")

	server := server.GRPCServerHandler{}
	server.StartGRPCServer(port)
}
