#!/bin/bash
go build
./splits-go-schema-codegen ../splits-go-api $@
