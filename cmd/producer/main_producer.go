package main

import (
	"log"
	"os"
	"context"
	"os/signal"
	"syscall"
	"time"

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

	producerService := adapter.NewProducerService(config)

	done := make(chan string)
	go post(ctx, *producerService ,done)
	
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	
	log.Println("Encerrando...")
	defer func() {
		log.Println("Encerrado !!!")
	}()
}

func post(ctx context.Context, producerService adapter.ProducerService, done chan string){
	for i := 0 ; i < 3600; i++ {
		producerService.Producer(ctx, i)
		time.Sleep(time.Millisecond * time.Duration(10000))
	}
	done <- "END"
}