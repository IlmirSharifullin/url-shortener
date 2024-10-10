FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080

ENV CONFIG_PATH "config/prod.yaml"

RUN go build -o url-shortener ./cmd/url-shortener/main.go

CMD ["./url-shortener"]