package types

type License struct {
	// 元数据
	PubKey    string   `json:"pubkey,omitempty"`
	BlockHash string   `json:"block_hash,omitempty"`
	BlockHeight int64    `json:"block_height,omitempty"`
	Signature string   `json:"signature,omitempty" binding:"required"`
	ID        string   `json:"id,omitempty" binding:"required"`

	// 数据类型
	Type      string   `json:"type,omitempty" binding:"required"`
	Params    map[string]string   `json:"params,omitempty" binding:"required"`
	Desc      string   `json:"desc,omitempty"`
	Created   string   `json:"created,omitempty"`
	Extra     string   `json:"extra,omitempty"`
}

func (a *License)Dumps() []byte {
	d, _ := json.Marshal(a)
	return d
}
