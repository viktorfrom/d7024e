#!/bin/bash

# Run this from root NOT inside /scripts
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
