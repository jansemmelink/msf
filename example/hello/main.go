package main

import (
	"github.com/jansemmelink/msf/lib/micro"
	"github.com/jansemmelink/msf/lib/mq"
	_ "github.com/jansemmelink/msf/lib/mq/nats"
	_ "github.com/jansemmelink/msf/lib/mq/rabbit"
	_ "github.com/jansemmelink/msf/lib/mq/redis"
	_ "github.com/jansemmelink/msf/lib/mq/rest"
)

func main() {
	//log.DebugOn()
	//log.Debugf("Starting...")
	mq.Listen(micro.Root())
}
