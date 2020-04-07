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
	// ErrRequired 表示缺少了必填欄位。
	ErrRequired = errors.New("tavern: missing required value")
	// ErrLength 表示大小超出指定長度。
	ErrLength = errors.New("tavern: out of length")
	// ErrRange 表示大小超出指定範圍。
	ErrRange = errors.New("tavern: out of range")
	// ErrDatetime 表示不符合指定格式的日期字串。
	ErrDatetime = errors.New("tavern: invalid datetime format")
	// ErrEmail 表示錯誤的電子郵件地址格式。
	ErrEmail = errors.New("tavern: invalid email format")

	// ErrAddress 表示無法解析的位置。
	ErrAddress = errors.New("tavern: unresolvable address")
	// ErrURL 表示錯誤的 URL 格式。
	ErrURL = errors.New("tavern: invalid url format")
	// ErrJSON 表示錯誤的 JSON 格式。
	ErrJSON = errors.New("tavern: invalid json format")
)

var (
	// ErrWrongType 表示建立驗證規則時，傳入錯誤的格式到驗證器。
	ErrWrongType = errors.New("tavern: passed wrong type to validator")
)

const (
	// KeyRequired 是必填的鍵值上下文資料。
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

// WithDatetime 會檢查字串內容是否符合指定的日期格式。
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

// WithEmail 會檢查字串是否符合 Email 格式。
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

// WithRegExp 會驗證指定字串是否通過 RegExp 正規表達式。
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

// WithPrefix 會檢查字串是否開頭帶有指定字元。
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

// WithSuffix 會檢查字串結尾是否以特定字元結束。
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

// WithAlpha 會檢查字串是否為基本大小寫英文字母。
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

// WithAlphanumeric 會檢查字串是否為大小寫英文字母與數字。
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

// WithAlphaUnicode 會檢查字串是否為標準的 Unicode 語系文字。
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

// WithAlphanumericUnicode 會檢查字串是否為標準的 Unicode 語系文字與數字。
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

// WithNumeric 會檢查字串是否為數字或帶有小數點的格式。
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

// WithRGB 會檢查字串是否為 `rgb(0,0,0)` 格式。
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

// WithRGBA 會檢查字串是否為 `rgba(0,0,0,0)` 格式。
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

// WithHSL 會檢查字串是否為 `hsl(0,0,0)` 格式。
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

// WithHSLA 會檢查字串是否為 `hsla(0,0,0,0)` 格式。
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

// WithJSON 會驗證指定字串是否為正規的 JSON 格式。
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

// WithBase64 會檢查字串是否為 Base64 格式。
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

// WithBase64URL 會檢查字串是否為帶有 Base64 資料的網址格式。
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

// WithBitcoinAddress 會檢查字串是否為比特幣地址。
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

// WithISBN10 會檢查字串是否為 ISBN10 格式。
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

// WithISBN13 會檢查字串是否為 ISBN13 格式。
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

// WithUUID 會檢查字串是否為 UUID 格式。
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

// WithUUID3 會檢查字串是否為 UUID3 格式。
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

// WithUUID4 會檢查字串是否為 UUID4 格式。
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

// WithUUID5 會檢查字串是否為 UUID5 格式。
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

// WithASCII 會檢查字串是否為 ASCII 字元。
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

// WithASCIIPrintable 會檢查字串是否為 ASCII 可列印字元。
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

// WithMultiByte 會檢查字串是否為雙重位元組字元（如：符號、中日文）。
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

// WithDataURI 會檢查字串是否為 DataURI 格式。
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

// WithLatitude 會檢查傳入的字串格式是否為座標緯度。
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

// WithLongitude 會檢查傳入的字串格式是否為座標經度。
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

// WithTCPAddress 會驗證 TCP 地址是否可供解析。
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

// WithTCPv4Address 會驗證 TCPv4 地址是否可供解析。
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

// WithTCPv6Address 會驗證 TCPv6 地址是否可供解析。
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

// WithUDPAddress 會驗證 UDP 地址是否可供解析。
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

// WithUDPv4Address 會驗證 UDPv4 地址是否可供解析。
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

// WithUDPv6Address 會驗證 UDPv6 地址是否可供解析。
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

// WithIPAddress 會驗證一個 IP 地址是否可供解析。
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

// WithIPv4Address 會驗證一個 IPv4 地址是否可供解析。
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

// WithIPv6Address 會驗證一個 IPv6 地址是否可供解析。
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

// WithUnixAddress 會驗證一個 Unix 地址是否可供解析。
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

// WithMAC 會驗證一個字串是否為正規的 MAC 地址。
func WithMAC() {

}

// WithHTML 會驗證字串是否為正規的 HTML 格式。
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

// WithHostname 會驗證指定的主機名稱是否可供解析。
func WithHostname() {

}
