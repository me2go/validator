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
func TestValidate(t *testing.T) {
	validator := NewValidator()
	t.Run("test >", func(tv *testing.T) {
		data := struct {
			a int     `validator:">10"`
			b float32 `validator:">10.1"`
			c string  `validator:">5"`
		}{}
		tv.Run("Success", func(tu *testing.T) {
			data.a = 11
			data.b = 11.0
			data.c = "hello,world"
			err := validator.Validate(data)
			if err == nil {

			} else {
				tu.Error(err.Error())
			}
		})
		tv.Run("Fail", func(tu *testing.T) {
			data.a = 10
			data.b = 10.1
			data.c = "me123"
			err := validator.Validate(data)
			if err == LESS_THAN_MIN {
				tu.Logf(err.Error())
			} else {
				tu.Fail()
			}
		})
	})

	t.Run("test >=", func(tv *testing.T) {
		data := struct {
			a int     `validator:">=10"`
			b float32 `validator:">=10.1"`
			c string  `validator:">=5"`
		}{}
		tv.Run("Success", func(tu *testing.T) {
			data.a = 10
			data.b = 10.1
			data.c = "hello"
			err := validator.Validate(data)
			if err == nil {

			} else {
				tu.Error(err.Error())
			}
		})
		tv.Run("Fail", func(tu *testing.T) {
			data.a = 9
			data.b = 10.0
			data.c = "me12"
			err := validator.Validate(data)
			if err == LESS_THAN_MIN {
				tu.Logf(err.Error())
			} else {
				tu.Fail()
			}
		})
	})

	t.Run("test <", func(tv *testing.T) {
		data := struct {
			a int     `validator:"<10"`
			b float32 `validator:"<10.1"`
			c string  `validator:"<5"`
		}{}
		tv.Run("Success", func(tu *testing.T) {
			data.a = 9
			data.b = 9.0
			data.c = "me02"
			err := validator.Validate(data)
			if err != nil {
				tu.Error(err.Error())
			}
		})
		tv.Run("Fail", func(tu *testing.T) {
			data.a = 10
			data.b = 10.1
			data.c = "me120"
			err := validator.Validate(data)
			if err == GREATER_THAN_MAX {
				tu.Logf(err.Error())
			} else {
				tu.Fail()
			}
		})
	})

	t.Run("test <=", func(tv *testing.T) {
		data := struct {
			a int     `validator:"<=10"`
			b float32 `validator:"<=10.1"`
			c string  `validator:"<=5"`
		}{}
		tv.Run("Success", func(tu *testing.T) {
			data.a = 10
			data.b = 10.1
			data.c = "me023"
			err := validator.Validate(data)
			if err != nil {
				tu.Error(err.Error())
			}
		})
		tv.Run("Fail", func(tu *testing.T) {
			data.a = 11
			data.b = 10.2
			data.c = "me1022"
			err := validator.Validate(data)
			if err == GREATER_THAN_MAX {
				tu.Logf(err.Error())
			} else {
				tu.Fail()
			}
		})
	})
	t.Run("test regexp", func(tv *testing.T) {
		data := struct {
			phone string `validator:"regexp:^1(31|32|33|35|36|37|38|39|51|52|55|56|58|59|77|89)\\d{8}$"`
		}{}
		tv.Run("Success", func(tu *testing.T) {
			data.phone = "13133841509"
			err := validator.Validate(data)
			if err != nil {
				tu.Error(err.Error())
			}
		})
		tv.Run("Fail", func(tu *testing.T) {
			data.phone = "1234567890"
			err := validator.Validate(data)
			if err == ERROR_REGEXP {
				tu.Logf(err.Error())
			} else {
				if err != nil {
					tu.Logf(err.Error())
				}
				tu.Fail()
			}
		})
	})
	t.Run("test non-zero", func(tv *testing.T) {
		data := struct {
			str string      `validator:"nonzero"`
			a   int         `validator:"nonzero"`
			i   interface{} `validator:"nonzero"`
		}{}
		tv.Run("Success", func(tu *testing.T) {
			data.str = "123"
			data.a = 100
			data.i = 10

			err := validator.Validate(data)
			if err != nil {
				tu.Logf(err.Error())
				tu.Fail()
			}
		})
		tv.Run("Fail", func(tu *testing.T) {
			data.str = ""
			data.a = 0
			data.i = nil
			err := validator.Validate(data)
			if err != ERROR_ZERO {
				if err != nil {
					tu.Logf(err.Error())
				}
				tu.Fail()
			} else {
				tu.Logf(err.Error())
			}
		})
	})
}
