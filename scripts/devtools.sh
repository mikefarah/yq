#!/bin/sh
set -ex
go mod download golang.org/x/tools@v0.44.0
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/6008b81b81c690c046ffc3fd5bce896da715d5fd/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.11.3
curl -sSfL https://raw.githubusercontent.com/securego/gosec/424fc4cd9c82ea0fd6bee9cd49c2db2c3cc0c93f/install.sh | sh -s v2.22.11

TYPOS_VERSION=v1.47.2
TYPOS_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
TYPOS_ARCH=$(uname -m)
case "${TYPOS_ARCH}" in
  x86_64) TYPOS_ARCH=x86_64 ;;
  aarch64|arm64) TYPOS_ARCH=aarch64 ;;
  *)
    echo "unsupported architecture for typos: ${TYPOS_ARCH}"
    exit 1
    ;;
esac
case "${TYPOS_OS}" in
  linux) TYPOS_TARGET="${TYPOS_ARCH}-unknown-linux-musl" ;;
  darwin) TYPOS_TARGET="${TYPOS_ARCH}-apple-darwin" ;;
  *)
    echo "unsupported OS for typos: ${TYPOS_OS}"
    exit 1
    ;;
esac
TYPOS_ARCHIVE="typos-${TYPOS_VERSION}-${TYPOS_TARGET}.tar.gz"
TYPOS_URL="https://github.com/crate-ci/typos/releases/download/${TYPOS_VERSION}/${TYPOS_ARCHIVE}"
case "${TYPOS_TARGET}" in
  aarch64-apple-darwin) TYPOS_SHA256=23ca24a9186b5cb395b5f6c8eea8cdb02911c8980833e016454b56e90c3bd474 ;;
  aarch64-unknown-linux-musl) TYPOS_SHA256=596d5c6b9ecf34307f68bea649178c5b45a4398fe3a1fcef9598e85aa2ccb742 ;;
  x86_64-apple-darwin) TYPOS_SHA256=469a2d9fc894b0cdcec6e4fa3719b4c4638e195feee6517d4845450f8e8985c6 ;;
  x86_64-unknown-linux-musl) TYPOS_SHA256=7aef58932fc123b4cf4b40d86468e89a3297d80169051d7cfd13a235e05fc426 ;;
  *)
    echo "unsupported typos target: ${TYPOS_TARGET}"
    exit 1
    ;;
esac
TYPOS_TMPDIR=$(mktemp -d)
curl -sSfL "${TYPOS_URL}" -o "${TYPOS_TMPDIR}/${TYPOS_ARCHIVE}"
TYPOS_ACTUAL_SHA256=$(sha256sum "${TYPOS_TMPDIR}/${TYPOS_ARCHIVE}" 2>/dev/null | cut -d' ' -f1)
if [ -z "${TYPOS_ACTUAL_SHA256}" ]; then
  TYPOS_ACTUAL_SHA256=$(shasum -a 256 "${TYPOS_TMPDIR}/${TYPOS_ARCHIVE}" | cut -d' ' -f1)
fi
if [ "${TYPOS_ACTUAL_SHA256}" != "${TYPOS_SHA256}" ]; then
  echo "typos archive checksum mismatch: expected ${TYPOS_SHA256}, got ${TYPOS_ACTUAL_SHA256}"
  exit 1
fi
tar xzf "${TYPOS_TMPDIR}/${TYPOS_ARCHIVE}" -C "${TYPOS_TMPDIR}"
install -m 755 "${TYPOS_TMPDIR}/typos" "$(go env GOPATH)/bin/typos"
rm -rf "${TYPOS_TMPDIR}"
