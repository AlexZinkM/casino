# Casino Transaction Management System

I completed the test assignment as requested. Kafka is handled asynchronously, with a simple retry mechanism for db connection issues without DQL. I also added Swagger, a simple middleware, and basic asynchronous logging to make the project look more polished and structured.

The menu was required to be covered with tests by 85%. As I understand it, this refers not to total code coverage, but specifically to covering parts that contain actual business logic.

## Architecture

The system follows Clean Architecture principles with the following structure:

- **Domain Layer**: Contains entities and use case implementations
- **Boundary Layer**: Contains interfaces, DTOs, and repository models
- **Adapter Layer**: Contains HTTP handlers and request/response structs
- **Infrastructure Layer**: Contains PostgreSQL repository and Kafka consumer implementations

## Prerequisites

- Go 1.24.4 or higher


### Setup
   For your convenience, I leave the docker-compose I used for testing the app.


## Features

- Clean Architecture implementation
- Asynchronous message processing with Kafka
- PostgreSQL database storage
- RESTful API with JSON responses
- Comprehensive unit and integration tests
- Transaction filtering by user and type
- Proper error handling and logging
