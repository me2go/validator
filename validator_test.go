package validator

import "testing"

func TestRegisterTag(t *testing.T) {
	v := NewValidator()
	tag := NewTag("test", WORD, nil)
	v.RegisterTag(tag)
	if !v.HasTag(tag) {
		t.FailNow()
	}
}
func TestUnregisterTag(t *testing.T) {
	v := NewValidator()
	v.UnregisterTag(GREATEROREQUAL)
	if v.HasTag(GREATEROREQUAL) {
		t.FailNow()
	}
}

func TestHasTag(t *testing.T) {
	v := NewValidator()
	if !v.HasTag(GREATEROREQUAL) {
		t.Error(">= exists")
	}
	tag := NewTag("test", WORD, nil)
	if v.HasTag(tag) {
		t.Error("the test tag does not exist")
	}
}
