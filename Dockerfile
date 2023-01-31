FROM golang:1.19.5

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . .

RUN go test ./...

RUN go build -o main ./cmd

EXPOSE 8080

CMD ["./main"]