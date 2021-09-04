FROM alpine:3.14
WORKDIR /workbench

# https://www.terraform.io/docs/cli/config/environment-variables.html
ENV TF_IN_AUTOMATION="yesplz"

RUN apk add --no-cache git terraform curl bash
RUN apk add --no-cache kubectl helm --repository=https://dl-cdn.alpinelinux.org/alpine/edge/testing

WORKDIR /go/src/app

RUN curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash && \
    mv /go/src/app/kustomize /usr/bin/kustomize

COPY server terraform-gha-controller

EXPOSE 8080

CMD ["/go/src/app/terraform-gha-controller"]
