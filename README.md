---
description: yq is a lightweight and portable command-line YAML processor
---

# yq

 ![Build](https://github.com/mikefarah/yq/workflows/Build/badge.svg) ![Docker Pulls](https://img.shields.io/docker/pulls/mikefarah/yq.svg) ![Github Releases \(by Release\)](https://img.shields.io/github/downloads/mikefarah/yq/total.svg) ![Go Report](https://goreportcard.com/badge/github.com/mikefarah/yq)

## Install

`yq` has pre-built binaries for most platforms - checkout the [releases page](https://github.com/mikefarah/yq/releases) for the latest build. Alternatively - you can use one of the methods below:

### On MacOS:

```bash
brew install yq
```

### On Windows:

```bash
choco install yq
```

Kindly maintained by @chillum \([https://github.com/chillum/choco-packages/tree/master/yq](https://github.com/chillum/choco-packages/tree/master/yq)\)

### On Ubuntu and other Linux distributions supporting `snap` packages:

```bash
snap install yq
```

#### Snap notes

`yq` installs with with [_strict confinement_](https://docs.snapcraft.io/snap-confinement/6233) in snap, this means it doesn't have direct access to root files. To read root files you can:

```bash
sudo cat /etc/myfile | yq -r - somecommand
```

And to write to a root file you can either use [sponge](https://linux.die.net/man/1/sponge):

```bash
sudo cat /etc/myfile | yq -r - somecommand | sudo sponge /etc/myfile
```

or write to a temporary file:

```bash
sudo cat /etc/myfile | yq -r - somecommand | sudo tee /etc/myfile.tmp
sudo mv /etc/myfile.tmp /etc/myfile
rm /etc/myfile.tmp
```

### On Ubuntu 16.04 or higher from Debian package:

```bash
sudo add-apt-repository ppa:rmescandon/yq
sudo apt update
sudo apt install yq -y
```

Kindly maintained by @rmescandon

### go get:

```text
GO111MODULE=on go get github.com/mikefarah/yq/v3
```

## Docker

Oneshot use:

```bash
docker run --rm -v ${PWD}:/workdir mikefarah/yq yq [flags] <command> FILE...
```

Run commands interactively:

```bash
docker run --rm -it -v ${PWD}:/workdir mikefarah/yq sh
```

It can be useful to have a bash function to avoid typing the whole docker command:

```bash
yq() {
  docker run --rm -i -v ${PWD}:/workdir mikefarah/yq yq $@
}
```

