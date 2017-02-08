package validator

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	INVALID_VALUE    = errors.New("the value should be a struct")
	GREATER_THAN_MAX = errors.New("greater than the max")
	LOWER_THAN_MIN   = errors.New("lower than the min")
	NOT_SUPPORTED    = errors.New("not supported")
)

type ValidatorError map[string]error

func (ve ValidatorError) Error() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("{")
	for name, e := range ve {
		buf.WriteString(fmt.Sprintf("[%v: %v]", name, e.Error()))
	}
	buf.WriteString("}")
	return buf.String()
}
