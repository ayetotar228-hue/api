FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_NAME=api
ARG SERVICE_PATH=./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build -o /service ${SERVICE_PATH}

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /service .
COPY --from=builder /app/migrations ./migrations

ARG SERVICE_PORT=8080
EXPOSE ${SERVICE_PORT}

CMD ["./service"]