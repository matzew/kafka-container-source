package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	sarama "github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/knative/eventing-sources/pkg/kncloudevents"
)

type consumerGroupHandler struct {
	ceClient client.Client
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		//go post(msg.Value, c)

		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				Type:   "kafka-event",
				Source: *types.ParseURLRef(topic),
			}.AsV02(),
			Data: msg.Value,
		}
		go h.ceClient.Send(context.TODO(), event)
		sess.MarkMessage(msg, "")
	}
	return nil
}

var (
	sink      string
	topic     string
	groupID   string
	bootstrap string
)

func init() {
	flag.StringVar(&sink, "sink", "", "the host url to send messages to")
	flag.StringVar(&topic, "topic", "", "the topic to read messages from")
	flag.StringVar(&groupID, "groupId", uuid.New().String(), "the consumer group")
	flag.StringVar(&bootstrap, "bootstrap", "localhost:9092", "Apache Kafka bootstrap servers")
}

func main() {

	flag.Parse()

	brokers := []string{bootstrap}
	topics := []string{topic}

	// kafka
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kafkaConfig.Version = sarama.V2_0_0_0
	kafkaConfig.Consumer.Return.Errors = true

	// Start with a client
	client, err := sarama.NewClient(brokers, kafkaConfig)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Close() }()

	// init consumer group
	group, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	c, err := kncloudevents.NewDefaultClient(sink)
	if err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		handler := consumerGroupHandler{
			ceClient: c,
		}

		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
