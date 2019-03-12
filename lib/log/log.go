package log

import (
	"fmt"
	"os"

	"github.com/jansemmelink/msf/lib/log/level"
)

var (
	//default logger writes to stderr
	logger = NewFileWriter(os.Stderr)
)

//function that operate of default logger
func log(l level.Enum, f string, a ...interface{}) {
	h := Header{}
	h.Set(3, "logger", l)
	logger.Write(h, fmt.Sprintf(f, a...))
}

//Tracef ...
func Tracef(f string, a ...interface{}) {
	log(level.Trace, f, a...)
}

//Debugf ...
func Debugf(f string, a ...interface{}) {
	log(level.Debug, f, a...)
}

//Infof ...
func Infof(f string, a ...interface{}) {
	log(level.Info, f, a...)
}

//Warnf ...
func Warnf(f string, a ...interface{}) {
	log(level.Warn, f, a...)
}

//Errorf ...
func Errorf(f string, a ...interface{}) {
	log(level.Error, f, a...)
}

//Notef ...
func Notef(f string, a ...interface{}) {
	log(level.Note, f, a...)
}

//Fatalf ...
func Fatalf(f string, a ...interface{}) {
	log(level.Fatal, f, a...)
	os.Exit(1)
}

//SetLevel ...
func SetLevel(l level.Enum) {
	logger.SetLevel(l)
}

//DebugOn is shorthand to set level to debug
func DebugOn() {
	logger.SetLevel(level.Debug)
}

//VerboseOn is shorthand to set level to info
func VerboseOn() {
	logger.SetLevel(level.Info)
}
