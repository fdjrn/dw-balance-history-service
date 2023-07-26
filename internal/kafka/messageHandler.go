package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/fdjrn/dw-balance-history-service/internal/db/entity"
	"github.com/fdjrn/dw-balance-history-service/internal/handlers/consumer"
	"github.com/fdjrn/dw-balance-history-service/internal/kafka/topic"
	"github.com/fdjrn/dw-balance-history-service/internal/utilities"
)

func HandleMessages(message *sarama.ConsumerMessage) {
	var (
		handler = consumer.NewTransactionHandler()
		history = new(entity.BalanceHistory)
		//err       error
		transType string
	)

	utilities.Log.SetPrefix("[CONSUMER] ")

	switch message.Topic {
	case topic.DeductResult:
		transType = "deduct"
	case topic.TopUpResult:
		transType = "topup"
	case topic.DistributionResult:
		transType = "merchant"
	case topic.DistributionResultMembers:
		transType = "member"
	default:
		utilities.Log.Println("| Unknown topic message")
		return
	}

	history, err := handler.DoHandleTransaction(message)
	if err != nil {
		utilities.Log.Printf("| failed to create %s history: %s\n", transType, err.Error())
	} else {
		utilities.Log.Printf("| %s transaction history with receipt number %s , has been created successfully\n",
			transType,
			history.ReceiptNumber)
	}

}
