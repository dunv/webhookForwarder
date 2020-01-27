#!/bin/bash

go get -u github.com/golang/protobuf/protoc-gen-go

# Generate .go files
files=proto/*
for fileName in $files
do
    protoc \
        -I=proto \
        --go_out=plugins=grpc:./src \
        $fileName
done

# Move files to the correct location, so they can be imported
cd src

files=$(find . -maxdepth 1 | grep "pb\.go")
mkdir -p pb
for fileName in $files
do
    mv $fileName pb/$fileName
    echo "moved $fileName"
done

cd ..
