package validator

import (
	"errors"
	"reflect"
)

var (
	INVALID_TYPE = errors.New("Unsupported type, the type must be struct")
)

func NewTag(tag string, kind TagKind, p TagHandler) *Tag {
	return &Tag{
		ID:      tag,
		Kind:    kind,
		Handler: p,
	}
}

type TagKind int

const (
	UNKNOWN TagKind = iota
	OP
	KEY
	WORD
)

type TagHandler interface {
	Match(string) bool
	Parse(string) (string, TagKind, string)
	Check(reflect.Value, string) error
}

type Tag struct {
	ID      string
	Kind    TagKind
	Handler TagHandler
}

func (t *Tag) Equal(v *Tag) bool {
	return v != nil && v.ID == t.ID && v.Kind == t.Kind
}
