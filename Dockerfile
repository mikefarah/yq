FROM golang:1.13 as builder

WORKDIR /go/src/mikefarah/yq

# cache devtools
COPY ./scripts/devtools.sh /go/src/mikefarah/yq/scripts/devtools.sh
RUN ./scripts/devtools.sh

COPY . /go/src/mikefarah/yq

RUN CGO_ENABLED=0 make local build

# Choose alpine as a base image to make this useful for CI, as many
# CI tools expect an interactive shell inside the container
FROM alpine:3.8 as production

COPY --from=builder /go/src/mikefarah/yq/yq /usr/bin/yq
RUN chmod +x /usr/bin/yq

ARG VERSION=none
LABEL version=${VERSION}

WORKDIR /workdir
