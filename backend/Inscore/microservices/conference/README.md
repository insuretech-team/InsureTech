# Conference Microservice

## Overview

The Conference microservice provides WebRTC signaling and room management capabilities for real-time video/audio conferencing. It implements a full-featured gRPC API with bidirectional streaming for signaling.

## Architecture

### Services Implemented

1. **RoomService** - Room management and lifecycle
   - CreateRoom, GetRoom, UpdateRoom, CloseRoom, ListRooms

2. **PeerService** - Peer connection management
   - JoinRoom, LeaveRoom, GetPeer, UpdatePeer, ListPeers

3. **TrackService** - Media track management
   - PublishTrack, UnpublishTrack, MuteTrack, GetTrack, ListTracks, UpdateTrack

4. **SignalingService** - WebRTC signaling with bidirectional streaming
   - Connect (bidirectional stream)
   - SendOffer, SendAnswer, SendICECandidate

5. **StatsService** - Connection statistics and analytics
   - GetConnectionStats, StreamStats, GetRoomAnalytics, GetPeerAnalytics

## Features

### Bidirectional Streaming Signaling

The SignalingService implements full bidirectional streaming for real-time WebRTC signaling:

- **Persistent Connections**: Each peer maintains a single bidirectional stream
- **Room Broadcasting**: Automatic event broadcasting to all peers in a room
- **Connection Management**: Automatic cleanup on disconnect
- **Event Types**:
  - Peer joined/left notifications
  - SDP offer/answer exchange
  - ICE candidate exchange
  - Track published/unpublished/muted events
  - Peer state changes

### Stream Management

The service maintains:
- `streams` map: peer_id → streamInfo (channel, roomID, peerID)
- `roomPeers` map: room_id → set of peer_ids
- Thread-safe operations with sync.Map and sync.RWMutex

### Message Flow

#### Joining a Room
1. Client calls PeerService.JoinRoom()
2. Client establishes SignalingService.Connect() stream
3. Service stores stream and notifies other peers
4. Client receives existing peers list

#### Signaling Flow
1. Peer A sends offer through bidirectional stream
2. Service validates and forwards to Peer B
3. Peer B receives offer via stream
4. Peer B sends answer through stream
5. Service forwards answer to Peer A
6. ICE candidates exchanged similarly

### Configuration

Server configuration via `ServerConfig`:
- Port (default: 50052)
- Max message sizes (10MB)
- Connection keepalive settings
- Reflection enabled for development

## Usage

### Starting the Service

```bash
# Default port 50052
go run lpc/microservices/conference/main.go

# Custom port
CONFERENCE_PORT=50053 go run lpc/microservices/conference/main.go
```

### Database Configuration

The service requires a PostgreSQL database. Configure in `lpc/configs/database.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: yourpassword
  dbname: lifepluscore
  sslmode: disable
```

### Client Integration

#### 1. Create a Room

```go
import (
    "github.com/fahara02/lifepluscore/gen/go/lifepluscore/webrtc/v1/service"
    "github.com/fahara02/lifepluscore/gen/go/lifepluscore/common/v1"
)

conn, _ := grpc.Dial("localhost:50052", grpc.WithInsecure())
roomClient := service.NewRoomServiceClient(conn)

resp, _ := roomClient.CreateRoom(ctx, &service.CreateRoomRequest{
    Name: "My Conference",
    Config: &entity.RoomConfig{
        MaxParticipants: 10,
        RequireToken: false,
    },
})
room := resp.Room
joinToken := resp.JoinToken
```

#### 2. Join Room as Peer

```go
peerClient := service.NewPeerServiceClient(conn)

resp, _ := peerClient.JoinRoom(ctx, &service.JoinRoomRequest{
    RoomId: room.RoomId,
    JoinToken: joinToken,
    DisplayName: "John Doe",
})
peer := resp.Peer
existingPeers := resp.ExistingPeers
```

#### 3. Establish Signaling Stream

