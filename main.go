package tavern

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrLength 表示欄位的長度低於或高於指定預期。
	ErrLength = errors.New("tavern: 長度錯誤")
	// ErrRange 表示欄位的數字低於或高於指定範圍。
	ErrRange = errors.New("tavern: 超出範圍")
	// ErrMin 表示欄位的長度或範圍低於指定值。
	ErrMin = errors.New("tavern: 低於最小值")
	// ErrMax 表示欄位的長度或範圍高於指定值。
	ErrMax = errors.New("tavern: 高於最大值")
	// ErrDate 表示欄位的日期格式不正確。
	ErrDate = errors.New("tavern: 日期格式錯誤")
	// ErrEmail 表示欄位的內容不符合電子郵件地址格式。
	ErrEmail = errors.New("tavern: 電子郵件地址格式錯誤")
	// ErrIn 表示欄位內容不在預期清單中。
	ErrIn = errors.New("tavern: 內容不在預期清單內")
	// ErrIP 表示欄位並不符合指定的 IP 位置格式。
	ErrIP = errors.New("tavern: IP 位置格式錯誤")
	// ErrURL 表示欄位不符合指定的網址格式。
	ErrURL = errors.New("tavern: 網址格式錯誤")
	// ErrEqual 表示欄位不符合指定內容。
	ErrEqual = errors.New("tavern: 內容與預期不相符")
	// ErrRegEx 表示欄位無法通過正規表達式驗證。
	ErrRegEx = errors.New("tavern: 無法通過正規表達式驗證")
	// ErrRequired 表示欄位必填但卻缺少內容或僅有空白。
	ErrRequired = errors.New("tavern: 缺少必填內容")
)

// Add 會增加新的值供 Tavern 檢查。
func Add(value interface{}) *Tavern {
	t := &Tavern{}
	return t.Add(value)
}

// E 呈現了自訂的錯誤訊息。
type E struct {
	Length, Range, Min, Max, Date, Email, In, IP, URL, Equal, RegExp, Required error
}

// Tavern 是最主要的檢查結構體，能夠透過函式互動增加新的檢查機制。
type Tavern struct {
	rules []rule
}

// lastRule 會回傳最後一個增加的規則指針。
func (t *Tavern) lastRule() *rule {
	return &t.rules[len(t.rules)-1]
}

// Add 會增加一個新的檢查目標內容。
func (t *Tavern) Add(value interface{}) *Tavern {
	r := rule{}
	switch v := value.(type) {
	case int:
		r.typ = "int"
		r.intValue = v
		r.stringValue = strconv.Itoa(v)
	case int8:
		r.typ = "int"
		r.intValue = int(v)
		r.stringValue = strconv.Itoa(int(v))
	case int16:
		r.typ = "int"
		r.intValue = int(v)
		r.stringValue = strconv.Itoa(int(v))
	case int32:
		r.typ = "int"
		r.intValue = int(v)
		r.stringValue = strconv.Itoa(int(v))
	case int64:
		r.typ = "int"
		r.intValue = int(v)
		r.stringValue = strconv.Itoa(int(v))
	case uint:
		r.typ = "int"
		r.intValue = int(v)
		r.stringValue = strconv.Itoa(int(v))
	case uint8:
		r.typ = "int"
		r.intValue = int(v)
		r.stringValue = strconv.Itoa(int(v))
	case uint16:
		r.typ = "int"
		r.intValue = int(v)
		r.stringValue = strconv.Itoa(int(v))
	case uint32:
		r.typ = "int"
		r.intValue = int(v)
		r.stringValue = strconv.Itoa(int(v))
	case uint64:
		r.typ = "int"
		r.intValue = int(v)
		r.stringValue = strconv.Itoa(int(v))
	case string:
		r.typ = "string"
		r.intValue, _ = strconv.Atoi(v)
		r.stringValue = v
	}
	r.value = value
	t.rules = append(t.rules, r)
	return t
}

// Length 會檢查數字、字串的長度。
func (t *Tavern) Length(min, max int) *Tavern {
	r := t.lastRule()
	r.length = [2]int{min, max}
	r.hasLength = true
	return t
}

// Range 會檢查數字的範圍。
func (t *Tavern) Range(min, max int) *Tavern {
	r := t.lastRule()
	r.rangeList = [2]int{min, max}
	r.hasRange = true
	return t
}

// Min 會檢查字串的最小長度，若是套用到數字則是最小範圍。
func (t *Tavern) Min(min int) *Tavern {
	r := t.lastRule()
	r.min = min
	r.hasMin = true
	return t
}

// Max 會檢查字串的最大長度，若是套用到數字則是最大範圍。
func (t *Tavern) Max(max int) *Tavern {
	r := t.lastRule()
	r.max = max
	r.hasMax = true
	return t
}

// Date 會檢查字串的日期格式，可以傳入多個格式（格式為 Golang 日期）。
func (t *Tavern) Date(formats ...string) *Tavern {
	r := t.lastRule()
	r.date = formats
	r.hasDate = true
	return t
}

