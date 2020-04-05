FROM golang:1.14

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV SLEEP_TIMEOUT=10

CMD ["sh", "-c", "sleep ${SLEEP_TIMEOUT} && go test"]
