package tavern

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	assert := assert.New(t)
	err := Add("123456").Check()
	assert.NoError(err)
	err = Add("").Check()
	assert.NoError(err)
	err = Add("0").Check()
	assert.NoError(err)

	err = Add("123456").Required().Check()
	assert.NoError(err)
	err = Add("").Required().Check()
	assert.Error(err)
	err = Add(0).Required().Check()
	assert.NoError(err)

	var s interface{}
	s = ""
	err = Add(s).Required().Check()
	assert.Error(err)
	s = nil
	err = Add(s).Required().Check()
	assert.Error(err)
	s = nil
	err = Add(s).Check()
	assert.NoError(err)
	s = 0
	err = Add(s).Required().Check()
	assert.NoError(err)
}

func TestLength(t *testing.T) {
	assert := assert.New(t)
	err := Add("123456").Length(1, 3).Check()
	assert.Error(err)
	err = Add("123456").Length(1, 6).Check()
	assert.NoError(err)
	err = Add("123456").Length(0, 1).Check()
	assert.Error(err)
	err = Add("123456").Length(0, 0).Check()
	assert.Error(err)
	err = Add("").Length(0, 0).Check()
	assert.NoError(err)
	err = Add("").Length(1, 2).Check()
	assert.Error(err)

	err = Add(123456).Length(1, 3).Check()
	assert.Error(err)
	err = Add(123456).Length(1, 6).Check()
	assert.NoError(err)
	err = Add(123456).Length(0, 1).Check()
	assert.Error(err)
	err = Add(123456).Length(0, 0).Check()
	assert.Error(err)
}

func TestRange(t *testing.T) {
	assert := assert.New(t)
	err := Add(1234).Range(1, 1234).Check()
	assert.NoError(err)
	err = Add(1234).Range(1, 1233).Check()
	assert.Error(err)
	err = Add(-1).Range(1, 1233).Check()
	assert.Error(err)
	err = Add(1).Range(0, 1).Check()
	assert.NoError(err)
}

func TestMin(t *testing.T) {
	assert := assert.New(t)
	err := Add(1).Min(100).Check()
	assert.Error(err)
	err = Add(1).Min(1).Check()
	assert.NoError(err)
	err = Add(-1).Min(1).Check()
	assert.Error(err)
}

func TestMax(t *testing.T) {
	assert := assert.New(t)
	err := Add(101).Max(100).Check()
	assert.Error(err)
	err = Add(100).Max(100).Check()
	assert.NoError(err)
	err = Add(10).Max(100).Check()
	assert.NoError(err)
}

func TestDate(t *testing.T) {
	assert := assert.New(t)
	err := Add("1998-07-13").Date("2006-01-02").Check()
	assert.NoError(err)
	err = Add("1998-07-32").Date("2006-01-02").Check()
	assert.Error(err)
	err = Add("13-07-1998").Date("2006-01-02").Check()
	assert.Error(err)
	err = Add("12:59:59").Date("03:04:05").Check()
	assert.NoError(err)
	err = Add("13:59:59").Date("03:04:05").Check()
	assert.Error(err)
	err = Add("13:59:59").Date("15:04:05").Check()
	assert.NoError(err)
	err = Add("135959").Date("15:04:05").Check()
	assert.Error(err)
}

func TestEmail(t *testing.T) {
	assert := assert.New(t)
	err := Add("yamiodymel@gmail.com").Email().Check()
	assert.NoError(err)
	err = Add("").Email().Check()
	assert.NoError(err)
	err = Add("").Email().Required().Check()
	assert.Error(err)
	err = Add("xxxxxxx").Email().Check()
	assert.Error(err)
}

func TestIn(t *testing.T) {
	assert := assert.New(t)
	err := Add("Female").In("Male", "Female", "Transgender").Check()
	assert.NoError(err)
	err = Add("Female").In("Male", "Transgender").Check()
	assert.Error(err)
	//err = Add("").In("Male", "Transgender").Check()
	//assert.NoError(err)
	err = Add("").In("Male", "Transgender").Required().Check()
	assert.Error(err)
}

