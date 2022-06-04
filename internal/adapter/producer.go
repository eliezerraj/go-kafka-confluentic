package adapter

import (
	"log"
	"os"
	"context"
	"time"
	"encoding/json"
	"strconv"
	"math/rand"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/go-kafka-confluentic/internal/core"

)
const producer_timeout = 10

type Message struct {
    ID          int     `json:"id"`
    Description string  `json:"description"`
    Status      bool    `json:"status"`
}

type ProducerService struct{
	configurations  *core.Configurations
	producer        *kafka.Producer
}

func NewProducerService(configurations *core.Configurations) *ProducerService {

	kafkaBrokerUrls := configurations.KafkaConfig.Brokers1 + "," + configurations.KafkaConfig.Brokers2 + "," + configurations.KafkaConfig.Brokers3
	log.Printf(kafkaBrokerUrls)
	config := &kafka.ConfigMap{	"bootstrap.servers":            kafkaBrokerUrls,
								"security.protocol":            configurations.KafkaConfig.Protocol, //"SASL_SSL",
								"sasl.mechanisms":              configurations.KafkaConfig.Mechanisms, //"SCRAM-SHA-256",
								"sasl.username":                configurations.KafkaConfig.Username,
								"sasl.password":                configurations.KafkaConfig.Password,
								"group.id":                     configurations.KafkaConfig.Groupid,
								"client.id": 					configurations.KafkaConfig.Clientid,
								"acks": 						"all",
								"default.topic.config":         kafka.ConfigMap{"auto.offset.reset": "earliest"},
								//"debug":                           "generic,broker,security",
								}

	p, err := kafka.NewProducer(config)
	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Created Producer %v\n", p)

	return &ProducerService{ configurations : configurations,
							producer : p,
	}
}

func (p *ProducerService) Close(ctx context.Context) error{
	/*if err := p.producer.Close(); err != nil {
		log.Printf("failed to close producer:", err)
		return err
	}*/
	return nil
}

func (p *ProducerService) Producer(ctx context.Context, i int) {
	log.Printf("kafka Producer")

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 4
	key := "key-"+ strconv.Itoa(rand.Intn(max-min+1) + min)
	//message := "teste-(" + strconv.Itoa(rand.Intn(max-min+1) + min)  + ")-" + strconv.Itoa(i)

	message := Message{}
	message.ID = rand.Intn(max-min+1) + min
	message.Description = "teste-(" + strconv.Itoa(rand.Intn(max-min+1) + min)  + ")-" + strconv.Itoa(i)
	message.Status = true
	res, _ := json.Marshal(message)
	msg := kafka.Message{
			Key:    []byte(key),
			Value:  []byte(string(res)),
	}

	log.Println("----------------------------------------")
	log.Printf("==> Message Topic %s Key %s Value %s \n", p.configurations.KafkaConfig.Topic ,string(msg.Key), string(msg.Value))
	log.Println("----------------------------------------")

	producer := p.producer
	deliveryChan := make(chan kafka.Event)

	err := producer.Produce(&kafka.Message{
							TopicPartition: kafka.TopicPartition{Topic: &p.configurations.KafkaConfig.Topic, Partition: kafka.PartitionAny}, 
							Value: 	[]byte(res), 
							Headers:  []kafka.Header{{Key: "key", Value: []byte(key)}},
							},deliveryChan)
	if err != nil {
		log.Printf("Failed to producer message: %s\n", err)
		os.Exit(1)
	}
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	close(deliveryChan)
}