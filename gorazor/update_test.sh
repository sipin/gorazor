#!/bin/sh
for file in test/layout/_*.go
do
	mv -f "$file" "${file/\/_//}"
done

for file in test/_*.go
do
	mv -f "$file" "${file/\/_//}"
done

cp -f test_opt_gen/*.go test_opt/
