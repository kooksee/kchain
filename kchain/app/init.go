package app

import (
	tlog "github.com/tendermint/tmlibs/log"
	kcfg "kchain/types/cfg"

	"github.com/json-iterator/go"
	"github.com/tendermint/tendermint/types"
)

var (
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	cfg    = kcfg.GetConfig()
	logger tlog.Logger
	pvfs   *types.PrivValidatorFS
)