// Email 會透過較為簡單的正規表達式（RegExp）檢查字串是否為電子郵件地址形式。但有些神奇字串仍然可以繞過這個驗證，因此將驗證信件傳送至目標電子郵件地址來得到最佳保護是不二選擇。
func (t *Tavern) Email() *Tavern {
	r := t.lastRule()
	r.email = true
	return t
}

// In 會檢查指定內容是否在清單裡。
func (t *Tavern) In(values ...interface{}) *Tavern {
	r := t.lastRule()
	r.in = values
	r.hasIn = true
	return t
}

// IP 會檢查字串是否為 IPv4 或 IPv6 格式，傳遞參數指定 `v4` 或 `v6`，留白則是兩者其中一個都可以。
func (t *Tavern) IP(typ ...string) *Tavern {
	r := t.lastRule()
	if len(typ) == 1 {
		r.ip = typ[0]
	}
	r.hasIP = true
	return t
}

// URL 會檢查字串是否為合法的網址，傳遞參數用以確保網址是某個開頭（如：`https://`）。
func (t *Tavern) URL(contains ...string) *Tavern {
	r := t.lastRule()
	r.url = contains
	r.hasURL = true
	return t
}

// Equal 會檢查內容是否與期望兩者完全相符。
func (t *Tavern) Equal(compare interface{}) *Tavern {
	r := t.lastRule()
	r.equal = compare
	r.hasEqual = true
	return t
}

// RegExp 會透過自訂的正規表達式檢查內容是否通過。
func (t *Tavern) RegExp(reg string) *Tavern {
	r := t.lastRule()
	r.regexp = reg
	r.hasRegExp = true
	return t
}

// Required 會要求內容是必填的且禁止僅有空白。
func (t *Tavern) Required() *Tavern {
	r := t.lastRule()
	r.required = true
	return t
}

// Error 允許你傳入自訂的錯誤，當該欄位驗證錯誤發生時則會回傳此錯誤訊息。
// 你也能夠傳入 `tavern.E` 結構體來自訂一個更為詳細的錯誤訊息。
func (t *Tavern) Error(err interface{}) *Tavern {
	r := t.lastRule()
	switch v := err.(type) {
	case E:
		r.errMessages = v
	case error:
		r.err = v
	}
	return t
}

