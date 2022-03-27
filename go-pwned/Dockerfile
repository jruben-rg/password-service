FROM golang:1.18-alpine AS build

WORKDIR /src/go-pwned/
ADD config /src/go-pwned/config
ADD handlers /src/go-pwned/handlers
ADD metric /src/go-pwned/metric
ADD middleware /src/go-pwned/middleware
ADD password /src/go-pwned/password
ADD pwned /src/go-pwned/pwned
COPY go.* main.go /src/go-pwned/
RUN go mod tidy && CGO_ENABLED=0 go build -o /bin/go-pwned

ENTRYPOINT ["/bin/go-pwned"]