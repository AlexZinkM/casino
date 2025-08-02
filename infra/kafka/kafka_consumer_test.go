package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	boundarydto "casino/boundary/dto"
	"casino/utils"

	"github.com/segmentio/kafka-go"
)

type MockTransactionUseCase struct {
	processError error
	processCount int
}

func (m *MockTransactionUseCase) ProcessTransaction(dto *boundarydto.CreateTransactionDTO) error {
	m.processCount++
	return m.processError
}

func (m *MockTransactionUseCase) GetUserTransactions(userID string, filter *boundarydto.TransactionFilterDTO) ([]*boundarydto.TransactionDTO, error) {
	return nil, nil
}

func (m *MockTransactionUseCase) GetAllTransactions(filter *boundarydto.TransactionFilterDTO) ([]*boundarydto.TransactionDTO, error) {
	return nil, nil
}

type MockLogger struct {
	errorCount int
	infoCount  int
	errors     []error
	messages   []string
}

func (m *MockLogger) Error(ctx context.Context, errs ...error) {
	m.errorCount++
	m.errors = append(m.errors, errs...)
}

func (m *MockLogger) Info(ctx context.Context, messages ...string) {
	m.infoCount++
	m.messages = append(m.messages, messages...)
}

func TestNewKafkaConsumer(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}

	consumer := NewKafkaConsumer(
		[]string{"localhost:9092"},
		"test-topic",
		mockUseCase,
		mockLogger,
	)

	if consumer == nil {
		t.Error("Expected consumer to be created")
	}

	if consumer.useCase != mockUseCase {
		t.Error("Expected use case to be set")
	}

	if consumer.logger != mockLogger {
		t.Error("Expected logger to be set")
	}

	if consumer.reader == nil {
		t.Error("Expected reader to be created")
	}
}

func TestNewKafkaConsumer_WithMultipleBrokers(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}

	consumer := NewKafkaConsumer(
		[]string{"localhost:9092", "localhost:9093"},
		"test-topic",
		mockUseCase,
		mockLogger,
	)

	if consumer == nil {
		t.Error("Expected consumer to be created with multiple brokers")
	}
}

func TestKafkaConsumer_Close(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}

	consumer := NewKafkaConsumer(
		[]string{"localhost:9092"},
		"test-topic",
		mockUseCase,
		mockLogger,
	)

	err := consumer.Close()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = consumer.Close()
	if err != nil {
		t.Errorf("Expected no error on second close, got %v", err)
	}
}

func TestTransactionMessage_JSON(t *testing.T) {
	message := TransactionMessage{
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Errorf("Expected no error marshaling to JSON, got %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("Expected JSON data to be generated")
	}

	var unmarshaledMessage TransactionMessage
	err = json.Unmarshal(jsonData, &unmarshaledMessage)
	if err != nil {
		t.Errorf("Expected no error unmarshaling from JSON, got %v", err)
	}

	if unmarshaledMessage.UserID != message.UserID {
		t.Errorf("Expected UserID %s, got %s", message.UserID, unmarshaledMessage.UserID)
	}

	if unmarshaledMessage.TransactionType != message.TransactionType {
		t.Errorf("Expected TransactionType %s, got %s", message.TransactionType, unmarshaledMessage.TransactionType)
	}

	if unmarshaledMessage.Amount != message.Amount {
		t.Errorf("Expected Amount %d, got %d", message.Amount, unmarshaledMessage.Amount)
	}
}

func TestTransactionMessage_JSON_EdgeCases(t *testing.T) {
	testCases := []struct {
		name            string
		userID          string
		transactionType string
		amount          uint
	}{
		{"Empty UserID", "", "bet", 100},
		{"Empty TransactionType", "user123", "", 100},
		{"Zero Amount", "user123", "bet", 0},
		{"Large Amount", "user123", "win", 999999},
		{"Special Characters", "user-123_456", "bet", 100},
		{"Unicode Characters", "user-123_456-测试", "bet", 100},
		{"Very Long UserID", "very-long-user-id-that-exceeds-normal-length-limits-for-testing-purposes", "bet", 100},
		{"Maximum Amount", "user123", "win", 4294967295},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			message := TransactionMessage{
				UserID:          tc.userID,
				TransactionType: tc.transactionType,
				Amount:          tc.amount,
			}

			jsonData, err := json.Marshal(message)
			if err != nil {
				t.Errorf("Expected no error marshaling %s, got %v", tc.name, err)
			}

			var unmarshaledMessage TransactionMessage
			err = json.Unmarshal(jsonData, &unmarshaledMessage)
			if err != nil {
				t.Errorf("Expected no error unmarshaling %s, got %v", tc.name, err)
			}

			if unmarshaledMessage.UserID != message.UserID {
				t.Errorf("UserID mismatch in %s: %s != %s", tc.name, unmarshaledMessage.UserID, message.UserID)
			}

			if unmarshaledMessage.TransactionType != message.TransactionType {
				t.Errorf("TransactionType mismatch in %s: %s != %s", tc.name, unmarshaledMessage.TransactionType, message.TransactionType)
			}

			if unmarshaledMessage.Amount != message.Amount {
				t.Errorf("Amount mismatch in %s: %d != %d", tc.name, unmarshaledMessage.Amount, message.Amount)
			}
		})
	}
}

