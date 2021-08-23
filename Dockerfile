# BUILD IMAGE
FROM golang:1.16 as builder

WORKDIR /go/src/app/

COPY . .

ARG version
ARG githubToken

RUN git config --global url."${githubToken}".insteadOf "https://github.com/"
RUN GO111MODULE=on go get github.com/swaggo/swag/cmd/swag
RUN GO111MODULE=on swag init --generalInfo app.go
RUN pkgPath="$(GO111MODULE=on go list -m)/conf.Version=${version}" && GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -ldflags="-X $pkgPath" -a -installsuffix cgo -o app

# PROD/DEV IMAGE
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/app/

COPY --from=builder /go/src/app/app .
COPY --from=builder /go/pkg/mod /go/pkg/mod
COPY --from=builder /go/src/app/env/ env/

CMD ./app -e $env

EXPOSE 8080
