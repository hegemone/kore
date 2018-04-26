FROM golang:1.8.3-alpine3.6

COPY ./ /go/src/github.com/hegemone/kore
WORKDIR /go/src/github.com/hegemone/kore
ENV KORE_CONFIG /etc/kore.yaml
ENV GOOGLE_SERVICE_ACCOUNT="/auth/jwt.json"

RUN apk update && apk add git make gcc musl-dev && go get github.com/golang/dep/cmd/dep
RUN make build
RUN install -m +x /go/src/github.com/hegemone/kore/build/kore /usr/bin/kore
RUN mkdir -p /usr/lib/kore && install /go/src/github.com/hegemone/kore/build/*.so /usr/lib/kore
RUN install /go/src/github.com/hegemone/kore/config.yaml /etc/kore.yaml

CMD ["/usr/bin/kore"]
