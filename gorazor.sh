#!/usr/bin/env sh
node gorazor.js $1 $2
gofmt -w $2
