#!/bin/bash -e

cd MfkWeb
ant
cd ..
cd MfkMaker
ant
cd ..

rm -rf marry-fuck-kill/static/gwt

mkdir -p marry-fuck-kill/static/gwt

cp -r MfkWeb/war/mfkweb/ marry-fuck-kill/static/gwt/mfkweb/
cp -r MfkMaker/war/mfkmaker/ marry-fuck-kill/static/gwt/mfkmaker/

