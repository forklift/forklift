package engine

const (
	ERR int = iota
	WARN
	INFO
)

type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Print(args ...interface{})
}
