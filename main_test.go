package main

import (
	"flag"
	"fmt"
	kafkago_batch "github.com/gguridi/benchmark-kafka-go-clients/kafkago-batch"
	sarama_sync "github.com/gguridi/benchmark-kafka-go-clients/sarama-sync"

	"github.com/gguridi/benchmark-kafka-go-clients/confluent"
	"github.com/gguridi/benchmark-kafka-go-clients/kafkago"
	"github.com/gguridi/benchmark-kafka-go-clients/sarama"
	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var _ = Describe("Benchmarks", func() {

	It("checks lag", func() {
		flag.Parse()
		consumer := confluent.NewConsumer(viper.GetString("kafka.brokers"), true)
		confluent.PreparePoll(consumer, 1)()
		Expect(confluent.Lag(consumer)).To(Equal(NumMessages))
	})

	Measure("producer", func(b Benchmarker) {
		flag.Parse()
		name := fmt.Sprintf("%s producing %d messages of %d bytes size", Library, NumMessages, MessageSize)
		switch Library {
		case "confluent":
			producer := confluent.NewProducer(viper.GetString("kafka.brokers"))
			process := confluent.Prepare(producer, GenMessage(), NumMessages)
			b.Time(name, func() {
				process()
				<-confluent.Done
			})
			producer.Close()
			break
		case "sarama":
			producer := sarama.NewProducer(viper.GetString("kafka.brokers"))
			process := sarama.Prepare(producer, GenMessage(), NumMessages)
			b.Time(name, func() {
				process()
				<-sarama.Done
			})
			if err := producer.Close(); err != nil {
				log.WithError(err).Panic("Unable to close the producer")
			}
			break
		case "sarama-sync":
			producer := sarama_sync.NewProducer(viper.GetString("kafka.brokers"))
			process := sarama_sync.Prepare(producer, GenMessage(), NumMessages)
			b.Time(name, func() {
				process()
			})
			if err := producer.Close(); err != nil {
				log.WithError(err).Panic("Unable to close the producer")
			}
			break
		case "kafkago":
			writer := kafkago.NewProducer(viper.GetString("kafka.brokers"))
			process := kafkago.Prepare(writer, GenMessage(), NumMessages)
			b.Time(name, func() {
				process()
			})
			if err := writer.Close(); err != nil {
				log.WithError(err).Panic("Unable to close the producer")
			}
			break
		case "kafkago-batch":
			writer := kafkago_batch.NewProducer(viper.GetString("kafka.brokers"))
			process := kafkago_batch.Prepare(writer, GenMessage(), NumMessages)
			b.Time(name, func() {
				process()
			})
			if err := writer.Close(); err != nil {
				log.WithError(err).Panic("Unable to close the producer")
			}
			break
		default:
			log.Panicf("Unable to find the libray %+v", Library)
		}
	}, 3)

	Context("prepopulate", func() {

		var (
			initialised = false
		)

		BeforeEach(func() {
			if !initialised {
				log.Infof("Prepopulating Kafka with %d messages of %d bytes", NumMessages, MessageSize)
				producer := sarama.NewProducer(viper.GetString("kafka.brokers"))
				process := sarama.Prepare(producer, GenMessage(), int(float64(NumMessages)*1.5))
				process()
				<-sarama.Done
				producer.Close()
				log.Infof("Finished prepopulating Kafka")
			}
			initialised = true
		})

		Measure("consumer", func(b Benchmarker) {
			flag.Parse()
			name := fmt.Sprintf("%s consuming %d messages", Library, NumMessages)
			switch Library {
			case "confluent-poll":
				consumer := confluent.NewConsumer(viper.GetString("kafka.brokers"), true)
				process := confluent.PreparePoll(consumer, NumMessages)
				b.Time(name, func() {
					process()
				})
				consumer.Close()
				break
			case "confluent-channel":
				consumer := confluent.NewConsumer(viper.GetString("kafka.brokers"), false)
				process := confluent.PrepareChannel(consumer, NumMessages)
				b.Time(name, func() {
					process()
				})
				consumer.Close()
				break
			case "sarama":
				consumer := sarama.NewConsumer(viper.GetString("kafka.brokers"), NumMessages)
				process := sarama.PrepareConsume(consumer)
				b.Time(name, func() {
					process()
				})
				if err := consumer.Client.Close(); err != nil {
					log.WithError(err).Panic("Unable to close the consumer")
				}
				break
			case "kafkago":
				reader := kafkago.NewReader(viper.GetString("kafka.brokers"))
				process := kafkago.PrepareReader(reader, NumMessages)
				b.Time(name, func() {
					process()
				})
				if err := reader.Close(); err != nil {
					log.WithError(err).Panic("Unable to close the consumer")
				}
				break
			default:
				log.Panicf("Unable to find the libray %+v", Library)
			}
		}, 1)
	})
})
