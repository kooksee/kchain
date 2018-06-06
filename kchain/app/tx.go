package app

import (
	"github.com/tendermint/go-crypto"

	"github.com/pkg/errors"
	"encoding/hex"
	"fmt"
	"kchain/types/cnst"
)

// FromBytes 解析Transaction
func (t *Transaction) FromBytes(bs []byte) error {
	return json.Unmarshal(bs, t)
}

// ToBytes Marshal
func (t *Transaction) ToBytes() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Transaction) DecodeValues() map[string]interface{} {
	d := map[string]interface{}{}
	json.UnmarshalFromString(t.Values, &d)
	return d
}

func (t *Transaction) Verify() error {

	if t.Signature == "" || t.PubKey == "" {
		return errors.New("sign or pubkey is null")
	}
	// 检查发送tx的节点有没有在区块链中
	if !state.db.Has([]byte(cnst.ValidatorPrefix + t.PubKey)) {
		return errors.New(f("the node %s does not exist", t.PubKey))
	}

	// 区块签名验证
	d, _ := hex.DecodeString(t.PubKey)
	if pk, err := crypto.PubKeyFromBytes(d); err != nil {
		return err
	} else {
		d, _ := hex.DecodeString(t.Signature)
		if sig, err := crypto.SignatureFromBytes(d); err != nil {
			return err
		} else {
			signMsg := crypto.Ripemd160([]byte(fmt.Sprintf("%s%s%d", t.Values, t.Path, t.Timestamp)))
			if !pk.VerifyBytes(signMsg, sig) {
				return errors.New("transaction verify false")
			}
		}
	}
	return nil
}
