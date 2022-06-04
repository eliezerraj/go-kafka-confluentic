package adapter

import (
	"log"
	"os"
	"context"
	"time"
	//"os/signal"
	//"syscall"

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
	kafkaBrokerUrls := configurations.KafkaConfig.Brokers1 + "," + configurations.KafkaConfig.Brokers2 + "," + configurations.KafkaConfig.Brokers3
	
	log.Printf(kafkaBrokerUrls)

	config := &kafka.ConfigMap{	"bootstrap.servers":            kafkaBrokerUrls,
								"security.protocol":            configurations.KafkaConfig.Protocol, //"SASL_SSL",
								"sasl.mechanisms":              configurations.KafkaConfig.Mechanisms, //"SCRAM-SHA-512",
								"sasl.username":                configurations.KafkaConfig.Username,
								"sasl.password":                configurations.KafkaConfig.Password,
								"group.id":                     configurations.KafkaConfig.Groupid,
								"broker.address.family": 		"v4",
								"client.id": 					configurations.KafkaConfig.Clientid,
								"session.timeout.ms":    		6000,
								"auto.offset.reset":     		"earliest"}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		log.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	
	log.Printf("Created Consumer : %v\n", c)

	return &ConsumerService{ 	configurations : configurations,
								consumer : c,
	}
}

func (c *ConsumerService) Close(ctx context.Context) error{
	log.Printf("Encerrando Consumer....")
	if err := c.consumer.Close(); err != nil {
		log.Printf("===>>> Failed to close reader:", err)
		return err
	}
	log.Printf("Encerrado Consumer !!!!")
	return nil
}

func (c *ConsumerService) Consumer(ctx context.Context) {
	log.Printf("kafka Consumer")

	topics := []string{c.configurations.KafkaConfig.Topic}

	//sigchan := make(chan os.Signal, 1)
	//signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	err := c.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Printf("Failed to subscriber topic: %s\n", err)
		os.Exit(1)
	}

	run := true
	for run {
		select {
		//case sig := <-sigchan:
		//	log.Printf("Caught signal %v: terminating\n", sig)
		//	run = false
		default:
			ev := c.consumer.Poll(100)

			if ev == nil {
				continue
			}

			switch e := ev.(type) {
				case kafka.AssignedPartitions:
					c.consumer.Assign(e.Partitions)
				case kafka.RevokedPartitions:
					c.consumer.Unassign()	
				case kafka.PartitionEOF:
					log.Printf("%% Reached %v\n", e)
				case *kafka.Message:
					log.Printf("Topic %s:\n",e.TopicPartition)
					log.Print("----------------------------------")
					log.Print("Value : " ,string(e.Value))
					log.Print("-----------------------------------")
					if e.Headers != nil {
						log.Printf("Headers: %v\n", e.Headers)
					}
					c.consumer.Commit()
				case kafka.Error:
					log.Printf("%% Error: %v\n", e)
					if e.Code() == kafka.ErrAllBrokersDown {
						run = false
					}
				default:
					log.Print("Ignored %v\n", e)
			}
		}
		if ( lag_consumer > 0){
			log.Printf("Waiting for %v seconds", lag_consumer)
			time.Sleep(time.Second * time.Duration(lag_consumer))
		}
	}

	log.Printf("Closing consumer\n")
	c.consumer.Close()

}
