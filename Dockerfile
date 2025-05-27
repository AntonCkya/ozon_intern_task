FROM golang:alpine 

WORKDIR /app

COPY . .

RUN go build -o ozon_habr ./cmd/main.go

EXPOSE 8080

CMD ["./ozon_habr", "-s", "p", "-d", "d"]
