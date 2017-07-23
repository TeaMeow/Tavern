package tavern

import "errors"

func Add(value interface{}) *Tavern {

}

type E struct {
	Min      error
	Max      error
	Length   error
	Range    error
	Email    error
	URL      error
	IP       error
	Date     error
	In       error
	RegExp   error
	Equal    error
	Required error
}

type Tavern struct {
	rules []rule
}

type rule struct {
	value       interface{}
	intValue    int
	stringValue string
	errors      E
	typ         string
	min         int
	max         int
	length      [2]int
	numberRange [2]int
	required    bool
	isEmail     bool
	urlFormat   []string
	ipFormat    string
	in          []interface{}
	dateFormat  []string
	regExp      []string
	equal       interface{}
}

func (t *Tavern) latestRule() *rule {
	return &t.rules[len(t.rules)-1]
}

func (r *rule) checkLength() error {
	if len(r.length) == 0 {
		return nil
	}
	if len(r.stringValue) < r.length[0] || len(r.stringValue) > r.length[0] {
		if r.errors.Length != nil {
			return r.errors.Length
		}
		return errors.New("The length is incorrect.")
	}
	return nil
}

func (r *rule) checkRange() error {
	if r.typ != "int" {
		return nil
	}
	if r.intValue < r.numberRange[0] || r.intValue > r.numberRange[1] {
		if r.errors.Range != nil {
			return r.errors.Range
		}
		return errors.New("The number doesn't in the range.")
	}
	return nil
}

func (r *rule) checkMinMax() error {
	if r.min == 0 {
		return nil
	}
}

func (r *rule) assert() {
	switch r.value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		r.typ = "int"
	case string:
		r.typ = "string"
	}
	panic("The value must be a int or a string.")
}

func (t *Tavern) Add(value interface{}) *Tavern {
	t.rules = append(t.rules, rule{})
	return t
}

func (t *Tavern) Min(min int) *Tavern {
	r := t.latestRule()
	r.min = min
	return t
}

func (t *Tavern) Max(max int) *Tavern {
	r := t.latestRule()
	r.max = max
	return t
}

func (t *Tavern) Length(min, max int) *Tavern {
	r := t.latestRule()
	r.length = [2]int{min, max}
	return t
}

func (t *Tavern) Range(min, max int) *Tavern {
	r := t.latestRule()
	r.numberRange = [2]int{min, max}
	return t
}

func (t *Tavern) Email() *Tavern {
	r := t.latestRule()
	r.isEmail = true
	return t
}

func (t *Tavern) URL(must ...[]string) *Tavern {
	r := t.latestRule()
	if len(must) == 1 {
		r.urlFormat = must[0]
	}
	return t
}

func (t *Tavern) IP(format ...string) *Tavern {
	r := t.latestRule()
	if len(format) == 1 {
		if format[0] != "v4" && format[0] != "v6" {
			panic("The IP format should be `v4` or either `v6`.")
		}
		r.ipFormat = format[0]
	}
	return t
}

func (t *Tavern) Date(format ...string) *Tavern {
	r := t.latestRule()
	r.dateFormat = format
	return t
}

func (t *Tavern) In(values ...interface{}) *Tavern {
	r := t.latestRule()
	r.in = values
	return t
}

func (t *Tavern) RegExp(regexps ...string) *Tavern {
	r := t.latestRule()
	r.regExp = regexps
	return t
}

func (t *Tavern) Equal(expect interface{}) *Tavern {
	r := t.latestRule()
	t.equal = compare
	return t
}

func (t *Tavern) Required() *Tavern {
	r := t.latestRule()
	r.required = true
	return t
}

func (t *Tavern) Error(e E) *Tavern {
	r := t.latestRule()
	r.errors = e
	return t
}

func (t *Tavern) Check() error {
	for _, v := range t.rules {
		v.assert()

	}
}

func (t *Tavern) CheckAll() []error {

}
