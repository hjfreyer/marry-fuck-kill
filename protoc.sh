#!/bin/bash -e

SRC="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

(cd $SRC/proto && protoc mfk.proto --go_out=$SRC/go/mfklib)
