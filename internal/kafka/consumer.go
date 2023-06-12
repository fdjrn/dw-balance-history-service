package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/fdjrn/dw-balance-history-service/configs"
	"github.com/fdjrn/dw-balance-history-service/internal/db/entity"
	"github.com/fdjrn/dw-balance-history-service/internal/handlers"
	"log"
	"os"
	"strings"
	"sync"
)

type MessageConsumer struct {
	ready chan bool
}

var Logger = log.New(os.Stdout, "[CONSUMER] ", log.LstdFlags)

const (
	DeductTopic = "mdw.transaction.deduct.created"
	TopUpTopic  = "mdw.transaction.topup.created"
)

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *MessageConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			//log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			HandleMessages(message)

			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *MessageConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *MessageConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

func initConsumer() (sarama.ConsumerGroup, MessageConsumer, error) {
	splitBrokers := strings.Split(configs.MainConfig.Kafka.Brokers, ",")

	conf := configs.NewSaramaConfig()

	switch configs.MainConfig.Kafka.Consumer.Assignor {
	case "sticky":
		conf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	case "roundRobin":
		conf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	case "range":
		conf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", configs.MainConfig.Kafka.Consumer.Assignor)
	}

	if configs.MainConfig.Kafka.Consumer.Oldest {
		conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	/**
	 * Set up a new Sarama consumer group
	 */
	consumer := MessageConsumer{
		ready: make(chan bool),
	}

	//ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(splitBrokers, configs.MainConfig.Kafka.Consumer.ConsumerGroupName, conf)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	return client, consumer, err

}

func StartConsumer() {

	var err error

	client, consumer, err := initConsumer()
	if err != nil {
		log.Fatalln(err)
	}

	topicMsg := strings.Split(configs.MainConfig.Kafka.Consumer.ConsumerTopics, ",")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {

			if err = client.Consume(context.Background(), topicMsg, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			//if ctx.Err() != nil {
			//	return
			//}
			consumer.ready = make(chan bool)
		}
	}()
	// wait till the consumer has been set up
	<-consumer.ready

	log.Println("[INIT] consumer >> up and running!...")
	wg.Done()
}

func HandleMessages(message *sarama.ConsumerMessage) {

	var history = new(entity.BalanceHistory)
	var err error

	switch message.Topic {
	case DeductTopic:
		history, err = handlers.InsertDeductHistory(message)
		if err != nil {
			Logger.Println("failed to create deduct transaction history: ", err.Error())
			return
		}
		Logger.Printf("deduct transaction history with receipt number %s , has been created successfully\n",
			history.ReceiptNumber)

	case TopUpTopic:
		history, err = handlers.InsertTopUpHistory(message)
		if err != nil {
			Logger.Println("failed to create topup history: ", err.Error())
			return
		}

		Logger.Printf("topup transaction history with receipt number %s , has been created successfully\n",
			history.ReceiptNumber)

	default:
		Logger.Println("Unknown topic message")
		return
	}

}
