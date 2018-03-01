package services

import (
	"encoding/hex"
	"net/http"
	"github.com/gin-gonic/gin"

	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/tendermint/types"

	kuc "kchain/utils/crypto"
)

func InitUrls(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 生成节点账号
	router.GET("/gen_validator", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": "ok",
			"data":types.GenPrivValidatorFS(""),
		})
		return
	})

	// 签名服务
	router.POST("/sign", func(c *gin.Context) {

		if d, err := c.GetRawData(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": "error",
				"msg":err.Error(),
			})
		} else {

			s := map[string]string{}
			json.Unmarshal(d, &s)

			d, ks := kuc.MapHasher(s)

			c.JSON(http.StatusOK, gin.H{
				"code": "ok",
				"data":gin.H{
					"hash": hex.EncodeToString(pvfs.PrivKey.Sign(d).Bytes()),
					"keys":ks,
				},
			})
		}

		return
	})


	// 签名服务
	router.POST("/sign_filter", func(c *gin.Context) {

		if d, err := c.GetRawData(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": "error",
				"msg":err.Error(),
			})
		} else {

			s := map[string]string{}
			json.Unmarshal(d, &s)

			d, ks := kuc.MapHasher(s)

			c.JSON(http.StatusOK, gin.H{
				"code": "ok",
				"data":gin.H{
					"hash": hex.EncodeToString(pvfs.PrivKey.Sign(d).Bytes()),
					"keys":ks,
				},
			})
		}

		return
	})

	// 签名验证服务
	router.POST("/verify", func(c *gin.Context) {
		if d, err := c.GetRawData(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": "error",
				"msg":err.Error(),
			})
		} else {

			s := map[string]string{}
			json.Unmarshal(d, &s)

			msg := s["msg"]
			sign_msg := s["sign_msg"]

			if d, e := hex.DecodeString(sign_msg); e != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": "error",
					"msg":err.Error(),
				})
			} else {
				if sig, err := crypto.SignatureFromBytes(d); err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code": "error",
						"msg":err.Error(),
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"code": "ok",
						"data":pvfs.PubKey.VerifyBytes([]byte(msg), sig),
					})
				}
			}
		}

		return
	})

	// 签名验证服务
	router.POST("/hasher", func(c *gin.Context) {
		if d, err := c.GetRawData(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": "error",
				"msg":err.Error(),
			})
		} else {

			s := map[string]string{}
			json.Unmarshal(d, &s)

			d, ks := kuc.MapHasher(s)

			c.JSON(http.StatusOK, gin.H{
				"code": "ok",
				"data":gin.H{
					"hash": d,
					"keys":ks,
				},
			})
		}
		return
	})
}
