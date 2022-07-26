#!/usr/bin/env bash
set -Eeuo pipefail
shopt -s expand_aliases


# check that docker is installed
if ! command -v docker &> /dev/null
then
    echo "This script requires docker to run go-swagger, please install it or use any of the installation methods described in https://goswagger.io/ to generate files manually"
    exit
fi

# pull go-swagger from its registry and create an alias for it
docker pull quay.io/goswagger/swagger
alias swagger="docker run --rm -it -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger"

# script directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# generate server files
swagger generate server -A grants-program -t ${DIR}/gen -f ${DIR}/open-api-specifications/grants-program/grants-program.yaml -P models.Principal --skip-tag-packages
