package mq

import (
	"fmt"
	"sync"

	"github.com/jansemmelink/msf/lib/config"
	"github.com/jansemmelink/msf/lib/log"
	"github.com/jansemmelink/msf/lib/micro"
)

//Add a listener implementation that can be configured
func Add(name string, implementation IListener, doc string) {
	log.Debugf("Adding %s for %s ...", name, doc)
	implementationsMutex.Lock()
	defer implementationsMutex.Unlock()

	if _, ok := implementations[name]; ok {
		panic(fmt.Sprintf("Duplicate mq.IListener(%s)", name))
	}

	config.Register("mq."+name, implementation, "Configure this to use "+doc+" for message queue processing.")

	implementations[name] = implementation
	log.Debugf("Registered mq.IListener(%s)", name)
}

var (
	implementationsMutex = sync.Mutex{}
	implementations      = make(map[string]IListener)
	defaultListener      IListener
)

//IListener ...
type IListener interface {
	config.IConfigurable
	Listen(d micro.IDomain)
}

//Listener ...
type Listener struct{}

//Listen using configured listener
func Listen(d micro.IDomain) {
	names := ""
	if defaultListener == nil {
		log.Debugf("Looking for one of %d listeners in config", len(implementations))
		for name /*,listenerImplementation*/ := range implementations {
			log.Debugf("  Trying %s ...", name)
			names += "|" + name

			configuredListener, err := config.Get("mq." + name) //, listenerImplementation)
			if err != nil {
				log.Debugf("    mq.%s not available: %v", name, err)
				continue
			}

			err = configuredListener.Validate()
			if err != nil {
				log.Debugf("    mq.%s is not valid: %v", name, err)
				continue
			}

			log.Infof("Using mq.%s", name)
			defaultListener = configuredListener.(IListener)
			break
		}
	}

	if defaultListener == nil {
		if len(names) > 0 {
			names = names[1:]
		}
		log.Fatalf("No mq listener configured, expecting mq.%s", names)
	}
	defaultListener.Listen(d)
}
