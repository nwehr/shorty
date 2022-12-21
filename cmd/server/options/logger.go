package options

import (
	"fmt"
	"os"
	"time"
)

type severity int

const (
	DEBUG severity = iota
	WARNNING
	ERROR
	FATAL
)

type logger struct {
	severity severity
}

func (l logger) Debugf(format string, args ...any) {
	if l.severity < WARNNING {
		fmt.Printf("%s [DEBUG] %s", time.Now().Format(time.Stamp), fmt.Sprintf(format, args...))
	}
}

func (l logger) Warningf(format string, args ...any) {
	if l.severity < ERROR {
		fmt.Printf("%s [WARN] %s", time.Now().Format(time.Stamp), fmt.Sprintf(format, args...))
	}
}

func (l logger) Errorf(format string, args ...any) {
	if l.severity < FATAL {
		fmt.Printf("%s [ERROR] %s", time.Now().Format(time.Stamp), fmt.Sprintf(format, args...))
	}
}

func (l logger) Fatalf(format string, args ...any) {
	fmt.Printf("%s [FATAL] %s", time.Now().Format(time.Stamp), fmt.Sprintf(format, args...))
	os.Exit(1)
}
