package kafka

import (
	"context"
	"encoding/json"

	"casino/boundary/dto"
	"casino/boundary/logging"
	"casino/boundary/usecase"

	"github.com/segmentio/kafka-go"
)

type KafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type KafkaConsumer struct {
	reader  KafkaReader
	useCase usecase.TransactionUseCase
	logger  logging.Logger
}

type TransactionMessage struct {
	UserID          string `json:"user_id"`
	TransactionType string `json:"transaction_type"`
	Amount          uint   `json:"amount"`
}

func NewKafkaConsumer(brokers []string, topic string, useCase usecase.TransactionUseCase, logger logging.Logger) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "casino-transaction-consumer",
	})

	return &KafkaConsumer{
		reader:  reader,
		useCase: useCase,
		logger:  logger,
	}
}

func (kc *KafkaConsumer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			kc.logger.Info(ctx, "Kafka consumer stopping due to context cancellation")
			return
		default:
			message, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				kc.logger.Error(ctx, err)
				continue
			}

			var transactionMsg TransactionMessage
			if err := json.Unmarshal(message.Value, &transactionMsg); err != nil {
				kc.logger.Error(ctx, err)
				continue
			}

			createDto := &dto.CreateTransactionDTO{
				UserID:          transactionMsg.UserID,
				TransactionType: transactionMsg.TransactionType,
				Amount:          transactionMsg.Amount,
			}

			if err := kc.useCase.ProcessTransaction(createDto); err != nil {
				kc.logger.Error(ctx, err)
				continue
			}
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}
