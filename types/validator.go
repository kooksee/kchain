package types

type Validator struct {
	PubKey string `json:"pub_key,omitempty" binding:"required"`
	Power  string  `json:"power,omitempty" binding:"required"`
}

func (v *Validator)Dumps() []byte {
	d, _ := json.Marshal(v)
	return d
}
