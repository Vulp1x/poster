ARG CONFIG_MODE=prod

FROM golang:1.20-alpine3.17 AS builder
ARG GITHUB_PATH=github.com/inst-api/poster
RUN apk add --no-cache  --update make git curl tzdata
COPY . /home/${GITHUB_PATH}

WORKDIR /home/${GITHUB_PATH}
#ENV GOPROXY=direct
RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@v1.9.1
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags="all=-N -l" -o ./bin/rest-api ./cmd/rest-api

FROM alpine:latest as server
ARG CONFIG_MODE
ARG GITHUB_PATH=github.com/inst-api/poster
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/rest-api .
COPY --from=builder /home/${GITHUB_PATH}/deploy ./deploy
COPY --from=builder /home/${GITHUB_PATH}/gen/http ./gen/http
COPY --from=builder /go/bin/dlv /

EXPOSE 8090 40000

ENV TZ=Europe/Moscow
ENV CONFIG_MODE=${CONFIG_MODE}

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./rest-api", "debug"]
#CMD [ "/bin/sh"]