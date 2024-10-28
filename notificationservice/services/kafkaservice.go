package Services

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type KafkaService struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func NewKafkaService(brokers []string) (*KafkaService, error) {
	// Configure Kafka producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	// Create producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	// Create consumer
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaService{
		producer: producer,
		consumer: consumer,
	}, nil
}

func (ks *KafkaService) PublishMessage(topic string, message interface{}) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonMessage),
	}

	_, _, err = ks.producer.SendMessage(msg)
	return err
}

func (ks *KafkaService) Subscribe(topic string, handler func([]byte)) error {
	partitionConsumer, err := ks.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	go func() {
		for message := range partitionConsumer.Messages() {
			handler(message.Value)
		}
	}()

	return nil
}
