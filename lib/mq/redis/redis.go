package redis

import (
	"fmt"
	"sync"
	"time"

	"github.com/jansemmelink/msf/lib/log"
	"github.com/jansemmelink/msf/lib/micro"
	"github.com/jansemmelink/msf/lib/mq"
	"github.com/pkg/errors"
)

func init() {
	mq.Add("redis", &popper{}, "REDIS")
}

type popper struct {
	mq.Listener
	Server        string `json:"server" doc:"REDIS Server address or hostname. Defaults to localhost."`
	Port          int    `json:"port" doc:"REDIS Server TCP port number. Defaults to 6379."`
	NrConn        int    `json:"nrConn" doc:"Nr of connections to make to the server. Defaults to 1."`
	QName         string `json:"qname" doc:"REDIS queue name to consume."`
	Limit         int    `json:"limit" doc:"Terminate after popping this nr of messages. Defaults to -1 = unlimited."`
	MaxConcurrent int    `json:"maxConcurrent" doc:"Max concurrent transactions. Defaults to 100."`

	//runtime data
	remain int
}

func (p *popper) Validate() error {
	if p.Server == "" {
		p.Server = "localhost"
	}
	if p.Port <= 0 {
		p.Port = 6379
	}
	if p.NrConn < 1 {
		p.NrConn = 1
	}
	if p.QName == "" {
		return fmt.Errorf("missing qname=... ")
	}
	if p.Limit < 1 {
		p.Limit = -1
	}
	if p.MaxConcurrent < 1 {
		p.MaxConcurrent = 100
	}
	log.Debugf("popper validated: %+v", p)
	return nil
}

//Listen ...
func (p popper) Listen(d micro.IDomain) {
	log.Debugf("REDIS Listening to %s ...", p.QName)
	//	p := Popper{pool: nil, stopped: false, count: 0, limit: limit}

	pool, err := NewRedis("tcp", fmt.Sprintf("%s:%d", p.Server, p.Port), p.NrConn)
	if err != nil {
		panic(errors.Wrapf(err, "Failed to create Redis pool for pop"))
	}

	//atomic counter of the current nr of decoders running
	//every decoder in the end result in another context being created
	//(except in the case of response popper) so the nr of decoders
	//are deducted from the capacity of the server
	//note: without this, under load, messages are popped very fast and
	//request decoders are spun up faster than they can complete and
	//server capacity does not decrease, causing the popper to take on too much
	//and use thousands of goroutines concurrently, eventually killing the server
	//var nrDecoders int32 //atomic.Value
	//nrDecoders = 0

	log.Debugf("Creating context channel...")
	ctxChannel := make(chan ctx, p.MaxConcurrent)
	for i := 0; i < p.MaxConcurrent; i++ {
		ctxChannel <- ctx{id: i}
		log.Debugf("Added context[%d] ...", i)
	}

	//start a go routine for each connection to pop messages
	wg := sync.WaitGroup{}
	p.remain = p.Limit
	for i := 0; i < p.NrConn; i++ {
		wg.Add(1)
		go func(conn int) {
			p.pop(pool, conn, ctxChannel)
			wg.Done()
		}(i)
	}
	//wait for all to terminate
	wg.Wait()
	log.Infof("Popper terminated")
}

func (p popper) pop(pool Redis, conn int, ctxPool chan ctx) {
	for {
		if p.Limit > 0 && p.remain <= 0 {
			log.Debugf("Popper(%s) conn[%d] terminating after %d pops.", p.QName, conn, p.Limit)
			break
		}

		//wait for and get next available context
		//or terminate when the channel is closed
		timeout := time.Duration(10) * time.Second
		select {
		case ctx := <-ctxPool:
			log.Debugf("Conn[%d]: Got context: %+v", conn, ctx)
			ctx.Pop(pool, p.QName, ctxPool)
		case <-time.After(timeout):
			log.Errorf("Popper(%s) conn[%d]: No available contexts...", p.QName, conn)
		}
	} //until stop
	log.Debugf("Popper(%s) conn[%d]: Stopped", p.QName, conn)
}

type ctx struct {
	id int
}

