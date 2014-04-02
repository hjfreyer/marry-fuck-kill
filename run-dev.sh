
trap 'kill $(jobs -pr)' SIGINT SIGTERM EXIT

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

BIN=$DIR/bin-dev

rm -rf $BIN
mkdir $BIN
mkdir $BIN/marry-fuck-kill
ln -s $DIR/backend/* $BIN/marry-fuck-kill
ln -s $DIR/frontend/app $BIN/marry-fuck-kill/jsapp

sass --watch frontend/app/css/app.scss &
dev_appserver.py $BIN/marry-fuck-kill
