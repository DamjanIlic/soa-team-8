FROM golang:1.25.1-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o purchase-service main.go

EXPOSE 8080

CMD ["./purchase-service"]