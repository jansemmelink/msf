package main

import (
	"testing"

	"github.com/jansemmelink/msf/lib/micro"
)

func Test1(t *testing.T) {
	micro.Test(&greeterService{Name: ""}, greeterResponse{"Hello anonymous!"}, greeterAudit{Len: 9})
	micro.Test(&greeterService{Name: "Jan"}, greeterResponse{"Hello Jan!"}, greeterAudit{Len: 3})
	micro.Test(&greeterService{greeting: "Goodbye", Name: "Jan"}, greeterResponse{"Goodbye Jan!"}, greeterAudit{Len: 3})
}
