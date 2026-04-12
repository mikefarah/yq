FROM golang:1.26.2@sha256:2a2b4b5791cea8ae09caecba7bad0bd9631def96e5fe362e4a5e67009fe4ae61 AS builder

WORKDIR /go/src/mikefarah/yq

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" .
# RUN ./scripts/test.sh -- this too often times out in the github pipeline.
RUN ./scripts/acceptance.sh

# Choose alpine as a base image to make this useful for CI, as many
# CI tools expect an interactive shell inside the container
FROM alpine:3@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659 AS production
LABEL maintainer="Mike Farah <mikefarah@users.noreply.github.com>"

COPY --from=builder /go/src/mikefarah/yq/yq /usr/bin/yq

WORKDIR /workdir

RUN set -eux; \
  addgroup -g 1000 yq; \
  adduser -u 1000 -G yq -s /bin/sh -h /home/yq -D yq

RUN chown -R yq:yq /workdir

USER yq

ENTRYPOINT ["/usr/bin/yq"]
