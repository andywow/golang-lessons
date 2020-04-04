FROM golang:1.14 AS builder

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/sender cmd/sender/main.go

FROM alpine:3.11.5

WORKDIR /app

COPY --from=builder /go/src/app/build/sender /app/

ENTRYPOINT ["/app/sender"]
CMD ["--configfile", "config.yaml"]
