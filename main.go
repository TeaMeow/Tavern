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

func Add(value interface{}) *Tavern {
	t := &Tavern{}
	return t.Add(value)
}

type E struct {
	Length, Range, Min, Max, Date, Email, In, IP, URL, Equal, RegExp, Required error
}

type Tavern struct {
	rules []rule
}

func (t *Tavern) lastRule() *rule {
	return &t.rules[len(t.rules)-1]
}

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

func (t *Tavern) Length(min, max int) *Tavern {
	r := t.lastRule()
	r.length = [2]int{min, max}
	r.hasLength = true
	return t
}

func (t *Tavern) Range(min, max int) *Tavern {
	r := t.lastRule()
	r.rangeList = [2]int{min, max}
	r.hasRange = true
	return t
}

func (t *Tavern) Min(min int) *Tavern {
	r := t.lastRule()
	r.min = min
	r.hasMin = true
	return t
}

func (t *Tavern) Max(max int) *Tavern {
	r := t.lastRule()
	r.max = max
	r.hasMax = true
	return t
}

func (t *Tavern) Date(formats ...string) *Tavern {
	r := t.lastRule()
	r.date = formats
	r.hasDate = true
	return t
}

func (t *Tavern) Email() *Tavern {
	r := t.lastRule()
	r.email = true
	return t
}

func (t *Tavern) In(values ...interface{}) *Tavern {
	r := t.lastRule()
	r.in = values
	r.hasIn = true
	return t
}

func (t *Tavern) IP(typ ...string) *Tavern {
	r := t.lastRule()
	r.ip = typ[0]
	r.hasIP = true
	return t
}

func (t *Tavern) URL(contains ...string) *Tavern {
	r := t.lastRule()
	r.url = contains
	r.hasURL = true
	return t
}

func (t *Tavern) Equal(compare interface{}) *Tavern {
	r := t.lastRule()
	r.equal = compare
	r.hasEqual = true
	return t
}

func (t *Tavern) RegExp(reg string) *Tavern {
	r := t.lastRule()
	r.regexp = reg
	r.hasRegExp = true
	return t
}

func (t *Tavern) Required() *Tavern {
	r := t.lastRule()
	r.required = true
	return t
}

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

func (t *Tavern) Check() error {
	for _, v := range t.rules {
		err := v.check()
		if err != nil {
			return err
		}
	}
	return nil
}

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
		return errors.New("The length of the value is.")
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
		return errors.New("The length of the value is.")
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
		return errors.New("The length of the value is.")
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
		return errors.New("The length of the value is.")
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
		return errors.New("Required value but it's empty.")
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
		return errors.New("Required value but it's empty.")
	}
	return nil
}

func (r *rule) checkURL() error {
	if !r.hasURL {
		return nil
	}

	var isErr bool
	_, err := url.ParseRequestURI(r.stringValue)
	if err != nil {
		isErr = true
	} else {
		for _, v := range r.url {
			if !strings.HasPrefix(r.stringValue, v) {
				isErr = true
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
		return errors.New("Required value but it's empty.")
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
		return errors.New("Required value but it's empty.")
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
		return errors.New("Required value but it's empty.")
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
		return errors.New("Required value but it's empty.")
	}
	return nil
}
