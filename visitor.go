package validator

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

const TAG_KEY = "validator"

type visitor struct {
	valid *validator
	Err   ValidatorError
}

func (o *visitor) Skip(f *reflect.StructField) bool {
	_, ok := f.Tag.Lookup(TAG_KEY)
	return !ok && !f.Anonymous
}

func (o *visitor) Visit(f *reflect.StructField, v *reflect.Value) {
	if !v.CanInterface() {
		return
	}
	log.Printf("%v", f.Name)
	if f.Type.Kind() == reflect.Struct {
		if err := o.valid.Validate(v.Interface()); err != nil {
			o.Err[f.Type.Name()] = err
		}
		return
	}
	cris := o.Criterias(f)
	for c, p := range cris {
		validFunc, ok := o.valid.tags[c]
		if !ok || validFunc == nil {
			continue
		}
		if err := validFunc(v.Interface(), p); err != nil {
			o.Err[fmt.Sprintf("%s:%s", f.Name, c)] = err
		}
	}
}
func (o *visitor) Criterias(f *reflect.StructField) map[string]string {
	tag, ok := f.Tag.Lookup(TAG_KEY)
	if !ok || tag == "" {
		return nil
	}
	ret := map[string]string{}
	conds := strings.Split(tag, ";")
	for _, c := range conds {
		if c == "" {
			continue
		}
		cp := strings.Split(c, ":")
		if len(cp) == 1 {
			ret[strings.Trim(cp[0], " ")] = ""
		} else if len(cp) == 2 {
			ret[strings.Trim(cp[0], " ")] = strings.Trim(cp[1], " ")
		} else {
			//not supported
			continue
		}
	}
	return ret
}
