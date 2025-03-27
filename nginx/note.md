# Nginx Technical Documentation

## Overview

Nginx serves as the central load balancer and reverse proxy in our distributed chat application, managing both HTTP and WebSocket connections across multiple API instances.

## Architecture

```plaintext
                        ┌─→ [API1:6300]
[Client] → [Nginx:80] ──┼─→ [API2:6400] 
                        └─→ [API3:6500]
```

## Core Components

### 1. Load Balancer Configuration
```nginx
upstream backend {
    server go_chat_1:6300;  # API instance 1
    server go_chat_2:6400;  # API instance 2
    server go_chat_3:6500;  # API instance 3
    keepalive 32;          # Maintains 32 idle connections
}
```

- **Load Distribution**: Round-robin distribution by default
- **Server Names**: Match Docker container names
- **Connection Pool**: 32 keepalive connections per upstream server
- **Health Checks**: Automatic removal of failed servers

### 2. Connection Flow

#### HTTP Request Path
```plaintext
1. Client → Nginx:80
2. Nginx selects API instance (round-robin)
3. Request forwarded to API:6300/6400/6500
4. API processes request
5. Response returns through same path
```

#### WebSocket Connection Path
```plaintext
1. Client initiates WS connection
2. Nginx upgrades connection
3. Persistent WS connection established with API
4. Messages flow bidirectionally
5. Redis handles cross-server communication
```

## Integration Points

### 1. Docker Network Integration
```plaintext
[Nginx Container] ←→ [Docker Network: go_net] ←→ [API Containers]
```

- Nginx container name: `nginxproxy`
- Network name: `go_net`
- Internal DNS resolution using container names
- Automated container discovery

### 2. API Server Communication

#### Headers Management
```nginx
proxy_set_header Host $host;
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
```

- Preserves client information
- Enables request tracing
- Maintains origin data

#### WebSocket Support
```nginx
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";
```

- Handles WebSocket protocol upgrade
- Maintains persistent connections
- Enables real-time communication

## Detailed Message Flow

### 1. Private Chat Flow
```plaintext
Client A → Nginx → API1 → Redis → API2 → Nginx → Client B
```

1. Client A connects through Nginx
2. Message reaches assigned API server
3. API publishes to Redis
4. Target API receives from Redis
5. Message delivered via WebSocket

### 2. Group Chat Flow
```plaintext
                    ┌─→ API1 → Nginx → Client B
Client A → Nginx → API2 → Redis ─┼─→ API2 → Nginx → Client C
                                └─→ API3 → Nginx → Client D
```

1. Group message reaches any API
2. API fans out via Redis
3. Each recipient's API server receives
4. Parallel delivery to online members

## Performance Configurations

### 1. Timeout Settings
```nginx
proxy_connect_timeout 60;
proxy_send_timeout 60;
proxy_read_timeout 60;
```

- Connect timeout: 60s
- Send timeout: 60s
- Read timeout: 60s

### 2. Buffer Settings
```nginx
proxy_buffering on;
proxy_buffer_size 4k;
proxy_buffers 8 4k;
```

- Optimized for WebSocket traffic
- Minimizes memory usage
- Balances throughput

## Scaling Capabilities

### 1. Horizontal Scaling
- Add new API servers to upstream
- No client-side changes needed
- Automatic load distribution

### 2. High Availability
- Automatic failed server removal
- Connection retry mechanisms
- Session persistence support

## Monitoring & Debug Points

### 1. Log Locations
```nginx
error_log /var/log/nginx/error.log debug;
access_log /var/log/nginx/access.log;
```

### 2. Health Metrics
- Connection status
- Error rates
- Request latencies
- Backend availability

## Best Practices

1. **Security**
   - Header sanitization
   - Connection limits
   - Timeout management

2. **Performance**
   - Keepalive connections
   - Buffer optimization
   - Connection pooling

3. **Reliability**
   - Health checks
   - Graceful failure handling
   - Connection retry logic
