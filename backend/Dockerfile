FROM golang:1.17.2-alpine3.14

WORKDIR /go/src/app
COPY . .

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

RUN go mod tidy
RUN go build -ldflags "-s -w" -o chefbook-server ./cmd/main.go
CMD /wait && ./chefbook-server