package main

import (
  "github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaRepository struct {
  Broker    string
  producer  *kafka.Producer
}


func NewKafkaRepository(broker string) (Repository, error) {
  p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

  if err != nil {
    return nil, err
  }

  kr := KafkaRepository{
    Broker : broker,
    producer : p,
  }

  return kr, nil
}

func (r KafkaRepository) Push(topic string, message []byte) error {
  r.producer.ProduceChannel() <- &kafka.Message{
    TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
    Value: message,
    }
}

type KafkaRepositoryStats struct {
  Status    string
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
