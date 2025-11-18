FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o hifiy ./main.go

FROM alpine:latest

WORKDIR /backend

COPY --from=builder /app/hifiy /backend/hifiy

EXPOSE 8082

ENTRYPOINT ["/backend/hifiy"]
