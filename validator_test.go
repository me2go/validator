package validator

import (
	"log"
	"testing"
)

type MySubStruct struct {
	A int `validator:"max:10"`
	B int `validator:"min:1"`
}

type MyStruct struct {
	A int    `validator:"max:10;min:10"`
	B string `validator:"required"`
	MySubStruct
}

func TestValidate(t *testing.T) {
	i := &MyStruct{
		A: 9,
		B: "",
		MySubStruct: MySubStruct{
			A: 11,
			B: 0,
		},
	}
	v := NewValidator()
	err := v.Validate(i)
	log.Printf("%v", err)
	if err == nil {
		t.Fail()
	}
}
