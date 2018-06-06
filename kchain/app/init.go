package app

import (
	"github.com/json-iterator/go"
	"github.com/tendermint/tmlibs/log"

	kcfg "kchain/types/cfg"
	kts "kchain/types"
	"fmt"
)

var (
	cfg    = kcfg.GetConfig()
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	state  State
	logger log.Logger
)

func f(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

type Transaction kts.Transaction

func (t *Transaction) Dumps() []byte {
	d, _ := json.Marshal(t)
	return d
}

func NewTransaction() *Transaction {
	return &Transaction{}
}
