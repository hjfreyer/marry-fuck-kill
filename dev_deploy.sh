#!/bin/bash -e

trap 'kill $(jobs -pr)' SIGINT SIGTERM EXIT

DEV_APPSERVER="/Users/hjfreyer/tools/google_appengine/dev_appserver.py"

SRC="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BIN="$SRC/bin/dev"
TMP=$(mktemp -d -t mfk)

rm -rf $BIN

mkdir -p $BIN

function lns {
    ln -s $SRC/$1 $BIN/$2
}

lns app.yaml ''
lns go gomfk
lns templates templates

mkdir $BIN/generated_templates
lns templates/js_include_debug.html generated_templates/js_include.html

mkdir -p $BIN/static
lns js static/js
lns third_party/closure-library static/closure-library
lns assets static/assets

sass \
    --watch $SRC/stylesheets/:$BIN/static/ \
    --cache-location $TMP > $TMP/sass.log &

$DEV_APPSERVER $BIN