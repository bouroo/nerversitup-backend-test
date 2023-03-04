package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig(configName string) (err error) {
	// load config from env
	viper.AutomaticEnv()

	// load config from file
	viper.SetConfigName(configName)  // name of config file (without extension)
	viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./configs") // path to look for the config file inworking directory
	err = viper.ReadInConfig()       // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		log.Panicf("fatal error config file: %+v", err)
	}

	return
}
