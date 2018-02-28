package cfg

import (
	"sync"
	cfg "github.com/tendermint/tendermint/config"
	nm "github.com/tendermint/tendermint/node"
	dbm "github.com/tendermint/tmlibs/db"
	tlog "github.com/tendermint/tmlibs/log"
	ttypes "github.com/tendermint/tendermint/types"
	c "github.com/tendermint/tendermint/rpc/client"
)

var (
	once sync.Once
	instance *services
)

type appConfig struct {
	Name string
	Addr string
}

type services struct {
	Config *cfg.Config
	App    *appConfig
	Node   *nm.Node
	store  *dbm.GoLevelDB
	log    tlog.Logger
	pk     *ttypes.PrivValidatorFS
	client *c.HTTP
}