```go
signalingClient := service.NewSignalingServiceClient(conn)
stream, _ := signalingClient.Connect(ctx)

// Start receiving messages
go func() {
    for {
        msg, err := stream.Recv()
        if err != nil {
            break
        }
        handleSignalResponse(msg)
    }
}()

// Send offer
stream.Send(&service.SignalRequest{
    PeerId: peer.PeerId,
    RoomId: room.RoomId,
    Payload: &service.SignalRequest_Offer{
        Offer: &service.SendOfferRequest{
            FromPeerId: myPeerId,
            ToPeerId: targetPeerId,
            Sdp: offerSdp,
        },
    },
})
```

#### 4. Handle Signal Responses

```go
func handleSignalResponse(msg *service.SignalResponse) {
    switch payload := msg.Payload.(type) {
    case *service.SignalResponse_OfferReceived:
        // Handle incoming offer
        handleOffer(payload.OfferReceived)
        
    case *service.SignalResponse_AnswerReceived:
        // Handle incoming answer
        handleAnswer(payload.AnswerReceived)
        
    case *service.SignalResponse_IceCandidateReceived:
        // Handle ICE candidate
        handleICECandidate(payload.IceCandidateReceived)
        
    case *service.SignalResponse_PeerJoined:
        // New peer joined room
        onPeerJoined(payload.PeerJoined)
        
    case *service.SignalResponse_PeerLeft:
        // Peer left room
        onPeerLeft(payload.PeerLeft)
        
    case *service.SignalResponse_TrackPublished:
        // Peer published new track
        onTrackPublished(payload.TrackPublished)
    }
}
```

## Server Configuration Options

### Keepalive Settings

Configured for reliable long-lived connections:
- MaxConnectionIdle: 5 minutes
- MaxConnectionAge: 30 minutes
- Keepalive interval: 30 seconds
- Keepalive timeout: 10 seconds

### Message Size Limits

- MaxRecvMsgSize: 10MB (for large SDP messages)
- MaxSendMsgSize: 10MB

### Channel Buffering

- Stream channels: 100 message buffer
- Prevents blocking on busy connections

## Monitoring and Health

### Health Check Endpoint

```go
server.HealthCheck(ctx)
```

Verifies:
- Database connection
- Database ping response

### Logging

The service uses structured logging:
- Service lifecycle events
- Connection events
- Error conditions
- Performance metrics

## Development

### Running Tests

```bash
cd lpc/pkg/webrtc
go test -v ./...
```

### Using gRPC Reflection

The server enables reflection for development tools:

```bash
# List services
grpcurl -plaintext localhost:50052 list

# Describe service
grpcurl -plaintext localhost:50052 describe lifepluscore.webrtc.v1.service.SignalingService

# Call method
grpcurl -plaintext -d '{"name": "Test Room"}' localhost:50052 \
  lifepluscore.webrtc.v1.service.RoomService/CreateRoom
```

## Production Considerations

### Security
- [ ] Add TLS/SSL certificates
- [ ] Implement JWT-based authentication
- [ ] Add rate limiting per peer
- [ ] Validate all user inputs

### Scalability
- [ ] Add Redis for distributed stream management
- [ ] Implement horizontal scaling with load balancer
- [ ] Add metrics collection (Prometheus)
- [ ] Implement connection pooling

### Reliability
- [ ] Add circuit breakers
- [ ] Implement retry logic
- [ ] Add request timeouts
- [ ] Monitor memory usage for stream channels

## Troubleshooting

### Connection Issues

1. **Stream disconnects frequently**
   - Check network stability
   - Review keepalive settings
   - Check client-side timeout handling

2. **Messages not delivered**
   - Verify peer is still connected
   - Check channel buffer size
   - Review error logs

3. **Database connection errors**
   - Verify database.yaml configuration
   - Check PostgreSQL is running
   - Verify network connectivity

### Performance Issues

1. **High latency**
   - Check network RTT
   - Review message size
   - Monitor server CPU/memory

2. **Memory growth**
   - Check for stream leaks
   - Review cleanup logic
   - Monitor goroutine count

## API Reference

See proto files in `proto/lifepluscore/webrtc/v1/service/` for complete API documentation.

## License

Proprietary - LifePlusCore
