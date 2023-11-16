package log

type ILogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

var logger ILogger

func Init() error {
	if logger == nil {
		l, err := NewZapLogger()
		if err != nil {
			return err
		}

		logger = l.Sugar()
	}

	return nil
}

func Logger() ILogger {
	return logger
}
