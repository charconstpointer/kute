package kute

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type PipeLogger struct {
	prefix string
}

func (l *PipeLogger) Infof(format string, args ...interface{}) {
	log.Infof(fmt.Sprintf("[%s] %s", l.prefix, format), args)
}
func (l *PipeLogger) Errorf(format string, args ...interface{}) {
	log.Errorf(fmt.Sprintf("[%s] %s", l.prefix, format), args)
}
