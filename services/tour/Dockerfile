FROM golang:1.25.1-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o tour-service main.go

EXPOSE 8080

CMD ["./tour-service"]