package input

type Input interface {
	Read() ([]byte, error)
}
