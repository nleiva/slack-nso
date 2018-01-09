FROM golang:1.9.2 as builder

# Needs full path to work with Container's $GOPATH
RUN mkdir -p /go/src/github.com/nleiva/slack-nso/
WORKDIR /go/src/github.com/nleiva/slack-nso/
COPY ./. .

# Disable crosscompiling 
# ENV CGO_ENABLED=0

# Ccompile Linux only
# ENV GOOS=linux

# Build the binary with debug information removed
# RUN go build -ldflags '-w -s' -a -installsuffix cgo -o cli client/main.go

# Not sure if it is worthwhile to build with flags in multi-stage builds
RUN go build -o serv server/main.go

# Could be any of these base images I guess
# FROM scratch
# FROM alpine
FROM gcr.io/distroless/base
COPY ./cert.pem /
COPY ./key.pem /
COPY --from=builder /go/src/github.com/nleiva/slack-nso/serv /

EXPOSE 50051

CMD ["/serv"]