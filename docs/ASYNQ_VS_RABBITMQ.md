# Asynq vs RabbitMQ: Kapan Menggunakan Yang Mana?

## Overview

Aplikasi ini menggunakan **dua sistem asynchronous** yang berbeda namun saling melengkapi:
- **Asynq** - Background job processing
- **RabbitMQ** - Event-driven messaging

Keduanya **BUKAN pengganti**, melainkan **saling melengkapi** untuk membuat arsitektur yang robust dan scalable.

## 🔄 Asynq (Background Jobs)

### Karakteristik
- Backend: Redis
- Pattern: Producer → Queue → Worker
- Consumers: 1 worker per task
- Fokus: Job reliability & retry mechanism

### Use Cases
✅ **Delayed/Scheduled Tasks**
- Send welcome email setelah registrasi
- Generate monthly reports
- Cleanup expired data

✅ **Heavy Processing**
- Image/video processing
- PDF generation
- Data export/import

✅ **Retry-Critical Operations**
- Payment processing
- API calls to external services
- Email notifications

### Contoh di Codebase

```go
// File: internal/app/user_service.go
// Enqueue welcome email task asynchronously
if s.taskClient != nil {
    task, err := tasks.NewEmailWelcomeTask(user.ID.String(), user.Email, user.Name)
    if err != nil {
        log.Printf("failed to create welcome email task: %v", err)
    } else {
        info, err := s.taskClient.Enqueue(task)
        if err != nil {
            log.Printf("failed to enqueue welcome email task: %v", err)
        } else {
            log.Printf("enqueued welcome email task: id=%s queue=%s", info.ID, info.Queue)
        }
    }
}
```

### Features
- ✅ Automatic retry with exponential backoff
- ✅ Task prioritization (high, medium, low)
- ✅ Scheduled/delayed execution
- ✅ Task inspection & monitoring
- ✅ Dead letter queue
- ✅ Unique tasks (prevent duplicates)

### Monitoring
Access Asynq web UI (optional):
```bash
# Install asynqmon
go install github.com/hibiken/asynq/tools/asynqmon@latest

# Run monitoring UI
asynqmon --redis-addr=localhost:6379

# Access at http://localhost:8080
```

## 📨 RabbitMQ (Message Broker)

### Karakteristik
- Backend: RabbitMQ Server
- Pattern: Publisher → Exchange → Consumers
- Consumers: Multiple consumers per event
- Fokus: Event broadcasting & service decoupling

### Use Cases
✅ **Event Broadcasting**
- User created → notify analytics, CRM, email service
- Order placed → update inventory, send notification, log audit
- Payment received → update order status, send receipt

✅ **Microservice Communication**
- Decouple services
- Asynchronous inter-service communication
- Real-time notifications

✅ **Complex Routing**
- Topic-based routing
- Fanout to multiple consumers
- Direct/topic exchanges

### Contoh di Codebase

```go
// File: internal/app/user_service.go
// Publish user created event
if s.eventPublisher != nil {
    event := domain.NewUserCreatedEvent(user.ID, user.Email, user.Name)
    if err := s.eventPublisher.PublishUserCreated(ctx, event); err != nil {
        // Log error but don't fail the operation
        fmt.Printf("failed to publish user created event: %v\n", err)
    }
}
```

### Features
- ✅ Pub/Sub pattern (1 event → N consumers)
- ✅ Message routing (direct, topic, fanout, headers)
- ✅ Message persistence
- ✅ Acknowledgment & confirms
- ✅ Dead letter exchanges
- ✅ Priority queues

### Monitoring
Access RabbitMQ Management UI:
```
URL: http://localhost:15672
Username: guest
Password: guest
```

## 🤝 Penggunaan Bersamaan

### Contoh: User Registration Flow

Ketika user melakukan registrasi di `internal/app/user_service.go`:

```go
func (s *UserService) CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.LoginResponse, error) {
    // 1. Create user in database
    user := domain.NewUser(req.Email, req.Name, hashedPassword)
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // 2. RabbitMQ: Broadcast event to multiple services
    //    → Analytics service: track new user
    //    → CRM service: sync customer data
    //    → Audit service: log event
    if s.eventPublisher != nil {
        event := domain.NewUserCreatedEvent(user.ID, user.Email, user.Name)
        s.eventPublisher.PublishUserCreated(ctx, event)
    }

    // 3. Asynq: Queue specific job for processing
    //    → Only email service processes this
    //    → Guaranteed execution with retry
    if s.taskClient != nil {
        task, err := tasks.NewEmailWelcomeTask(user.ID.String(), user.Email, user.Name)
        s.taskClient.Enqueue(task)
    }

    return response, nil
}
```

