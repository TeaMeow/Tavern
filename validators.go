package tavern

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrRequired
	ErrRequired = errors.New("tavern: missing required value")
	// ErrLength
	ErrLength = errors.New("tavern: out of the length")
	//
	ErrRange = errors.New("")
	//
	ErrDatetime = errors.New("")
	//
	ErrEmail = errors.New("")

	//
	ErrIP = errors.New("")
	//
	ErrURL = errors.New("")
	//
	ErrJSON = errors.New("")
)

var (
	// ErrWrongType
	ErrWrongType = errors.New("tavern: passed wrong type to validator")
)

const (
	KeyRequired = "KEY_REQUIRED"
)

// WithRequired 表示該內容值必須有內容而非零值（如：0、""）。
func WithRequired() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		value := reflect.ValueOf(v)
		if value.IsZero() {
			return ErrRequired, ctx
		}
		return nil, ctx
	}
}

// WithLength 會檢查切片、字串或正整數的長度。
func WithLength(min, max int) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		var err error
		err, ctx = WithMinLength(min)(v, ctx)
		if err != nil {
			return err, ctx
		}
		err, ctx = WithMaxLength(max)(v, ctx)
		if err != nil {
			return err, ctx
		}
		return nil, ctx
	}
}

// WithMaxLength 會檢查切片、字串、正整數的最大長度。
func WithMaxLength(max int) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		// IS NOT REQUIRD

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			if value.Len() > max {
				return ErrLength, ctx
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			str := strconv.Itoa(int(value.Int()))
			if len(str) > max {
				return ErrLength, ctx
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			str := strconv.Itoa(int(value.Uint()))
			if len(str) > max {
				return ErrLength, ctx
			}
		case reflect.Float32, reflect.Float64:
			str := fmt.Sprintf("%g", value.Float())
			if len(str) > max {
				return ErrLength, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

// WithMinLength 會檢查切片、字串、正整數的最小長度。
func WithMinLength(min int) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		// IS NOT REQUIRD

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			if value.Len() < min {
				return ErrLength, ctx
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			str := strconv.Itoa(int(value.Int()))
			if len(str) < min {
				return ErrLength, ctx
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			str := strconv.Itoa(int(value.Uint()))
			if len(str) < min {
				return ErrLength, ctx
			}
		case reflect.Float32, reflect.Float64:
			str := fmt.Sprintf("%g", value.Float())
			if len(str) < min {
				return ErrLength, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

// WithFixedLength 會要求切片、字串、正整數必須符合指定長度。
func WithFixedLength(length int) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		// IS NOT REQUIRD

		var err error
		err, ctx = WithMinLength(length)(v, ctx)
		if err != nil {
			return err, ctx
		}
		err, ctx = WithMaxLength(length)(v, ctx)
		if err != nil {
			return err, ctx
		}
		return nil, ctx
	}
}

// WithRange 會檢查正整數的數值是否在指定範圍內。
func WithRange(min, max int) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		// IS NOT REQUIRD

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() < int64(min) || value.Int() > int64(max) {
				return ErrRange, ctx
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value.Uint() < uint64(min) || value.Uint() > uint64(max) {
				return ErrRange, ctx
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() < float64(min) || value.Float() > float64(max) {
				return ErrRange, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

// WithMaxRange 會檢查正整數的數值是否小於某個範圍內。
func WithMaxRange(max int) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		// IS NOT REQUIRD

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() > int64(max) {
				return ErrRange, ctx
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value.Uint() > uint64(max) {
				return ErrRange, ctx
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() > float64(max) {
				return ErrRange, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

// WithMinRange 會檢查正整數的數值是否小於某個範圍內。
func WithMinRange(min int) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		// IS NOT REQUIRD

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() < int64(min) {
				return ErrRange, ctx
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value.Uint() < uint64(min) {
				return ErrRange, ctx
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() < float64(min) {
				return ErrRange, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

// WithMaximum 會要求切片、字串、正整數必須小於指定長度或範圍內。
func WithMaximum(max int) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		// IS NOT REQUIRD

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			if value.Len() > max {
				return ErrLength, ctx
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() > int64(max) {
				return ErrRange, ctx
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value.Uint() > uint64(max) {
				return ErrRange, ctx
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() > float64(max) {
				return ErrRange, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

// WithMinimum 會要求切片、字串、正整數必須符小於指定長度或範圍內。
func WithMinimum(min int) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		// IS NOT REQUIRD

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			if value.Len() < min {
				return ErrLength, ctx
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() < int64(min) {
				return ErrRange, ctx
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value.Uint() < uint64(min) {
				return ErrRange, ctx
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() < float64(min) {
				return ErrRange, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithDatetime(f string) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			t, err := time.Parse(f, k)
			if err != nil {
				return err, ctx
			}
			if t.Format(f) != k {
				return ErrDatetime, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithEmail() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpEmailRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithOneOf() {

}

//
func WithNotOneOf() {

}

//
func WithIP() {

}

//
func WithIPv4() {

}

//
func WithIPv6() {

}

//
func WithURL() {

}

//
func WithEqual() {

}

//
func WithNotEqual() {

}

//
func WithTrue() {

}

//
func WithFalse() {

}

//
func WithRegExp(r string) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			m, err := regexp.Match(r, []byte(k))
			if !m || err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithPrefix(p string) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !strings.HasPrefix(k, p) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithSuffix(s string) Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !strings.HasSuffix(k, s) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithAlpha() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpAlphaRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithAlphanumeric() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpAlphaNumericRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithAlphaUnicode() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpAlphaUnicodeRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithAlphanumericUnicode() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpAlphaUnicodeNumericRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithNumeric() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpNumericRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithHexadecimal() {

}

//
func WithHexColor() {

}

//
func WithLowercase() {

}

//
func WithUppercase() {

}

//
func WithRGB() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpRgbRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithRGBA() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpRgbaRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithHSL() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpHslRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithHSLA() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpHslaRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithJSON() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !json.Valid([]byte(k)) {
				return ErrEmail, ctx
			}
		case []byte:
			if !json.Valid(k) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithFilePath() {

}

//
func WithURI() {

}

//
func WithURNRFC2141() {

}

//
func WithBase64() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpBase64Regex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithBase64URL() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpBase64URLRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithBitcoinAddress() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpBtcAddressRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithBitcoinAddressBech32() {

}

//
func WithEthereumAddress() {

}

//
func WithContains() {

}

//
func WithNotContains() {

}

//
func WithISBN() {

}

//
func WithISBN10() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpISBN10Regex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithISBN13() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpISBN13Regex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithUUID() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpUUIDRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithUUID3() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpUUID3Regex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.String:
			if !regExpUUID3Regex.Match([]byte(value.String())) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithUUID4() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpUUID4Regex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithUUID5() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpUUID5Regex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithASCII() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpASCIIRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithASCIIPrintable() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpASCIIPrintableRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithMultiByte() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpMultibyteRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithDataURI() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpDataURIRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithLatitude() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpLatitudeRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithLongitude() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpLongitudeRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithTCPAddress() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			_, err := net.ResolveTCPAddr("tcp", k)
			if err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithTCPv4Address() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			_, err := net.ResolveTCPAddr("tcp4", k)
			if err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithTCPv6Address() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			_, err := net.ResolveTCPAddr("tcp6", k)
			if err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithUDPAddress() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			_, err := net.ResolveTCPAddr("udp", k)
			if err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithUDPv4Address() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			_, err := net.ResolveTCPAddr("udp4", k)
			if err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithUDPv6Address() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			_, err := net.ResolveTCPAddr("udp6", k)
			if err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithIPAddress() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			_, err := net.ResolveIPAddr("ip", k)
			if err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithIPv4Address() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			_, err := net.ResolveIPAddr("ip4", k)
			if err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithIPv6Address() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			_, err := net.ResolveIPAddr("ip6", k)
			if err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithUnixAddress() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			_, err := net.ResolveUnixAddr("unix", k)
			if err != nil {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithMAC() {

}

//
func WithHTML() Validator {
	return func(v interface{}, ctx context.Context) (error, context.Context) {
		switch k := v.(type) {
		case string:
			if !regExpHTMLRegex.Match([]byte(k)) {
				return ErrEmail, ctx
			}
		default:
			panic(ErrWrongType)
		}
		return nil, ctx
	}
}

//
func WithHostname() {

}
