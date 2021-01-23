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

func TestUnixAddress(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "helloworld",
			Validators: []Validator{WithUnixAddress()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "192.168.1.123",
			Validators: []Validator{WithUnixAddress()},
		},
	})
	a.NoError(err)
}

func TestHTML(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "hello",
			Validators: []Validator{WithHTML()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "<bhello",
			Validators: []Validator{WithHTML()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "<b>hello</b>",
			Validators: []Validator{WithHTML()},
		},
	})
	a.NoError(err)
}

func TestIPAddress(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "abcdefg",
			Validators: []Validator{WithIPAddress()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "0",
			Validators: []Validator{WithIPAddress()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "localhost",
			Validators: []Validator{WithIPAddress()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234",
			Validators: []Validator{WithIPAddress()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "192.168.1.123",
			Validators: []Validator{WithIPAddress()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "::0",
			Validators: []Validator{WithIPAddress()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			Validators: []Validator{WithIPAddress()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "127.0.0.1",
			Validators: []Validator{WithIPAddress()},
		},
	})
	a.NoError(err)
}

func TestIPv4Address(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "abcdefg",
			Validators: []Validator{WithIPv4Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			Validators: []Validator{WithIPv4Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "0",
			Validators: []Validator{WithIPv4Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "192.168.1.123:1234",
			Validators: []Validator{WithIPv4Address()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "192.168.1.123",
			Validators: []Validator{WithIPv4Address()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "127.0.0.1",
			Validators: []Validator{WithIPv4Address()},
		},
	})
	a.NoError(err)
}

func TestIPv6Address(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "abcdefg",
			Validators: []Validator{WithIPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "192.168.1.123",
			Validators: []Validator{WithIPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "127.0.0.1",
			Validators: []Validator{WithIPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "0",
			Validators: []Validator{WithIPv6Address()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234",
			Validators: []Validator{WithIPv6Address()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			Validators: []Validator{WithIPv6Address()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "::0",
			Validators: []Validator{WithIPv6Address()},
		},
	})
	a.NoError(err)
}

func TestUDPAddress(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "abcdefg",
			Validators: []Validator{WithUDPAddress()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "0",
			Validators: []Validator{WithUDPAddress()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "localhost",
			Validators: []Validator{WithUDPAddress()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "192.168.1.123",
			Validators: []Validator{WithUDPAddress()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "::0",
			Validators: []Validator{WithUDPAddress()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "192.168.1.123:1234",
			Validators: []Validator{WithUDPAddress()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "[::0]:1234",
			Validators: []Validator{WithUDPAddress()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234",
			Validators: []Validator{WithUDPAddress()},
		},
	})
	a.NoError(err)
}

func TestUDPv4Address(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "abcdefg",
			Validators: []Validator{WithUDPv4Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234",
			Validators: []Validator{WithUDPv4Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "0",
			Validators: []Validator{WithUDPv4Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "192.168.1.123",
			Validators: []Validator{WithUDPv4Address()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "127.0.0.1:1234",
			Validators: []Validator{WithUDPv4Address()},
		},
	})
	a.NoError(err)
}

func TestUDPv6Address(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "abcdefg",
			Validators: []Validator{WithUDPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "192.168.1.123",
			Validators: []Validator{WithUDPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "127.0.0.1",
			Validators: []Validator{WithUDPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "0",
			Validators: []Validator{WithUDPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "127.0.0.1:1234",
			Validators: []Validator{WithUDPv6Address()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234",
			Validators: []Validator{WithUDPv6Address()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "[::0]:1234",
			Validators: []Validator{WithUDPv6Address()},
		},
	})
	a.NoError(err)
}

func TestTCPAddress(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "abcdefg",
			Validators: []Validator{WithTCPAddress()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "0",
			Validators: []Validator{WithTCPAddress()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "localhost",
			Validators: []Validator{WithTCPAddress()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "192.168.1.123",
			Validators: []Validator{WithTCPAddress()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "::0",
			Validators: []Validator{WithTCPAddress()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "192.168.1.123:1234",
			Validators: []Validator{WithTCPAddress()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "[::0]:1234",
			Validators: []Validator{WithTCPAddress()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234",
			Validators: []Validator{WithTCPAddress()},
		},
	})
	a.NoError(err)
}

func TestTCPv4Address(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "abcdefg",
			Validators: []Validator{WithTCPv4Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234",
			Validators: []Validator{WithTCPv4Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "0",
			Validators: []Validator{WithTCPv4Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "192.168.1.123",
			Validators: []Validator{WithTCPv4Address()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "127.0.0.1:1234",
			Validators: []Validator{WithTCPv4Address()},
		},
	})
	a.NoError(err)
}

func TestTCPv6Address(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "abcdefg",
			Validators: []Validator{WithTCPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "192.168.1.123",
			Validators: []Validator{WithTCPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "127.0.0.1",
			Validators: []Validator{WithTCPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "0",
			Validators: []Validator{WithTCPv6Address()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "127.0.0.1:1234",
			Validators: []Validator{WithTCPv6Address()},
		},
	})
	a.Error(err)

	err = Validate([]Rule{
		{
			Value:      "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234",
			Validators: []Validator{WithTCPv6Address()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "[::0]:1234",
			Validators: []Validator{WithTCPv6Address()},
		},
	})
	a.NoError(err)
}

func TestLatitude(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "1234.92967312345678",
			Validators: []Validator{WithLatitude()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "35.",
			Validators: []Validator{WithLatitude()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "12345678",
			Validators: []Validator{WithLatitude()},
		},
	})

	err = Validate([]Rule{
		{
			Value:      "35.929673",
			Validators: []Validator{WithLatitude()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "35",
			Validators: []Validator{WithLatitude()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "-78.948237",
			Validators: []Validator{WithLatitude()},
		},
	})
	a.NoError(err)
}

func TestLongitude(t *testing.T) {
	a := assert.New(t)
	err := Validate([]Rule{
		{
			Value:      "1234.92967312345678",
			Validators: []Validator{WithLongitude()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "35.",
			Validators: []Validator{WithLongitude()},
		},
	})
	a.Error(err)
	err = Validate([]Rule{
		{
			Value:      "12345678",
			Validators: []Validator{WithLongitude()},
		},
	})

	err = Validate([]Rule{
		{
			Value:      "35.929673",
			Validators: []Validator{WithLongitude()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "35",
			Validators: []Validator{WithLongitude()},
		},
	})
	a.NoError(err)
	err = Validate([]Rule{
		{
			Value:      "-78.948237",
			Validators: []Validator{WithLongitude()},
		},
	})
	a.NoError(err)
}
