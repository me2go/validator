/*
 scheam is a common library, which is designed to reflect a struct type easily
*/
package schema

import "reflect"

type FieldVisitor interface {
	Skip(*reflect.StructField) bool
	Visit(*reflect.StructField, *reflect.Value)
}
type MethodVisitor interface {
	Skip(m *reflect.Method, in []reflect.Type, out []reflect.Type) bool
	Visit(m *reflect.Method, in []reflect.Type, out []reflect.Type)
}
type StructSchema struct {
	F FieldVisitor
	M MethodVisitor

	C Cache
}

func (ss *StructSchema) meta(i interface{}) *StructMeta {
	if ss.C != nil {
		key := reflect.TypeOf(i).String()
		if b := ss.C.Get(key); b != nil {
			if v, ok := b.(*StructMeta); ok {
				return v
			}
		}
	}
	return From(i)
}

/*
 return the string representation of the type of i
*/
func (ss *StructSchema) Type(i interface{}) string {
	sm := ss.meta(i)
	if sm == nil {
		return ""
	}
	return sm.Name
}

/*
 Fields will visit all the struct fields
*/
func (ss *StructSchema) Fields(i interface{}) {
	if ss.F == nil {
		return
	}
	sm := ss.meta(i)
	if sm == nil {
		return
	}
	value := reflect.ValueOf(i)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	for _, f := range sm.Fields {
		if ss.F.Skip(f) {
			continue
		}
		fv := value.FieldByIndex(f.Index)
		ss.F.Visit(f, &fv)
	}
}

/*
 Methods will visit all the struct methods
*/
func (ss *StructSchema) Methods(i interface{}) {
	if ss.M == nil {
		return
	}
	sm := ss.meta(i)
	if sm == nil {
		return
	}
	for _, mi := range sm.Methods {
		if ss.M.Skip(mi.m, mi.in, mi.out) {
			continue
		}
		ss.M.Visit(mi.m, mi.in, mi.out)
	}
}
