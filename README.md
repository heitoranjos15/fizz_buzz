# FizzBuzz API

A customizable FizzBuzz API built with Go and Gin.

## Getting Started

### Build & Run

```sh
make build
make run

The server will start on `localhost:8080`.

### Run Tests

```sh
make test
```

## API Endpoints

### POST /fizzbuzz

Generate a FizzBuzz sequence with custom parameters.

**Example:**
```sh
curl --location --request POST 'localhost:8080/fizzbuzz?multiples=3%2C5&words=teste%2CFizz&limit=5'
```

### GET /stats

Get usage statistics.

**Example:**
```sh
curl --location 'localhost:8080/stats'
```
