#!/bin/sh
go install github.com/tuxounet/k2@v0.9.12
$(go env GOPATH)/bin/k2 apply