func (ctx ctx) Pop(pool Redis, qname string, ctxPool chan ctx) int {
	defer func() {
		//put context back in the pool
		ctxPool <- ctx
	}()

	data, err := pool.BRPOP(qname, 1)
	if err != nil {
		if err.Error() == "EOF" {
			panic("%s: Redis Connection dropped: " + err.Error()) //TODO: Handle with re-connect
		}
		log.Errorf("%s: BRPOP Error: %v", qname, err)
		return 0
	}

	if len(data) <= 0 {
		//blocking popped timed out - queue is idle
		return 0
	}

	//popped a message
	log.Errorf("%+v: NOT YET PROCESSING: %v", ctx, data)
	return 1

	//decode and handle in a separate go-routine
	//so that this routine can immediately pop again
	// go func() {
	// 	atomic.AddInt32(&nrDecoders, 1)
	// 	defer atomic.AddInt32(&nrDecoders, -1)

	// 	messageHeader := MessageHeader{}
	// 	messageHeader.json = data
	// 	if err := messageHeader.FromJSON(data); err != nil {
	// 		log.Error.Printf("%s: Discard: Invalid JSON: %v: %v", qname, err, data)
	// 		return
	// 	}

	// 	if messageHeader.Header == nil {
	// 		log.Error.Printf("%s: Discard: Missing header in %v", qname, data)
	// 		return
	// 	}

	// 	if len(messageHeader.Header.IntGUID) == 0 {
	// 		log.Error.Printf("%s: Discard: Missing header.int_guid in %v", qname, data)
	// 		return
	// 	}

	// 	if len(messageHeader.Header.Timestamp) == 0 {
	// 		log.Error.Printf("%s: Discard: Missing header.timestamp in %v", qname, data)
	// 		return
	// 	}

	// 	if len(messageHeader.Header.Timestamp) == len("2017-06-07 11:37:58") {
	// 		messageHeader.Header.Timestamp = messageHeader.Header.Timestamp + ".000"
	// 	}

	// 	var err error
	// 	if messageHeader.Header._ts, err = time.ParseInLocation("2006-01-02 15:04:05.000", messageHeader.Header.Timestamp, time.Local); err != nil {
	// 		log.Error.Printf("%s: Discard: Invalid header.timestamp %v: %v", qname, err, data)
	// 		return
	// 	}

	// 	if messageHeader.Header.TTL < 0 {
	// 		log.Error.Printf("%s: Discard: Negative Ttl: %v", qname, data)
	// 		return
	// 	}

	// 	//zero TTL is used for messages that should not expire
	// 	//they can sit long in a queue, but now that we popped the messageHeader,
	// 	//we give it a sensible TTL starting from now + maxDuration in this process
	// 	//so we use TTL = timeNow + TTL0 - timestamp
	// 	if messageHeader.Header.TTL == 0 {
	// 		if c.Pop.TTL0 == 0 {
	// 			log.Error.Printf("%s: Discard: Zero Ttl (TTL0 not configured): %v", qname, data)
	// 			return
	// 		}
	// 		messageHeader.Header.TTL = int(time.Now().Add(time.Duration(c.Pop.TTL0)*time.Millisecond).Sub(messageHeader.Header._ts) / time.Millisecond)
	// 	}

	// 	messageHeader.Header._ttl = time.Duration(messageHeader.Header.TTL) * time.Millisecond
	// 	messageHeader.Header._expTime = messageHeader.Header._ts.Add(messageHeader.Header._ttl)

	// 	//set default empty provider|consumer when not specified to require fewer checks later
	// 	if messageHeader.Header.Provider == nil {
	// 		messageHeader.Header.Provider = &Provider{}
	// 	}
	// 	if messageHeader.Header.Consumer == nil {
	// 		messageHeader.Header.Consumer = &Consumer{}
	// 	}

	// 	//requests (messages without header.result) must have a valid provider name
	// 	//written as "/domain/oper"
	// 	if messageHeader.Header.Result == nil {
	// 		do := strings.Split(messageHeader.Header.Provider.Name, "/")
	// 		if len(do) != 3 {
	// 			log.Error.Printf("%s: Discard: Provider.Name=\"%v\" not /domain/oper", qname, messageHeader.Header.Provider.Name)
	// 			return
	// 		}
	// 		messageHeader.Header.Provider._domainName = do[1]
	// 		messageHeader.Header.Provider._operName = do[2]
	// 	} //if request

	// 	// push into messageHeader channel then return
	// 	log.Trace.Printf("%s: popped: %v", qname, data)
	// 	log.Trace.Printf("          header=%v", *messageHeader.Header)
	// 	log.Trace.Printf("          provider=%v", *messageHeader.Header.Provider)
	// 	if messageHeader.Header.Consumer != nil {
	// 		log.Trace.Printf("          consumer=%v", *messageHeader.Header.Consumer)
	// 	}
	// 	if messageHeader.Header.Result != nil {
	// 		log.Trace.Printf("          result=%v", *messageHeader.Header.Result)
	// 	}

	// 	//send to the msgChannel
	// 	//(blocking if necessary, i.e. do not use select{} around this)
	// 	msgChannel <- messageHeader
	// 	return

}
