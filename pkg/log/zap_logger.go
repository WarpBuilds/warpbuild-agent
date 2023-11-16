package log

import (
	"go.uber.org/zap"
)

func NewZapLogger() (*zap.Logger, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return l, nil
}
