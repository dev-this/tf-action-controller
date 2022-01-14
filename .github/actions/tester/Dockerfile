FROM golang:1.17-alpine

# https://www.terraform.io/docs/cli/config/environment-variables.html
ENV TF_IN_AUTOMATION="yesplz"

RUN apk add --no-cache git terraform curl bash
RUN apk add --no-cache kubectl helm --repository=https://dl-cdn.alpinelinux.org/alpine/edge/testing

EXPOSE 8080

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
