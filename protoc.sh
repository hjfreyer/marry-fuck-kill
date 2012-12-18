#!/bin/bash -e

SRC="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

(cd $SRC/proto && protoc mfk.proto --go_out=$SRC/go/mfklib)

sed -i.bak 's/code.google.com\/p\/goprotobuf\/proto/github.com\/hjfreyer\/marry-fuck-kill\/go\/third_party\/proto/g' "$SRC/go/mfklib/mfk.pb.go"
rm $SRC/go/mfklib/mfk.pb.go.bak