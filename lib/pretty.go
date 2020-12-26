package homed

// Pretty is an interface to display a pretty output
type Pretty interface {
	Pretty() string
}

// PrettyPrint pretty prints an struct
func PrettyPrint(p Pretty) string {
	return p.Pretty()
}
