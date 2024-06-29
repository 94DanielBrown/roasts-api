FROM golang:1.23 as builder

RUN mkdir /app && COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o webApp ./cmd/webApp && chmod +x /app/webApp

FROM alpine:latest

RUN apk --no-cache add \
    curl \
    aws-cli

RUN mkdir /app

COPY --from=builder /app/webApp /app/webApp

EXPOSE 8000

CMD ["/app/webApp"]