# Build
FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Set CGO off for static binary; adjust if you need CGO
ENV CGO_ENABLED=0
RUN go build -o /bin/ytd ./cmd/ytd

# Runtime
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /bin/ytd /usr/local/bin/ytd
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/ytd"]
