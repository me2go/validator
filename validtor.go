package validator

import (
	"sync"

	"github.com/me2go/schema"
)

type ValidateFunc func(interface{}, string) error

func NewValidator() *validator {
	v := &validator{
		tags: make(map[string]ValidateFunc),
	}
	for tag, f := range builtin {
		v.RegisterTag(tag, f)
	}
	return v
}

type validator struct {
	sync.RWMutex
	tags map[string]ValidateFunc
}

func (vor *validator) RegisterTag(tag string, f ValidateFunc) {
	if tag == "" || f == nil {
		return
	}
	vor.Lock()
	vor.tags[tag] = f
	vor.Unlock()
}
func (vor *validator) UnregisterTag(tag string) {
	if tag == "" {
		return
	}
	vor.Lock()
	delete(vor.tags, tag)
	vor.Unlock()
}

func (vor *validator) Validate(v interface{}) error {
	if v == nil {
		return nil
	}
	ve := make(ValidatorError)
	sm := &schema.StructSchema{
		F: &visitor{
			valid: vor,
			Err:   ve,
		},
	}
	sm.Fields(v)
	if len(ve) != 0 {
		return ve
	}
	return nil
}
