package messager

type Messager interface {
	Encode() (*Package, error)
	Decode([]byte) error
	Title() string
}
