#!/bin/bash
go build || { echo 'failed to build' ; exit 1; }
./splits-go-schema-codegen ../splits-go-api $@
