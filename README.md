# Casino Transaction Management System

I completed the test assignment as requested. Kafka is handled asynchronously. I also added Swagger, a simple middleware, and basic asynchronous logging to make the project look more polished and structured.

The menu was required to be covered with tests by 85%. As I understand it, this refers not to total code coverage, but specifically to covering parts that contain actual business logic.

## Architecture

The system follows Clean Architecture principles with the following structure:

- **Domain Layer**: Contains entities and use case implementations
- **Boundary Layer**: Contains interfaces, DTOs, and repository models
- **Adapter Layer**: Contains HTTP handlers and request/response structs
- **Infrastructure Layer**: Contains PostgreSQL repository and Kafka consumer implementations

## Prerequisites

- Go 1.24.4 or higher


### PostgreSQL Setup
1. Install PostgreSQL 12 or higher
2. Create a database named `casino_db`
3. Use the following credentials:
   - Username: `login`
   - Password: `password`
   - Database: `casino_db`
   - Host: `localhost`
   - Port: `5432`

### Kafka Setup
1. Install Kafka 2.8 or higher
2. Start Zookeeper and Kafka brokers
3. Create topic: `casino-transactions-stream`


## Features

- Clean Architecture implementation
- SOLID principles adherence
- Asynchronous message processing with Kafka
- PostgreSQL database storage
- RESTful API with JSON responses
- Comprehensive unit and integration tests
- Transaction filtering by user and type
- Proper error handling and logging
