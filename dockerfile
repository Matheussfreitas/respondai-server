FROM golang:1.25.1-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/api

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates curl
COPY --from=builder /build/app .
EXPOSE 8080
CMD [ "./app" ]
