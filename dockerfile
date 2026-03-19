FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .

FROM alpine:3.x
WORKDIR /app
COPY --from=builder /build/app .
EXPOSE 8080
CMD [ "./app" ]