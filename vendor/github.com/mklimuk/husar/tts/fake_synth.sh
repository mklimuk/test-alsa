#!/bin/sh

BASEDIR="$( cd "$( dirname "$0" )" && pwd )"
TEST_WAV="$BASEDIR/fake.wav"

echo "Syth params:"
echo "text: $1"
echo "file: $2"
echo "tempo: $3"
echo "ogg: $4"

cp "$TEST_WAV" "$2"

exit 0
