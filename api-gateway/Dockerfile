FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GODEBUG=netdns=go+1

ENV CGO_ENABLED=0

RUN go build -o api-gateway ./cmd/main.go

EXPOSE 8088

CMD ["./api-gateway"]
