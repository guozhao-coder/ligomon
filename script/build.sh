#!/bin/bash
echo "-----build start-------"
go build -o _output/ligo ../cmd/main.go
echo "-----build over-------"

sudo  ./_output/ligo

