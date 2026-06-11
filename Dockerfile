FROM golang:1.26.4@sha256:68cb6d68bed024785b69195b89af7ac7a444f27791435f98647edff595aa0479 AS builder

WORKDIR /go/src/mikefarah/yq

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" .
# RUN ./scripts/test.sh -- this too often times out in the github pipeline.
RUN ./scripts/acceptance.sh

# Choose alpine as a base image to make this useful for CI, as many
# CI tools expect an interactive shell inside the container
FROM alpine:3@sha256:a2d49ea686c2adfe3c992e47dc3b5e7fa6e6b5055609400dc2acaeb241c829f4 AS production
LABEL maintainer="Mike Farah <mikefarah@users.noreply.github.com>"

COPY --from=builder /go/src/mikefarah/yq/yq /usr/bin/yq

WORKDIR /workdir

RUN set -eux; \
  addgroup -g 1000 yq; \
  adduser -u 1000 -G yq -s /bin/sh -h /home/yq -D yq

RUN chown -R yq:yq /workdir

USER yq

ENTRYPOINT ["/usr/bin/yq"]
