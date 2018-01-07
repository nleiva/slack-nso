# Slack-NSO

Small PoC to control NSO via Slack

![slack-nso](slack-nso.png)

## Compiling protocol buffers

´´´shell
protoc --go_out=plugins=grpc:. comm.proto
´´´

## Generating SSL Certificates

- Local testing

´´´shell
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj '/CN=localhost'
´´´

- Online testing

CN needs to be replace by your gRPC Server

´´´shell
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj '/CN=grpc.nleiva.com'
´´´

## Generating binaries for Client and Server

- gRPC client: `go build -o cli client/main.go`
- gRPC server: `go build -o serv server/main.go`

## Enviroment variables requiered

- SLACK_TOKEN
- NSO_SERVER
- NSO_DEVICE
- NSO_USER
- NSO_PASSWORD