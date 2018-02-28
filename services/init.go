package services

import (
	kcfg "kchain/types/cfg"
	tlog "github.com/tendermint/tmlibs/log"

	"github.com/json-iterator/go"

	"github.com/tendermint/tendermint/types"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	cfg = kcfg.GetConfig()
	logger tlog.Logger
	pvfs *types.PrivValidatorFS
)
