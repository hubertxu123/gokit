package zapi

type Option func(l *Logger)

func TraceIdKey(key string) Option {
	return func(l *Logger) {
		l.traceIdKey = key
	}
}
