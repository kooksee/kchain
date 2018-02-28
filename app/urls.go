package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
	//"kchain/utils/validation"

)

func InitUrls(router *gin.Engine) {

	//binding.Validator.RegisterValidation("bookabledate", validation.BookableDate)]

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := router.Group("/v1")
	v1.POST("/metadata", metadata_post)
	v1.GET("/metadata/:dna", metadata_get)
	v1.POST("/license", license_post)
	v1.GET("/license/:name", license_get)
	v1.POST("/validator", validator_post)

}
