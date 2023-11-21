package log

type ILoggerManager interface {
	Logger() ILogger
	Sync() error
}

type ILogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

var loggerManager ILoggerManager

type InitOptions struct {
	StdoutFile string
	StderrFile string
}

func Init(opts *InitOptions) (ILoggerManager, error) {
	if loggerManager == nil {
		l, err := NewZapLogger(opts)
		if err != nil {
			return nil, err
		}

		lm := NewZapLoggerManager(l)

		loggerManager = lm
	}

	return loggerManager, nil
}

func Logger() ILogger {
	return loggerManager.Logger()
}
