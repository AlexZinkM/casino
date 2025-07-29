# Casino Transaction Management System

A clean architecture implementation of a casino transaction management system using Go, PostgreSQL, and Kafka.

## Architecture

The system follows Clean Architecture principles with the following structure:

- **Domain Layer**: Contains entities and use case implementations
- **Boundary Layer**: Contains interfaces, DTOs, and repository models
- **Adapter Layer**: Contains HTTP handlers and request/response structs
- **Infrastructure Layer**: Contains PostgreSQL repository and Kafka consumer implementations

## Prerequisites

- Go 1.24.4 or higher
- PostgreSQL 12 or higher
- Kafka 2.8 or higher

## Database Setup

1. Create a PostgreSQL database named `casino_db`
2. Use the following credentials:
   - Username: `login`
   - Password: `password`
   - Database: `casino_db`
   - Host: `localhost`
   - Port: `5432`

## Kafka Setup

1. Start Kafka with a topic named `casino-transactions`
2. The system expects Kafka to be running on `localhost:9092`

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Application

1. Start the application:
   ```bash
   go run main.go
   ```

The server will start on port 8080.

## API Endpoints

### Get All Transactions
```
GET /transactions?transaction_type=bet
```

Query Parameters:
- `transaction_type` (optional): Filter by transaction type (`bet` or `win`)

### Get User Transactions
```
GET /transactions/user?user_id=1&transaction_type=win
```

Query Parameters:
- `user_id` (required): User ID to filter transactions
- `transaction_type` (optional): Filter by transaction type (`bet` or `win`)

## Kafka Message Format

The system consumes messages from the `casino-transactions` topic with the following JSON format:

```json
{
  "user_id": 1,
  "transaction_type": "bet",
  "amount": 100.0
}
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Project Structure

```
casino/
├── domain/
│   ├── entities/
│   │   └── transaction.go
│   └── usecases/
│       ├── transaction_usecase_impl.go
│       └── transaction_usecase_test.go
├── boundary/
│   ├── usecases/
│   │   └── transaction_usecase.go
│   ├── repositories/
│   │   └── transaction_repository.go
│   ├── dtos/
│   │   └── transaction_dto.go
│   └── repo_models/
│       └── transaction_model.go
├── adapter/
│   ├── handlers/
│   │   ├── transaction_handler.go
│   │   └── transaction_handler_test.go
│   └── structs/
│       └── transaction_structs.go
├── infra/
│   ├── repository/
│   │   ├── postgres_repository.go
│   │   └── postgres_repository_test.go
│   └── kafka/
│       └── kafka_consumer.go
├── main.go
├── go.mod
└── README.md
```

## Features

- Clean Architecture implementation
- SOLID principles adherence
- Domain-Driven Design (DDD)
- Asynchronous message processing with Kafka
- PostgreSQL database storage
- RESTful API with JSON responses
- Comprehensive unit and integration tests
- Transaction filtering by user and type
- Proper error handling and logging

## Database Schema

The system creates the following table:

```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    transaction_type VARCHAR(10) NOT NULL CHECK (transaction_type IN ('bet', 'win')),
    amount DECIMAL(10,2) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);
```

Indexes are created for:
- `user_id`
- `transaction_type`
- `timestamp` 