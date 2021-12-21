#!/bin/bash -eu
#
# Copyright (C) 2021 Roberto Mier Escandón <rmescandon@gmail.com>
#
# This script creates a .deb package file with yq valid for ubuntu 20.04 by default
# You can pass 

DOCKER_IMAGE_NAME=yq-deb-builder
DOCKER_IMAGE_TAG=$(git describe --always --tags)
OUTPUT=
GOVERSION="1.17.4"
KEYID=
MAINTAINER=
PPA=
VERSION=
DISTRIBUTION=
DO_SIGN=
PASSPHRASE=

show_help() {
  echo "  usage: $(basename "$0") VERSION [options...]"
  echo ""
  echo "  positional arguments"
  echo "    VERSION"
  echo ""
  echo "  optional arguments:"
  echo "    -h, --help		              Shows this help"
  echo "    -d, --distribution DISTRO   The distribution to use for the changelog generation. If not provided, last changelog entry"
  echo "                                  distribution is considered"
  echo "    --goversion VERSION         The version of golang to use. Default to $GOVERSION"
  echo "    -k, --sign-key KEYID        Sign the package sources with the provided gpg key id (long format). When not provided this"
  echo "                                  paramater, the generated sources are not signed"
  echo "    -s, --sign                  Sign the package sources with a gpg key of the maintainer"
  echo "    -m, --maintainer WHO        The maintainer used as author of the changelog. git.name and git.email (see git config) is"
  echo "                                  the considered format"
  echo "    -o DIR, --output DIR        The path where leaving the generated debian package. Default to a temporary folder if not set"
  echo "    -p, --push PPA              Push resultant files to indicated ppa. This option should be given along with a signing key."
  echo "                                  Otherwise, the server could reject the package building"
  echo "    --passphrase PASSPHRASE     Passphrase to decrypt the signage key"
  exit 1
}
# read input args
while [ $# -ne 0 ]; do
  case $1 in
    -h|--help)
      show_help
      ;;
    -d|--distribution)
      shift
      DISTRIBUTION="$1"
      ;;
    --goversion)
      shift
      GOVERSION="$1"
      ;;
    -k|--sign-key)
      shift
      DO_SIGN='y'
      KEYID="$1"
      ;;
    -s|--sign)
      DO_SIGN='y'
      ;;
    -m|--maintainer)
      shift
      MAINTAINER="$1"
      ;;
    -o|--output)
      shift
      OUTPUT="$1"
      ;;
    -p|--push)
      shift
      PPA="$1"
      ;;
    --passphrase)
      shift
      PASSPHRASE="$1"
      ;;
    *)
      if [ -z "$VERSION" ]; then
        VERSION="$1"
      else
        show_help
      fi
  esac
  shift
done

[ -n "$VERSION" ] || (echo "error - you have to provide a version" && show_help)

if [ -n "$OUTPUT" ]; then
  OUTPUT="$(realpath "$OUTPUT")"
  mkdir -p "$OUTPUT"
else
  # Temporary folder where leaving the built deb package in case that output folder is not provided
  OUTPUT="$(mktemp -d)"
fi 

# Create temporary folder with all the artifacts to create and deploy the docker image
srcdir="$(realpath "$(dirname "$0")"/..)"
blddir="$(cd "${srcdir}" && mkdir -p build && cd build && echo "$(pwd)")"

cleanup() {
  rm -f "${blddir}/build.sh" || true
  rm -f "${blddir}/Dockerfile" || true
  docker rmi "${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}" -f > /dev/null 2>&1 || true
}
trap cleanup EXIT INT

cat << EOF > ${blddir}/build.sh
#!/bin/bash 
set -e -o pipefail

PATH=$PATH:/usr/local/go/bin
export GPG_TTY=$(tty)

go mod vendor

### bump debian/changelog

# maintainer
export DEBEMAIL="$MAINTAINER"
if [ -z "$MAINTAINER" ]; then
  export DEBEMAIL="\$(dpkg-parsechangelog -S Maintainer)"
fi

# prepend a 'v' char to complete the tag name from where calculating the changelog
SINCE="v\$(dpkg-parsechangelog -S Version)"

