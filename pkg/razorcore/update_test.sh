#!/bin/sh
for file in testdata/layout/_*.go
do
	mv -f "$file" "${file/\/_//}"
done

for file in testdata/_*.go
do
	mv -f "$file" "${file/\/_//}"
done

cp -f testoptimizegen/*.go testoptimize/
