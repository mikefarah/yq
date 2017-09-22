#!/bin/sh

go get -u github.com/alecthomas/gometalinter
go get -u golang.org/x/tools/cmd/goimports
go get -u github.com/mitchellh/gox
go get -u github.com/kardianos/govendor

# install all the linters
gometalinter --install --update
