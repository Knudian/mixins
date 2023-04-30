#!/usr/bin/env bash

# Run tests for the package

go test -coverprofile cover.out .

# Exports an html based coverage

go tool cover -html=cover.out