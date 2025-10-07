# 🧠 CRUD Project

A lightweight RESTful API written in **Go (Golang)** that demonstrates a clean architecture for CRUD operations using **PostgreSQL**. The project is containerized using **Docker** and includes live reload support via **Air**.

---

## ⚙️ Tech Stack

- **Language:** Go 1.25+
- **Database:** PostgreSQL
- **Hot Reload:** [Air](https://github.com/cosmtrek/air)
- **Containerization:** Docker & Docker Compose
- **Environment Management:** .env file (loaded automatically)

---

## 📁 Project Structure

```
CRUD PROJECT
├── cmd/api/              # Application entry point
│   ├── main.go           # App configuration and startup
│   ├── handler.go        # HTTP handlers for CRUD endpoints
│   └── route.go          # Route setup
│
├── internal/             # Internal packages (not exposed externally)
│   ├── database/         # Database connection logic
│   │   └── db.go
│   ├── env/              # Environment variable helpers
│   ├── store/            # Data models and queries
│   │   ├── note.go
│   │   └── error.go
│   └── migrations/       # SQL migrations (if any)
│
├── tmp/                  # Temporary build or runtime files
│
├── .air.toml             # Air configuration (for live reload)
├── .dockerignore         # Ignore rules for Docker
├── .env                  # Environment configuration
├── .gitignore
├── docker-compose.yml    # Docker Compose setup
├── Dockerfile            # Docker build instructions
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── Makefile              # Helper commands (run, build, etc.)
└── README.md             # Project documentation
```

---

## 🚀 Quick Start

### 1️⃣ Clone the repo

```bash
git clone https://github.com/dmc0001/crud-project.git
cd crud-project
```

### 2️⃣ Set up environment variables

Edit `.env`:

```env
PORT=:8080
DB_DSN=postgres://postgres:password@localhost:5432/note_db?sslmode=disable
```

### 3️⃣ Run locally with Air (hot reload)

```bash
air
```

Visit: [http://localhost:8080/notes](http://localhost:8080/notes)

### 4️⃣ Run with Docker Compose

```bash
docker compose up --build
```

This will start both the **API** and the **PostgreSQL** service.

---

## 📡 API Endpoints

| Method | Endpoint       | Description             |
| ------ | -------------- | ----------------------- |
| GET    | `/notes`       | Get latest notes        |
| GET    | `/note?id=1`   | Get a note by ID        |
| POST   | `/create`      | Create a new note       |
| PUT    | `/update?id=1` | Update an existing note |
| DELETE | `/delete?id=1` | Delete a note by ID     |

### Example Request (POST /create)

```json
{
  "title": "My first note",
  "content": "Hello from Go!"
}
```

### Example Response

```json
{
  "message": "Note with id 1 has been created",
  "id": 1
}
```

---

## 🐳 Docker Setup

**Dockerfile** builds the Go API image:

```dockerfile
FROM golang:1.25-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o crud cmd/api/main.go
EXPOSE 8080
CMD ["./crud"]
```

**docker-compose.yml** spins up API + DB:

```yaml
version: "3.9"
services:
  db:
    image: postgres:16
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: note_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  api:
    build: .
    ports:
      - "8080:8080"
    env_file: .env
    depends_on:
      - db

volumes:
  postgres_data:
```

---

## 🧩 Makefile Commands

```makefile
run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

docker-up:
	docker compose up --build

docker-down:
	docker compose down
```

---

## 🧪 Testing the API

You can use Postman or curl:

```bash
curl -X GET http://localhost:8080/notes
curl -X POST http://localhost:8080/create -d '{"title":"Test","content":"Example"}' -H 'Content-Type: application/json'
```

---

## 🧱 Future Improvements

- Add request validation middleware
- Integrate sqlc for fully type-safe idiomatic Go code from SQL.
- Implement JWT authentication
- Add automated tests (unit + integration)

---

## 👨‍💻 Author

**Haitham Attab**

---