// Check 會檢查所有的規則並在遇到錯誤時立即停下並回傳。
func (t *Tavern) Check() error {
	for _, v := range t.rules {
		err := v.check()
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckAll 會檢查所有的規則，這會比起 `Check` 還要更慢。因為檢查遇到錯誤時並不會停下而會檢查下一個規則。當你希望收集所有錯誤訊息的時候就能用上這個函式。
func (t *Tavern) CheckAll() []error {
	var errs []error
	for _, v := range t.rules {
		err := v.check()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// rule 定義了一個欄位的檢查規則與其內容。
type rule struct {
	value interface{}
	typ   string

	stringValue string
	intValue    int

	required bool
	email    bool

	url    []string
	hasURL bool

	ip    string
	hasIP bool

	in    []interface{}
	hasIn bool

	date    []string
	hasDate bool

	rangeList [2]int
	hasRange  bool

	length    [2]int
	hasLength bool

	min    int
	hasMin bool

	max    int
	hasMax bool

	equal    interface{}
	hasEqual bool

	regexp    string
	hasRegExp bool

	errMessages E
	err         error
}

type number struct {
	value int
	has   bool
}

func (r *rule) check() error {
	var err error

	err = r.checkRequired()
	if err != nil {
		return err
	}
	err = r.checkLength()
	if err != nil {
		return err
	}
	err = r.checkRange()
	if err != nil {
		return err
	}
	err = r.checkMin()
	if err != nil {
		return err
	}
	err = r.checkMax()
	if err != nil {
		return err
	}
	err = r.checkDate()
	if err != nil {
		return err
	}
	err = r.checkEmail()
	if err != nil {
		return err
	}
	err = r.checkIn()
	if err != nil {
		return err
	}
	err = r.checkIP()
	if err != nil {
		return err
	}
	err = r.checkURL()
	if err != nil {
		return err
	}
	err = r.checkEqual()
	if err != nil {
		return err
	}
	err = r.checkRegExp()
	if err != nil {
		return err
	}
	return nil
}

func (r *rule) checkLength() error {
	if !r.hasLength {
		return nil
	}

	var isErr bool
	switch r.typ {
	case "int", "string":
		if len(r.stringValue) < r.length[0] || len(r.stringValue) > r.length[1] {
			isErr = true
		}
	}
	if isErr {
		if r.errMessages.Length != nil {
			return r.errMessages.Length
		}
		if r.err != nil {
			return r.err
		}
		return ErrLength
	}
	return nil
}

func (r *rule) checkRange() error {
	if !r.hasRange {
		return nil
	}

	var isErr bool
	if r.typ == "int" {
		if r.intValue < r.rangeList[0] || r.intValue > r.rangeList[1] {
			isErr = true
		}
	}
	if isErr {
		if r.errMessages.Range != nil {
			return r.errMessages.Range
		}
		if r.err != nil {
			return r.err
		}
		return ErrRange
	}
	return nil
}

func (r *rule) checkMin() error {
	if !r.hasMin {
		return nil
	}

	var isErr bool
	switch r.typ {
	case "int":
		if r.intValue < r.min {
			isErr = true
		}
	case "string":
		if len(r.stringValue) < r.min {
			isErr = true
		}
	}
	if isErr {
		if r.errMessages.Min != nil {
			return r.errMessages.Min
		}
		if r.err != nil {
			return r.err
		}
		return ErrMin
	}
	return nil
}

func (r *rule) checkMax() error {
	if !r.hasMax {
		return nil
	}

	var isErr bool
	switch r.typ {
	case "int":
		if r.intValue > r.max {
			isErr = true
		}
	case "string":
		if len(r.stringValue) > r.max {
			isErr = true
		}
	}
	if isErr {
		if r.errMessages.Max != nil {
			return r.errMessages.Max
		}
		if r.err != nil {
			return r.err
		}
		return ErrMax
	}
	return nil
}

func (r *rule) checkDate() error {
	if !r.hasDate {
		return nil
	}

	var isErr bool
	isErr = true
	for _, v := range r.date {
		t, err := time.Parse(v, r.stringValue)
		if err == nil && t.Format(v) == r.stringValue {
			isErr = false
			break
		}
	}
	if isErr {
		if r.errMessages.Date != nil {
			return r.errMessages.Date
		}
		if r.err != nil {
			return r.err
		}
		return ErrDate
	}
	return nil
}

func (r *rule) checkEmail() error {
	if !r.email || (!r.required && r.stringValue == "") {
		return nil
	}

	var isErr bool
	m, err := regexp.Match(`\S+@\S+`, []byte(r.stringValue))
	if !m || err != nil {
		isErr = true
	}
	if isErr {
		if r.errMessages.Email != nil {
			return r.errMessages.Email
		}
		if r.err != nil {
			return r.err
		}
		return ErrEmail
	}
	return nil
}

func (r *rule) checkIn() error {
	if !r.hasIn {
		return nil
	}
	var isErr bool
	isErr = true
	for _, v := range r.in {

		if r.value == v {
			isErr = false
			break
		}
	}
	if isErr {
		if r.errMessages.In != nil {
			return r.errMessages.In
		}
		if r.err != nil {
			return r.err
		}
		return ErrIn
	}
	return nil
}

func (r *rule) checkIP() error {
	if !r.hasIP {
		return nil
	}

	var isErr bool
	switch r.ip {
	case "":
		ip := net.ParseIP(r.stringValue)
		if ip == nil {
			isErr = true
		}
	case "v4":
		ip := net.ParseIP(r.stringValue)
		if ip == nil || !strings.Contains(r.stringValue, ".") {
			isErr = true
		}
	case "v6":
		ip := net.ParseIP(r.stringValue)
		if ip == nil || !strings.Contains(r.stringValue, ":") {
			isErr = true
		}
	}
	if isErr {
		if r.errMessages.IP != nil {
			return r.errMessages.IP
		}
		if r.err != nil {
			return r.err
		}
		return ErrIP
	}
	return nil
}

func (r *rule) checkURL() error {
	if !r.hasURL {
		return nil
	}

	var isErr bool
	isErr = true
	_, err := url.ParseRequestURI(r.stringValue)
	if err == nil {
		isErr = false
	}
	if len(r.url) != 0 {
		isErr = true
		for _, v := range r.url {
			if strings.HasPrefix(r.stringValue, v) {
				isErr = false
				break
			}
		}
	}
	if isErr {
		if r.errMessages.URL != nil {
			return r.errMessages.URL
		}
		if r.err != nil {
			return r.err
		}
		return ErrURL
	}
	return nil
}

func (r *rule) checkEqual() error {
	if !r.hasEqual {
		return nil
	}

	var isErr bool
	if r.value != r.equal {
		isErr = true
	}
	if isErr {
		if r.errMessages.Equal != nil {
			return r.errMessages.Equal
		}
		if r.err != nil {
			return r.err
		}
		return ErrEqual
	}
	return nil
}

func (r *rule) checkRegExp() error {
	if !r.hasRegExp {
		return nil
	}

	var isErr bool
	m, err := regexp.Match(r.regexp, []byte(r.stringValue))
	if !m || err != nil {
		isErr = true
	}
	if isErr {
		if r.errMessages.RegExp != nil {
			return r.errMessages.RegExp
		}
		if r.err != nil {
			return r.err
		}
		return ErrRegEx
	}
	return nil
}

func (r *rule) checkRequired() error {
	if r.required == false {
		return nil
	}

	var isErr bool
	if r.value == nil {
		isErr = true
	}
	if r.typ == "string" {
		if r.stringValue == "" {
			isErr = true
		}
		if strings.Trim(r.stringValue, " ") == "" {
			isErr = true
		}
	}
	if isErr {
		if r.errMessages.Required != nil {
			return r.errMessages.Required
		}
		if r.err != nil {
			return r.err
		}
		return ErrRequired
	}
	return nil
}
