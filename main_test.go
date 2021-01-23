package tavern

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "",
			Validators: []Validator{WithRequired()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      0,
			Validators: []Validator{WithRequired()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      []string{},
			Validators: []Validator{WithRequired()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "ABC",
			Validators: []Validator{WithRequired()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      123,
			Validators: []Validator{WithRequired()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      []string{"wow"},
			Validators: []Validator{WithRequired()},
		},
	})
	a.NoError(err)
}

func TestLength(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "",
			Validators: []Validator{WithRequired(), WithLength(1, 10)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      10,
			Validators: []Validator{WithLength(3, 10)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      []string{},
			Validators: []Validator{WithLength(1, 10)},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "",
			Validators: []Validator{WithLength(1, 10)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "A",
			Validators: []Validator{WithLength(1, 10)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      1000,
			Validators: []Validator{WithLength(1, 10)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      []string{"wow"},
			Validators: []Validator{WithLength(1, 10)},
		},
	})
	a.NoError(err)
}

func TestMaxLength(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "ABCDEF",
			Validators: []Validator{WithMaxLength(5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      100000,
			Validators: []Validator{WithMaxLength(5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      []string{"A", "B", "C", "D", "E", "F"},
			Validators: []Validator{WithMaxLength(5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "ABC",
			Validators: []Validator{WithMaxLength(3)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      10,
			Validators: []Validator{WithMaxLength(3)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      []string{"A", "B"},
			Validators: []Validator{WithMaxLength(3)},
		},
	})
	a.NoError(err)
}

func TestMinLength(t *testing.T) {
	a := assert.New(t)

	err := Validate([]Rule{
		{
			Value:      "ABCD",
			Validators: []Validator{WithMinLength(5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      1000,
			Validators: []Validator{WithMinLength(5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      []string{"A", "B", "C"},
			Validators: []Validator{WithMinLength(5)},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "",
			Validators: []Validator{WithMinLength(5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "ABC",
			Validators: []Validator{WithMinLength(3)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      1000,
			Validators: []Validator{WithMinLength(3)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      []string{"A", "B", "D"},
			Validators: []Validator{WithMinLength(3)},
		},
	})
	a.NoError(err)
}

func TestFixedLength(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "ABCDEF",
			Validators: []Validator{WithFixedLength(5)},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      1000,
			Validators: []Validator{WithFixedLength(5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      []string{"A", "B", "C"},
			Validators: []Validator{WithFixedLength(5)},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "",
			Validators: []Validator{WithFixedLength(3)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "ABC",
			Validators: []Validator{WithFixedLength(3)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      100,
			Validators: []Validator{WithFixedLength(3)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      []string{"A", "B", "D"},
			Validators: []Validator{WithFixedLength(3)},
		},
	})
	a.NoError(err)
}

func TestRange(t *testing.T) {
	a := assert.New(t)

	err := Validate([]Rule{
		{
			Value:      -1,
			Validators: []Validator{WithRange(0, 5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      6,
			Validators: []Validator{WithRange(0, 5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      -0.3,
			Validators: []Validator{WithRange(0, 5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      0,
			Validators: []Validator{WithRequired(), WithRange(1, 5)},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      0,
			Validators: []Validator{WithRange(1, 5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      3,
			Validators: []Validator{WithRange(0, 5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      0,
			Validators: []Validator{WithRange(0, 5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      0.3,
			Validators: []Validator{WithRange(0, 5)},
		},
	})
	a.NoError(err)
}

func TestMaxRange(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      6,
			Validators: []Validator{WithMaxRange(5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      5.3,
			Validators: []Validator{WithMaxRange(5)},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      3,
			Validators: []Validator{WithMaxRange(5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      5,
			Validators: []Validator{WithMaxRange(5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      4.9,
			Validators: []Validator{WithMaxRange(5)},
		},
	})
	a.NoError(err)
}

func TestMinRange(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      4,
			Validators: []Validator{WithMinRange(5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      4.3,
			Validators: []Validator{WithMinRange(5)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      0,
			Validators: []Validator{WithRequired(), WithMinRange(5)},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      0,
			Validators: []Validator{WithMinRange(5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      6,
			Validators: []Validator{WithMinRange(5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      5,
			Validators: []Validator{WithMinRange(5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      5.1,
			Validators: []Validator{WithMinRange(5)},
		},
	})
	a.NoError(err)
}

func TestMaximum(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      11,
			Validators: []Validator{WithMaximum(10)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "ABC",
			Validators: []Validator{WithMaximum(2)},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      []string{"A", "B", "C"},
			Validators: []Validator{WithMaximum(2)},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      3,
			Validators: []Validator{WithMaximum(5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "AB",
			Validators: []Validator{WithMaximum(5)},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      []string{"A", "B"},
			Validators: []Validator{WithMaximum(2)},
		},
	})
	a.NoError(err)
}

func TestDatetime(t *testing.T) {
	a := assert.New(t)

	err := Validate([]Rule{
		{
			Value:      "2009/01/23",
			Validators: []Validator{WithDatetime("2006-01-02")},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "2018/04/39",
			Validators: []Validator{WithDatetime("2006/01/02")},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "14:32",
			Validators: []Validator{WithDatetime("03:04")},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "2009/01/23",
			Validators: []Validator{WithDatetime("2006/01/02")},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "2018-02-14",
			Validators: []Validator{WithDatetime("2006-01-02")},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "12:32",
			Validators: []Validator{WithDatetime("03:04")},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "14:32",
			Validators: []Validator{WithDatetime("15:04")},
		},
	})
	a.NoError(err)
}

func TestEmail(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "yamiodymel@",
			Validators: []Validator{WithEmail()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "yamiodymel",
			Validators: []Validator{WithEmail()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "yamiodymel@xx@xx.com",
			Validators: []Validator{WithEmail()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "yamiodymel@x",
			Validators: []Validator{WithEmail()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "yamiodymel@x.",
			Validators: []Validator{WithEmail()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "yamiodymel@xx.com",
			Validators: []Validator{WithEmail()},
		},
	})
	a.NoError(err)
}
