package app

import (
	"github.com/gin-gonic/gin"
	kcfg "kchain/types/cfg"
	"github.com/tendermint/tendermint/types"
)

func Run() {

	logger = kcfg.GetLogWithKeyVals("module", "app")
	pvfs = types.LoadOrGenPrivValidatorFS(cfg().Config.PrivValidatorFile())

	app := gin.Default()

	logger.Info("init urls", "init", "urls")
	InitUrls(app)
	app.Run(cfg().App.Addr)
}