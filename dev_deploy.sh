#!/bin/bash -e

trap 'kill $(jobs -pr)' SIGINT SIGTERM EXIT

APPENGINE_DIR=${APPENGINE_DIR-"$HOME/tools/google_appengine"}
DEV_APPSERVER="$APPENGINE_DIR/dev_appserver.py"

SRC="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BIN="$SRC/bin/dev"
TMP=$(mktemp -d -t mfk.XXXX)

GODIR="$BIN/github.com/hjfreyer/marry-fuck-kill/go"

rm -rf $BIN

mkdir -p $BIN
mkdir -p $(dirname $GODIR)
ln -s $SRC/go $GODIR

function lns {
    ln -s $SRC/$1 $BIN/$2
}

lns main.go github.com/

lns app.yaml ''
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
guard -i &

$DEV_APPSERVER --use_sqlite $BIN
