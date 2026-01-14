# Graceful Shutdown in Go

A practical example demonstrating graceful shutdown of an HTTP server in Go. This implementation ensures that in-progress requests are completed before the server terminates, while also handling timeout scenarios.

## Overview

This example shows how to properly shut down an HTTP server in Go by:
- Listening for OS signals (SIGINT, SIGTERM)
- Allowing in-progress requests to complete
- Implementing a shutdown timeout to prevent indefinite waiting
- Gracefully handling timeout scenarios with fallback to force close

## Features

- **Signal Handling**: Listens for `SIGINT` (Ctrl+C) and `SIGTERM` signals
- **Graceful Shutdown**: Uses `http.Server.Shutdown()` to allow in-progress requests to complete
- **Timeout Protection**: Configurable shutdown timeout prevents the server from hanging indefinitely
- **Fallback Mechanism**: If shutdown times out, the server is forcefully closed

## How It Works

The server includes a `/slow` endpoint that simulates a long-running request (8 seconds). When a shutdown signal is received:

1. The server stops accepting new requests
2. In-progress requests are allowed to complete
3. If all requests complete within the timeout period, the server exits gracefully
4. If the timeout is exceeded, the server is forcefully closed

## Running the Example

### Start the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

### Test the Graceful Shutdown

1. In one terminal, start the server:
   ```bash
   go run main.go
   ```

2. In another terminal, make a request to the slow endpoint:
   ```bash
   curl http://localhost:8080/slow
   ```

3. While the request is processing, send a shutdown signal (Ctrl+C) to the server

4. Observe that the request completes before the server exits

### Run Tests

```bash
go test -v
```

The test suite includes:
- **TestServerGracefulShutdown**: Verifies that in-progress requests complete during shutdown
- **TestServerTimeoutDuringShutdown**: Verifies that the server handles timeout scenarios correctly

## Code Structure

- `main.go`: Contains the server implementation and graceful shutdown logic
- `main_test.go`: Test cases for graceful shutdown scenarios

### Key Components

- `createServer()`: Creates an HTTP server with a `/slow` endpoint
- `runServer()`: Manages server lifecycle, signal handling, and graceful shutdown
- `main()`: Entry point that initializes and runs the server

## Configuration

The shutdown timeout is currently set to 3 seconds in `main()`:

```go
runServer(context.Background(), server, 3*time.Second)
```

You can adjust this value based on your application's requirements.

## Learning Resources

This example was coded along with the following YouTube videos:
- [Graceful Shutdown](https://www.youtube.com/watch?v=UPVSeZXBTxI&list=PL10piHcP2kVJOxO18iPHsq8IqTArbvFe0&index=43) - Implementation of graceful shutdown
- [Integration Tests](https://www.youtube.com/watch?v=J9yHXJC8aBg&list=PL10piHcP2kVJOxO18iPHsq8IqTArbvFe0&index=44) - Writing integration tests for graceful shutdown

## Requirements

- Go 1.25.4 or later

## License

This is an educational example. Feel free to use it as a reference for implementing graceful shutdown in your own Go applications.
