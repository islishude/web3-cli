# syntax=docker/dockerfile:1
FROM golang:1.22.1-alpine as compiler
RUN apk add --no-cache make gcc musl-dev linux-headers git ca-certificates
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go install

FROM alpine:3.19.1
RUN apk add --no-cache jq
COPY --from=compiler /go/bin/web3-cli /usr/local/bin/
COPY --from=compiler /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT [ "web3-cli" ]
