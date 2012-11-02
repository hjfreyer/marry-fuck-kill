#!/bin/bash -e

trap 'kill $(jobs -pr)' SIGINT SIGTERM EXIT

DEV_APPSERVER="/Users/hjfreyer/tools/google_appengine/dev_appserver.py"

DIR="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BIN="$DIR/bin/dev"
TMP=$(mktemp -d -t mfk)

rm -rf $BIN

mkdir -p $BIN
mkdir -p $BIN/static

ln -s $DIR/app.yaml $BIN

ln -s $DIR/go $BIN/gomfk
ln -s $DIR/third_party/closure-library $BIN/static/closure-library

ln -s $DIR/js/script.js $BIN/static/

sass \
    --watch $DIR/stylesheets/:$BIN/static/ \
    --cache-location $TMP > $TMP/sass.log &

$DEV_APPSERVER $BIN