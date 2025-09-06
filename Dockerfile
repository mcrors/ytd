# ---- Build your Go app -------------------------------------------------------
FROM golang:1.23 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /bin/ytd ./cmd/ytd

# ---- Fetch yt-dlp and ffmpeg (robust, multi-arch) ---------------------------
FROM debian:bookworm-slim AS media
ARG YTDLP_VERSION=2025.08.11     # or whatever you're pinning
ARG FFMPEG_VERSION=6.1.1
ARG TARGETARCH

RUN set -eux; \
  apt-get update; \
  # add tools we actually use here: ldd (libc-bin) + readelf (binutils)
  apt-get install -y --no-install-recommends ca-certificates curl xz-utils libc-bin binutils; \
  update-ca-certificates; \
  arch="${TARGETARCH:-amd64}"; \
  case "$arch" in \
    amd64)  ytdlp_url="https://github.com/yt-dlp/yt-dlp/releases/download/${YTDLP_VERSION}/yt-dlp_linux"; \
            ffmpeg_release_url="https://www.johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"; \
            ffmpeg_pinned_url="https://www.johnvansickle.com/ffmpeg/releases/ffmpeg-${FFMPEG_VERSION}-amd64-static.tar.xz" ;; \
    arm64)  ytdlp_url="https://github.com/yt-dlp/yt-dlp/releases/download/${YTDLP_VERSION}/yt-dlp_linux_aarch64"; \
            ffmpeg_release_url="https://www.johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz"; \
            ffmpeg_pinned_url="https://www.johnvansickle.com/ffmpeg/releases/ffmpeg-${FFMPEG_VERSION}-arm64-static.tar.xz" ;; \
    *) echo "Unsupported TARGETARCH: $arch"; exit 1 ;; \
  esac; \
  # yt-dlp (self-contained PyInstaller binary)
  curl -fL --retry 5 --connect-timeout 15 "$ytdlp_url" -o /usr/local/bin/yt-dlp; \
  chmod +x /usr/local/bin/yt-dlp; \
  \
  # ---- collect yt-dlp runtime libs for distroless ----
  mkdir -p /opt/ytlibs; \
  # list deps; ignore "not found" lines; copy each absolute path with its parent dirs
  ldd /usr/local/bin/yt-dlp | awk '{print $3}' | grep -E '^/' | xargs -r -I{} cp -v --parents {} /opt/ytlibs; \
  # also copy the ELF interpreter (dynamic loader), in case ldd didn’t include it
  interp="$(readelf -l /usr/local/bin/yt-dlp | awk '/Requesting program interpreter/ {print $NF}' | tr -d "]")"; \
  [ -n "$interp" ] && cp -v --parents "$interp" /opt/ytlibs || true; \
  \
  # ffmpeg (static, includes ffprobe)
  ff_url="$ffmpeg_pinned_url"; \
  if ! curl -fL --retry 5 --connect-timeout 15 "$ff_url" -o /tmp/ffmpeg.tar.xz; then \
    ff_url="$ffmpeg_release_url"; \
    curl -fL --retry 5 --connect-timeout 15 "$ff_url" -o /tmp/ffmpeg.tar.xz; \
  fi; \
  xz -t /tmp/ffmpeg.tar.xz; \
  mkdir -p /opt/ffmpeg; \
  tar -xJf /tmp/ffmpeg.tar.xz -C /opt/ffmpeg --strip-components=1; \
  install -m 0755 /opt/ffmpeg/ffmpeg /usr/local/bin/ffmpeg; \
  install -m 0755 /opt/ffmpeg/ffprobe /usr/local/bin/ffprobe; \
  cp /etc/ssl/certs/ca-certificates.crt /ca-certificates.crt; \
  ls -lh /usr/local/bin/yt-dlp /usr/local/bin/ffmpeg /usr/local/bin/ffprobe /ca-certificates.crt

# ---- Tiny final image (glibc present) ---------------------------------------
FROM gcr.io/distroless/cc-debian12:nonroot
WORKDIR /app

COPY --from=build /bin/ytd /usr/local/bin/ytd
COPY --from=media /usr/local/bin/yt-dlp /usr/local/bin/yt-dlp
COPY --from=media /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=media /usr/local/bin/ffprobe /usr/local/bin/ffprobe

# copy CA bundle and yt-dlp’s required shared libs
COPY --from=media /ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=media /opt/ytlibs/ /

ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt \
    PATH=/usr/local/bin:$PATH
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/ytd"]
