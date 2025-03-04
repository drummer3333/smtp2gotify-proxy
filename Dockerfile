FROM golang:1.22 as builder
RUN git clone https://github.com/drummer3333/smtp2gotify-proxy /go/src/build
WORKDIR /go/src/build
RUN go mod vendor
ENV CGO_ENABLED=0
RUN GOOS=linux go build -mod vendor -a -o smtp2gotify-proxy .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/build/smtp2gotify-proxy /usr/bin/smtp2gotify-proxy
ENTRYPOINT ["smtp2gotify-proxy"]
