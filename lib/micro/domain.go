package micro

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/jansemmelink/msf/lib/log"
)

//IDomain ...
type IDomain interface {
	Sub(n string) IDomain
	GetSub(n string) IDomain
	GetSubs() map[string]IDomain

	Add(m IMicro)
	AddName(n string, m IMicro)
	Get(n string) IMicro
	Opers() map[string]IMicro
}

type domain struct {
	name  string
	mutex sync.Mutex
	sub   map[string]IDomain
	oper  map[string]oper
}

type oper struct {
	req                IMicro
	responseStructType reflect.Type
	auditStructType    reflect.Type
}

//New default domain
func newDomain(n string) IDomain {
	return &domain{
		name: n,
		sub:  make(map[string]IDomain),
		oper: make(map[string]oper),
	}
}

func (d *domain) Sub(n string) IDomain {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if existing, ok := d.sub[n]; ok {
		return existing
	}
	new := newDomain(n)
	d.sub[n] = new
	return new
}

func (d *domain) GetSub(n string) IDomain {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if existing, ok := d.sub[n]; ok {
		return existing
	}
	return nil
}

func (d *domain) GetSubs() map[string]IDomain {
	return d.sub
}

func (d *domain) Add(m IMicro) {
	d.AddName("", m)
}

func (d *domain) AddName(n string, m IMicro) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	t := reflect.TypeOf(m)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Errorf("micro.Add(%T) must use &struct{}", m))
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("micro.Add(%T) is not &struct{}", m))
	}

	if len(n) == 0 {
		n = t.Name()
	}

	if _, ok := d.oper[n]; ok {
		panic(fmt.Sprintf("Duplicate name=\"%s\" in micro.Add(%T)", n, m))
	}

	//registered pointer to struct
	operCopy := m
	if err := operCopy.Validate(); err != nil {
		panic(fmt.Sprintf("micro.Add(%s,%T): invalid oper: %v", n, m, err))
	}
	operResponseStruct, operAuditStruct := operCopy.Handle()
	newOper := oper{
		req:                m,
		responseStructType: reflect.TypeOf(operResponseStruct),
		auditStructType:    reflect.TypeOf(operAuditStruct),
	}
	d.oper[n] = newOper
	log.Infof("Registered %s/%s: %T=%+v -> %v + %v\n", d.name, n, newOper.req, newOper.req, newOper.responseStructType, newOper.auditStructType)
}

func (d *domain) Get(n string) IMicro {
	if existing, ok := d.oper[n]; ok {
		return existing.req
	}
	return nil
}

func (d *domain) Opers() map[string]IMicro {
	log.Debugf("Listing %d opers", len(d.oper))
	opers := make(map[string]IMicro)
	for name, oper := range d.oper {
		opers[name] = oper.req
	}
	return opers
}
