# SusyDB Release Guide

This document provides step-by-step instructions for publishing SusyDB as an open-source project on GitHub.

---

## Pre-Release Checklist

Before releasing, ensure the following are complete:

- [x] All code compiles without errors (`go build ./...`)
- [x] All static analysis passes (`go vet ./...`)
- [x] README.md is complete with features, benchmarks, and usage examples
- [x] LICENSE file exists (MIT)
- [x] ARCHITECTURE.md provides technical documentation
- [x] Dockerfile is present for containerized deployment
- [x] CI workflow (`.github/workflows/go.yml`) is configured

---

## Step 1: Initialize Git Repository

```bash
# Navigate to project directory
cd path/to/SusyDB

# Initialize Git (if not already done)
git init

# Add all files
git add .

# Create initial commit
git commit -m "feat: Initial release of SusyDB v1.0.0"

# Rename branch to main
git branch -M main
```

---

## Step 2: Create GitHub Repository

1. Go to [github.com/new](https://github.com/new)
2. **Repository name**: `susydb`
3. **Description**: "A tiny, high-performance, in-memory key-value store written in Go"
4. **Visibility**: Public
5. **Do NOT** initialize with README, .gitignore, or license (we already have them)
6. Click **Create repository**

---

## Step 3: Push to GitHub

```bash
# Add remote origin (replace with your actual URL)
git remote add origin https://github.com/Syed-Suhaan/susydb.git

# Push code to main branch
git push -u origin main
```

---

## Step 4: Create Version Tag

Go modules require semantic version tags for `go install` to work properly.

```bash
# Create version tag
git tag v1.0.0

# Push tag to GitHub
git push origin v1.0.0
```

---

## Step 5: Create GitHub Release

1. Navigate to your repository on GitHub
2. Click **Releases** (right sidebar)
3. Click **Draft a new release**
4. Configure the release:
   - **Choose a tag**: Select `v1.0.0`
   - **Release title**: `v1.0.0 - Initial Release`
   - **Description**: Use the template below

### Release Description Template

```markdown
# SusyDB v1.0.0

A tiny, high-performance, in-memory key-value store written in Go.

## Highlights

- **154k RPS** Mixed Workload performance
- Thread-safe in-memory store with `sync.RWMutex`
- Zero external dependencies
- RESP-compatible protocol (works with `redis-cli`)
- Pub/Sub messaging support
- Embeddable as a Go library

## Installation

```bash
go install github.com/Syed-Suhaan/susydb/cmd/susydb@latest
```

## Docker

```bash
docker build -t susydb .
docker run -d -p 7379:7379 susydb
```

## Documentation

- [README](https://github.com/Syed-Suhaan/susydb#readme)
- [Architecture](https://github.com/Syed-Suhaan/susydb/blob/main/ARCHITECTURE.md)

## What's Included

- Core KV operations: SET, GET, DEL
- TTL support: SETEX with lazy and active expiry
- Hash operations: HSET, HGET, HDEL, HGETALL
- Atomic counters: INCR, INCRBY
- Pub/Sub: PUBLISH, SUBSCRIBE
- Embedded library mode
```

5. Click **Publish release**

---

## Step 6: Post-Release Tasks

### 6.1 Go Report Card
1. Visit [goreportcard.com](https://goreportcard.com)
2. Enter your repository URL: `github.com/Syed-Suhaan/susydb`
3. Click **Generate Report**

### 6.2 pkg.go.dev Indexing
1. Visit `https://pkg.go.dev/github.com/Syed-Suhaan/susydb`
2. If not indexed, click **Request** to trigger indexing

### 6.3 Add Topics on GitHub
1. Go to repository settings
2. Add topics: `golang`, `key-value-store`, `in-memory-database`, `redis-alternative`, `cache`

### 6.4 Social Announcement (Optional)
Share on LinkedIn/Twitter:

> Just released SusyDB - a tiny, high-performance in-memory key-value store in Go that hits 154k RPS. Built for session storage, caching, and rate limiting.
>
> GitHub: github.com/Syed-Suhaan/susydb
>
> #golang #opensource #database

---

## Version History

| Version | Date | Notes |
|---------|------|-------|
| v1.0.0 | 2024-12-10 | Initial release with KV, TTL, Hash, Pub/Sub, Embedded Mode |

---

## Future Roadmap

- [x] Phase 1: Core Engine (KV, TTL, Hash, Counters)
- [x] Phase 2: Pub/Sub & Event Notification
- [ ] Phase 3: Persistence (AOF/RDB Snapshot) - *Out of scope for simplicity*
- [x] Phase 4: Embedded Library Mode
