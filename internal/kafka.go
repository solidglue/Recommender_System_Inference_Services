package internal

import (
	"context"
	"infer-microservices/internal/flags"
	"infer-microservices/internal/logs"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

var (
	kafkaWriter *kafka.Writer
	kafkaURL    string
	kafkaTopic  string
	kafkaGroup  string
)

type CallbackFunc func(string, string)

func init() {
	flagFactory := flags.FlagFactory{}
	flagKafka := flagFactory.CreateFlagKafka()
	kafkaURL = flagKafka.GetKafkaUrl()
	kafkaTopic = flagKafka.GetKafkaTopic()
	kafkaGroup = flagKafka.GetKafkaGroup()

}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e7, // 100MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset, //  kafka.FirstOffset,
	})
}

// 从消息队列里监听来自用户管理后台传来的信息
func KafkaConsumer(callback CallbackFunc) {
	reader := getKafkaReader(kafkaURL, kafkaTopic, kafkaGroup)
	defer reader.Close()

	index0 := 0
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			logs.Error(err)
			continue
		}

		//do something
		callback(string(m.Key), string(m.Value))

		if index0%100000 == 0 {
			logs.Info(" topic, partition, offset, time, key ,value = ", m.Topic, m.Partition, m.Offset, m.Time.Unix(), string(m.Key), string(m.Value))
		}

		index0 += 1
	}
}

func KafkaProducer(msgKey string, msgValue string) {
	kafkaWriter = getKafkaWriter(kafkaURL, kafkaTopic)
	defer kafkaWriter.Close()

	msg := kafka.Message{
		Key:   []byte(msgKey),
		Value: []byte(msgValue),
	}

	err := kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		logs.Error(err)
	}
}
