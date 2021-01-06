FROM golang:1.15-alpine as compiler
ENV GO111MODULE=on
WORKDIR /app
RUN apk add --no-cache make gcc musl-dev linux-headers git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install

FROM alpine:3.12
COPY --from=compiler /go/bin/web3-cli /usr/local/bin/
COPY --from=compiler /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT [ "web3-cli" ]
