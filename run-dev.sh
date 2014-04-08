
trap 'kill $(jobs -pr)' SIGINT SIGTERM EXIT

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

BIN=$DIR/bin-dev


rm -rf $BIN
mkdir -p $BIN/

ln -s $DIR/backend/* $BIN/

mkdir -p $BIN/static2/app/
ln -s $DIR/frontend/app/* $BIN/static2/app/
rm $BIN/static2/app/css
mkdir -p $BIN/static2/app/css
sass --watch $DIR/frontend/app/css/app.scss:$BIN/static2/app/css/app.css &

cd $DIR
bower update
cd -
ln -s $DIR/bower_components $BIN/static2/bower

dev_appserver.py $BIN
