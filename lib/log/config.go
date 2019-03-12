package log

import (
	"fmt"

	"github.com/jansemmelink/msf/lib/log/level"
)

//Config for log levels
type Config struct {
	Global  level.Enum            `json:"global" doc:"This is the default level for packages that are not configured."`
	Package map[string]level.Enum `json:"package" doc:"This level only applies to the specified package."`
}

//Validate ...
func (c *Config) Validate() error {
	if c.Global <= level.None || c.Global >= level.Trace {
		c.Global = level.Error
	}
	for pkgName, pkgLevel := range c.Package {
		if len(pkgName) < 1 {
			return fmt.Errorf("log.package configured without a name")
		}
		if pkgLevel <= level.None || pkgLevel >= level.Trace {
			c.Package[pkgName] = level.Error
		}
	}
	return nil
}
