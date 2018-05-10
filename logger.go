package pqutil

type logger interface {
	Debugf(string, ...interface{})
	Printf(string, ...interface{})
}

var lg logger

// SetLogger sets a logger on the package that will print
// messages. Must have Printf and Debugf.
func SetLogger(l logger) {
	lg = l
}

func debugf(f string, a ...interface{}) {
	if lg == nil {
		return
	}
	lg.Debugf(f, a...)
}

func logf(f string, a ...interface{}) {
	if lg == nil {
		return
	}
	lg.Printf(f, a...)
}
