package tavern

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("", WithRequired()))
	a.Error(err)
	err = Validate(NewRule(0, WithRequired()))
	a.Error(err)

	err = Validate(NewRule([]string{}, WithRequired()))
	a.NoError(err)
	err = Validate(NewRule("ABC", WithRequired()))
	a.NoError(err)
	err = Validate(NewRule(123, WithRequired()))
	a.NoError(err)
	err = Validate(NewRule([]string{"wow"}, WithRequired()))
	a.NoError(err)
}

func TestLength(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("", WithRequired(), WithLength(1, 10)))
	a.Error(err)
	err = Validate(NewRule(10, WithLength(3, 10)))
	a.Error(err)
	err = Validate(NewRule([]string{}, WithLength(1, 10)))
	a.Error(err)

	err = Validate(NewRule("", WithLength(1, 10)))
	a.NoError(err)
	err = Validate(NewRule("A", WithLength(1, 10)))
	a.NoError(err)
	err = Validate(NewRule(1000, WithLength(1, 10)))
	a.NoError(err)
	err = Validate(NewRule([]string{"wow"}, WithLength(1, 10)))
	a.NoError(err)
}

func TestMaxLength(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("ABCDEF", WithMaxLength(5)))
	a.Error(err)
	err = Validate(NewRule(100000, WithMaxLength(5)))
	a.Error(err)
	err = Validate(NewRule([]string{"A", "B", "C", "D", "E", "F"}, WithMaxLength(5)))
	a.Error(err)
	err = Validate(NewRule("ABC", WithMaxLength(3)))

	a.NoError(err)
	err = Validate(NewRule(10, WithMaxLength(3)))
	a.NoError(err)
	err = Validate(NewRule([]string{"A", "B"}, WithMaxLength(3)))
	a.NoError(err)
}

func TestMinLength(t *testing.T) {
	a := assert.New(t)

	err := Validate(NewRule("ABCD", WithMinLength(5)))
	a.Error(err)
	err = Validate(NewRule(1000, WithMinLength(5)))
	a.Error(err)
	err = Validate(NewRule([]string{"A", "B", "C"}, WithMinLength(5)))
	a.Error(err)

	err = Validate(NewRule("", WithMinLength(5)))
	a.NoError(err)
	err = Validate(NewRule("ABC", WithMinLength(3)))
	a.NoError(err)
	err = Validate(NewRule(1000, WithMinLength(3)))
	a.NoError(err)
	err = Validate(NewRule([]string{"A", "B", "D"}, WithMinLength(3)))
	a.NoError(err)
}

func TestFixedLength(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("ABCDEF", WithFixedLength(5)))
	a.Error(err)
	err = Validate(NewRule(1000, WithFixedLength(5)))
	a.Error(err)
	err = Validate(NewRule([]string{"A", "B", "C"}, WithFixedLength(5)))
	a.Error(err)

	err = Validate(NewRule("", WithFixedLength(3)))
	a.NoError(err)
	err = Validate(NewRule("ABC", WithFixedLength(3)))
	a.NoError(err)
	err = Validate(NewRule(100, WithFixedLength(3)))
	a.NoError(err)
	err = Validate(NewRule([]string{"A", "B", "D"}, WithFixedLength(3)))
	a.NoError(err)
}

func TestRange(t *testing.T) {
	a := assert.New(t)

	err := Validate(NewRule(-1, WithRange(0, 5)))
	a.Error(err)
	err = Validate(NewRule(6, WithRange(0, 5)))
	a.Error(err)
	err = Validate(NewRule(-0.3, WithRange(0, 5)))
	a.Error(err)
	err = Validate(NewRule(0, WithRequired(), WithRange(1, 5)))
	a.Error(err)

	err = Validate(NewRule(0, WithRange(1, 5)))
	a.NoError(err)
	err = Validate(NewRule(3, WithRange(0, 5)))
	a.NoError(err)
	err = Validate(NewRule(0, WithRange(0, 5)))
	a.NoError(err)
	err = Validate(NewRule(0.3, WithRange(0, 5)))
	a.NoError(err)
}

func TestMaxRange(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule(6, WithMaxRange(5)))
	a.Error(err)
	err = Validate(NewRule(5.3, WithMaxRange(5)))
	a.Error(err)

	err = Validate(NewRule(3, WithMaxRange(5)))
	a.NoError(err)
	err = Validate(NewRule(5, WithMaxRange(5)))
	a.NoError(err)
	err = Validate(NewRule(4.9, WithMaxRange(5)))
	a.NoError(err)
}

func TestMinRange(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule(4, WithMinRange(5)))
	a.Error(err)
	err = Validate(NewRule(4.3, WithMinRange(5)))
	a.Error(err)
	err = Validate(NewRule(0, WithRequired(), WithMinRange(5)))
	a.Error(err)

	err = Validate(NewRule(0, WithMinRange(5)))
	a.NoError(err)
	err = Validate(NewRule(6, WithMinRange(5)))
	a.NoError(err)
	err = Validate(NewRule(5, WithMinRange(5)))
	a.NoError(err)
	err = Validate(NewRule(5.1, WithMinRange(5)))
	a.NoError(err)
}

