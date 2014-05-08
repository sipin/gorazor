#!/bin/sh
go build vash.go
if [ $? -eq 0 ]; then
    ./vash
else
    echo "Build error"
fi
