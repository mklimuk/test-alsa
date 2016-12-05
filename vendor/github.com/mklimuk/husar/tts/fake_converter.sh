#!/bin/sh

BASEDIR="$( cd "$( dirname "$0" )" && pwd )"
TEST_OGG="$BASEDIR/fake.ogg"

echo "converter params:"
echo "in: $1"
echo "out: $2"

cp "$TEST_OGG" "$2"

exit 0
