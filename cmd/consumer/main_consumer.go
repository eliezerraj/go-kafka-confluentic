package main

import (
	"log"
	"os"
	"sync"
	//"os/signal"
	//"syscall"

	"github.com/spf13/viper"

	"github.com/go-kafka-confluentic/internal/core"
	"github.com/go-kafka-confluentic/internal/adapter"

)

func LoadConfig() (*core.Configurations, error) {
	log.Println("-> Loading config.yaml")
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

	config, err := LoadConfig()
	if err != nil{
		log.Println("* FATAL ERROR load config.yaml *", err)
		os.Exit(3)
	}
	log.Println("-> Config", config)

	consumerService := adapter.NewConsumerService(config)
	
    var wg sync.WaitGroup
	wg.Add(1)
	go consumerService.Consumer(&wg)
	wg.Wait()

}