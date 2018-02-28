#!/usr/bin/env bash

a="""
github.com/gin-gonic/gin
github.com/json-iterator/go
github.com/shopspring/decimal
github.com/spf13/cobra
github.com/spf13/viper
github.com/tendermint/abci
github.com/tendermint/tendermint
github.com/tendermint/tmlibs
golang.org/x/crypto
github.com/ebuchman/fail-test
github.com/fsnotify/fsnotify
github.com/gin-contrib/sse
github.com/go-kit/kit/log
github.com/go-logfmt/logfmt
github.com/gogo/protobuf/proto
github.com/gogo/protobuf
github.com/gorilla/websocket
github.com/hashicorp/hcl
github.com/magiconair/properties
github.com/rcrowley/go-metrics
github.com/pelletier/go-toml
github.com/mitchellh/mapstructure
github.com/mattn/go-isatty
github.com/syndtr/goleveldb
github.com/tendermint/go-wire
github.com/ugorji/go
golang.org/x/net/context
gopkg.in/yaml.v2

github.com/go-stack/stack
github.com/golang/protobuf/proto
github.com/golang/snappy
github.com/spf13/afero
github.com/spf13/cast
github.com/spf13/jwalterweatherman
github.com/spf13/pflag
github.com/tendermint/iavl
golang.org/x/sys/unix
gopkg.in/go-playground/validator.v8
github.com/tendermint/go-crypto
golang.org/x/text/transform
google.golang.org/grpc
google.golang.org/genproto
github.com/tendermint/ed25519
github.com/btcsuite/btcd/btcec
gopkg.in/go-playground/validator.v9
github.com/pkg/errors
github.com/go-playground/universal-translator
github.com/go-playground/locales
"""

echo $a | xargs gopm get -l