### Arsitektur Flow

```
┌─────────────────┐
│  User Register  │
└────────┬────────┘
         │
         v
┌────────────────────┐
│  CreateUser()      │
│  - Save to DB      │
└────────┬───────────┘
         │
         ├─────────────────────┐
         │                     │
         v                     v
┌────────────────┐    ┌──────────────────┐
│   RabbitMQ     │    │     Asynq        │
│   (Events)     │    │     (Jobs)       │
└────────┬───────┘    └────────┬─────────┘
         │                     │
         ├──────┬──────┬       │
         │      │      │       │
         v      v      v       v
    ┌─────┐ ┌─────┐ ┌─────┐ ┌──────┐
    │Analy│ │ CRM │ │Audit│ │Email │
    │tics │ │     │ │ Log │ │Worker│
    └─────┘ └─────┘ └─────┘ └──────┘
```

## 📊 Comparison Table

| Aspect | Asynq (Jobs) | RabbitMQ (Events) |
|--------|--------------|-------------------|
| **Backend** | Redis | RabbitMQ Server |
| **Pattern** | Producer → Queue → Worker | Publisher → Exchange → Consumers |
| **Consumers** | 1 worker per task | Multiple per event |
| **Main Purpose** | Job processing | Event broadcasting |
| **Retry** | Built-in automatic | Manual implementation |
| **Scheduling** | ✅ Delayed/cron jobs | ❌ Not built-in |
| **Priority** | ✅ Task priority | ✅ Queue priority |
| **Monitoring** | Asynq inspector | RabbitMQ UI |
| **Complexity** | Simple (Redis only) | Complex (separate server) |
| **Best For** | Background jobs | Microservices |

## 💡 Decision Tree: Kapan Pakai Apa?

### Gunakan Asynq Jika:
- ✅ Perlu process task yang berat/lama (image processing, reports)
- ✅ Butuh retry otomatis dengan backoff
- ✅ Perlu scheduled/delayed execution
- ✅ Simple job queue sudah cukup
- ✅ Hanya 1 worker yang perlu process task
- ✅ Perlu monitoring task progress

**Contoh:**
- Send email notifications
- Generate PDF reports
- Process uploaded images
- Cleanup expired sessions
- Export data to CSV

### Gunakan RabbitMQ Jika:
- ✅ Event-driven architecture
- ✅ Multiple services perlu tahu tentang 1 event
- ✅ Microservice inter-communication
- ✅ Complex routing patterns (topic, fanout)
- ✅ Message durability & persistence penting
- ✅ Perlu decouple services

**Contoh:**
- User created → notify multiple services
- Order placed → trigger multiple workflows
- Payment received → update multiple systems
- Real-time notifications
- Audit logging across services

### Gunakan Keduanya Jika:
- ✅ Production-ready application
- ✅ Butuh event broadcasting DAN background processing
- ✅ Scalable microservice architecture
- ✅ Separation of concerns penting

## 🚀 Setup & Configuration

### Asynq Setup

**1. Pastikan Redis Running**
```bash
# Via docker-compose
make docker-up

# Atau manual
docker run -d -p 6379:6379 redis:7-alpine
```

**2. Worker akan auto-start**
```bash
# Via docker-compose (automatic)
make docker-up

# Atau manual
go run cmd/worker/main.go
```

**3. Verify**
```bash
# Check worker logs
docker logs gohexaclean-worker

# Should see:
# "Asynq worker started (Redis: redis:6379, Concurrency: 10)"
```

### RabbitMQ Setup

**1. Start RabbitMQ** (Optional - disabled by default)
```bash
docker-compose -f docker-compose.rabbitmq.yml up -d
```

**2. Enable in Config**
```yaml
# config/app.yaml
broker:
  enabled: true
  type: rabbitmq
  rabbitmq:
    url: amqp://guest:guest@localhost:5672/
```

**3. Verify**
```bash
# Access management UI
http://localhost:15672
# Username: guest
# Password: guest
```

## 📝 Adding New Tasks/Events

### Adding New Asynq Task