func TestTransactionMessage_JSON_InvalidData(t *testing.T) {
	invalidJSON := []byte(`{"user_id": "test", "transaction_type": "bet", "amount": "invalid"}`)

	var message TransactionMessage
	err := json.Unmarshal(invalidJSON, &message)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid JSON")
	}
}

func TestTransactionMessage_JSON_MissingFields(t *testing.T) {
	partialJSON := []byte(`{"user_id": "test"}`)

	var message TransactionMessage
	err := json.Unmarshal(partialJSON, &message)
	if err != nil {
		t.Errorf("Expected no error when unmarshaling partial JSON, got %v", err)
	}

	if message.UserID != "test" {
		t.Errorf("Expected UserID 'test', got %s", message.UserID)
	}

	if message.TransactionType != "" {
		t.Errorf("Expected empty TransactionType, got %s", message.TransactionType)
	}

	if message.Amount != 0 {
		t.Errorf("Expected Amount 0, got %d", message.Amount)
	}
}

func TestKafkaConsumer_ProcessTransaction_Success(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}

	consumer := NewKafkaConsumer(
		[]string{"localhost:9092"},
		"test-topic",
		mockUseCase,
		mockLogger,
	)

	message := TransactionMessage{
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
	}

	createDto := &boundarydto.CreateTransactionDTO{
		UserID:          message.UserID,
		TransactionType: message.TransactionType,
		Amount:          message.Amount,
	}

	err := consumer.useCase.ProcessTransaction(createDto)
	if err != nil {
		t.Errorf("Expected no error processing transaction, got %v", err)
	}

	if mockUseCase.processCount != 1 {
		t.Errorf("Expected process count 1, got %d", mockUseCase.processCount)
	}
}

func TestKafkaConsumer_ProcessTransaction_Error(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		processError: fmt.Errorf("processing error"),
	}
	mockLogger := &MockLogger{}

	consumer := NewKafkaConsumer(
		[]string{"localhost:9092"},
		"test-topic",
		mockUseCase,
		mockLogger,
	)

	message := TransactionMessage{
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
	}

	createDto := &boundarydto.CreateTransactionDTO{
		UserID:          message.UserID,
		TransactionType: message.TransactionType,
		Amount:          message.Amount,
	}

	err := consumer.useCase.ProcessTransaction(createDto)
	if err == nil {
		t.Error("Expected error processing transaction")
	}

	if mockUseCase.processCount != 1 {
		t.Errorf("Expected process count 1, got %d", mockUseCase.processCount)
	}
}

