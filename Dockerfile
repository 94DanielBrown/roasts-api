FROM golang:1.23-rc-bookworm as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o webApp ./cmd/app && chmod +x webApp

FROM alpine:latest

RUN apk --no-cache add \
    curl \
    aws-cli

RUN mkdir /app

COPY --from=builder /app/webApp /app/webApp

EXPOSE 8000

CMD ["/app/webApp"]
