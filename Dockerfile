FROM golang:1.19.5

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]