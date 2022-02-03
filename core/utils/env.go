package utils

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func GetEnvValue(key string) string {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//TODO: solve it
			panic("Config file not found.")
		} else {
			log.Fatalf("Config file was found but another error was produced: %s\n", err)
		}
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf(fmt.Sprintf("Invalid type assertion: %s", key))
	}

	return value
}
