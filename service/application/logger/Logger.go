package logger

type Logger interface {
	Info(message ...interface{})
	Error(message ...interface{})
	Fatal(message ...interface{})
}
