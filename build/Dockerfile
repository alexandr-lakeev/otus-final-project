FROM golang:1.16.13 as build

ENV BIN_FILE /opt/previewer/previewer-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/main*

FROM alpine:3.9

ENV BIN_FILE "/opt/previewer/previewer-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/previewer/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

CMD sleep 3; ${BIN_FILE} -config ${CONFIG_FILE}
