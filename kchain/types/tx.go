package types

import (
	ktx "kchain/types/tx"
	"github.com/tendermint/go-crypto"
	"github.com/json-iterator/go"
)

type Transaction1 struct {
	SignPubKey string      `json:"pk,omitempty"`
	Signature  string      `json:"sign,omitempty"`
	IsSign     int         `json:"is_sign,omitempty"`
	Method     string      `json:"method,omitempty"`
	Params     interface{} `json:"params,omitempty"`
	Db         *ktx.Db
	Account    *ktx.Account
	Validator  *ktx.Validator
}

type Transaction struct {
	PubKey    string `json:"pubkey,omitempty"`
	Signature string `json:"sign,omitempty"`
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
	Path      string `json:"path,omitempty"`
	Timestamp int64  `json:"time,omitempty"`
}

func (t *Transaction) Dumps() []byte {
	d, _ := json.Marshal(t)
	return d
}

type PrivValidator struct {
	PubKey  crypto.PubKey  `json:"pub_key" mapstructure:"pubkey"`
	PrivKey crypto.PrivKey `json:"priv_key" mapstructure:"pubkey"`
}

func (p *PrivValidator) FromPubKeyString(key string) (crypto.PubKey, error) {
	a := map[string]map[string]string{
		"pub_key": map[string]string{"type": "ed25519", "data": key},
	}

	d, _ := jsoniter.Marshal(a)
	if err := jsoniter.Unmarshal(d, p); err != nil {
		return p.PubKey, err
	} else {
		return p.PubKey, nil
	}
}

func (p *PrivValidator) FromPriKeyString(key string) (crypto.PrivKey, error) {
	a := map[string]map[string]string{
		"priv_key": map[string]string{"type": "ed25519", "data": key},
	}

	d, _ := jsoniter.Marshal(a)
	if err := jsoniter.Unmarshal(d, p); err != nil {
		return p.PrivKey, err
	} else {
		return p.PrivKey, nil
	}
}
