# URL Shortening Service

A high-performance URL shortener with an analytics dashboard, built with **Go**, **Next.js/TypeScript**, **PostgreSQL**, and **Redis**.

## Features

- Shorten long URLs into unique 7-character short codes
- Fast redirects via Redis cache-aside pattern
- Async click logging with buffered channel + batch inserts
- Analytics dashboard with interactive charts (clicks over time, browser & country breakdown)
- Period filtering (24h, 7d, 30d, 90d)
- Full REST API with strict TypeScript interfaces

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go (chi router, pgx, go-redis) |
| Frontend | Next.js 16, TypeScript, Tailwind CSS, Recharts |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.21+
- [Node.js](https://nodejs.org/) 20+
- [Docker](https://www.docker.com/) (for PostgreSQL & Redis)

### 1. Start databases

```bash
docker-compose up -d
```

### 2. Start the backend

```bash
cd backend
cp ../.env.example .env    # or set DATABASE_URL env var
DATABASE_URL="postgres://urlshortener:devpassword@localhost:5432/urlshortener?sslmode=disable" go run ./cmd/server
```

The server starts on `http://localhost:8080`. Migrations run automatically on startup.

### 3. Start the frontend

```bash
cd frontend
npm install
npm run dev
```

The app is available at `http://localhost:3000`. API calls are proxied to the Go backend.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/urls` | Create a short URL |
| `GET` | `/api/urls` | List all URLs (paginated) |
| `GET` | `/api/urls/{id}` | Get a single URL |
| `DELETE` | `/api/urls/{id}` | Delete a URL |
| `GET` | `/r/{shortCode}` | Redirect to original URL |
| `GET` | `/api/urls/{id}/analytics/clicks?period=7d` | Clicks over time |
| `GET` | `/api/urls/{id}/analytics/sources?period=7d` | Browser & country breakdown |

## Architecture

```
Browser → GET /r/{code} → Go Backend → Redis (cache) → PostgreSQL
                              ↓
                    Async click logging (buffered channel → batch insert)

Browser → Next.js (:3000) → /api/* proxy → Go Backend (:8080)
```

**Backend** follows a 3-layer architecture: Handler → Service → Repository. Click events are logged asynchronously via a buffered Go channel to keep redirect latency minimal.

**Frontend** uses the Next.js App Router with typed API wrappers and Recharts for data visualization.
