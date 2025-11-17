# Background Jobs dengan Asynq

Aplikasi ini menggunakan [Asynq](https://github.com/hibiken/asynq) untuk menangani background job processing secara asynchronous dan scalable.

## Arsitektur

```
┌─────────────┐         ┌────────┐         ┌────────────┐
│   HTTP/gRPC │────────▶│ Redis  │◀────────│   Worker   │
│   Server    │  Enqueue│ Queue  │ Process │  Process   │
└─────────────┘         └────────┘         └────────────┘
```

- **HTTP/gRPC Server**: Mengirim task ke Redis queue
- **Redis**: Menyimpan task queue dan metadata
- **Worker**: Memproses task dari queue secara concurrent

## Task yang Tersedia

### 1. Welcome Email Task
Dikirim otomatis setelah user berhasil registrasi.

**Payload:**
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "name": "John Doe"
}
```

**Location:** `internal/infrastructure/asynq/tasks/email_task.go`

## Menjalankan Worker

### Development (Local)

Pastikan Redis sudah running:
```bash
# Menggunakan Docker
docker run -d -p 6379:6379 redis:7-alpine

# Atau menggunakan docker-compose
docker-compose up -d redis
```

Jalankan worker:
```bash
go run cmd/worker/main.go
```

### Production (Docker)

Worker sudah termasuk dalam `docker-compose.yml`:
```bash
docker-compose up -d worker
```

## Konfigurasi

Worker menggunakan environment variable:

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `REDIS_ADDR` | `localhost:6379` | Alamat Redis server |

Concurrency default: **10 workers**

## Queue Priority

Asynq menggunakan weighted priority untuk memproses task:

| Queue | Priority | Deskripsi |
|-------|----------|-----------|
| `critical` | 60% | Task dengan prioritas tinggi |
| `default` | 30% | Task dengan prioritas normal |
| `low` | 10% | Task dengan prioritas rendah |

## Menambahkan Task Baru

### 1. Buat Task Definition

Buat file baru di `internal/infrastructure/asynq/tasks/`:

```go
package tasks

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "github.com/hibiken/asynq"
)

const (
    TypeNewTask = "task:new"
)

type NewTaskPayload struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}

func NewNewTask(field1 string, field2 int) (*asynq.Task, error) {
    payload, err := json.Marshal(NewTaskPayload{
        Field1: field1,
        Field2: field2,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to marshal payload: %w", err)
    }
    return asynq.NewTask(TypeNewTask, payload), nil
}

func HandleNewTask(ctx context.Context, t *asynq.Task) error {
    var payload NewTaskPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return fmt.Errorf("failed to unmarshal payload: %w", err)
    }

    // Process task
    log.Printf("Processing task: %+v", payload)

    return nil
}
```

### 2. Register Handler di Worker

Edit `cmd/worker/main.go`:

```go
// Register task handlers
mux.HandleFunc(tasks.TypeEmailWelcome, tasks.HandleEmailWelcomeTask)
mux.HandleFunc(tasks.TypeNewTask, tasks.HandleNewTask) // Tambahkan ini
```

### 3. Enqueue Task dari Service

```go
import (
    "github.com/gieart87/gohexaclean/internal/infrastructure/asynq/tasks"
)

// Di dalam service method
task, err := tasks.NewNewTask("value1", 123)
if err != nil {
    log.Printf("failed to create task: %v", err)
} else {
    info, err := s.taskClient.Enqueue(task)
    if err != nil {
        log.Printf("failed to enqueue task: %v", err)
    } else {
        log.Printf("enqueued task: id=%s", info.ID)
    }
}
```

## Advanced Features

### Task Options

Asynq mendukung berbagai opsi untuk task:

```go
// Delay task execution
task, _ := tasks.NewEmailWelcomeTask(userID, email, name)
s.taskClient.Enqueue(task, asynq.ProcessIn(5*time.Minute))

// Set max retry
s.taskClient.Enqueue(task, asynq.MaxRetry(5))

// Set task priority
s.taskClient.Enqueue(task, asynq.Queue("critical"))

// Set deadline
s.taskClient.Enqueue(task, asynq.Deadline(time.Now().Add(1*time.Hour)))
```

### Retry Policy

Secara default, Asynq akan retry task yang gagal dengan exponential backoff.

### Monitoring

Asynq menyediakan web UI untuk monitoring (opsional):
```bash
go install github.com/hibiken/asynq/tools/asynq@latest
asynq dash
```

Akses di: http://localhost:8080

## Best Practices

1. **Idempotent Tasks**: Pastikan task dapat dijalankan berulang kali tanpa side effect
2. **Error Handling**: Selalu return error jika task gagal agar bisa di-retry
3. **Timeout**: Set timeout yang reasonable untuk task yang berjalan lama
4. **Monitoring**: Monitor queue size dan failure rate di production
5. **Graceful Shutdown**: Worker sudah implement graceful shutdown secara default

## Troubleshooting

### Task tidak diproses
- Pastikan Redis sudah running
- Cek log worker untuk error
- Pastikan task handler sudah diregister

### Task gagal terus
- Cek log untuk error detail
- Pastikan payload valid
- Periksa apakah ada dependency yang tidak tersedia

### Redis connection error
- Pastikan REDIS_ADDR sudah benar
- Cek network connectivity
- Pastikan Redis accessible dari worker

## Referensi

- [Asynq Documentation](https://github.com/hibiken/asynq)
- [Asynq Wiki](https://github.com/hibiken/asynq/wiki)
