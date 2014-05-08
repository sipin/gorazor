#!/bin/sh
go build vash.go
echo $?
if [ $? -eq 0 ]; then
    echo "Building finished"
    ./vash
else
    echo "Build error"
fi
