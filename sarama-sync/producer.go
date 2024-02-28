package sarama_sync

import (
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewProducer returns a new Sarama async producer.
func NewProducer(brokers string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Producer.Return.Successes = true
	config.Producer.Flush.Frequency = time.Duration(100) * time.Millisecond
	sarama.MaxRequestSize = 999000

	log.Infof("Connecting to %s", brokers)
	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		log.WithError(err).Panic("Unable to start the producer")
	}
	return producer
}

// Prepare returns a function that can be used during the benchmark as it only
// performs the sending of messages, checking that the sending was successful.
func Prepare(producer sarama.SyncProducer, message []byte, numMessages int) func() {
	log.Infof("Preparing to send message of %d bytes %d times", len(message), numMessages)

	return func() {
		for j := 0; j < numMessages; j++ {
			_, _, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic:     viper.GetString("kafka.topic"),
				Partition: kafka.PartitionAny,
				Value:     sarama.ByteEncoder(message),
			})
			if err != nil {
				log.WithError(err).Panic("Unable to deliver the message")
			}

		}
	}
}
