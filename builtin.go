package validator

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

var builtin = map[string]ValidateFunc{
	"required": required,
	"in":       in,
	"max":      max,
	"min":      min,
}

func required(v interface{}, param string) error {
	if v == nil {
		return errors.New("value is nil")
	}
	t := reflect.TypeOf(v)
	if v == reflect.Zero(t).Interface() {
		return errors.New("value is zero")
	}
	return nil
}

/**
 * 仅支持数字类型
 */
func in(v interface{}, param string) error {
	vtype := reflect.TypeOf(v)
	slicePtr := reflect.New(reflect.SliceOf(vtype))
	err := json.Unmarshal([]byte(param), slicePtr.Interface())
	if err != nil {
		return err
	}
	for i, num := 0, slicePtr.Elem().Len(); i < num; i++ {
		if v == slicePtr.Elem().Index(i).Interface() {
			return nil
		}
	}
	return errors.New("not in")
}
func max(v interface{}, param string) error {
	ret, err := compare(v, param)
	if err != nil {
		return err
	}
	if ret > 0 {
		return GREATER_THAN_MAX
	}
	return nil
}
func min(v interface{}, param string) error {
	ret, err := compare(v, param)
	if err != nil {
		return err
	}
	if ret < 0 {
		return LOWER_THAN_MIN
	}
	return nil
}
func compare(v interface{}, param string) (int, error) {
	switch v.(type) {
	case uint, uint8, uint16, uint32, uint64:
		p, err := strconv.ParseUint(param, 0, 64)
		if err != nil {
			return 0, err
		}
		temp := reflect.ValueOf(v).Uint()
		if temp > p {
			return 1, nil
		} else if temp < p {
			return -1, nil
		} else {
			return 0, nil
		}
	case int, int8, int16, int32, int64:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return 0, err
		}
		temp := reflect.ValueOf(v).Int()
		if temp > p {
			return 1, nil
		} else if temp < p {
			return -1, nil
		} else {
			return 0, nil
		}

	case string:
		str := v.(string)
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return 0, nil
		}
		temp := int64(len([]rune(str)))
		if temp > p {
			return 1, nil
		} else if temp < p {
			return -1, nil
		} else {
			return 0, nil
		}
	default:
		return 0, NOT_SUPPORTED
	}
}
