#!/bin/bash

find . \( -path ./vendor \) -prune -o -name "*.go" -exec goimports -w {} \;
