# kp-grpc-gw-example

An example project for grpc with json gateway with tracing enabled using opentelemettry.

## Frameworks used.

* opentelemetry for tracing which traces from LB(traefik) all throw gateway and spring grpc and jdbc connection.
* gofiber based gateway, note json to protobuf conversion is written manuually,  alternatively you use grpc gateway, since fasthttp doesn't support `h2c`.
  

## How to build.

### Dev

* go to kp-grpc-server and build bakend
```
$ gradle clean build
```

* Replace `kp-dev.local`, with your hostname in `compose-dev.yml`.

* bring up docker container.
```
$ docker compose -f compose-dev.yml up -d
```

### Quick  Build

* Replace `kp-dev.local`, with your hostname in `compose.yml`.
* 
* bring up docker container.
```
$ docker compose -f compose.yml up -d
```

### How  to Validate

* Add  `user`, replace `kp-dev.local`  with hostname.
```
curl --location 'http://kp-dev.local/kp-grpc-gw-example/api/users/1' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "name": "kp",
    "age": 99
}'
```
* List `user`, replace `kp-dev.local`  with hostname.
```
curl --location 'http://kp-dev.local/kp-grpc-gw-example/api/users'
```
* check the tracing in zipkin by opening url `http://<docker-host>:9411/` 
