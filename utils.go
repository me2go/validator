package validator

var global = NewValidator()

func Validate(v interface{}) error {
	return global.Validate(v)
}
