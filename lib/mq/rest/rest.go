package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/jansemmelink/msf/lib/log"
	"github.com/jansemmelink/msf/lib/micro"
	"github.com/jansemmelink/msf/lib/mq"
	"github.com/pkg/errors"
)

func init() {
	mq.Add("rest", &rest{}, "HTTP REST")
}

type rest struct {
	mq.Listener
	Addr string `json:"addr" doc:"HTTP Server Address (defaults to localhost)"`
	Port int    `json:"port" doc:"TCP Port number to listen on (defaults to 8000)"`
}

func (p *rest) Validate() error {
	if p.Addr == "" {
		p.Addr = "localhost"
	}
	if p.Port <= 0 {
		p.Port = 8000
	}
	log.Debugf("rest validated: %+v", p)
	return nil
}

//Listen ...
func (p rest) Listen(d micro.IDomain) {
	addr := fmt.Sprintf("%s:%d", p.Addr, p.Port)
	log.Debugf("REDIS Listening to %s ...", addr)
	if err := http.ListenAndServe(addr, router{d: d}); err != nil {
		panic(errors.Wrapf(err, "Failed to serve HTTP REST"))
	}
}

type router struct {
	d micro.IDomain
}

func (r router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//get "/<domain>/<oper>" from the URL
	path := strings.Split(req.URL.Path, "/")
	if len(path) < 3 {
		http.Error(res, "invalid domain", http.StatusBadRequest)
		return
	}
	domainName := path[1]
	operName := path[2]
	log.Debugf("domain=%s oper=%s", domainName, operName)

	domain := r.d.GetSub(domainName)
	if domain == nil {
		subNames := ""
		for name := range r.d.GetSubs() {
			subNames += "|" + name
		}
		if len(subNames) > 0 {
			subNames = subNames[1:]
		}
		http.Error(res, "unknown domain, expecting "+subNames, http.StatusBadRequest)
		return
	}

	oper := domain.Get(operName)
	if oper == nil {
		operNames := ""
		log.Debugf("names...")
		for name := range domain.Opers() {
			log.Debugf("name=%s", name)
			operNames += "|" + name
		}
		if len(operNames) > 0 {
			operNames = operNames[1:]
		}
		http.Error(res, "unknown oper, expecting "+operNames, http.StatusBadRequest)
		return
	}

	//allocate a new operation request structure
	operStructType := reflect.TypeOf(oper)
	if operStructType.Kind() == reflect.Ptr {
		operStructType = operStructType.Elem()
	}
	log.Debugf("t=%v", operStructType)

	//allocate a new copy of the operation (request) struct
	operRequest := reflect.New(operStructType).Interface()
	//copy operation values from registered oper
	operRequest = oper
	log.Debugf("Allocated %T=%+v", operRequest, operRequest)

	//parse operation request from body
	if req.Body != nil {
		if err := json.NewDecoder(req.Body).Decode(&operRequest); err != nil && err != io.EOF {
			http.Error(res, "invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	log.Debugf("Request Body: %+v", operRequest)

	//todo: parse URL params
	operValue := reflect.ValueOf(operRequest).Elem()
	log.Debugf("v=%v", operValue)
	for paramName, paramValues := range req.URL.Query() {
		found := false
		for fti := 0; fti < operStructType.NumField(); fti++ {
			ft := operStructType.Field(fti)
			log.Debugf("Check p=%s against f[%d]=%s", paramName, fti, ft.Name)
			if paramName == ft.Name || paramName == ft.Tag.Get("json") {
				log.Debugf("Got match on ft[%d]=%s", fti, ft.Name)
				fieldValue := operValue.Field(fti)
				if !fieldValue.CanSet() {
					http.Error(res, "URL param not allowed: "+paramName, http.StatusBadRequest)
					return
				}
				fieldValue.Set(reflect.ValueOf(paramValues[0]))
				found = true
				break
			}
		}
		if !found {
			http.Error(res, "unknown URL param "+paramName, http.StatusBadRequest)
			return
		}
	}
	log.Debugf("Request Params: %+v", operRequest)

	//execute
	operResponse, operAudit := operRequest.(micro.IMicro).Handle()
	log.Debugf("Res: %+v", operResponse)
	log.Debugf("Audit: %+v", operAudit)

	jsonRes, _ := json.Marshal(operResponse)
	res.Write(jsonRes)
}
