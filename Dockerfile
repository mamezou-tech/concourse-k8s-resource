FROM golang:alpine as builder

ADD cmd /go/src/cmd
ADD pkg /go/src/pkg
ADD go.mod /go/src

ENV CGO_ENABLED 0

WORKDIR /go/src
RUN go test -cover github.com/mamezou-tech/concourse-k8s-resource/pkg/k8s/...
RUN go build -o /assets/out /go/src/cmd/out
RUN go build -o /assets/in /go/src/cmd/in
RUN go build -o /assets/check /go/src/cmd/check

################ thats our production image

FROM alpine:edge AS resource
#RUN apk --no-cache add \
COPY --from=builder /assets /opt/resource

################ thats our release image

FROM resource AS release
