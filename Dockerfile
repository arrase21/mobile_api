FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/mobile_api ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot AS runtime

COPY --from=builder /app/mobile_api /mobile_api

ENV LOG_LEVEL=info
ENV LOG_FILE=""
ENV PORT=8080

EXPOSE 8080

USER nonroot

ENTRYPOINT ["/mobile_api"]
