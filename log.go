package tgbotapi

import (
	"errors"
	"io"
	stdlog "log"
	"os"
)

// BotLogger is an interface that represents the required methods to log data.
//
// Instead of requiring the standard logger, we can just specify the methods we
// use and allow users to pass anything that implements these. We use a subset
// of standard logging library APIs. A recommended library is k8s.io/klog.
type BotLogger interface {
	Infoln(args ...interface{})
	Infof(format string, args ...interface{})
	Errorln(args ...interface{})
	Errorf(format string, args ...interface{})
}

var log = newStdLogAdapter(os.Stderr, "", stdlog.LstdFlags)

// SetLogger specifies the logger that the package should use.
func SetLogger(logger BotLogger) error {
	if logger == nil {
		return errors.New("logger is nil")
	}
	log = logger
	return nil
}

// Implement standard log adapter
func newStdLogAdapter(out io.Writer, prefix string, flag int) BotLogger {
	return &stdlogAdapter{*stdlog.New(out, prefix, flag)}
}

type stdlogAdapter struct {
	stdlog.Logger
}

func (slw *stdlogAdapter) Infoln(args ...interface{}) {
	slw.Println(args...)
}

func (slw *stdlogAdapter) Infof(format string, args ...interface{}) {
	slw.Printf(format, args...)
}
func (slw *stdlogAdapter) Errorln(args ...interface{}) {
	slw.Println(args...)
}

func (slw *stdlogAdapter) Errorf(format string, args ...interface{}) {
	slw.Printf(format, args...)
}
