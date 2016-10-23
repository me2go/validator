package validator

import (
	"reflect"
	"strings"
	"sync"
)

const FIELD_TAG_KEY = "validator"

func NewValidator() *Validator {
	v := &Validator{
		tags: make([]*Tag, 0),
	}
	v.RegisterTag(GREATEROREQUAL)
	v.RegisterTag(GREATER)
	v.RegisterTag(LESSOREQUAL)
	v.RegisterTag(LESS)
	v.RegisterTag(EQUAL)
	v.RegisterTag(NONEQUAL)
	v.RegisterTag(REGEXP)
	v.RegisterTag(NONZERO)
	return v
}

type Validator struct {
	sync.RWMutex
	tags []*Tag
}

func (v *Validator) Validate(i interface{}) error {
	if i == nil {
		return nil
	}
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Struct {
		return INVALID_TYPE
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	vtype := value.Type()
	var err error = nil
	for i, nums := 0, value.NumField(); i < nums; i++ {
		vf := vtype.Field(i)
		vv := value.Field(i)
		if vv.Kind() != reflect.Struct {
			err = v.validateField(&vf, &vv)
		} else if vv.CanInterface() {
			err = v.Validate(vv.Interface())
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) validateField(f *reflect.StructField, fv *reflect.Value) error {
	tagString := f.Tag.Get(FIELD_TAG_KEY)
	if tagString == "" {
		return nil
	}
	validateElems := strings.Split(tagString, ";")
	for _, e := range validateElems {
		v.RLock()
		for _, tag := range v.tags {
			if tag.Handler == nil || !tag.Handler.Match(e) {
				continue
			}
			id, tagKind, tmp := tag.Handler.Parse(e)
			if id != tag.ID || tagKind != tag.Kind || tmp == "" {
				continue
			}
			if err := tag.Handler.Check(*fv, tmp); err != nil {
				return err
			}
		}
		v.RUnlock()
	}
	return nil
}

func (v *Validator) HasTag(t *Tag) bool {
	v.RLock()
	defer v.RUnlock()
	_, tt := v.findTagById(t.ID)
	return tt != nil
}

func (v *Validator) RegisterTag(t *Tag) error {
	if t == nil {
		return nil
	}
	v.Lock()
	defer v.Unlock()

	i, _ := v.findTagById(t.ID)
	if i == -1 {
		v.tags = append(v.tags, t)
	} else {
		v.tags[i] = t
	}

	return nil
}
func (v *Validator) UnregisterTag(t *Tag) error {
	if t == nil {
		return nil
	}
	v.Lock()
	defer v.Unlock()

	i, _ := v.findTagById(t.ID)
	if i < 0 {
		return nil
	}
	v.tags = append(v.tags[0:i], v.tags[i+1:]...)
	return nil
}

func (v *Validator) findTagById(id string) (int, *Tag) {
	for i, tag := range v.tags {
		if tag.ID == id {
			return i, tag
		}
	}
	return -1, nil
}
