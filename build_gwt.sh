#!/bin/bash -e

cd MfkMaker
ant
cd ..

rm -rf marry-fuck-kill/static/gwt

mkdir -p marry-fuck-kill/static/gwt

cp -r MfkMaker/war/mfkmaker/ marry-fuck-kill/static/gwt/mfkmaker/

