package adapter

import (
	"log"
	"os"
	"time"
	"os/signal"
	"syscall"
	"sync"
	"strconv"
	"math/rand"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/go-kafka-confluentic/internal/core"

)

var lag_commit = 0
var lag_consumer = 0

type ConsumerService struct{
	configurations		*core.Configurations
	consumer 			*kafka.Consumer
}

func NewConsumerService(configurations *core.Configurations) *ConsumerService {	
	lag_consumer = configurations.KafkaConfig.Lag
	lag_commit = configurations.KafkaConfig.LagCommit

	kafkaBrokerUrls := configurations.KafkaConfig.Brokers1 + "," + configurations.KafkaConfig.Brokers2 + "," + configurations.KafkaConfig.Brokers3
	

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 4
	client_id := configurations.KafkaConfig.Clientid +"-" + strconv.Itoa(rand.Intn(max-min+1) + min)

	log.Printf(kafkaBrokerUrls)

	config := &kafka.ConfigMap{	"bootstrap.servers":            kafkaBrokerUrls,
								"security.protocol":            configurations.KafkaConfig.Protocol, //"SASL_SSL",
								"sasl.mechanisms":              configurations.KafkaConfig.Mechanisms, //"SCRAM-SHA-512",
								"sasl.username":                configurations.KafkaConfig.Username,
								"sasl.password":                configurations.KafkaConfig.Password,
								"group.id":                     configurations.KafkaConfig.Groupid,
								"enable.auto.commit":           false, //true,
								"broker.address.family": 		"v4",
								"client.id": 					client_id,
								"session.timeout.ms":    		6000,
								"enable.idempotence":			true,
								"auto.offset.reset":     		"latest", //"earliest",
							}

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

func (c *ConsumerService) Consumer(wg *sync.WaitGroup) {
	log.Printf("kafka Consumer")

	topics := []string{c.configurations.KafkaConfig.Topic}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	err := c.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Printf("Failed to subscriber topic: %s\n", err)
		os.Exit(1)
	}

	run := true
	banner := 0
	count_rollback := 5
	count := 10
	is_abort := false
//	lag_commit = 1

	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
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
					if banner == 0 {
						log.Print("* * * *")
						banner = 1
					} else {
						log.Print("- - - -")
						banner = 0
					}
					log.Printf("Topic %s:\n",e.TopicPartition)
					log.Print("----------------------------------")
					if e.Headers != nil {
						log.Printf("Headers: %v\n", e.Headers)
					}
					log.Print("Value : " ,string(e.Value))
					log.Print("-----------------------------------")

					if ( lag_commit > 0){
						log.Print(" >>>> LAG COMMIT <<<")
						log.Printf("Waiting for %v seconds", lag_commit)
						time.Sleep(time.Millisecond * time.Duration(lag_commit))
					}
					
					if is_abort == true{
						count++
						if count%count_rollback == 0{
							log.Print("===> ABORTING !!!!!")
							os.Exit(3)
						} else {
							log.Print("===> COMMIT !!!!!")
							c.consumer.Commit()
						}
					} else {
						c.consumer.Commit()
					}

				case kafka.Error:
					log.Printf("%% Error: %v\n", e)
					if e.Code() == kafka.ErrAllBrokersDown {
						run = false
					}
				default:
					log.Printf("Ignored %v\n", e)
			}
		}
		if ( lag_consumer > 0){
			log.Print(" >>>> LAG MESSAGE <<<")
			log.Printf("Waiting for %v seconds", lag_consumer)
			time.Sleep(time.Second * time.Duration(lag_consumer))
		}
	}

	log.Printf("Closing consumer waiting please !!!! \n")
	c.consumer.Close()
	defer wg.Done()
}