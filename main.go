package main

import (
	"errors"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need list of brokers")
	}
	brokers := os.Args[1:]

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	const topic = "test-1"

	replicaID, replicaRack, err := topicReplica(brokers, topic)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("replica ID", replicaID, "rack", replicaRack)

	cfg := sarama.NewConfig()
	cfg.RackID = replicaRack
	cfg.Version = sarama.V2_6_0_0
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewClient(brokers, cfg)
	if err != nil {
		log.Fatal(err)
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			key := uuid.New().String()

			_, _, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(key),
				Value: sarama.StringEncoder(key),
			})
			if err != nil {
				log.Fatal(err)
			}

			// log.Println("produced", key, pa, of)
		}
	}()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatal(err)
	}

	pc, err := consumer.ConsumePartition(topic, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	for range pc.Messages() {
	}
}

func topicReplica(brokers []string, topic string) (id int32, rack string, err error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_6_0_0

	client, err := sarama.NewClient(brokers, cfg)
	if err != nil {
		return 0, "", err
	}
	defer client.Close()

	leader, err := client.Leader(topic, 0)
	if err != nil {
		return 0, "", err
	}

	replicas, err := client.Replicas(topic, 0)
	if err != nil {
		return 0, "", err
	}

	for _, b := range client.Brokers() {
		for _, r := range replicas {
			if r != leader.ID() {
				return r, b.Rack(), nil
			}
		}
	}

	return 0, "", errors.New("no non-leader replicas found")
}
