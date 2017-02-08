package validator

import (
	"log"
	"testing"
)

type SubStruct struct {
	S1 int    `validator:"max:10"`
	S2 string `validator:"required"`
}
type TestStruct struct {
	SubStruct
	F1 string `validator:"required;max: 10;min:2"`
	F2 int    `validator:"in:[2,4,6,7]"`
	F3 string `validator:"required"`
}

func TestValidateFail(t *testing.T) {
	s := &TestStruct{
		F1: "",
		F2: 3,
		F3: "fdas",
	}
	err := Validate(s)
	if err != nil {
		log.Printf(err.Error())
	} else {
		t.Error()
	}
}
func TestValidateSuccess(t *testing.T) {
	s := &TestStruct{
		SubStruct: SubStruct{
			S1: 9,
			S2: "hello",
		},
		F1: "大家好，我是谁谁方便",
		F2: 2,
		F3: "fdas",
	}
	err := Validate(s)
	if err == nil {
		log.Printf("ok")
	} else {
		t.Errorf(err.Error())
	}
}
func BenchmarkValidateSuccess(b *testing.B) {
	s := &TestStruct{
		F1: "dasf",
		F2: 2,
		F3: "fdas",
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = Validate(s)
	}
	b.StopTimer()
}
func BenchmarkValidateFail(b *testing.B) {
	s := &TestStruct{
		F1: "",
		F2: 3,
		F3: "fdas",
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = Validate(s)
	}
	b.StopTimer()
}
