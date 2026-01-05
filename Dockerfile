FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/exporter ./cmd/exporter

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /app/exporter /exporter
EXPOSE 8080
ENTRYPOINT ["/exporter"]
