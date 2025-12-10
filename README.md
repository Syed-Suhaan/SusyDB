# SusyDB

> **Tagline**: SusyDB is a tiny, high-performance, in-memory key-value store written in Go, optimized for simple cache, sessions, and counters.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Syed-Suhaan/SusyDB)](https://goreportcard.com/report/github.com/Syed-Suhaan/SusyDB)

---

## Why SusyDB?
Teams often deploy **Redis** just to store a few session tokens or cache API responses. 

SusyDB is for when you want:
- **Zero Ops**: Single binary, no dependencies.
- **Simplicity**: Just `SET`, `GET`, `TTL`, `Hashes`, and `Counters`.
- **Pure Go**: Easy to hack on, embed, or deploy alongside Go microservices.

**What it's NOT:**
- A durable database (No persistence *yet*).
- A distributed cluster (No sharding).

---

## Architecture & Design
SusyDB follows a shared-nothing, single-process architecture. It handles concurrency using Go's native synchronization primitives.

```mermaid
graph LR
    Client -->|TCP| Server
    Server -->|Goroutine| Parser
    Parser -->|RWMutex| Memory
```

- **Concurrency**: Uses `sync.RWMutex` for thread-safe access to the global map.
    > *Note: While a global lock is simple, SusyDB minimizes contention by keeping critical sections purely in-memory and extremely short. For 99% of session/cache workloads, the network I/O is the bottleneck, not the lock.*
- **Eviction**: Hybrid approach (Lazy check on Get + Active background sweeper).
- **Communication**: Simple text-based TCP protocol.

**[Read the Full Architecture Doc](ARCHITECTURE.md)** for deep dives into specific modules.

---

## Features
- **Protocol**: 100% RESP-compatible (works with `redis-cli`).
- **Data Structures**:
    - **Strings**: `SET`, `GET`, `DEL`, `SETEX` (legacy TTL command).
    - **Hashes**: `HSET`, `HGET`, `HDEL`, `HGETALL` (Perfect for sessions).
    - **Counters**: `INCR`, `INCRBY` (Rate limiting ready).
- **Pub/Sub**: Lightweight Message Broker (`PUBLISH`, `SUBSCRIBE`).
- **Embedded Mode**: Use as a library `import "github.com/Syed-Suhaan/SusyDB/pkg/core"` in your Go apps.
- **Hybrid Expiry**: Lazy + Active TTL implementation.
- **Architecture**: Thread-safe design using `sync.RWMutex`.
- **Observability**: `INFO` command for stats.

---

## Quick Start
### Installation
```bash
go install github.com/Syed-Suhaan/SusyDB/cmd/susydb@latest
```

### Running
```bash
susydb --addr :7379
```

### Docker
```bash
# Build the image
docker build -t susydb .

# Run it
docker run -d -p 7379:7379 susydb
```

### Usage
Connect with `telnet` or `netcat`:
```bash
$ telnet localhost 7379
Trying 127.0.0.1...
Connected to localhost.

SET user:1 Suhaan
+OK

GET user:1
Suhaan

INCR hits
:1

INCR hits
:2

INFO
$55
# Stats
keys:2
uptime_in_seconds:12
```

### Embedded Library
SusyDB can be imported directly into your Go applications, bypassing the network layer entirely for ultra-low latency.

```go
package main

import (
    "fmt"
    "github.com/Syed-Suhaan/SusyDB/pkg/core"
)

func main() {
    // Zero setup, pure Go
    db := core.NewKVStore()

    // Use it just like the server
    db.Set("session:123", "active", 3600)
    
    val, _, _ := db.Get("session:123")
    fmt.Println(val) // Output: active
}
```

---

## Benchmarks
### Fairness Disclosure
These tests compare **Redis (v5.0.14)** and **SusyDB (v1.0.0)** running natively on the same **Windows 11** machine. Both used default configurations and the same Go-based benchmarking client.

**Machine Specs:** Native Windows 11, Localhost.
**Workload:** 50 Concurrent Clients, 100,000 Requests, Small (16B) payloads.

| Workload | Redis (RPS / P99) | SusyDB (RPS / P99) | Gain |
| :--- | :--- | :--- | :--- |
| **A. Mixed** (50/50 SET/GET) | 99k / 5.2ms | **154k / 1.8ms** | **1.5x** |
| **B. TTL** (SETEX 60s) | 17k / 24ms | **170k / 0.9ms** | **9.5x** |
| **C. Counters** (INCR Hot Key) | 22k / 4.8ms | **163k / 1.1ms** | **7.3x** |
| **D. Hash** (Session HSET/HGET) | 19k / 5.1ms | **168k / 2.3ms** | **8.4x** |
| **E. Pub/Sub** (Message Blast) | 18k / 5.4ms | **143k / 1.2ms** | **7.6x** |

> **Context**: Redis is optimized for Linux; its Windows port is known to be slower due to OS-level threading differences. However, this demonstrates SusyDB's efficiency as a **native Go solution** for Windows-based microservices.

**Scope Note**: SusyDB is specialized for lightweight, in-memory session/cache workloads. It does **not** replace Redis features like replication, clustering, or disk persistence.

### Reproduce this Bench:
```bash
# Mixed Workload
go run cmd/susy-bench/main.go -h localhost:7379 -c 50 -n 100000 -test mixed

# TTL Workload
go run cmd/susy-bench/main.go -h localhost:7379 -c 50 -n 100000 -test setex
```

---

## License
MIT License. See [LICENSE](LICENSE) for details.
