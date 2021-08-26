FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/petstore/main.go

EXPOSE 5000

ENTRYPOINT ["./main"]
