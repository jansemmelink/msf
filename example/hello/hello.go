package main

import (
	"github.com/jansemmelink/msf/lib/micro"
)

func init() {
	//micro-service must be registered with a domain name ("greet" in this case)
	//and an operation name, that defaults to the struct type name:
	//this registers domain="greet" oper="greeterService"
	//this registration uses a pointer to the struct, which allows Validate() to
	//modify the contents before Handle() is called. If validate is a constant
	//method, the operation may be registered as a struct.
	micro.Domain("greet").Add(&greeterService{})
	//this registers domain="greet" oper="goodbye"
	micro.Domain("greet").AddName("goodbye", &greeterService{greeting: "Goodbye"})
}

//each micro-service operation is defined as a struct that implements IMicro
type greeterService struct {
	micro.Service

	//private struct members can modify the service definition to use the same
	//struct for different outcomes, e.g. here we have different greetings
	//and the user of the service cannot change this:
	greeting string

	//Public members are request data that are specified by the user
	//when the operation is invoked.

	//todo: To define default values, specify them also during registration.
	Name string
}

type greeterResponse struct {
	Message string
}

type greeterAudit struct {
	Len int
}

func (h *greeterService) Validate() error {
	if h.greeting == "" {
		h.greeting = "Hello"
	}
	if len(h.Name) <= 0 {
		h.Name = "anonymous"
	}
	return nil
}

func (h greeterService) Handle() (res interface{}, a interface{}) {
	return greeterResponse{h.greeting + " " + h.Name + "!"}, greeterAudit{Len: len(h.Name)}
}
