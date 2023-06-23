package broker

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
)

const (
	OrderCreatedTopic string = "order_created"
)

func IniProducer(brokerAddr string) sarama.SyncProducer {
	brokerCfg := sarama.NewConfig()
	brokerCfg.Producer.RequiredAcks = sarama.WaitForAll
	brokerCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{brokerAddr}, brokerCfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return producer
}
