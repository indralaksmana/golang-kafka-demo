FROM golang:1.17-alpine
RUN apk --no-cache add gcc g++ make ca-certificates curl git openssh tzdata

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

WORKDIR /go/src/${SERVICE_NAME}
COPY . /go/src/${SERVICE_NAME}
ARG CMD_PATH
ENV CMD_PATH=${CMD_PATH}

RUN go get github.com/githubnemo/CompileDaemon

RUN go mod download
RUN go mod verify

ENTRYPOINT CompileDaemon -exclude-dir=".git" -build="go build -v -o /go/bin/app /go/src/${SERVICE_NAME}/${CMD_PATH}/..." -command="/go/bin/app" 