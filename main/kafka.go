package main

import (
  "github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaRepository struct {
	Broker   string
	producer *kafka.Producer
	channel  chan kafka.Event
}

func NewKafkaRepository(broker string) (Repository, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return KafkaRepository{}, err
	}
	deliveryChan := make(chan kafka.Event)

	kr := KafkaRepository{
		Broker:   broker,
		producer: p,
		channel:  deliveryChan,
	}

	return kr, nil
}

func (r KafkaRepository) Push(topic string, message []byte) error {
	err := r.producer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: message}, r.channel)
	if err != nil {
		return err
	}

	e := <-r.channel
	m := e.(*kafka.Message)

	return m.TopicPartition.Error
}

type KafkaRepositoryStats struct {
	Status string
}

func (r KafkaRepository) Stat() interface{} {
	s := KafkaRepositoryStats{
		"Ok",
	}

	return s
}

func (r KafkaRepository) Close() {
	r.producer.Close()
}

func (r KafkaRepository) Health() error {
	//TODO: Check health
	return nil
}
