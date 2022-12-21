package main

import (
	"log"
	"os"
	"sync"
	//"os/signal"
	//"syscall"
	"strconv"

	"github.com/spf13/viper"

	"github.com/go-kafka-confluentic/internal/core"
	"github.com/go-kafka-confluentic/internal/adapter"

)

var app_kafka_config core.Configurations

func init(){
	log.Printf("==========================")
	log.Printf("Init")
	getEnvfromFile()
	getEnv()
	log.Printf("==========================")
}

func getEnvfromFile() {
	log.Println("-> Loading config.yaml")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.ReadInConfig()

	err := viper.Unmarshal(&app_kafka_config)
    if err != nil {
		log.Print("FATAL ERROR load application.yaml", err)
		os.Exit(3)
	}
	log.Println("Variables", app_kafka_config)
}

func getEnv() {
	log.Printf("-------------------")
	log.Printf("Overiding variables from enviroment")

	if os.Getenv("KAFKA_USER") !=  "" {
		app_kafka_config.KafkaConfig.Username = os.Getenv("KAFKA_USER")
	}
	if os.Getenv("KAFKA_PASSWORD") !=  "" {
		app_kafka_config.KafkaConfig.Password = os.Getenv("KAFKA_PASSWORD")
	}
	if os.Getenv("KAFKA_PROTOCOL") !=  "" {
		app_kafka_config.KafkaConfig.Protocol = os.Getenv("KAFKA_PROTOCOL")
	}
	if os.Getenv("KAFKA_MECHANISM") !=  "" {
		app_kafka_config.KafkaConfig.Mechanisms = os.Getenv("KAFKA_MECHANISM")
	}
	if os.Getenv("KAFKA_CLIENT_ID") !=  "" {
		app_kafka_config.KafkaConfig.Clientid = os.Getenv("KAFKA_CLIENT_ID")
	}
	if os.Getenv("KAFKA_BROKER_1") !=  "" {
		app_kafka_config.KafkaConfig.Brokers1 = os.Getenv("KAFKA_BROKER_1")
	}
	if os.Getenv("KAFKA_BROKER_2") !=  "" {
		app_kafka_config.KafkaConfig.Brokers2 = os.Getenv("KAFKA_BROKER_2")
	}
	if os.Getenv("KAFKA_BROKER_3") !=  "" {
		app_kafka_config.KafkaConfig.Brokers3 = os.Getenv("KAFKA_BROKER_3")
	}
	if os.Getenv("KAFKA_GROUP_ID") !=  "" {
		app_kafka_config.KafkaConfig.Groupid = os.Getenv("KAFKA_GROUP_ID")
	}
	if os.Getenv("KAFKA_TOPIC") !=  "" {
		app_kafka_config.KafkaConfig.Topic = os.Getenv("KAFKA_TOPIC")
	}
	if os.Getenv("KAFKA_PARTITION") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("KAFKA_PARTITION"))
		app_kafka_config.KafkaConfig.Partition = intVar
	}
	if os.Getenv("KAFKA_REPLICATION") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("KAFKA_REPLICATION"))
		app_kafka_config.KafkaConfig.ReplicationFactor = intVar
	}

	if os.Getenv("APP_LAG") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("APP_LAG"))
		app_kafka_config.KafkaConfig.Lag = intVar
	}
	if os.Getenv("APP_LAG_COMMIT") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("APP_LAG_COMMIT"))
		app_kafka_config.KafkaConfig.LagCommit = intVar
	}

	log.Println("Variables", app_kafka_config)
}

func main(){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting kafka")

	consumerService := adapter.NewConsumerService(&app_kafka_config)
	
    var wg sync.WaitGroup
	wg.Add(1)
	go consumerService.Consumer(&wg)
	wg.Wait()
}