func TestMaximum(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule(11, WithMaximum(10)))
	a.Error(err)
	err = Validate(NewRule("ABC", WithMaximum(2)))
	a.Error(err)
	err = Validate(NewRule([]string{"A", "B", "C"}, WithMaximum(2)))
	a.Error(err)

	err = Validate(NewRule(3, WithMaximum(5)))
	a.NoError(err)
	err = Validate(NewRule("AB", WithMaximum(5)))
	a.NoError(err)
	err = Validate(NewRule([]string{"A", "B"}, WithMaximum(2)))
	a.NoError(err)
}

func TestDatetime(t *testing.T) {
	a := assert.New(t)

	err := Validate(NewRule("2009/01/23", WithDatetime("2006-01-02")))
	a.Error(err)
	err = Validate(NewRule("2018/04/39", WithDatetime("2006/01/02")))
	a.Error(err)
	err = Validate(NewRule("14:32", WithDatetime("03:04")))
	a.Error(err)

	err = Validate(NewRule("2009/01/23", WithDatetime("2006/01/02")))
	a.NoError(err)
	err = Validate(NewRule("2018-02-14", WithDatetime("2006-01-02")))
	a.NoError(err)
	err = Validate(NewRule("12:32", WithDatetime("03:04")))
	a.NoError(err)
	err = Validate(NewRule("14:32", WithDatetime("15:04")))
	a.NoError(err)
}

func TestEmail(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("yamiodymel@", WithEmail()))
	a.Error(err)
	err = Validate(NewRule("yamiodymel", WithEmail()))
	a.Error(err)
	err = Validate(NewRule("yamiodymel@xx@xx.com", WithEmail()))
	a.Error(err)
	err = Validate(NewRule("yamiodymel@x", WithEmail()))
	a.Error(err)
	err = Validate(NewRule("yamiodymel@x.", WithEmail()))
	a.Error(err)

	err = Validate(NewRule("yamiodymel@xx.com", WithEmail()))
	a.NoError(err)
}

func TestUnixAddress(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("helloworld", WithUnixAddress()))
	a.NoError(err)
	err = Validate(NewRule("192.168.1.123", WithUnixAddress()))
	a.NoError(err)
}

func TestHTML(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("hello", WithHTML()))
	a.Error(err)
	err = Validate(NewRule("<bhello", WithHTML()))
	a.Error(err)

	err = Validate(NewRule("<b>hello</b>", WithHTML()))
	a.NoError(err)
}

func TestIPAddress(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("abcdefg", WithIPAddress()))
	a.Error(err)
	err = Validate(NewRule("0", WithIPAddress()))
	a.Error(err)

	err = Validate(NewRule("localhost", WithIPAddress()))
	a.NoError(err)
	err = Validate(NewRule("[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234", WithIPAddress()))
	a.NoError(err)
	err = Validate(NewRule("192.168.1.123", WithIPAddress()))
	a.NoError(err)
	err = Validate(NewRule("::0", WithIPAddress()))
	a.NoError(err)
	err = Validate(NewRule("2001:0db8:85a3:0000:0000:8a2e:0370:7334", WithIPAddress()))
	a.NoError(err)
	err = Validate(NewRule("127.0.0.1", WithIPAddress()))
	a.NoError(err)
}

func TestIPv4Address(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("abcdefg", WithIPv4Address()))
	a.Error(err)
	err = Validate(NewRule("2001:0db8:85a3:0000:0000:8a2e:0370:7334", WithIPv4Address()))
	a.Error(err)
	err = Validate(NewRule("0", WithIPv4Address()))
	a.Error(err)
	err = Validate(NewRule("192.168.1.123:1234", WithIPv4Address()))
	a.Error(err)

	err = Validate(NewRule("192.168.1.123", WithIPv4Address()))
	a.NoError(err)
	err = Validate(NewRule("127.0.0.1", WithIPv4Address()))
	a.NoError(err)
}

func TestIPv6Address(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("abcdefg", WithIPv6Address()))
	a.Error(err)
	err = Validate(NewRule("192.168.1.123", WithIPv6Address()))
	a.Error(err)
	err = Validate(NewRule("127.0.0.1", WithIPv6Address()))
	a.Error(err)
	err = Validate(NewRule("0", WithIPv6Address()))
	a.Error(err)

	err = Validate(NewRule("[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234", WithIPv6Address()))
	a.NoError(err)
	err = Validate(NewRule("2001:0db8:85a3:0000:0000:8a2e:0370:7334", WithIPv6Address()))
	a.NoError(err)
	err = Validate(NewRule("::0", WithIPv6Address()))
	a.NoError(err)
}

