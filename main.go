package main

import (
	"fmt"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/config"
	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/server"
	"github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

var opts struct {
	ConfigFile string `short:"c" long:"configfile" description:"File of the config file" default:".config/dev/config.yaml"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		TimestampFormat: "02-01-2006 15:04:05",
	},
	)
	log.SetReportCaller(true)

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
