package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	GREATER_THAN_MAX = errors.New("greater than the max")
	LESS_THAN_MIN    = errors.New("less than the min")
	UNSUPPORTED_TYPE = errors.New("unsupported type")
	ERROR_EQUAL      = errors.New("equal")
	ERROR_NONEQUAL   = errors.New("non-equal")
	ERROR_REGEXP     = errors.New("unmatching pattern")
	ERROR_ZERO       = errors.New("the value should not be zero")
)
var (
	GREATEROREQUAL = NewTag(">=", OP, &greaterOrEqualHandler{})
	GREATER        = NewTag(">", OP, &greaterHandler{})
	LESSOREQUAL    = NewTag("<=", OP, &lessOrEqualHandler{})
	LESS           = NewTag("<", OP, &lessHandler{})
	EQUAL          = NewTag("==", OP, &equalHandler{})
	NONEQUAL       = NewTag("!=", OP, &nonequalHandler{})

	REGEXP = NewTag("regexp", KEY, &regexpHandler{})

	NONZERO = NewTag("nonzero", KEY, &nonequalHandler{})
)

func compare(v reflect.Value, c string) (int, error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intV := v.Int()
		intC, err := strconv.ParseInt(c, 0, 64)
		if err != nil {
			return 0, err
		} else {
			if intV < intC {
				return -1, nil
			} else if intV == intC {
				return 0, nil
			} else {
				return 1, nil
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintV := v.Uint()
		uintC, err := strconv.ParseUint(c, 0, 64)
		if err != nil {
			return 0, err
		} else {
			if uintV < uintC {
				return -1, nil
			} else if uintV == uintC {
				return 0, nil
			} else {
				return 1, nil
			}
		}
	case reflect.Float32, reflect.Float64:
		floatV := v.Float()
		floatC, err := strconv.ParseFloat(c, 32)
		if err != nil {
			return 0, err
		} else {
			if floatV < floatC {
				return -1, nil
			} else if floatV == floatC {
				return 0, nil
			} else {
				return 1, nil
			}
		}
	case reflect.String:
		strV := v.String()
		num, err := strconv.Atoi(c)
		if err != nil {
			return 0, err
		}
		if len([]rune(strV)) < num {
			return -1, nil
		} else if len([]rune(strV)) == num {
			return 0, nil
		} else {
			return 1, nil
		}
	}
	return 0, UNSUPPORTED_TYPE
}

type greaterHandler struct{}

func (g *greaterHandler) Match(s string) bool {
	return strings.HasPrefix(s, ">") && []rune(s)[1] >= '0' && []rune(s)[1] <= '9'
}

func (g *greaterHandler) Parse(s string) (string, TagKind, string) {
	return ">", OP, strings.Trim(s[1:], " ")
}

func (g *greaterHandler) Check(v reflect.Value, c string) error {
	ret, err := compare(v, c)
	if err != nil {
		return err
	}
	if ret <= 0 {
		return LESS_THAN_MIN
	}
	return nil
}

type greaterOrEqualHandler struct {
}

func (g *greaterOrEqualHandler) Match(s string) bool {
	return strings.HasPrefix(s, ">=")
}
func (g *greaterOrEqualHandler) Parse(s string) (string, TagKind, string) {
	return ">=", OP, strings.Trim(s[2:], " ")
}

func (g *greaterOrEqualHandler) Check(v reflect.Value, c string) error {
	ret, err := compare(v, c)
	if err != nil {
		return err
	}
	if ret < 0 {
		return LESS_THAN_MIN
	}
	return nil
}

type lessHandler struct{}

func (l *lessHandler) Match(s string) bool {
	return strings.HasPrefix(s, "<") && []rune(s)[1] >= '0' && []rune(s)[1] <= '9'
}

func (l *lessHandler) Parse(s string) (string, TagKind, string) {
	return "<", OP, strings.Trim(s[1:], " ")
}

func (l *lessHandler) Check(v reflect.Value, c string) error {
	ret, err := compare(v, c)
	if err != nil {
		return err
	}
	if ret >= 0 {
		return GREATER_THAN_MAX
	}
	return nil
}

type lessOrEqualHandler struct {
}

func (l *lessOrEqualHandler) Match(s string) bool {
	return strings.HasPrefix(s, "<=")
}

func (l *lessOrEqualHandler) Parse(s string) (string, TagKind, string) {
	return "<=", OP, strings.Trim(s[2:], " ")
}

func (l *lessOrEqualHandler) Check(v reflect.Value, c string) error {
	ret, err := compare(v, c)
	if err != nil {
		return err
	}
	if ret > 0 {
		return GREATER_THAN_MAX
	}
	return nil
}

type equalHandler struct{}

func (e *equalHandler) Match(s string) bool {
	return strings.HasPrefix(s, "==")
}
func (e *equalHandler) Parse(s string) (string, TagKind, string) {
	return "==", OP, strings.Trim(s[2:], " ")
}
func (e *equalHandler) Check(v reflect.Value, c string) error {
	ret, err := compare(v, c)
	if err != nil {
		return err
	}
	if ret != 0 {
		return ERROR_NONEQUAL
	}
	return nil
}

type nonequalHandler struct{}

func (e *nonequalHandler) Match(s string) bool {
	return strings.HasPrefix(s, "!=")
}
func (e *nonequalHandler) Parse(s string) (string, TagKind, string) {
	return "!=", OP, strings.Trim(s[2:], " ")
}
func (e *nonequalHandler) Check(v reflect.Value, c string) error {
	ret, err := compare(v, c)
	if err != nil {
		return err
	}
	if ret == 0 {
		return ERROR_EQUAL
	}
	return nil
}

type regexpHandler struct{}

func (r *regexpHandler) Match(s string) bool {
	return strings.HasPrefix(s, "regexp")
}
func (r *regexpHandler) Parse(s string) (string, TagKind, string) {
	parts := strings.Split(s, ":")
	return "regexp", KEY, strings.Trim(parts[1], " ")
}
func (r *regexpHandler) Check(v reflect.Value, c string) error {
	if v.Kind() != reflect.String {
		return UNSUPPORTED_TYPE
	}
	reg, err := regexp.Compile(c)
	if err != nil {
		return err
	}
	if !reg.MatchString(v.String()) {
		return ERROR_REGEXP
	}
	return nil
}

type nonzeroHandler struct{}

func (r *nonzeroHandler) Match(s string) bool {
	return s == "nonzero"
}
func (r *nonzeroHandler) Parse(s string) (string, TagKind, string) {
	return "nonzero", WORD, ""
}
func (r *nonzeroHandler) Check(v reflect.Value, c string) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return ERROR_ZERO
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() == 0 {
			return ERROR_ZERO
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0 {
			return ERROR_ZERO
		}
	case reflect.String:
		if v.String() == "" {
			return ERROR_ZERO
		}

	case reflect.Slice, reflect.Map, reflect.Interface, reflect.Func, reflect.Ptr, reflect.Chan:
		if v.IsNil() {
			return ERROR_ZERO
		}
	}
	return nil
}
