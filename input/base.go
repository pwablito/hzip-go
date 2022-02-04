package input

type Input interface {
	GetData() ([]byte, error)
}
