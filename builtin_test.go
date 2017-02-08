package validator

import (
	"log"
	"testing"
)

func TestRequired(t *testing.T) {
	err := required("", "")
	if err == nil {
		t.Error()
	} else {
		log.Printf(err.Error())
	}
	err = required("fdas", "")
	if err != nil {
		t.Error()
	} else {
		log.Printf("string ok")
	}
}
func TestNotIn(t *testing.T) {
	var a int16 = 10
	err := in(a, "[1,2,3,4]")
	if err == nil {
		t.Error()
	} else {
		log.Printf(err.Error())
	}
	var b uint16 = 10
	err = in(b, "[1,2,3,4]")
	if err == nil {
		t.Error()
	} else {
		log.Printf(err.Error())
	}
	var c float64 = 1.332
	err = in(c, "[1.31, 1.29, 1.10]")
	if err == nil {
		t.Error()
	} else {
		log.Printf(err.Error())
	}
}
func TestIn(t *testing.T) {
	var a int = 1
	err := in(a, "[1,2,3,4]")
	if err != nil {
		log.Printf(err.Error())
		t.FailNow()
	} else {
		log.Printf("int ok")
	}
	var b uint = 1
	err = in(b, "[1,2,3,4]")
	if err != nil {
		log.Printf(err.Error())
		t.FailNow()
	} else {
		log.Printf("uint ok")
	}
	var c float64 = 1.31
	err = in(c, "[1.31, 1.29, 1.10]")
	if err != nil {
		log.Printf(err.Error())
		t.FailNow()
	} else {
		log.Printf("float ok")
	}
}
