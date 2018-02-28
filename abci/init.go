package abci

import (
	"github.com/json-iterator/go"
	"github.com/tendermint/tmlibs/log"
	"github.com/tendermint/iavl"

	kts "kchain/types"
	kcfg "kchain/types/cfg"
)

var (
	cfg = kcfg.GetConfig()
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	state  *iavl.VersionedTree
	logger log.Logger
)



type Transaction kts.Transaction

func NewTransaction() *Transaction {
	return &Transaction{
	}
}
