# Event-Driven Architecture with RabbitMQ

This project demonstrates an **Event-Driven Architecture** (EDA) using **RabbitMQ** as the message broker. It showcases the flow of events between different microservices, promoting decoupled and scalable systems.

---

## Table of Contents

1. [Overview](#overview)
2. [Project Architecture](#project-architecture)
3. [Requirements](#requirements)
4. [Installation](#installation)
5. [Running the Services](#running-the-services)
6. [How It Works](#how-it-works)
7. [Folder Structure](#folder-structure)
8. [Tech Stack](#tech-stack)
9. [Endpoints & Message Flow](#endpoints--message-flow)
10. [Next Steps](#next-steps)

---

## Overview

The project implements a distributed system where services communicate via events using RabbitMQ. This architecture enables:

- Loose coupling between services
- Scalability and extensibility
- Asynchronous communication

---

## Project Architecture

### Key Components:
- **Producer Service**: Publishes events to RabbitMQ.
- **RabbitMQ**: Message broker that routes events to appropriate queues.
- **Consumer Service**: Listens to queues and processes events.

**Workflow:**
1. Producer sends events to RabbitMQ.
2. RabbitMQ routes events to relevant queues.
3. Consumers process events from queues asynchronously.

### High-Level Diagram:
```
[Producer Service] ---> [RabbitMQ Broker] ---> [Consumer Services]
```

---

## Requirements

Ensure the following tools are installed:
- [Go 1.20+](https://golang.org/)
- [RabbitMQ](https://www.rabbitmq.com/)
- Docker & Docker Compose (optional, for local RabbitMQ setup)

---

## Installation

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/event-driven-rabbitmq.git
cd event-driven-rabbitmq
```

### 2. Start RabbitMQ
Use Docker Compose for local setup:

```bash
docker-compose up -d
```

RabbitMQ will be available at:
- **Web UI**: `http://localhost:15672`
- Default credentials: `guest/guest`

### 3. Install Dependencies
```bash
go mod tidy
```

---

## Running the Services

### 1. Start Producer Service
```bash
go run ./producer/main.go
```

### 2. Start Consumer Service
```bash
go run ./consumer/main.go
```

---

## How It Works

- **Producer Service**
  - Sends events/messages to a RabbitMQ exchange.
- **RabbitMQ**
  - Routes the messages to one or more queues based on the routing key and exchange type.
- **Consumer Service**
  - Listens to the queue(s) and processes incoming messages asynchronously.

**Exchange Type**: Direct/Topic/Fanout can be configured based on use cases.

---

## Folder Structure
```plaintext
.
|-- producer/          # Producer service code
|   |-- main.go        # Publishes events to RabbitMQ
|-- consumer/          # Consumer service code
|   |-- main.go        # Processes messages from RabbitMQ
|-- config/            # Configuration files
|-- docker-compose.yml # Local RabbitMQ setup
|-- README.md          # Project documentation
|-- go.mod             # Go dependencies
|-- go.sum             # Go checksums
```

---

## Tech Stack

- **Golang**: Core programming language
- **RabbitMQ**: Message broker for event routing
- **Docker**: Containerized RabbitMQ setup

---

## Endpoints & Message Flow

### Producer Service
- **Publishes**: Events/messages to RabbitMQ exchange.
- **Message Structure**:
```json
{
  "event": "user_created",
  "data": {
    "id": "12345",
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### RabbitMQ Configuration
- **Exchange**: `events_exchange`
- **Queues**: 
  - `user_events`
- **Routing Key**: `user.*`

### Consumer Service
- **Consumes**: Messages from the `user_events` queue.
- **Processing Logic**:
  - Logs the message or triggers further processing.

---

## Next Steps

- Implement more exchange types (Topic, Fanout) for advanced routing.
- Add retry mechanisms for failed messages.
- Integrate monitoring tools for RabbitMQ queues.

---

**Contributions Welcome!**
Feel free to fork and submit a pull request.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

Happy Coding! 🚀
