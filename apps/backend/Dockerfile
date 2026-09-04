FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /echotalk ./cmd

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=builder --chmod=0755 /echotalk /echotalk

EXPOSE 8080

CMD ["/echotalk"]