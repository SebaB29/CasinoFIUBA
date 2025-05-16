FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "./cmd/server/main.go"]
