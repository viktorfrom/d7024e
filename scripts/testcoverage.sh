#!/bin/bash

# Run this from root NOT inside /scripts
# run on local machine and not on VM
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
