package ports

type LoggingService interface {
	Debug(msg string)
	Debugf(template string, msg string)
	Info(msg string)
	Infof(template string, msg string)
	Warn(msg string)
	Warnf(template string, msg string)
	Error(err error)
	Fatal(err error)
}
