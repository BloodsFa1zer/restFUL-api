# RESTful API

## Description

This project is a RESTful API built with Go (Golang). It provides endpoints to perform CRUD operations on various resources. It serves as a backend application for managing data and exposing it through a well-defined API.

## Features

- CRUD operations for various resources
- RESTful design principles
- JSON-based API responses
- Integration with a SQLite database
- Docker support for containerization

## Installation

### Prerequisites

- Go
- SQLite
- Redis
- Docker
- JWT Authorization
- Echo

### Steps

1. Clone the repository:
    ```sh
    git clone https://github.com/BloodsFa1zer/restFUL-api.git
    cd restFUL-api
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Set up the SQLite database:
    - Create a database named `restful_api`.
    - Configure your `config.yaml` with the appropriate database connection details.

4. Run the application:
    ```sh
    go run main.go
    ```

### Docker Setup

1. Build the Docker image:
    ```sh
    docker build -t restful-api .
    ```

2. Run the Docker container:
    ```sh
    docker run -p 8080:8080 restful-api
    ```

## Usage

### Endpoints

- **GET** `/resources` - List all resources
- **GET** `/resources/{id}` - Get a resource by ID
- **POST** `/resources` - Create a new resource
- **PUT** `/resources/{id}` - Update a resource by ID
- **DELETE** `/resources/{id}` - Delete a resource by ID

### Example Requests

- Get all resources:
    ```sh
    curl -X GET http://localhost:8080/resources
    ```

- Create a new resource:
    ```sh
    curl -X POST http://localhost:8080/resources -d '{"name": "New Resource"}' -H "Content-Type: application/json"
    ```

## Configuration

The application configuration is managed through a `config.yaml` file. Ensure that you set the appropriate values for your environment.

```yaml
server:
  port: 8080

database:
  DatabasePath: put_path_to_database_here
  DatabaseName: put_database_name_here
  SigningKey: put_signing_key_for_jwt_here
  RedisAddr: put_your_redis_address_here
  RedisPassword: put_your_redis_password_here
  RedisDB: put_your_redis_DB_here
