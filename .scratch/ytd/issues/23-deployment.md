Status: ready-for-agent
Phase: 6 — Deployment

Production-ready container and k3s deployment. Multi-stage Dockerfile with yt-dlp baked in. docker-compose for local dev. k3s manifests for the homelab cluster.

## Acceptance criteria

- [ ] `Dockerfile` — multi-stage: Go build stage → minimal runtime image with yt-dlp installed
- [ ] `docker-compose.yaml` — mounts `./data` for SQLite and media, exposes port 8080
- [ ] k3s `Deployment` manifest — single replica, env vars via ConfigMap, resource limits
- [ ] k3s `Service` manifest — ClusterIP
- [ ] k3s `PersistentVolumeClaim` for MEDIA_DIR (NFS) and DB_PATH
- [ ] k3s `ConfigMap` for non-secret env vars
- [ ] `/healthz` and `/readyz` probes wired in the Deployment manifest
- [ ] `docker compose up` starts a working local instance
