package services

import (
	"github.com/gin-gonic/gin"
	kcfg "kchain/types/cfg"
)

func Run() {

	logger = kcfg.GetLogWithKeyVals("module", "app")

	app := gin.Default()

	logger.Info("init urls", "init", "urls")
	InitUrls(app)

	app.Run(cfg().App.Addr)
}
