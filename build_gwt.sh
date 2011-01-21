#!/bin/bash -e

cd MfkWeb
ant
cd ..
cd MfkMaker
ant
cd ..

mkdir -p marry-fuck-kill/static/gwt

cp -r MfkWeb/war marry-fuck-kill/static/gwt/web
cp -r MfkMaker/war marry-fuck-kill/static/gwt/maker

