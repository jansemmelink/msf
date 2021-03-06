package log

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jansemmelink/msf/lib/log/level"
)

//function to get goroutine ID
//should not be used for anything except logging
//because golang explicitly do not want to identify go-routines,
//but we use it so that concurrent log entries in a file can be identified
//to which routine wrote them.
//Code was taken from https://blog.sgmansfield.com/2015/12/goroutine-ids/
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

//Header is details about the log event that is not part of the actual log message
type Header struct {
	Timestamp time.Time    `json:"timestamp"`
	Logger    string       `json:"logger"`
	GoRoutine uint64       `json:"goroutine"`
	Function  functionInfo `json:"function"`
	FileName  string       `json:"filename"`
	LineNr    int          `json:"linenr"`
	Level     level.Enum   `json:"level"`
}

//Set gathers info about the caller
// depth is the depth of call stack inside log module, so we can skip
// over those to get the log user
func (header *Header) Set(depth int, loggerName string, level level.Enum) error {
	//define timestamp and validate log level
	header.Timestamp = time.Now()
	header.Logger = loggerName
	header.Level = level

	//get file name and line number from runtime caller info
	//skip 3 entries for:
	//	  0: (this func)defineHeader()
	//    1: log()
	//    2: Debug()/Error()/...
	//    3: (user function that called Debug()/Error()/...)
	header.FileName = "<unknown>"
	header.LineNr = -1

	header.GoRoutine = getGID()

	for skip := depth; ; skip++ {
		_, fn, ln, gotit := runtime.Caller(skip)
		if !gotit {
			fmt.Println("DID NOT GET INFO")
			break
		}

		//stack sometimes has fn=<autogenerated>
		//skip to next level
		if fn[0:1] == "<" {
			//continue
		}

		//got info
		header.FileName = fn
		header.LineNr = ln
		break
	} /*for*/

	//get function from runtime caller stack
	pc := make([]uintptr, 10)
	n := runtime.Callers(0, pc)
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return fmt.Errorf("Unable to get call stack for log header")
	}

	var frame runtime.Frame
	if n >= 1+depth {
		pc = pc[1+depth : 3+depth] // pass only valid pcs to runtime.CallersFrames
		frames := runtime.CallersFrames(pc)
		frame, _ = frames.Next()

		//function starts with package name, then function name
		header.Function.Set(frame.Function)

		// Loop to get frames.
		// A fixed number of pcs can expand to an indefinite number of Frames.
		/*
			for {
				frame, more := frames.Next()
				// To keep this example's output stable
				// even if there are changes in the testing package,
				// stop unwinding when we leave package runtime.
				if strings.Contains(frame.File, "runtime/") {
					continue
				}
				if !more {
					break
				}
				fmt.Printf("- more:%v | %s (%d)\n", more, frame.Function, frame.Line)
			}
		*/
	}
	return nil
} //Header.Set()

type functionInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type,omitempty"`
	Package string `json:"package"`
}

//Set the function info from Caller.Function
func (info *functionInfo) Set(fn string) {
	//default when cannot interpret, put all in Name field
	info.Name = fn
	if fn == "" {
		return
	}

	//parse do notation from the back, expect at least one dot:
	//  "main.main"                                       -> func   "main"   in package "main"
	//	"github.com/jansemmelink/etl/lib-etl/logger.New"  -> func   "New"    in package "github.com/jansemmelink/etl/lib-etl"
	//	"bb.org/vs/lib-etl/logger.*static.newLog"         -> method "newLog" in type "static" in package "bb.org/vs/lib-etl/logger"
	// start with the last dot-notation element
	lastDot := strings.LastIndex(fn, ".")
	if lastDot < 0 {
		return
	}

	//got last '.', everything after the dot is the function/method name
	info.Name = fn[lastDot+1:]
	remain := fn[:lastDot]

	//before last dot we expect "<package>.<type>" or just "<package>"
	//package may have both '.' and '/'
	//if '.' is after last '/', it separates type from package and type is present
	lastDot = strings.LastIndex(remain, ".")
	lastSlash := strings.LastIndex(remain, "/")
	if lastDot > lastSlash {
		//has type
		info.Type = remain[lastDot+1:]
		info.Package = remain[0:lastDot]
	} else {
		//no dot, no type, just package
		info.Package = remain
	}
	return
}
