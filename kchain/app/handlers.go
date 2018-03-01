package app

import (
	"net/http"
	"github.com/gin-gonic/gin"
	tdata "github.com/tendermint/go-wire/data"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tendermint/types"

	kts "kchain/types"
	kcfg "kchain/types/cfg"
	"encoding/hex"
	"fmt"
	"time"
)

func metadata_post(c *gin.Context) {
	md := &kts.Metadata{}
	if err := c.ShouldBindJSON(md); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "error",
			"msg":  err.Error(),
		})
		return
	}

	// 生成dna信息
	md.DNA = hex.EncodeToString(crypto.Ripemd160([]byte(md.Signature)))

	// 构建tx
	tx := &kts.Transaction{
		PubKey:    hex.EncodeToString(pvfs.PubKey.Bytes()),
		Key:       md.DNA,
		Value:     string(md.Dumps()),
		Path:      "db",
		Timestamp: time.Now().Unix(),
	}
	// 签名
	tx.Signature = hex.EncodeToString(pvfs.PrivKey.Sign([]byte(fmt.Sprintf("%s%d%s", tx.Key, tx.Timestamp, tx.Value))).Bytes())

	if res, err := kcfg.Abci().BroadcastTxCommit(types.Tx(string(tx.Dumps()))); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "error",
			"msg":  err.Error(),
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"code": "ok",
			"data": gin.H{
				"dna": md.DNA,
				"tx":  res,
			},
		})
	}

	return
}

func metadata_get(c *gin.Context) {
	dna := c.Param("dna")
	if dna == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "error",
			"data": "id is null",
		})
		return
	}

	if res, err := kcfg.Abci().ABCIQuery("db", tdata.Bytes(dna)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "error",
			"msg":  err.Error(),
		})
	} else {
		t := map[string]interface{}{}
		json.Unmarshal(res.Response.Value, &t)
		c.JSON(http.StatusOK, gin.H{
			"code": "ok",
			"data": t,
		})
	}

	return
}

func validator_post(c *gin.Context) {
	t := &kts.Validator{}

	if err := c.ShouldBindJSON(t); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "error",
			"msg":  err.Error(),
		})
		return
	}

	// 构建tx
	tx := &kts.Transaction{
		PubKey:    hex.EncodeToString(pvfs.PubKey.Bytes()),
		Key:       "01" + t.PubKey,
		Value:     t.Power,
		Path:      "validator",
		Timestamp: time.Now().Unix(),
	}

	// 签名
	tx.Signature = hex.EncodeToString(pvfs.PrivKey.Sign([]byte(fmt.Sprintf("%s%d%s", tx.Key, tx.Timestamp, tx.Value))).Bytes())

	if res, err := kcfg.Abci().BroadcastTxCommit(types.Tx(tx.Dumps())); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "error",
			"msg":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": "ok",
			"data": res,
		})
	}
}

func license_post(c *gin.Context) {
	t := &kts.License{}
	if err := c.ShouldBindJSON(t); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": "error",
			"msg":  err.Error(),
		})
		return
	}

	// 构建tx
	tx := &kts.Transaction{
		PubKey:    hex.EncodeToString(pvfs.PubKey.Bytes()),
		Key:       t.Type,
		Value:     string(t.Dumps()),
		Path:      "const_db",
		Timestamp: time.Now().Unix(),
	}
	// 签名
	tx.Signature = hex.EncodeToString(pvfs.PrivKey.Sign([]byte(fmt.Sprintf("%s%d%s", tx.Key, tx.Timestamp, tx.Value))).Bytes())
	logger.Error(string(tx.Dumps()))
	if res, err := kcfg.Abci().BroadcastTxCommit(types.Tx(string(tx.Dumps()))); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "error",
			"msg":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": "ok",
			"data": res,
		})
	}
	return
}

func license_get(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "ok",
			"data": "id is null",
		})
		return
	}

	if res, err := kcfg.Abci().ABCIQuery("const_db", tdata.Bytes(name)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "error",
			"msg":  err.Error(),
		})
	} else {
		t := map[string]interface{}{}
		json.Unmarshal(res.Response.Value, &t)
		c.JSON(http.StatusOK, gin.H{
			"code": "ok",
			"data": t,
		})
	}
	return
}
