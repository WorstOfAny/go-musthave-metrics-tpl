package logger

type Logger interface {
	Log(message ...any)
}

type LoggerAdapter func(message ...any)

func (la LoggerAdapter) Log(message ...any) {
	la(message...)
}
