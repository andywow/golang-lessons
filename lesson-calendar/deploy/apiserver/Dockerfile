FROM golang:1.14 AS builder

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/apiserver cmd/apiserver/main.go

FROM alpine:3.11.5

WORKDIR /app

COPY --from=builder /go/src/app/build/apiserver /app/

EXPOSE 8080 9090

ENTRYPOINT ["/app/apiserver"]
CMD ["--configfile", "config.yaml"]
