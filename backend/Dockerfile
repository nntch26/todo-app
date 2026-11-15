# ---------- Stage 1: build the Go binary ----------
FROM golang:1.27.3-alpine AS build
WORKDIR /src

# Download modules first: this layer is cached until go.mod / go.sum change
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO off = a static binary that runs on plain Alpine; -s -w strips debug info (smaller)
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# ---------- Stage 2: small image with only the binary ----------
FROM alpine:3.23

# Run as a normal user, not root (-D: no password, -H: no home dir)
RUN adduser -D -H app
USER app

WORKDIR /app
COPY --from=build /out/api /app/api

# Config comes from environment variables set by docker compose (no .env file inside the image)
EXPOSE 8080
CMD ["/app/api"]
