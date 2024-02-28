package kafkago_batch

import (
	"context"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var batchSize = 1000

// NewProducer returns a new kafkago writer.
func NewProducer(brokers string) *kafkago.Writer {
	return kafkago.NewWriter(kafkago.WriterConfig{
		Brokers:       []string{brokers},
		Topic:         viper.GetString("kafka.topic"),
		Balancer:      &kafkago.Hash{},
		BatchTimeout:  time.Duration(100) * time.Millisecond,
		QueueCapacity: 10000,
		BatchSize:     1000000,
		// Async doesn't allow us to know if message has been successfully sent to Kafka.
		// Async:         true,
	})
}

// Prepare returns a function that can be used during the benchmark as it only
// performs the sending of messages.
func Prepare(writer *kafkago.Writer, message []byte, numMessages int) func() {
	log.Infof("Preparing to send message of %d bytes %d times", len(message), numMessages)
	return func() {
		var batch []kafkago.Message
		for j := 0; j < numMessages; j++ {
			batch = append(batch, kafkago.Message{Value: message})
			if (j+1)%batchSize == 0 || j == numMessages-1 {
				err := writer.WriteMessages(context.Background(), batch...)
				if err != nil {
					log.WithError(err).Panic("Unable to deliver the message")
				}
				batch = nil
			}
		}
	}
}
