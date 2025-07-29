package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"casino/adapter/handler"
	domainusecases "casino/domain/usecase"
	"casino/infra/kafka"
	infralogging "casino/infra/logging"
	"casino/infra/middleware"
	"casino/infra/repository"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Casino Transaction Management API
// @version 1.0
// @description A clean architecture implementation of a casino transaction management system
// @host localhost:8080
// @BasePath /
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

	sql, err := os.ReadFile("001_init.sql")
	if err != nil {
		asyncLogger.Error(context.Background(), err)
		log.Fatal("Failed to read SQL file:", err)
	}

	if err := db.Exec(string(sql)).Error; err != nil {
		asyncLogger.Error(context.Background(), err)
		log.Fatal("Failed to execute SQL:", err)
	}

	transactionRepo := repository.NewPostgresTransactionRepository(db)
	transactionUseCase := domainusecases.NewTransactionUseCaseImpl(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase, asyncLogger)

	kafkaConsumer := kafka.NewKafkaConsumer(
		[]string{"localhost:9092"},
		"casino-transactions",
		transactionUseCase,
		asyncLogger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go kafkaConsumer.Start(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", infra.LoggingMiddleware(transactionHandler.GetAllTransactions, asyncLogger))
	mux.HandleFunc("/transactions/user", infra.LoggingMiddleware(transactionHandler.GetUserTransactions, asyncLogger))
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  2 * time.Second,
	}

	go func() {
		asyncLogger.Info(context.Background(), "Server starting on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			asyncLogger.Error(context.Background(), err)
			log.Fatal("Server error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	asyncLogger.Info(context.Background(), "Shutting down server")
	if err := server.Shutdown(context.Background()); err != nil {
		asyncLogger.Error(context.Background(), err)
		log.Fatal("Server shutdown error:", err)
	}

	if err := kafkaConsumer.Close(); err != nil {
		asyncLogger.Error(context.Background(), err)
		log.Fatal("Kafka consumer close error:", err)
	}
}
