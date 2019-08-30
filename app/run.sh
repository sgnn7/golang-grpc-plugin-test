#!/bin/bash -ex

echo "Building plugin..."
go build -o ./testplugin.plugin ../plugin/main.go

echo "Running app..."
go run main.go
