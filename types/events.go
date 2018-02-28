package types

type Events struct {
	Desc    string
	Name    string
	Handler func(params string) error
}