func TestKafkaConsumer_MessageProcessing(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}

	consumer := NewKafkaConsumer(
		[]string{"localhost:9092"},
		"test-topic",
		mockUseCase,
		mockLogger,
	)

	testCases := []struct {
		name            string
		transactionType string
		amount          uint
	}{
		{"Bet Transaction", "bet", 100},
		{"Win Transaction", "win", 200},
		{"Large Amount", "bet", 1000},
		{"Small Amount", "win", 1},
		{"Zero Amount", "bet", 0},
		{"Maximum Amount", "win", 4294967295},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			message := TransactionMessage{
				UserID:          utils.GenerateUUID(),
				TransactionType: tc.transactionType,
				Amount:          tc.amount,
			}

			createDto := &boundarydto.CreateTransactionDTO{
				UserID:          message.UserID,
				TransactionType: message.TransactionType,
				Amount:          message.Amount,
			}

			err := consumer.useCase.ProcessTransaction(createDto)
			if err != nil {
				t.Errorf("Expected no error processing %s, got %v", tc.name, err)
			}
		})
	}
}

func TestKafkaConsumer_ErrorHandling(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		processError: fmt.Errorf("simulated error"),
	}
	mockLogger := &MockLogger{}

	consumer := NewKafkaConsumer(
		[]string{"localhost:9092"},
		"test-topic",
		mockUseCase,
		mockLogger,
	)

	message := TransactionMessage{
		UserID:          "",
		TransactionType: "invalid",
		Amount:          0,
	}

	createDto := &boundarydto.CreateTransactionDTO{
		UserID:          message.UserID,
		TransactionType: message.TransactionType,
		Amount:          message.Amount,
	}

	err := consumer.useCase.ProcessTransaction(createDto)
	if err == nil {
		t.Error("Expected error processing invalid transaction")
	}
}

func TestKafkaConsumer_LoggerIntegration(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		processError: fmt.Errorf("test error"),
	}
	mockLogger := &MockLogger{}

	consumer := NewKafkaConsumer(
		[]string{"localhost:9092"},
		"test-topic",
		mockUseCase,
		mockLogger,
	)

	ctx := context.Background()

	consumer.logger.Error(ctx, fmt.Errorf("test error"))
	consumer.logger.Info(ctx, "test message")

	if mockLogger.errorCount != 1 {
		t.Errorf("Expected error count 1, got %d", mockLogger.errorCount)
	}

	if mockLogger.infoCount != 1 {
		t.Errorf("Expected info count 1, got %d", mockLogger.infoCount)
	}

	if len(mockLogger.errors) != 1 {
		t.Errorf("Expected 1 error logged, got %d", len(mockLogger.errors))
	}

	if len(mockLogger.messages) != 1 {
		t.Errorf("Expected 1 message logged, got %d", len(mockLogger.messages))
	}
}

func TestKafkaConsumer_Constructor_EdgeCases(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}

	testCases := []struct {
		name    string
		brokers []string
		topic   string
	}{
		{"Single Broker", []string{"localhost:9092"}, "test-topic"},
		{"Multiple Brokers", []string{"localhost:9092", "localhost:9093", "localhost:9094"}, "test-topic"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			consumer := NewKafkaConsumer(tc.brokers, tc.topic, mockUseCase, mockLogger)
			if consumer == nil {
				t.Error("Expected consumer to be created")
			}
			if consumer.useCase != mockUseCase {
				t.Error("Expected use case to be set")
			}
			if consumer.logger != mockLogger {
				t.Error("Expected logger to be set")
			}
		})
	}
}

