FROM golang:alpine as builder
WORKDIR /go/src/github.com/beldpro-ci/subscriber/
RUN apk add --update git build-base bash && go get -u -v github.com/Masterminds/glide

ADD ./ ./
RUN make deps
RUN make linux -j4

FROM alpine:latest
COPY --from=builder /go/src/github.com/beldpro-ci/subscriber/subscriber/subscriber.linux.amd64 /usr/local/bin/subscriber

ENTRYPOINT ["subscriber"]
CMD [ "listen" ]
