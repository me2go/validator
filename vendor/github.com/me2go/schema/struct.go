package schema

import "reflect"

type MethodInfo struct {
	m   *reflect.Method
	in  []reflect.Type
	out []reflect.Type
}
type StructMeta struct {
	Name    string
	Fields  []*reflect.StructField
	Methods []*MethodInfo
}

func IsStruct(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr &&
		t.Elem().Kind() != reflect.Struct {
		return false
	}
	if t.Kind() != reflect.Ptr && t.Kind() != reflect.Struct {
		return false
	}
	return true
}
func From(i interface{}) *StructMeta {
	t := reflect.TypeOf(i)

	if !IsStruct(t) {
		return nil
	}
	fs := StructFields(t)
	ms := StructMethods(t)
	sm := &StructMeta{
		Fields:  fs,
		Methods: ms,
	}
	if t.Kind() == reflect.Ptr {
		sm.Name = t.Elem().Name()
	} else {
		sm.Name = t.Name()
	}
	return sm
}
func StructFields(t reflect.Type) []*reflect.StructField {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fs := make([]*reflect.StructField, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fs = append(fs, &f)
	}
	return fs
}
func StructMethods(t reflect.Type) []*MethodInfo {
	ms := make([]*MethodInfo, 0, t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		in := make([]reflect.Type, m.Func.Type().NumIn())
		for i := 0; i < m.Func.Type().NumIn(); i++ {
			in = append(in, m.Func.Type().In(i))
		}
		out := make([]reflect.Type, m.Func.Type().NumOut())
		for i := 0; i < m.Func.Type().NumOut(); i++ {
			out = append(out, m.Func.Type().Out(i))
		}
		mi := &MethodInfo{
			m:   &m,
			in:  in,
			out: out,
		}
		ms = append(ms, mi)
	}
	return ms
}
