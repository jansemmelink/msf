package nats

import (
	"fmt"

	"github.com/jansemmelink/msf/lib/log"
	"github.com/jansemmelink/msf/lib/micro"
	"github.com/jansemmelink/msf/lib/mq"
)

func init() {
	mq.Add("nats", &popper{}, "NATS")
}

type popper struct {
	mq.Listener
	Subject string
}

func (p popper) Validate() error {
	if p.Subject == "" {
		return fmt.Errorf("missing subject=... ")
	}
	log.Debugf("popper validated: %+v", p)
	return fmt.Errorf("NYI") //nil
}

//Listen ...
func (p popper) Listen(d micro.IDomain) {
	log.Debugf("NATS Listening to %s ...", p.Subject)
	for {
		log.Fatalf("NYI")
	}
}
