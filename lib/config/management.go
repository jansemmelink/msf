package config

import (
	"reflect"

	"github.com/jansemmelink/msf/lib/log"
)

//Describe implements IMicro to document config
type Describe struct {
	Name string `json:"name" doc:"Name of configuration to document"`
}

//Validate ...
func (oper Describe) Validate() error {
	return nil
}

//Schema of config
type Schema struct {
	Name  string `json:"name"`
	Doc   string `json:"doc"`
	Items []item `json:"items"`
}

type item struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Doc  string `json:"doc"`
}

//Handle ...
func (oper Describe) Handle() (res interface{}, audit interface{}) {
	log.Debugf("Documenting %+v", oper)
	s := Schema{Name: oper.Name}
	if cs != nil {
		if c, ok := cs.all[oper.Name]; ok {
			//schema for named config
			s.Doc = c.doc
			s.Items = make([]item, 0)

			//list struct fields
			t := reflect.TypeOf(c.data)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			for fti := 0; fti < t.NumField(); fti++ {
				ft := t.Field(fti)
				if ft.Anonymous {
					continue
				}
				if ft.Name[0] < 'A' || ft.Name[0] > 'Z' {
					continue
				}
				i := item{Name: ft.Name, Type: ft.Type.String(), Doc: ft.Tag.Get("doc")}
				s.Items = append(s.Items, i)
			}
			return s, nil
		}

		//list all config
		s.Name = ""
		s.Doc = "The following items can be configured."
		s.Items = make([]item, 0)
		for name, c := range cs.all {
			s.Items = append(s.Items, item{Name: name, Doc: c.doc, Type: reflect.TypeOf(c.data).Elem().String()})
		}
	}
	return s, nil
}
