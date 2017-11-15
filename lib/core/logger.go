package core

type Logger interface {
	Println(a ...interface{})
	Printf(format string, a ...interface{})
}
