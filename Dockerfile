FROM golang:1.27.3@sha256:e519d3b0a5b285e5aa53ab9eef02bc6bbb9cb04608389b743a136a9ce75f0d1e AS builder

WORKDIR /go/src/mikefarah/yq

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" .
# RUN ./scripts/test.sh -- this too often times out in the github pipeline.
RUN ./scripts/acceptance.sh

# Choose alpine as a base image to make this useful for CI, as many
# CI tools expect an interactive shell inside the container
FROM alpine:3@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf AS production
LABEL maintainer="Mike Farah <mikefarah@users.noreply.github.com>"

COPY --from=builder /go/src/mikefarah/yq/yq /usr/bin/yq

WORKDIR /workdir

RUN set -eux; \
  addgroup -g 1000 yq; \
  adduser -u 1000 -G yq -s /bin/sh -h /home/yq -D yq

RUN chown -R yq:yq /workdir

USER yq

ENTRYPOINT ["/usr/bin/yq"]