**1. Create Task Definition**
```go
// internal/infra/asynq/tasks/my_task.go
package tasks

import (
    "context"
    "encoding/json"
    "github.com/hibiken/asynq"
)

const TypeMyTask = "task:my_task"

type MyTaskPayload struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}

func NewMyTask(field1 string, field2 int) (*asynq.Task, error) {
    payload, err := json.Marshal(MyTaskPayload{
        Field1: field1,
        Field2: field2,
    })
    if err != nil {
        return nil, err
    }
    return asynq.NewTask(TypeMyTask, payload), nil
}

func HandleMyTask(ctx context.Context, t *asynq.Task) error {
    var payload MyTaskPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return err
    }

    // Process task logic here
    log.Printf("Processing task: %+v", payload)

    return nil
}
```

**2. Register Handler**
```go
// cmd/worker/main.go
mux.HandleFunc(tasks.TypeMyTask, tasks.HandleMyTask)
```

**3. Enqueue from Service**
```go
task, err := tasks.NewMyTask("value", 123)
if err != nil {
    return err
}
info, err := s.taskClient.Enqueue(task)
```

### Adding New RabbitMQ Event

**1. Create Domain Event**
```go
// internal/domain/events.go
type MyEvent struct {
    ID        uuid.UUID
    Field1    string
    Field2    int
    OccurredAt time.Time
}

func NewMyEvent(id uuid.UUID, field1 string, field2 int) *MyEvent {
    return &MyEvent{
        ID:         id,
        Field1:     field1,
        Field2:     field2,
        OccurredAt: time.Now(),
    }
}
```

**2. Publish Event**
```go
// In service
event := domain.NewMyEvent(id, "value", 123)
if err := s.eventPublisher.PublishMyEvent(ctx, event); err != nil {
    log.Printf("failed to publish event: %v", err)
}
```

**3. Create Consumer** (in consuming service)
```go
// internal/adapter/inbound/consumer/my_consumer.go
func (c *MyConsumer) HandleMyEvent(ctx context.Context, event *domain.MyEvent) error {
    log.Printf("Received event: %+v", event)
    // Process event
    return nil
}
```

## 🔍 Troubleshooting

### Asynq Issues

**Tasks not being processed:**
```bash
# Check worker is running
docker ps | grep worker

# Check Redis connection
docker logs gohexaclean-worker | grep -i redis

# Check for errors
docker logs gohexaclean-worker | grep -i error
```

**Redis connection failed:**
- Ensure Redis container is running
- Check Redis port (default: 6379)
- Verify network connectivity

### RabbitMQ Issues

**Events not received:**
```bash
# Check RabbitMQ is running
docker ps | grep rabbitmq

# Check connection
docker logs gohexaclean-http | grep -i rabbitmq

# Verify broker enabled in config
cat config/app.yaml | grep -A 5 broker
```

**Connection refused:**
- Ensure RabbitMQ container is running
- Check port 5672 is accessible
- Verify credentials in config

## 📚 Further Reading

- [Asynq Documentation](https://github.com/hibiken/asynq)
- [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
- [Message Broker Documentation](./MESSAGE_BROKER.md)
- [Async Jobs Documentation](./ASYNC_JOBS.md)

## 🎯 Best Practices

### Asynq Best Practices
1. ✅ Always handle errors gracefully
2. ✅ Use exponential backoff for retries
3. ✅ Set appropriate timeouts
4. ✅ Monitor task queues regularly
5. ✅ Use task priorities wisely
6. ✅ Implement idempotent handlers

### RabbitMQ Best Practices
1. ✅ Always acknowledge messages
2. ✅ Use persistent messages for critical events
3. ✅ Implement dead letter queues
4. ✅ Set appropriate prefetch count
5. ✅ Monitor queue lengths
6. ✅ Use proper exchange types

## 📊 When to Scale

### Scale Asynq Workers When:
- Queue backlog consistently high
- Task processing time increases
- Need parallel processing
- Different task priorities needed

**Solution:** Run multiple worker instances

### Scale RabbitMQ When:
- Message throughput exceeds single node
- Need high availability
- Geographic distribution required
- Message persistence critical

**Solution:** RabbitMQ clustering or federation

---

**Summary:** Asynq untuk **job processing**, RabbitMQ untuk **event broadcasting**. Gunakan keduanya untuk arsitektur yang robust!
