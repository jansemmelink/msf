package level

import (
	"fmt"
	"strings"
)

//Enum stores a enumerator value
type Enum int

//Enumerator values
const (
	//None is place holder
	None Enum = iota
	//Fatal is most critical level
	Fatal
	//Notes are always printed, cause most processes run at Error or higher level
	Note
	//Error is something wrong that does not necessarily terminate the process
	Error
	//Warn might require attention
	Warn
	//Info is non-transactional info, e.g. creation of something on startup
	Info
	//Debug is for flow of code
	Debug
	//Trace is for underlying value details, e.g. hex dumps
	Trace
)

var mapEnum2Text = map[Enum]string{
	None:  "none",
	Fatal: "fatal",
	Note:  "note",
	Error: "error",
	Warn:  "warn",
	Info:  "info",
	Debug: "debug",
	Trace: "trace",
}

var mapText2Enum = map[string]Enum{}

func init() {
	for level, text := range mapEnum2Text {
		mapText2Enum[text] = level
	}
}

//String converts enum to text
func (e Enum) String() string {
	text, ok := mapEnum2Text[e]
	if !ok {
		return mapEnum2Text[None]
	}
	return text
}

//Parse converts text into enum
func Parse(text string) Enum {
	if e, ok := mapText2Enum[strings.ToLower(text)]; ok {
		return e
	}
	return None
}

//MarshalJSON ...
func (e Enum) MarshalJSON() ([]byte, error) {
	return []byte(e.String()), nil
}

//UnmarshalJSON ...
func (e *Enum) UnmarshalJSON(b []byte) error {
	if value, ok := mapText2Enum[strings.Trim(strings.ToLower(string(b)), "\"")]; ok {
		*e = value
		return nil
	}
	return fmt.Errorf("\"%s\" not a log level", string(b))
}

// todo:
// - config is gotten repeatedly, including marshal & unmarshal of JSON is repeated - not good.
// - module level is not applied
