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
	// ErrRequired is missing the required value.
	ErrRequired = errors.New("tavern: missing required value")
	// ErrLength is out of the length.
	ErrLength = errors.New("tavern: out of length")
	// ErrRange is out of the range.
	ErrRange = errors.New("tavern: out of range")
	// ErrDatetime is invalid datetime format.
	ErrDatetime = errors.New("tavern: invalid datetime format")
	// ErrEmail is invalid email format.
	ErrEmail = errors.New("tavern: invalid email format")
	// ErrInvalidJSON is invalid JSON.
	ErrInvalidJSON = errors.New("tavern: invalid JSON")
	// ErrInvalidHTML is invalid HTML.
	ErrInvalidHTML = errors.New("tavern: invalid HTML")
	// ErrInvalidPattern is invalid pattern.
	ErrInvalidPattern = errors.New("invalid pattern")
	// ErrAddress is unresolvable address.
	ErrAddress = errors.New("tavern: unresolvable address")
	// ErrURL is invalid url format.
	ErrURL = errors.New("tavern: invalid url format")
	// ErrJSON is invalid json format.
	ErrJSON = errors.New("tavern: invalid json format")
)

var (
	// ErrWrongType is passed a wrong value type to validator.
	ErrWrongType = errors.New("tavern: passed wrong value type to validator")
)

// Key represents the keys in the context.
type Key int

const (
	// KeyRequired returns true if the value was set to required. The validators can rely on the value with it's own logic.
	KeyRequired Key = iota
)

// isNotRequiredAndZeroValue 表示這個欄位是不是非必要而且還零值。
func isNotRequiredAndZeroValue(ctx context.Context, v interface{}) bool {
	_, ok := ctx.Value(KeyRequired).(bool)
	return !ok && reflect.ValueOf(v).IsZero()
}

// WithRequired requires the value to not be a zero value (e.g. 0, "") nor an empty value.
func WithRequired() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		ctx = context.WithValue(ctx, KeyRequired, true)
		value := reflect.ValueOf(v)
		if value.IsZero() {
			return ctx, ErrRequired
		}
		return ctx, nil
	}
}

// WithLength requires the length of the value (e.g. slive, string, number) to be in a certian length. It counts the length of the number if the value was a number.
func WithLength(min, max int) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		var err error
		ctx, err = WithMinLength(min)(ctx, v)
		if err != nil {
			return ctx, err
		}
		ctx, err = WithMaxLength(max)(ctx, v)
		if err != nil {
			return ctx, err
		}
		return ctx, nil
	}
}

