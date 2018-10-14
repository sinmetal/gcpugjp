#!/bin/sh
cd `dirname $0`
npm run-script build
cp -pR dist/* ../static