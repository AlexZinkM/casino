package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"casino/adapter/handler"
	domainusecases "casino/domain/usecase"
	"casino/infra/kafka"
	infralogging "casino/infra/logging"
	"casino/infra/repository"
	"casino/infra/restserver/nethttp"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "casino/docs"
)

// @title Casino Transaction Management API
// @version 1.0
// @description A clean architecture implementation of a casino transaction management system
// @host localhost:8080
func main() {
	asyncLogger := infralogging.NewAsyncLogger("casino")
	simpleLogger := &infralogging.SimpleLogger{}
	asyncLogger.Register(simpleLogger)
	defer asyncLogger.Close()

	dsn := "host=localhost user=login password=password dbname=casino_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		asyncLogger.Error(context.Background(), err)
		log.Fatal("Failed to connect to database:", err)
	}

	transactionRepo := repository.NewPostgresTransactionRepository(db)
	transactionUseCase := domainusecases.NewTransactionUseCaseImpl(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase, asyncLogger)

	kafkaConsumer := kafka.NewKafkaConsumer(
		[]string{"localhost:9092"},
		"casino-transactions-stream",
		transactionUseCase,
		asyncLogger,
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go kafkaConsumer.Start(ctx)

	server := nethttp.NewNetHttpServer()

	server.RegisterPublicRoute("GET", "/transactions", transactionHandler.GetAllTransactions, asyncLogger)
	server.RegisterPublicRoute("GET", "/transactions/user", transactionHandler.GetUserTransactions, asyncLogger)
	server.RegisterSwaggerRoutes()

	asyncLogger.Info(context.Background(), "Server starting on port 8080")
	if err := server.Start(":8080"); err != nil {
		asyncLogger.Error(context.Background(), err)
		log.Fatal("Server error:", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	asyncLogger.Info(context.Background(), "Shutting down server")

	if err := kafkaConsumer.Close(); err != nil {
		asyncLogger.Error(context.Background(), err)
		log.Fatal("Kafka consumer close error:", err)
	}
}
