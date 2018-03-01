package main

import (
	"encoding/hex"
	"fmt"
	crypto "github.com/tendermint/go-crypto"
)

func main() {

	d, _ := hex.DecodeString("016A66FE139DE9CC71DE4B940775504F8EDCB6962C6C6B291A89D703962EFCD3B9")
	if pk, err := crypto.PubKeyFromBytes(d); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(pk)
	}
}
