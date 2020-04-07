package tavern

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	a := assert.New(t)
	err := New().Add("", WithRequired()).Validate()
	a.Error(err)
	err = New().Add(0, WithRequired()).Validate()
	a.Error(err)
	err = New().Add([]string{}, WithRequired()).Validate()
	a.Error(err)

	err = New().Add("ABC", WithRequired()).Validate()
	a.NoError(err)
	err = New().Add(123, WithRequired()).Validate()
	a.NoError(err)
	err = New().Add([]string{"wow"}, WithRequired()).Validate()
	a.NoError(err)
}

func TestLength(t *testing.T) {
	a := assert.New(t)
	err := New().Add("", WithLength(1, 10)).Validate()
	a.Error(err)
	err = New().Add(10, WithLength(3, 10)).Validate()
	a.Error(err)
	err = New().Add([]string{}, WithLength(1, 10)).Validate()
	a.Error(err)

	err = New().Add("A", WithLength(1, 10)).Validate()
	a.NoError(err)
	err = New().Add(1000, WithLength(1, 10)).Validate()
	a.NoError(err)
	err = New().Add([]string{"wow"}, WithLength(1, 10)).Validate()
	a.NoError(err)
}

func TestMaxLength(t *testing.T) {
	a := assert.New(t)
	err := New().Add("ABCDEF", WithMaxLength(5)).Validate()
	a.Error(err)
	err = New().Add(100000, WithMaxLength(5)).Validate()
	a.Error(err)
	err = New().Add([]string{"A", "B", "C", "D", "E", "F"}, WithMaxLength(5)).Validate()
	a.Error(err)

	err = New().Add("ABC", WithMaxLength(3)).Validate()
	a.NoError(err)
	err = New().Add(10, WithMaxLength(3)).Validate()
	a.NoError(err)
	err = New().Add([]string{"A", "B"}, WithMaxLength(3)).Validate()
	a.NoError(err)
}

func TestMinLength(t *testing.T) {
	a := assert.New(t)
	err := New().Add("ABCD", WithMinLength(5)).Validate()
	a.Error(err)
	err = New().Add(1000, WithMinLength(5)).Validate()
	a.Error(err)
	err = New().Add([]string{"A", "B", "C"}, WithMinLength(5)).Validate()
	a.Error(err)

	err = New().Add("ABC", WithMinLength(3)).Validate()
	a.NoError(err)
	err = New().Add(1000, WithMinLength(3)).Validate()
	a.NoError(err)
	err = New().Add([]string{"A", "B", "D"}, WithMinLength(3)).Validate()
	a.NoError(err)
}

func TestFixedLength(t *testing.T) {
	a := assert.New(t)
	err := New().Add("ABCDEF", WithFixedLength(5)).Validate()
	a.Error(err)
	err = New().Add(1000, WithFixedLength(5)).Validate()
	a.Error(err)
	err = New().Add([]string{"A", "B", "C"}, WithFixedLength(5)).Validate()
	a.Error(err)

	err = New().Add("ABC", WithFixedLength(3)).Validate()
	a.NoError(err)
	err = New().Add(100, WithFixedLength(3)).Validate()
	a.NoError(err)
	err = New().Add([]string{"A", "B", "D"}, WithFixedLength(3)).Validate()
	a.NoError(err)
}

func TestRange(t *testing.T) {
	a := assert.New(t)
	err := New().Add(-1, WithRange(0, 5)).Validate()
	a.Error(err)
	err = New().Add(6, WithRange(0, 5)).Validate()
	a.Error(err)
	err = New().Add(-0.3, WithRange(0, 5)).Validate()
	a.Error(err)

	err = New().Add(3, WithRange(0, 5)).Validate()
	a.NoError(err)
	err = New().Add(0, WithRange(0, 5)).Validate()
	a.NoError(err)
	err = New().Add(0.3, WithRange(0, 5)).Validate()
	a.NoError(err)
}

func TestMaxRange(t *testing.T) {
	a := assert.New(t)
	err := New().Add(6, WithMaxRange(5)).Validate()
	a.Error(err)
	err = New().Add(5.3, WithMaxRange(5)).Validate()
	a.Error(err)

	err = New().Add(3, WithMaxRange(5)).Validate()
	a.NoError(err)
	err = New().Add(5, WithMaxRange(5)).Validate()
	a.NoError(err)
	err = New().Add(4.9, WithMaxRange(5)).Validate()
	a.NoError(err)
}

func TestMinRange(t *testing.T) {
	a := assert.New(t)
	err := New().Add(4, WithMinRange(5)).Validate()
	a.Error(err)
	err = New().Add(4.3, WithMinRange(5)).Validate()
	a.Error(err)

	err = New().Add(6, WithMinRange(5)).Validate()
	a.NoError(err)
	err = New().Add(5, WithMinRange(5)).Validate()
	a.NoError(err)
	err = New().Add(5.1, WithMinRange(5)).Validate()
	a.NoError(err)
}

func TestMaximum(t *testing.T) {
	a := assert.New(t)
	err := New().Add(11, WithMaximum(10)).Validate()
	a.Error(err)
	err = New().Add("ABC", WithMaximum(2)).Validate()
	a.Error(err)
	err = New().Add([]string{"A", "B", "C"}, WithMaximum(2)).Validate()
	a.Error(err)

	err = New().Add(3, WithMaximum(5)).Validate()
	a.NoError(err)
	err = New().Add("AB", WithMaximum(5)).Validate()
	a.NoError(err)
	err = New().Add([]string{"A", "B"}, WithMaximum(2)).Validate()
	a.NoError(err)
}

func TestDatetime(t *testing.T) {
	a := assert.New(t)
	err := New().Add("2009/01/23", WithDatetime("2006-01-02")).Validate()
	a.Error(err)
	err = New().Add("2018/04/39", WithDatetime("2006/01/02")).Validate()
	a.Error(err)
	err = New().Add("14:32", WithDatetime("03:04")).Validate()
	a.Error(err)

	err = New().Add("2009/01/23", WithDatetime("2006/01/02")).Validate()
	a.NoError(err)
	err = New().Add("2018-02-14", WithDatetime("2006-01-02")).Validate()
	a.NoError(err)
	err = New().Add("12:32", WithDatetime("03:04")).Validate()
	a.NoError(err)
	err = New().Add("14:32", WithDatetime("15:04")).Validate()
	a.NoError(err)
}

func TestEmail(t *testing.T) {
	a := assert.New(t)
	err := New().Add("yamiodymel@", WithEmail()).Validate()
	a.Error(err)
	err = New().Add("yamiodymel", WithEmail()).Validate()
	a.Error(err)
	err = New().Add("yamiodymel@xx@xx.com", WithEmail()).Validate()
	a.Error(err)
	err = New().Add("yamiodymel@x", WithEmail()).Validate()
	a.Error(err)
	err = New().Add("yamiodymel@x.", WithEmail()).Validate()
	a.Error(err)

	err = New().Add("yamiodymel@xx.com", WithEmail()).Validate()
	a.NoError(err)
}

func T