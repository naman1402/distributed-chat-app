# Distributed Chat Application

A scalable, real-time chat application designed with distributed systems principles. This application demonstrates modern cloud-native architecture using microservices and event-driven communication.

## System Architecture

### Core Technologies & Their Purpose

1. **Golang (Backend)**
   - High performance and concurrent processing
   - Strong standard library for network programming
   - Efficient garbage collection

2. **Redis (Message Broker)**
   - Pub/Sub functionality for real-time message distribution
   - Built-in support for distributed systems
   - Low latency for real-time communication

3. **Apache Cassandra (Database)**
   - Distributed NoSQL database for high availability
   - Linear scalability for large datasets
   - Multi-datacenter replication

4. **Nginx (Load Balancer)**
   - Efficient load distribution across multiple servers
   - WebSocket connection handling
   - High-performance reverse proxy
   - Health checking and failover

5. **Docker (Containerization)**
   - Consistent development and production environments
   - Easy scaling and deployment
   - Resource isolation

### System Flow
1. Client connects via WebSocket through Nginx load balancer
2. Connection is routed to one of multiple Go servers
3. Messages are published to Redis channels
4. Redis distributes messages to appropriate server instances
5. Messages are persisted in Cassandra
6. Real-time delivery to recipients via WebSocket

### Scalability Features
- Horizontal scaling of Go servers
- Redis cluster for message distribution
- Cassandra's distributed nature for data storage
- Load balancing across multiple instances
- Containerized deployment for easy scaling

## Project Structure
```
distributed-chat-app/
├── config/
│   ├── redis.go     # Redis configuration and pub/sub
│   └── ws.go        # WebSocket handlers
├── controller/
│   ├── message.go   # Message handling logic
│   ├── room.go      # Room management
│   └── user.go      # User operations
├── database/
│   └── db.go        # Cassandra connection & queries
├── model/
│   ├── room.go      # Room data structures
│   └── user.go      # User data structures
├── nginx/
│   ├── Dockerfile   # Nginx container setup
│   └── default.conf # Load balancer configuration
├── router/
│   └── router.go    # API routes definition
├── .env             # Environment variables
├── db.cql           # Database schema
├── docker-compose.yaml
├── Dockerfile
├── go.mod
└── main.go
```

## Setup Instructions

### Initial Setup
```bash
# Clone the repository
git clone https://github.com/naman1402/distributed-chat-app.git
cd distributed-chat-app

# Clean existing containers (if any)
docker-compose down
docker system prune -f --volumes
```

### Start Services
```bash
# Start Cassandra first
docker-compose up -d cassandra
timeout /t 30  # Wait for Cassandra to initialize

# Initialize database schema
docker exec -i testCass cqlsh < db.cql

# Start remaining services
docker-compose up -d

# Verify all containers are running
docker ps
```

### Verification Commands
```bash
# Verify Cassandra schema
docker exec -it testCass cqlsh -e "USE chat; DESCRIBE TABLES;"

# Check service logs
docker logs nginxproxy     # Nginx logs
docker logs go_chat_1      # API1 logs
docker logs redis_chat_app # Redis logs
docker logs testCass       # Cassandra logs
```

## API Documentation

### Authentication Endpoints

#### Create User
```bash
curl -X POST http://localhost/signin \
  -H "Content-Type: application/json" \
  -d '{"username":"user1"}'
```
- Method: POST
- Endpoint: /signin
- Request Body: username (string)
- Response: 202 Accepted with user creation confirmation

#### User Login
```bash
curl -X POST http://localhost/login \
  -H "Content-Type: application/json" \
  -d '{"id":"user1"}'
```
- Method: POST
- Endpoint: /login
- Request Body: id (string)
- Response: 202 Accepted with user ID and name
- Sets cookie: uid

### Chat Room Operations

#### Create Room
```bash
curl -X POST http://localhost/create \
  -H "Content-Type: application/json" \
  -d '{"name":"room1"}'
```
- Method: POST
- Endpoint: /create
- Request Body: name (string)
- Response: 200 OK with confirmation

#### Join Room
```bash
curl -X POST http://localhost/join \
  -H "Content-Type: application/json" \
  -d '{"name":"room1","user":"user1"}'
```
- Method: POST
- Endpoint: /join
- Request Body: name (string), user (string)
- Response: 200 OK with confirmation

### WebSocket Connection
- Endpoint: ws://localhost/ws?id={userId}
- Query Parameter: id (user identifier)
- Authentication: Required via user ID

### Database Queries

#### Check Data
```bash
# View users
docker exec -it testCass cqlsh -e "USE chat; SELECT * FROM users;"

# View rooms
docker exec -it testCass cqlsh -e "USE chat; SELECT * FROM room;"

# View room members
docker exec -it testCass cqlsh -e "USE chat; SELECT * FROM room_members;"
```

## Cleanup Commands

### Graceful Shutdown
```bash
# Stop all services
docker-compose down

# Force shutdown if needed
docker-compose down --remove-orphans

# Clean volumes
docker volume prune -f

# Remove related containers
docker rm -f $(docker ps -a | grep 'go_chat\|redis\|nginx\|cass' | awk '{print $1}')
```

## Monitoring & Logs

### Server Logs
```bash
# API Server 1 logs
docker logs -f go_chat_1 
# With timestamp
docker logs -f --timestamps go_chat_1
# Last 100 lines
docker logs --tail 100 go_chat_1

# API Server 2 logs
docker logs -f go_chat_2
# With timestamp
docker logs -f --timestamps go_chat_2
# Last 100 lines
docker logs --tail 100 go_chat_2

# API Server 3 logs
docker logs -f go_chat_3
# With timestamp
docker logs -f --timestamps go_chat_3
# Last 100 lines
docker logs --tail 100 go_chat_3
```

### Infrastructure Logs
```bash
# Load Balancer logs
docker logs -f nginxproxy

# Redis logs
docker logs -f redis_chat_app

# Cassandra logs
docker logs -f testCass
```

## Learning Journey
This project is part of my learning journey into distributed systems, cloud-native architectures, and modern deployment practices. It has helped me understand:
- Distributed system patterns and challenges
- Real-time communication in scaled environments
- Container orchestration and service discovery
- Message broker patterns and event-driven architecture
- NoSQL database scaling and replication

Feel free to contribute! Whether it's improvements to the code, documentation updates, or bug fixes, all contributions are welcome. Please create a pull request or open an issue to discuss your ideas.