# distribution
DISTRIBUTION="$DISTRIBUTION"
if [ -z "$DISTRIBUTION" ]; then
  DISTRIBUTION="\$(dpkg-parsechangelog -S Distribution)"
fi

# generate changelog
gbp dch --ignore-branch --no-multimaint -N "$VERSION" -s "\$SINCE" -D "\$DISTRIBUTION"

# using -d to prevent failing when searching for golang dep on control file
params=("-d" "-S")

# add the -sa option for signing along with the key to use when provided key id
if [ -n "$DO_SIGN" ]; then
  params+=("-sa")

  # read from gpg the key id associated with the maintainer if not provided explicitly
  if [ -z "$KEYID" ]; then
    KEYID="\$(gpg --list-keys "\$(dpkg-parsechangelog -S Maintainer)" | head -2 | tail -1 | xargs)"
  else
    KEYID="$KEYID"
  fi
  
  params+=("--sign-key="\$KEYID"")

  if [ -n "$PASSPHRASE" ]; then
    gpg-agent --verbose --daemon --options /home/yq/.gnupg/gpg-agent.conf --log-file /tmp/gpg-agent.log --allow-preset-passphrase --default-cache-ttl=31536000
    KEYGRIP="\$(gpg --with-keygrip -k "\$KEYID" | grep 'Keygrip = ' | cut -d'=' -f2 | head -1 | xargs)"
    /usr/lib/gnupg/gpg-preset-passphrase --preset --passphrase "$PASSPHRASE" "\$KEYGRIP"
  fi

else
  params+=("-us" "-uc")
fi

debuild  \${params[@]}
mv ../yq_* /home/yq/output

echo ""
echo -e "\tfind resulting package at: "$OUTPUT""

# publish to ppa whether given
if [ -n "$PPA" ]; then
  dput ppa:"$PPA" "$OUTPUT"/yq_*.changes
fi
EOF
chmod +x "${blddir}"/build.sh

cat << EOF > ${blddir}/Dockerfile
FROM bitnami/minideb:bullseye as base
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV DEBIAN_FRONTEND noninteractive
ENV GO111MODULE on
ENV GOMODCACHE /home/yq/go

RUN set -e \
  && sed -i -- 's/# deb-src/deb-src/g' /etc/apt/sources.list \
  && apt-get -qq update

# install golang on its $GOVERSION
FROM base as golang
RUN apt-get -qq -y --no-install-recommends install \
    ca-certificates \
    wget
RUN wget "https://golang.org/dl/go${GOVERSION}.linux-amd64.tar.gz" -4
RUN tar -C /usr/local -xvf "go${GOVERSION}.linux-amd64.tar.gz"

FROM base
RUN apt-get -qq -y --no-install-recommends install \
    build-essential \
    debhelper \
    devscripts \
    fakeroot \
    git-buildpackage \
    gpg-agent \
    libdistro-info-perl \
    pandoc \
    rsync \
    sensible-utils && \
  apt-get clean && \
  rm -rf /tmp/* /var/tmp/*

COPY --from=golang /usr/local/go /usr/local/go

# build debian package as yq user
RUN useradd -ms /bin/bash yq && \
  mkdir /home/yq/src && chown -R yq: /home/yq/src && \
  mkdir /home/yq/output && chown -R yq: /home/yq/output

ADD ./build/build.sh /usr/bin/build.sh
RUN chmod +x /usr/bin/build.sh && chown -R yq: /usr/bin/build.sh

USER yq

WORKDIR /home/yq/src
VOLUME ["/home/yq/src"]

# dir where output packages are finally left
VOLUME ["/home/yq/output"]

CMD ["/usr/bin/build.sh"]
EOF

DOCKER_BUILDKIT=1 docker build --pull -f "${blddir}"/Dockerfile -t "${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}" .

docker run --rm -i \
  -v "${srcdir}":/home/yq/src:delegated \
  -v "${OUTPUT}":/home/yq/output \
  -v "${HOME}"/.gnupg:/home/yq/.gnupg:delegated \
  -u "$(id -u)" \
  "${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"
