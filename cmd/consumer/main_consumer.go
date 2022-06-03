package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/go-kafka/internal/core"

)

func LoadConfig() (*core.Configurations, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.ReadInConfig()

	var conf *core.Configurations
	err := viper.Unmarshal(&conf)
    if err != nil {
		return nil, err
	}
	return conf, nil
}

func main(){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting kafka")

	log.Println("-> Loading config.yaml")
	config, err := LoadConfig()
	if err != nil{
		log.Println("* FATAL ERROR load config.yaml *", err)
		os.Exit(3)
	}

	log.Println("-> Config", config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
}