func TestKafkaConsumer_MessageProcessing_EdgeCases(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}

	consumer := NewKafkaConsumer(
		[]string{"localhost:9092"},
		"test-topic",
		mockUseCase,
		mockLogger,
	)

	edgeCases := []struct {
		name            string
		userID          string
		transactionType string
		amount          uint
	}{
		{"Empty UserID", "", "bet", 100},
		{"Empty TransactionType", "user123", "", 100},
		{"Zero Amount", "user123", "bet", 0},
		{"Large Amount", "user123", "win", 999999},
		{"Special Characters in UserID", "user-123_456@test.com", "bet", 100},
		{"Unicode UserID", "user-测试-123", "win", 200},
		{"Very Long UserID", "very-long-user-id-that-exceeds-normal-length-limits-for-testing-purposes-and-should-still-work-correctly", "bet", 100},
		{"Maximum Amount", "user123", "win", 4294967295},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			message := TransactionMessage{
				UserID:          tc.userID,
				TransactionType: tc.transactionType,
				Amount:          tc.amount,
			}

			createDto := &boundarydto.CreateTransactionDTO{
				UserID:          message.UserID,
				TransactionType: message.TransactionType,
				Amount:          message.Amount,
			}

			err := consumer.useCase.ProcessTransaction(createDto)
			if err != nil {
				t.Errorf("Expected no error processing %s, got %v", tc.name, err)
			}
		})
	}
}

func TestKafkaConsumer_StructFields(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}

	consumer := NewKafkaConsumer(
		[]string{"localhost:9092"},
		"test-topic",
		mockUseCase,
		mockLogger,
	)

	if consumer.reader == nil {
		t.Error("Expected reader to be initialized")
	}

	if consumer.useCase == nil {
		t.Error("Expected use case to be initialized")
	}

	if consumer.logger == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestTransactionMessage_StructFields(t *testing.T) {
	message := TransactionMessage{
		UserID:          "test-user",
		TransactionType: "bet",
		Amount:          100,
	}

	if message.UserID != "test-user" {
		t.Errorf("Expected UserID 'test-user', got %s", message.UserID)
	}

	if message.TransactionType != "bet" {
		t.Errorf("Expected TransactionType 'bet', got %s", message.TransactionType)
	}

	if message.Amount != 100 {
		t.Errorf("Expected Amount 100, got %d", message.Amount)
	}
}

type MockKafkaReader struct {
	messages  [][]byte
	errOnRead error
	readIndex int
	closed    bool
}

func (m *MockKafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	if m.errOnRead != nil {
		return kafka.Message{}, m.errOnRead
	}
	if m.readIndex >= len(m.messages) {
		return kafka.Message{}, context.Canceled
	}
	msg := kafka.Message{Value: m.messages[m.readIndex]}
	m.readIndex++
	return msg, nil
}

func (m *MockKafkaReader) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return nil
}

func (m *MockKafkaReader) Close() error {
	m.closed = true
	return nil
}

func newTestKafkaConsumerWithMockReader(reader *MockKafkaReader, useCase *MockTransactionUseCase, logger *MockLogger) *KafkaConsumer {
	return &KafkaConsumer{
		reader:  reader,
		useCase: useCase,
		logger:  logger,
	}
}

func TestKafkaConsumer_Start_MockReader(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}

	validMsg := TransactionMessage{
		UserID:          "user1",
		TransactionType: "bet",
		Amount:          100,
	}
	validBytes, _ := json.Marshal(validMsg)
	invalidJSON := []byte(`{"user_id": "user2", "transaction_type": "bet", "amount": "bad"}`)

	reader := &MockKafkaReader{
		messages: [][]byte{validBytes, invalidJSON},
	}

	consumer := newTestKafkaConsumerWithMockReader(reader, mockUseCase, mockLogger)

	ctx, cancel := context.WithCancel(context.Background())

	go consumer.Start(ctx)

	time.Sleep(10 * time.Millisecond)
	cancel()

	consumer.Close()

	if mockUseCase.processCount != 1 {
		t.Errorf("Expected 1 valid message processed, got %d", mockUseCase.processCount)
	}

	if mockLogger.errorCount == 0 {
		t.Error("Expected at least one error to be logged for invalid JSON")
	}

	if !reader.closed {
		t.Error("Expected reader to be closed after Close")
	}
}