// WithMaxLength requires the length of the value (e.g. slive, string, number) cannot be too long. It counts the length of the number if the value was a number.
func WithMaxLength(max int) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			if value.Len() > max {
				return ctx, ErrLength
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			str := strconv.Itoa(int(value.Int()))
			if len(str) > max {
				return ctx, ErrLength
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			str := strconv.Itoa(int(value.Uint()))
			if len(str) > max {
				return ctx, ErrLength
			}
		case reflect.Float32, reflect.Float64:
			str := fmt.Sprintf("%g", value.Float())
			if len(str) > max {
				return ctx, ErrLength
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithMinLength requires the length of the value (e.g. slive, string, number) cannot be too short. It counts the length of the number if the value was a number.
func WithMinLength(min int) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			if value.Len() < min {
				return ctx, ErrLength
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			str := strconv.Itoa(int(value.Int()))
			if len(str) < min {
				return ctx, ErrLength
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			str := strconv.Itoa(int(value.Uint()))
			if len(str) < min {
				return ctx, ErrLength
			}
		case reflect.Float32, reflect.Float64:
			str := fmt.Sprintf("%g", value.Float())
			if len(str) < min {
				return ctx, ErrLength
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithFixedLength requires the length of the value (e.g. slive, string, number) to be the exact length. It counts the length of the number if the value was a number.
func WithFixedLength(length int) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		var err error
		ctx, err = WithMinLength(length)(ctx, v)
		if err != nil {
			return ctx, err
		}
		ctx, err = WithMaxLength(length)(ctx, v)
		if err != nil {
			return ctx, err
		}
		return ctx, nil
	}
}

// WithRange requires the number of the value to be in a certian range.
func WithRange(min, max int) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() < int64(min) || value.Int() > int64(max) {
				return ctx, ErrRange
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value.Uint() < uint64(min) || value.Uint() > uint64(max) {
				return ctx, ErrRange
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() < float64(min) || value.Float() > float64(max) {
				return ctx, ErrRange
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithMaxRange requires the number of the value to be equal or less than the specified number.
func WithMaxRange(max int) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() > int64(max) {
				return ctx, ErrRange
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value.Uint() > uint64(max) {
				return ctx, ErrRange
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() > float64(max) {
				return ctx, ErrRange
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithMinRange requires the number of the value to be least equal or greater than the specified number.
func WithMinRange(min int) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() < int64(min) {
				return ctx, ErrRange
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value.Uint() < uint64(min) {
				return ctx, ErrRange
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() < float64(min) {
				return ctx, ErrRange
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithMaximum requires the length of the slice, string and the range of the number to be least equal or less than the specified number.
func WithMaximum(max int) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			if value.Len() > max {
				return ctx, ErrLength
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() > int64(max) {
				return ctx, ErrRange
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value.Uint() > uint64(max) {
				return ctx, ErrRange
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() > float64(max) {
				return ctx, ErrRange
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithMinimum requires the length of the slice, string and the range of the number to be least equal or greater than the specified number.
func WithMinimum(min int) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			if value.Len() < min {
				return ctx, ErrLength
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() < int64(min) {
				return ctx, ErrRange
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value.Uint() < uint64(min) {
				return ctx, ErrRange
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() < float64(min) {
				return ctx, ErrRange
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithDatetime requires the date format to match the Golang date format. It validates via the `time.Parse` function.
func WithDatetime(f string) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			t, err := time.Parse(f, k)
			if err != nil {
				return ctx, err
			}
			if t.Format(f) != k {
				return ctx, ErrDatetime
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithEmail requires the value to be an email. It validates with a built-in email regexp pattern.
func WithEmail() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpEmailRegex.Match([]byte(k)) {
				return ctx, ErrEmail
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
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

// WithRegExp validates the valiue with specified regular expression.
func WithRegExp(r string) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			m, err := regexp.Match(r, []byte(k))
			if !m || err != nil {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithPrefix requires the value started with a specified sentence.
func WithPrefix(p string) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !strings.HasPrefix(k, p) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithSuffix requires the value ended with a specified sentence.
func WithSuffix(s string) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !strings.HasSuffix(k, s) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithAlpha requires the value to be alphabets only.
func WithAlpha() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpAlphaRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithAlphanumeric requires the value to be alphanumerics only.
func WithAlphanumeric() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpAlphaNumericRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithAlphaUnicode requires the value to be standard unicode characters.
func WithAlphaUnicode() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpAlphaUnicodeRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithAlphanumericUnicode requires the value to be numerics or standard unicode characters.
func WithAlphanumericUnicode() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpAlphaUnicodeNumericRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithNumeric requires the value to be numerics (includes the floating point).
func WithNumeric() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpNumericRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithHexadecimal 會檢查字串是否為十六進制格式。
func WithHexadecimal() {

}

// WithHexColor 會檢查字串是否為 # 井字開頭與結尾 3 或 6 個長度的十六進制格式。
func WithHexColor() {

}

// WithLowercase 會檢查字串是否僅有小寫英文字母。
func WithLowercase() {

}

// WithUppercase 會檢查字串是否僅有大寫英文字母。
func WithUppercase() {

}

// WithRGB requires the value to be a string RGB with `rgb(0,0,0)` format.
func WithRGB() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpRgbRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithRGBA requires the value to be a string RGBA with `rgba(0,0,0,0)` format.
func WithRGBA() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpRgbaRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithHSL requires the value to be a string HSL with `hsl(0,0,0)` format.
func WithHSL() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpHslRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithHSLA requires the value to be a string HSLA with `hsla(0,0,0,0)` format.
func WithHSLA() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpHslaRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithJSON requires the value to be a valid JSON. It validates via the `json.Valid` function.
func WithJSON() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !json.Valid([]byte(k)) {
				return ctx, ErrInvalidJSON
			}
		case []byte:
			if !json.Valid(k) {
				return ctx, ErrInvalidJSON
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
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

// WithBase64 requires the value to be a base64 string.
func WithBase64() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpBase64Regex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithBase64URL requires the value to be a URL base64.
func WithBase64URL() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpBase64URLRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithBitcoinAddress requires the value to be a Bitcoin address.
func WithBitcoinAddress() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpBtcAddressRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
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

// WithISBN10 requires the value to be a valid ISBN10 string.
func WithISBN10() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpISBN10Regex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithISBN13 requires the value to be a valid ISBN13 string.
func WithISBN13() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpISBN13Regex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithUUID requires the value to be a valid UUID string.
func WithUUID() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpUUIDRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithUUID3 requires the value to be a valid UUID3 string.
func WithUUID3() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpUUID3Regex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.String:
			if !regExpUUID3Regex.Match([]byte(value.String())) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithUUID4 requires the value to be a valid UUID4 string.
func WithUUID4() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpUUID4Regex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithUUID5 requires the value to be a valid UUID5 string.
func WithUUID5() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpUUID5Regex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithASCII requires the value to be a valid ASCII characters.
func WithASCII() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpASCIIRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithASCIIPrintable requires the value to be a valid printable ASCII characters.
func WithASCIIPrintable() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpASCIIPrintableRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithMultiByte requires the value to be multi-byte (e.g. Japanese, Chinese, Symbols).
func WithMultiByte() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpMultibyteRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithDataURI requires the value to be a data URI string.
func WithDataURI() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpDataURIRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithLatitude requires the value to be a valid latitude format.
func WithLatitude() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpLatitudeRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithLongitude requires the value to be a valid longitude format.
func WithLongitude() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpLongitudeRegex.Match([]byte(k)) {
				return ctx, ErrInvalidPattern
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithTCPAddress requires the value TCP address to be resolvable. It validates via the `net.ResolveTCPAddr` function.
func WithTCPAddress() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			_, err := net.ResolveTCPAddr("tcp", k)
			if err != nil {
				return ctx, ErrAddress
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithTCPv4Address requires the value TCPv4 address to be resolvable. It validates via the `net.ResolveTCPAddr` function.
func WithTCPv4Address() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			_, err := net.ResolveTCPAddr("tcp4", k)
			if err != nil {
				return ctx, ErrAddress
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithTCPv6Address requires the value TCPv6 address to be resolvable. It validates via the `net.ResolveTCPAddr` function.
func WithTCPv6Address() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			_, err := net.ResolveTCPAddr("tcp6", k)
			if err != nil {
				return ctx, ErrAddress
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithUDPAddress requires the value UDP address to be resolvable. It validates via the `net.ResolveTCPAddr` function.
func WithUDPAddress() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			_, err := net.ResolveUDPAddr("udp", k)
			if err != nil {
				return ctx, ErrAddress
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithUDPv4Address requires the value UDPv4 address to be resolvable. It validates via the `net.ResolveTCPAddr` function.
func WithUDPv4Address() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			_, err := net.ResolveUDPAddr("udp4", k)
			if err != nil {
				return ctx, ErrAddress
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithUDPv6Address requires the value UDPv6 address to be resolvable. It validates via the `net.ResolveTCPAddr` function.
func WithUDPv6Address() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			_, err := net.ResolveUDPAddr("udp6", k)
			if err != nil {
				return ctx, ErrAddress
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithIPAddress requires the value IP address to be resolvable. It validates via the `net.ResolveIPAddr` function.
func WithIPAddress() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			_, err := net.ResolveIPAddr("ip", k)
			if err != nil {
				return ctx, ErrAddress
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithIPv4Address requires the value IPv4 address to be resolvable. It validates via the `net.ResolveIPAddr` function.
//
// FIX: `::0` is IPv6 but resolvable, WTF GOLANG?
func WithIPv4Address() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			_, err := net.ResolveIPAddr("ip4", k)
			if err != nil {
				return ctx, ErrAddress
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithIPv6Address requires the value IPv6 address to be resolvable. It validates via the `net.ResolveIPAddr` function.
func WithIPv6Address() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			_, err := net.ResolveIPAddr("ip6", k)
			if err != nil {
				return ctx, ErrAddress
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithUnixAddress requires the value Unix address to be resolvable. It validates via the `net.ResolveUnixAddr` function.
func WithUnixAddress() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			_, err := net.ResolveUnixAddr("unix", k)
			if err != nil {
				return ctx, ErrAddress
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithMAC 會驗證一個字串是否為正規的 MAC 地址。
func WithMAC() {

}

// WithHTML requires the value to be a valid HTML.
func WithHTML() Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		if isNotRequiredAndZeroValue(ctx, v) {
			return ctx, nil
		}

		switch k := v.(type) {
		case string:
			if !regExpHTMLRegex.Match([]byte(k)) {
				return ctx, ErrInvalidHTML
			}
		default:
			panic(ErrWrongType)
		}
		return ctx, nil
	}
}

// WithHostname 會驗證指定的主機名稱是否可供解析。
func WithHostname() {

}

// WithCustomError accepts a validator with a custom error. It returns the custom error instead of the native Tavern error when the validator didn't pass it's validation. Useful if you are trying to create custom errors for each validation.
func WithCustomError(validator Validator, err error) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		ctx, originalErr := validator(ctx, v)
		if originalErr != nil {
			return ctx, err
		}
		return ctx, nil
	}
}
