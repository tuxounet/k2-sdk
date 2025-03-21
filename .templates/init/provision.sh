#!/bin/sh
K2_VERSION=$(curl --silent "https://api.github.com/repos/tuxounet/k2/tags" | jq -r '.[0].name')
go install github.com/tuxounet/k2@${K2_VERSION}
$(go env GOPATH)/bin/k2 apply