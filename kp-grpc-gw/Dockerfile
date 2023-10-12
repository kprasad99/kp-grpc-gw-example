FROM        golang:1.21-alpine3.18 as builder
WORKDIR     /usr/src/kp-grpc-gw
COPY        go.mod /usr/src/kp-grpc-gw/go.mod
COPY        go.sum /usr/src/kp-grpc-gw/go.sum
RUN         go mod download
COPY        api /usr/src/kp-grpc-gw/api
COPY        main.go /usr/src/kp-grpc-gw/main.go
RUN         CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o kp-grpc-gw

FROM        gcr.io/distroless/static:nonroot
COPY        --from=builder /usr/src/kp-grpc-gw/kp-grpc-gw /usr/local/bin/kp-grpc-gw

ENV WEB_PORT=8080

VOLUME /static
EXPOSE ${WEB_PORT}

USER 65532:65532
ENTRYPOINT  ["/usr/local/bin/kp-grpc-gw"]