func TestIP(t *testing.T) {
	assert := assert.New(t)
	//err := Add("").IP().Check()
	//assert.NoError(err)
	err := Add("").IP().Required().Check()
	assert.Error(err)
	err = Add("192.168.0.0").IP("v6").Check()
	assert.Error(err)
	err = Add("192.168.0.0").IP("v4").Check()
	assert.NoError(err)
	err = Add("::1").IP("v6").Check()
	assert.NoError(err)
	err = Add("2001:0db8:0a0b:12f0:0000:0000:0000:0001").IP("v6").Check()
	assert.NoError(err)
}

func TestURL(t *testing.T) {
	assert := assert.New(t)
	//err := Add("").URL().Check()
	//assert.NoError(err)
	err := Add("").URL().Required().Check()
	assert.Error(err)
	err = Add("abcde").URL().Check()
	assert.Error(err)
	err = Add("http://www.google.com/").URL().Check()
	assert.NoError(err)
	err = Add("google.com").URL().Check()
	assert.Error(err)
	err = Add("www.google.com").URL().Check()
	assert.Error(err)
	err = Add("ftp://google.com").URL().Check()
	assert.NoError(err)
	err = Add("ftp://google.com").URL("http://").Check()
	assert.Error(err)
	err = Add("http://google.com").URL("http://").Check()
	assert.NoError(err)
	err = Add("ftp://google.com").URL("http://", "ftp://").Check()
	assert.NoError(err)
}

func TestEqual(t *testing.T) {
	assert := assert.New(t)
	err := Add("").Equal("").Check()
	assert.NoError(err)
	err = Add("123").Equal("").Check()
	assert.Error(err)
	err = Add(12345).Equal("12345").Check()
	assert.Error(err)

	var a interface{}
	var b interface{}
	a = 1
	b = 1
	err = Add(a).Equal(b).Check()
	assert.NoError(err)
}

func TestRegExp(t *testing.T) {
	assert := assert.New(t)
	err := Add(12345).RegExp(`^[0-9]*$`).Check()
	assert.NoError(err)
	err = Add("12345").RegExp(`^[0-9]*$`).Check()
	assert.NoError(err)
	err = Add("ABCDEFG").RegExp(`^[0-9]*$`).Check()
	assert.Error(err)
}

func TestError(t *testing.T) {
	assert := assert.New(t)
	err := Add("123").Equal("").Error(errors.New("Fuck")).Check()
	assert.Equal(err.Error(), "Fuck")
	err = Add("123").Equal("123").Error(errors.New("Fuck")).Check()
	assert.NoError(err)
}

func TestErrors(t *testing.T) {
	assert := assert.New(t)
	err := Add("123").Length(5, 10).Error(E{
		Length: errors.New("Length"),
	}).Check()
	assert.Equal(err.Error(), "Length")
	err = Add("123").Length(5, 10).Equal("12345").Error(E{
		Length: errors.New("Length"),
		Equal:  errors.New("Equal"),
	}).Check()
	assert.Equal(err.Error(), "Length")
	err = Add("123").Length(0, 10).Equal("12345").Error(E{
		Length: errors.New("Length"),
		Equal:  errors.New("Equal"),
	}).Check()
	assert.Equal(err.Error(), "Equal")
}

func TestCheckAll(t *testing.T) {
	assert := assert.New(t)
	errs := Add("123").Length(5, 10).Equal("12345").Error(E{
		Length: errors.New("Length"),
		Equal:  errors.New("Equal"),
	}).Add("123").Length(5, 10).Equal("12345").Error(E{
		Length: errors.New("Length"),
		Equal:  errors.New("Equal"),
	}).CheckAll()
	assert.Len(errs, 2)
	assert.Equal(errs[0].Error(), "Length")
	assert.Equal(errs[1].Error(), "Length")
}
