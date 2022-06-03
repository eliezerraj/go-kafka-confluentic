package adapter

import (
	"log"
	"os"
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/go-kafka-confluentic/internal/core"

)

const consumer_timeout = 10
var lag_consumer = 0

type ConsumerService struct{
	configurations		*core.Configurations
	consumer 			*kafka.Consumer
}

func NewConsumerService(configurations *core.Configurations) *ConsumerService {	
	lag_consumer = configurations.KafkaConfig.Lag
	//kafkaBrokerUrls := []string {	configurations.KafkaConfig.Brokers1,
	//								configurations.KafkaConfig.Brokers2,
	//								configurations.KafkaConfig.Brokers3,}
	log.Printf(configurations.KafkaConfig.Brokers1)
	config := &kafka.ConfigMap{	"metadata.broker.list":            configurations.KafkaConfig.Brokers1,
								"security.protocol":               configurations.KafkaConfig.Protocol, //"SASL_SSL",
								"sasl.mechanisms":                 configurations.KafkaConfig.Mechanisms, //"SCRAM-SHA-256",
								"sasl.username":                   configurations.KafkaConfig.Username,
								"sasl.password":                   configurations.KafkaConfig.Password,
								"group.id":                        configurations.KafkaConfig.Groupid,
								"go.events.channel.enable":        true,
								"go.application.rebalance.enable": true,
								"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
								//"debug":                           "generic,broker,security",
								}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		log.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	
	log.Printf("Created Consumer %v\n", c)

	return &ConsumerService{ 	configurations : configurations,
								consumer : c,
	}
}

func (c *ConsumerService) Close(ctx context.Context) error{
	if err := c.consumer.Close(); err != nil {
		log.Printf("failed to close reader:", err)
		return err
	}
	return nil
}

func (c *ConsumerService) Consumer(ctx context.Context) {
	log.Printf("kafka Consumer")

	consumer := c.consumer

	err := consumer.Subscribe(c.configurations.KafkaConfig.Topic, nil)
	if err != nil {
		log.Printf("Failed to subscriber topic: %s\n", err)
		os.Exit(1)
	}
	
	run := true

	for run == true {
		ev := consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			log.Printf("%% Message on %s:\n%s\n",
				e.TopicPartition, string(e.Value))
		case kafka.PartitionEOF:
			log.Printf("%% Reached %v\n", e)
		case kafka.Error:
			log.Printf("%% Error: %v\n", e)
			run = false
		default:
			log.Printf("Ignored %v\n", e)
		}
		if ( lag_consumer > 0){
			//log.Println("Waiting for %s", lag_consumer)
			time.Sleep(time.Second * time.Duration(lag_consumer))
		}
	}

}
