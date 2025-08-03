package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"casino/boundary/dto"
	"casino/boundary/logging"
	"casino/boundary/usecase"
	"casino/utils"

	"github.com/segmentio/kafka-go"
)

type KafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type KafkaConsumer struct {
	reader  KafkaReader
	useCase usecase.TransactionUseCase
	logger  logging.Logger
}

type TransactionMessage struct {
	ID              string `json:"id"`
	UserID          string `json:"user_id"`
	TransactionType string `json:"transaction_type"`
	Amount          uint   `json:"amount"`
}

func NewKafkaConsumer(brokers []string, topic string, useCase usecase.TransactionUseCase, logger logging.Logger) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        "casino-transaction-consumer",
		CommitInterval: 0,
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
				kc.logger.Error(ctx, fmt.Errorf("failed to read message: %w", err))
				continue
			}

			var transactionMsg TransactionMessage
			if err := json.Unmarshal(message.Value, &transactionMsg); err != nil {
				kc.logger.Error(ctx, fmt.Errorf("failed to unmarshal message: %w", err))
				_ = kc.reader.CommitMessages(ctx, message)
				continue
			}

			createDto := &dto.CreateTransactionDTO{
				ID:              transactionMsg.ID,
				UserID:          transactionMsg.UserID,
				TransactionType: transactionMsg.TransactionType,
				Amount:          transactionMsg.Amount,
			}

			maxRetries := 3
			var processErr error

			for attempt := 1; attempt <= maxRetries; attempt++ {

				processErr = kc.useCase.ProcessTransaction(createDto)
				
				if utils.IsDatabaseConnectionError(processErr) {
					if attempt < maxRetries {
						time.Sleep(3 * time.Second)
						continue
					}
				}
				break
			}

			if processErr != nil {
				kc.logger.Error(ctx, processErr)
			} else {
				kc.logger.Info(ctx, "trans " + createDto.ID + " saved")
			}

			_ = kc.reader.CommitMessages(ctx, message)
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}
