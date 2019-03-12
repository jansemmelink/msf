package log

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jansemmelink/msf/lib/log/level"
)

//IWriter ...
type IWriter interface {
	Write(h Header, text string)
	Level() level.Enum
	SetLevel(l level.Enum)
}

//NewFileWriter ...
func NewFileWriter(f *os.File) IWriter {
	return &fileWriter{
		level: level.Error,
		f:     f,
	}
}

type fileWriter struct {
	level level.Enum
	f     *os.File
}

func (fw fileWriter) Level() level.Enum {
	return fw.level
}

func (fw *fileWriter) SetLevel(l level.Enum) {
	fw.level = l
}

func (fw fileWriter) Write(hdr Header, msg string) {
	if hdr.Level > fw.level {
		return
	}

	//-------------------------------------------------------------------------------
	//todo: get rid of name param in New(name) and define package name automatically
	//and see if can create logger automatically for each package, or just determine
	//package name and use to control log levels and throttle each package
	//-------------------------------------------------------------------------------
	fn := path.Base(hdr.Function.Package)
	if len(hdr.Function.Type) > 0 {
		fn += "." + hdr.Function.Type
	}
	fn += "." + hdr.Function.Name
	fl := len(fn)
	if fl > 30 {
		fn = fn[fl-30:]
	}
	//trim all trailing newline characters
	for strings.HasSuffix(msg, "\n") {
		msg = msg[:len(msg)-len("\n")]
	}

	fw.f.Write([]byte(fmt.Sprintf("%s %016X %5.5s %30.30s(%5d): %s\n",
		hdr.Timestamp.Format("2006-01-02 15:04:05.000"),
		hdr.GoRoutine,
		hdr.Level,
		//hdr.Logger, //omit logger.Name, rather use package base name to control log levels
		fn,
		hdr.LineNr,
		msg)))
}
