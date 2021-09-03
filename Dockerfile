FROM golang:1.16-alpine AS builder

WORKDIR /go/src/app
COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN CGO_ENABLED=0 go build -ldflags="-w -s" .

FROM scratch

COPY --from=builder /go/src/app/faksqldb-server /app/
COPY tmp/terraform /usr/bin/terraform

ENTRYPOINT ["/app/faksqldb-server"]
