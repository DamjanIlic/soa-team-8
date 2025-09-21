# Koristi Go 1.25.1 Alpine image
FROM golang:1.25.1-alpine

# Set radni direktorijum
WORKDIR /app

# Kopiraj go.mod i go.sum i preuzmi dependency-je
COPY go.mod go.sum ./
RUN go mod tidy

# Kopiraj ceo projekat
COPY . .

# Build aplikacije
RUN go build -o auth-service main.go

EXPOSE 8080 50051
CMD ["./auth-service"]

