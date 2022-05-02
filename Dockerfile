# syntax=docker/dockerfile:1.3
FROM --platform=${BUILDPLATFORM} golang:1.18-alpine as compiler
WORKDIR /app
RUN apk add --no-cache make gcc musl-dev linux-headers git ca-certificates
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go install

FROM --platform=${BUILDPLATFORM} alpine:3.15
COPY --from=compiler /go/bin/web3-cli /usr/local/bin/
COPY --from=compiler /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT [ "web3-cli" ]
