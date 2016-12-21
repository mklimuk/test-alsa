#!/bin/sh

# Compiles go sources inside docker container and installs both arm and x64 version into dist

# get script's location path
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd $DIR > /dev/null

docker-compose -f compile.yml up
docker-compose -f compile.yml rm -vf
