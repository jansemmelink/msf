package rabbit

import (
	"fmt"

	"github.com/jansemmelink/msf/lib/log"
	"github.com/jansemmelink/msf/lib/micro"
	"github.com/jansemmelink/msf/lib/mq"
)

func init() {
	mq.Add("rabbit", &consumer{}, "RabbitMQ")
}

type consumer struct {
	mq.Listener
	Subject string
}

func (p consumer) Validate() error {
	if p.Subject == "" {
		return fmt.Errorf("missing subject=... ")
	}
	log.Debugf("popper validated: %+v", p)
	return fmt.Errorf("NYI")
}

//Listen ...
func (p consumer) Listen(d micro.IDomain) {
	log.Debugf("Rebbit consuming %s ...", p.Subject)
	for {
		log.Fatalf("NYI")
	}
}
