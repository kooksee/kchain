package services

import (
	"github.com/gin-gonic/gin"
	kcfg "kchain/types/cfg"
	"github.com/tendermint/tendermint/types"
)

func Run() {

	logger = kcfg.GetLogWithKeyVals("module", "app")

	app := gin.Default()

	logger.Info("init urls", "init", "urls")
	InitUrls(app)

	pvfs = types.LoadOrGenPrivValidatorFS(cfg().Config.PrivValidatorFile())
	app.Run(cfg().App.Addr)
}