func TestUDPAddress(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("abcdefg", WithUDPAddress()))
	a.Error(err)
	err = Validate(NewRule("0", WithUDPAddress()))
	a.Error(err)
	err = Validate(NewRule("localhost", WithUDPAddress()))
	a.Error(err)
	err = Validate(NewRule("192.168.1.123", WithUDPAddress()))
	a.Error(err)
	err = Validate(NewRule("::0", WithUDPAddress()))
	a.Error(err)

	err = Validate(NewRule("192.168.1.123:1234", WithUDPAddress()))
	a.NoError(err)
	err = Validate(NewRule("[::0]:1234", WithUDPAddress()))
	a.NoError(err)
	err = Validate(NewRule("[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234", WithUDPAddress()))
	a.NoError(err)
}

func TestUDPv4Address(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("abcdefg", WithUDPv4Address()))
	a.Error(err)
	err = Validate(NewRule("[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234", WithUDPv4Address()))
	a.Error(err)
	err = Validate(NewRule("0", WithUDPv4Address()))
	a.Error(err)
	err = Validate(NewRule("192.168.1.123", WithUDPv4Address()))
	a.Error(err)

	err = Validate(NewRule("127.0.0.1:1234", WithUDPv4Address()))
	a.NoError(err)
}

func TestUDPv6Address(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("abcdefg", WithUDPv6Address()))
	a.Error(err)
	err = Validate(NewRule("192.168.1.123", WithUDPv6Address()))
	a.Error(err)
	err = Validate(NewRule("127.0.0.1", WithUDPv6Address()))
	a.Error(err)
	err = Validate(NewRule("0", WithUDPv6Address()))
	a.Error(err)
	err = Validate(NewRule("127.0.0.1:1234", WithUDPv6Address()))
	a.Error(err)

	err = Validate(NewRule("[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234", WithUDPv6Address()))
	a.NoError(err)
	err = Validate(NewRule("[::0]:1234", WithUDPv6Address()))
	a.NoError(err)
}

func TestTCPAddress(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("abcdefg", WithTCPAddress()))
	a.Error(err)
	err = Validate(NewRule("0", WithTCPAddress()))
	a.Error(err)
	err = Validate(NewRule("localhost", WithTCPAddress()))
	a.Error(err)
	err = Validate(NewRule("192.168.1.123", WithTCPAddress()))
	a.Error(err)
	err = Validate(NewRule("::0", WithTCPAddress()))
	a.Error(err)

	err = Validate(NewRule("192.168.1.123:1234", WithTCPAddress()))
	a.NoError(err)
	err = Validate(NewRule("[::0]:1234", WithTCPAddress()))
	a.NoError(err)
	err = Validate(NewRule("[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234", WithTCPAddress()))
	a.NoError(err)
}

func TestTCPv4Address(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("abcdefg", WithTCPv4Address()))
	a.Error(err)
	err = Validate(NewRule("[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234", WithTCPv4Address()))
	a.Error(err)
	err = Validate(NewRule("0", WithTCPv4Address()))
	a.Error(err)
	err = Validate(NewRule("192.168.1.123", WithTCPv4Address()))
	a.Error(err)

	err = Validate(NewRule("127.0.0.1:1234", WithTCPv4Address()))
	a.NoError(err)
}

func TestTCPv6Address(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("abcdefg", WithTCPv6Address()))
	a.Error(err)
	err = Validate(NewRule("192.168.1.123", WithTCPv6Address()))
	a.Error(err)
	err = Validate(NewRule("127.0.0.1", WithTCPv6Address()))
	a.Error(err)
	err = Validate(NewRule("0", WithTCPv6Address()))
	a.Error(err)
	err = Validate(NewRule("127.0.0.1:1234", WithTCPv6Address()))
	a.Error(err)

	err = Validate(NewRule("[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234", WithTCPv6Address()))
	a.NoError(err)
	err = Validate(NewRule("[::0]:1234", WithTCPv6Address()))
	a.NoError(err)
}

func TestLatitude(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("1234.92967312345678", WithLatitude()))
	a.Error(err)
	err = Validate(NewRule("35.", WithLatitude()))
	a.Error(err)
	err = Validate(NewRule("12345678", WithLatitude()))

	err = Validate(NewRule("35.929673", WithLatitude()))
	a.NoError(err)
	err = Validate(NewRule("35", WithLatitude()))
	a.NoError(err)
	err = Validate(NewRule("-78.948237", WithLatitude()))
	a.NoError(err)
}

func TestLongitude(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("1234.92967312345678", WithLongitude()))
	a.Error(err)
	err = Validate(NewRule("35.", WithLongitude()))
	a.Error(err)
	err = Validate(NewRule("12345678", WithLongitude()))

	err = Validate(NewRule("35.929673", WithLongitude()))
	a.NoError(err)
	err = Validate(NewRule("35", WithLongitude()))
	a.NoError(err)
	err = Validate(NewRule("-78.948237", WithLongitude()))
	a.NoError(err)
}

func TestCustomError(t *testing.T) {
	a := assert.New(t)
	err := Validate(NewRule("", WithRequired()))
	a.Equal(ErrRequired.Error(), err.Error())
	err = Validate(NewRule("", WithCustomError(WithRequired(), errors.New("hello"))))
	a.Equal("hello", err.Error())
}
