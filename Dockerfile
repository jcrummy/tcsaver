FROM golang:1.14 AS builder

WORKDIR /src
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .
WORKDIR /src/cmd
RUN go build -o /tcsaver

#tcsaver
FROM scratch

CMD ["/tcsaver"]
#VOLUME ["/certs", "/private", "/config.yaml", "/acme.json"]
VOLUME ["/certs", "/private"]
COPY --from=builder /tcsaver .
