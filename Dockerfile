FROM golang:alpine as builder

ADD cmd /go/src/cmd
ADD pkg /go/src/pkg
ADD go.mod /go/src

ENV CGO_ENABLED 0

WORKDIR /go/src
RUN go build -o /assets/out /go/src/cmd/out
RUN go build -o /assets/in /go/src/cmd/in
RUN go build -o /assets/check /go/src/cmd/check

################ thats our production image

FROM alpine:edge AS resource
#RUN apk --no-cache add \
COPY --from=builder /assets /opt/resource


################ thats our test image

#FROM resource AS test
#COPY tests/ /tests
#RUN /tests/integration/all_tests.sh

################ thats our release image

FROM resource AS release
