FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/api .

EXPOSE 8067

CMD ["./api"]
