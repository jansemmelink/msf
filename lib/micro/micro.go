package micro

import (
	"fmt"

	"github.com/jansemmelink/msf/lib/audit"
	"github.com/jansemmelink/msf/lib/config"
)

//IMicro is a micro-service
type IMicro interface {
	//Validate the parsed request data
	Validate() error

	//Handle the request to product a response and audit record
	Handle() (response interface{}, audit interface{})
}

//Service ...
type Service struct {
}

var (
	rootDomain IDomain
)

func init() {
	rootDomain = newDomain("")

	//add some management operations
	mgt := rootDomain.Sub("config")
	mgt.AddName("describe", &config.Describe{})
}

//Root domain
func Root() IDomain {
	return rootDomain
}

//Domain gets/creates the named domain
func Domain(n string) IDomain {
	return rootDomain.Sub(n)
}

//Test ...
func Test(req IMicro, expectedResponse interface{}, expectedAudit audit.IRecord) {
	if err := req.Validate(); err != nil {
		panic(fmt.Sprintf("Validation failed: %v", err))
	}
	res, ar := req.Handle()
	if res != expectedResponse {
		panic(fmt.Sprintf("Wrong response: %+v != %+v", res, expectedResponse))
	}
	if ar != expectedAudit {
		panic(fmt.Sprintf("Wrong audit: %+v != %+v", ar, expectedAudit))
	}
}
