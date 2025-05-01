FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o billing-service ./cmd/main.go

EXPOSE 8081

CMD ["./billing-service"]
