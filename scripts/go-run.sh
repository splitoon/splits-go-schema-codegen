#!/bin/bash
go build
if [ -d "$1" ]; then
  ./splits-go-schema-codegen $1
else
  ./splits-go-schema-codegen ../splits-go-api/db/
fi
