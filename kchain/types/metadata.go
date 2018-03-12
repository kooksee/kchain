package types

type Metadata struct {
	PubKey      string `json:"pubkey,omitempty"`
	BlockHash   string `json:"block_hash,omitempty"`
	BlockHeight int64  `json:"block_height,omitempty"`
	Signature   string `json:"signature,omitempty" binding:"required"`
	ID          string `json:"id,omitempty" binding:"required"`

	// 用逗号隔开
	Category    string `json:"category,omitempty" binding:"required"`
	ContentHash string `json:"content_hash,omitempty" binding:"required"`
	Type        string `json:"type,omitempty" binding:"required"`
	Title       string `json:"title,omitempty" binding:"required"`

	// 时间戳
	Created  string      `json:"created,omitempty"`
	Abstract string      `json:"abstract,omitempty"`
	DNA      string      `json:"dna,omitempty"`
	Language string      `json:"language,omitempty"`
	Source   string      `json:"source,omitempty"`
	Extra    interface{} `json:"extra,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	License struct {
		Type   string            `json:"type,omitempty" binding:"required"`
		Params map[string]string `json:"params,omitempty" binding:"required"`
	} `json:"license,omitempty" binding:"required"`
}

func (a *Metadata) Dumps() []byte {
	d, _ := json.Marshal(a)
	return d
}
