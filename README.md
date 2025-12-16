# go-image-web

A minimal lightweight image board server written in Go

<p align="center">
  <img src="https://img.shields.io/badge/Source Code-Go-00ADD8?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Container-Docker-2496ED?style=flat-square&logo=docker" alt="Docker">
</p>

## Quick Start (Development)

### Prerequisites

- Docker & Docker Compose
- Git

### Run the Container

```bash
docker run --name go-image-web -v ./data:/app/data -p 9191:9991 -d ghcr.io/irrisdev/go-image-web:latest
 ```

### Run Locally

```bash
# Clone the repository
git clone https://github.com/irrisdev/go-image-web.git

cd go-image-web

docker compose up --build

# Access the platform: http://localhost:9991
```

---