#! /bin/bash

go get -d github.com/tendermint/abci

# get the app commit used by tendermint
COMMIT=`bash scripts/glide/parse.sh abci`

echo "Checking out vendored commit for app: $COMMIT"

cd $GOPATH/src/github.com/tendermint/abci
git checkout $COMMIT
glide install
go install ./cmd/...
