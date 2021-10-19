FROM golang:1.17 as builder

WORKDIR /go/src/mikefarah/yq

COPY . /go/src/mikefarah/yq

RUN CGO_ENABLED=0 go build .
RUN ./scripts/test.sh
RUN ./scripts/acceptance.sh

# Choose alpine as a base image to make this useful for CI, as many
# CI tools expect an interactive shell inside the container
FROM alpine:3 as production

RUN mkdir /home/yq/
RUN addgroup -g 1000 yq && \
    adduser -u 1000 -G yq -s /bin/bash -h /home/yq -D yq
RUN chown -R yq:yq /home/yq/

COPY --from=builder /go/src/mikefarah/yq/yq /usr/bin/yq
RUN chmod +x /usr/bin/yq

ARG VERSION=none
LABEL version=${VERSION}

USER yq

WORKDIR /workdir

ENTRYPOINT ["/usr/bin/yq"]
