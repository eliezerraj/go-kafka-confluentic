package main

import (
	"log"
	"os"
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"

	"github.com/go-kafka-confluentic/internal/core"
	"github.com/go-kafka-confluentic/internal/adapter"

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

	log.Println("-> ctx", ctx)

	consumerService := adapter.NewConsumerService(config)

	go consumerService.Consumer(ctx)
	
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	
	log.Println("Encerrando...")
	defer func() {
		if err := consumerService.Close(ctx); err != nil {
			log.Printf("Failed to close connection: %s", err)
		}
		log.Println("Encerrado !!!")
	}()
}