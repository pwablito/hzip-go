package output

type Output interface {
	Write([]byte) error
	Open() error
	Close() error
}
