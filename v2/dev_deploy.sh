#!/bin/bash -e

trap 'kill $(jobs -pr)' SIGINT SIGTERM EXIT

DIR="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

mkdir -p $DIR/bin/dev
mkdir -p $DIR/bin/static



xeyes&
xeyes&

pause