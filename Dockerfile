ARG GITHUB_PATH=github.com/inst-api/poster
ARG CONFIG_MODE=prod

FROM golang:1.20-alpine3.17 AS builder
RUN apk add --no-cache  --update make git curl tzdata
COPY . /home/${GITHUB_PATH}
WORKDIR /home/${GITHUB_PATH}
ENV GOPROXY=direct
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/rest-api ./cmd/rest-api


FROM alpine:latest as server
ARG CONFIG_MODE
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/rest-api .
COPY --from=builder /home/${GITHUB_PATH}/deploy ./deploy
COPY --from=builder /home/${GITHUB_PATH}/gen/http ./gen/http

ENV TZ=Europe/Moscow
ENV CONFIG_MODE=${CONFIG_MODE}
CMD ["./rest-api", "--debug"]