# Technical Workflow: WebSocket and Redis Integration

## System Architecture Overview

```
[Client] <--WebSocket--> [Server Instance 1] <--Redis Pub/Sub--> [Server Instance 2] <--WebSocket--> [Client]
```

## Core Components

1. **WebSocket Layer**
   - Manages real-time bidirectional communication
   - Maintains persistent connections with clients
   - Handles connection state (connect/disconnect)

2. **Redis Layer**
   - Manages cross-server communication
   - Handles message distribution
   - Maintains server-client mappings

## Detailed Workflow

### 1. Connection Establishment
```plaintext
1. Client initiates WebSocket connection with user ID
2. Server validates user (WSHandler function)
3. On successful validation:
   - Maps user ID to WebSocket connection
   - Stores user-server mapping in Redis
   - Sends acknowledgment to client
```

### 2. Private Message Flow
```plaintext
1. Client sends message via WebSocket
2. Server processes message (ReceiveMessage function):
   - Generates unique message ID
   - Validates message content
   - Saves to database
   - Determines recipient's server ID
3. Message Distribution:
   - If recipient is on same server:
     → Direct WebSocket delivery
   - If recipient is on different server:
     → Publishes to Redis channel
     → Receiving server gets message from Redis
     → Forwards to recipient via WebSocket
```

### 3. Group Message Flow
```plaintext
1. Client sends group message
2. Server processes message:
   - Validates group existence
   - Saves message to database
   - Retrieves group members
3. Server creates server-member mapping:
   - Groups members by their server IDs
4. Message Distribution:
   - For each server:
     → Creates message with server-specific member list
     → Publishes to server's Redis channel
   - Receiving servers:
     → Get message from Redis
     → Distribute to online members via WebSocket
```

## Redis Channel Structure

1. **Channel Naming**
   - Each server has unique SERVERID
   - Servers subscribe to their SERVERID channel
   - Messages published to specific server channels

2. **Message Format**
```json
{
    "id": "unique-message-id",
    "msg": "message content",
    "sender": "sender-id",
    "receiver": "receiver-id",
    "is_group": boolean,
    "group_name": "group-name",
    "group_members": ["member1", "member2"],
    "server_id": "target-server-id"
}
```

## Error Handling

1. **WebSocket Failures**
   - Connection drops: Client removed from active connections
   - Message delivery failure: Connection closed, client mapping removed
   - Invalid messages: Error response sent to client

2. **Redis Failures**
   - Connection loss: Panic and require restart
   - Message publish failure: Error logged
   - Subscribe channel error: Panic and require restart

## Performance Considerations

1. **Memory Management**
   - WebSocket connections stored in memory
   - Client mappings cleaned on disconnection
   - Redis connections pooled for efficiency

2. **Message Distribution**
   - Server-specific channels reduce message broadcasting
   - Group messages optimized by server-based member grouping
   - Async message handling for better performance

## Scaling Considerations

1. **Horizontal Scaling**
   - Multiple server instances possible
   - Redis handles cross-server communication
   - Server-specific channels prevent message duplication

2. **Load Distribution**
   - Users can connect to any server instance
   - Redis ensures message delivery regardless of connection point
   - Server-client mapping enables efficient message routing
