# 📼 ytd

[![Go Version](https://img.shields.io/badge/go-1.22-blue)](https://golang.org/dl/)
[![Dockerized](https://img.shields.io/badge/docker-ready-blue)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**ytd** is a small service for downloading YouTube videos into organized folders on a media server.

---

## 🏗 Project Structure

```
ytd/
├── backend/          # Go backend (HTTP API)
├── frontend/         # Angular frontend (coming soon)
├── docker-compose.yml
├── helm/             # Helm chart for k3s deployment (coming soon)
└── README.md         # This file
```

---

## 🚀 Features

- 📥 Download YouTube videos to a media server.
- 📁 Organize videos into folders.
- 📝 Optionally rename videos on download.
- 🐳 Dockerized services.
- ☸️ Helm chart for easy k3s cluster deployment.

---

## 🛠 Backend API

| Method | Endpoint           | Description                  |
|--------|--------------------|------------------------------|
| POST   | `/api/download`     | Download a YouTube video.     |
| GET    | `/api/directories`  | List existing folders.        |
| POST   | `/api/directory`    | Create a new folder.          |

---

## 🎯 Quick Start (Backend)

```bash
cd backend
make build   # Build the backend
make run     # Run the backend locally
```

Or using Go directly:

```bash
go run main.go
```

Server will start on `localhost:8080`.

---

## 🐳 Running with Docker Compose

```bash
docker-compose up --build
```

---

## 📋 TODO

- Integrate frontend (Angular).
- Implement video download functionality.
- Handle recursive folder listing (optional).
- Add authentication.
- Configure persistent storage for downloads.
- Write unit tests for backend.

---

## 📄 License

MIT License.
