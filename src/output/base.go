package output

type Output interface {
	Write([]byte) error
}
