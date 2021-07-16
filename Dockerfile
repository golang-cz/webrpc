FROM golang:1.16.0-alpine3.13

RUN apk add --update make bash git

ADD ./ /webrpc

WORKDIR /webrpc
RUN go generate ./...
RUN go build -o /usr/bin/gen ./cmd/webrpc-gen

CMD ["echo", "docker", "run", "golangcz/webrpc", "gen", "-schema=./api.go", "-target=ts", "-client"]
