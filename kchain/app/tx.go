package app

import (
	"github.com/tendermint/go-crypto"

	//c1 "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"encoding/hex"
	"fmt"
)

// FromBytes 解析Transaction
func (t *Transaction) FromBytes(bs []byte) error {
	return json.Unmarshal(bs, t)
}

// ToBytes Marshal
func (t *Transaction) ToBytes() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Transaction) Verify() error {

	if t.Signature == "" || t.PubKey == "" {
		return errors.New("sign or pubkey is null")
	}

	// 检查发送tx的节点有没有在区块链中
	//if !state.db.Has([]byte(cnst.ValidatorPrefix + t.PubKey)) {
	//	return errors.New(f("the node %s does not exist", t.PubKey))
	//}

	// 事务超过一分钟没有被确认，则认定为超时
	//if time.Now().Unix()-t.Timestamp > int64(time.Minute*1) {
	//	return errors.New("transaction timeout")
	//}

	// 区块签名验证
	d, _ := hex.DecodeString(t.PubKey)
	if pk, err := crypto.PubKeyFromBytes(d); err != nil {
		return err
	} else {
		d, _ := hex.DecodeString(t.Signature)
		if sig, err := crypto.SignatureFromBytes(d); err != nil {
			return err
		} else {
			signMsg := crypto.Ripemd160([]byte(fmt.Sprintf("%s%d%s", t.Key, t.Timestamp, t.Value)))
			if !pk.VerifyBytes(signMsg, sig) {
				return errors.New("transaction verify false")
			}
		}
	}
	return nil
}
