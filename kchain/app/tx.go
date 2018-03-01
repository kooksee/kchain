package app

import ()

type Tx struct {
	SignPubKey string      `json:"pubkey,omitempty" binding:"required"`
	Signature  string      `json:"sign,omitempty" binding:"required"`
	ID         string      `json:"id,omitempty" binding:"required"`
	Data       interface{} `json:"data,omitempty" binding:"required"`
}
