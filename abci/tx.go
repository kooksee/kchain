package abci

import (
	crypto "github.com/tendermint/go-crypto"

	//c1 "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"strings"

	"fmt"
	"time"
	"encoding/hex"
)

func (t *Transaction) FromBytes(bs []byte) error {
	return json.Unmarshal(bs, t)
}

func (t *Transaction) ToBytes() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Transaction) Verify() error {

	if strings.Compare(t.Signature, "") == 0 || strings.Compare(t.PubKey, "") == 0 {
		return errors.New("sign or pubkey is null")
	}

	// 事务超过两分钟没有被确认，则认定为超时
	if time.Now().Unix() - t.Timestamp > int64(time.Minute * 1) {
		return errors.New("transaction timeout")
	}

	d, _ := hex.DecodeString(t.PubKey)
	if pk, err := crypto.PubKeyFromBytes(d); err != nil {
		return err
	} else {
		d, _ := hex.DecodeString(t.Signature)
		if sig, err := crypto.SignatureFromBytes(d); err != nil {
			return err
		} else {
			_sign_msg := []byte(fmt.Sprintf("%s%d%s", t.Key, t.Timestamp, t.Value))
			if !pk.VerifyBytes(_sign_msg, sig) {
				return errors.New("transaction verify false")
			}
		}
	}
	return nil
}
