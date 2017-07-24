package tavern

import "strconv"

func Add() *Tavern {
	return &Tavern{}
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
	return t
}

func (t *Tavern) Range(min, max int) *Tavern {
	r := t.lastRule()
	r.rangeList = [2]int{min, max}
	return t
}

func (t *Tavern) Min(min int) *Tavern {
	r := t.lastRule()
	r.min = number{
		value: min,
		has:   true,
	}
	return t
}

func (t *Tavern) Max(max int) *Tavern {
	r := t.lastRule()
	r.max = number{
		value: max,
		has:   true,
	}
	return t
}

func (t *Tavern) Date(formats ...string) *Tavern {
	r := t.lastRule()
	r.date = formats
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
	return t
}

func (t *Tavern) IP(typ ...string) *Tavern {
	r := t.lastRule()
	r.ip = typ
	return t
}

func (t *Tavern) URL(contains ...string) *Tavern {
	r := t.lastRule()
	r.url = contains
	return t
}

func (t *Tavern) Equal(compare interface{}) *Tavern {
	r := t.lastRule()
	r.equal = compare
	return t
}

func (t *Tavern) RegExp(reg string) *Tavern {
	r := t.lastRule()
	r.regexp = reg
	return t
}

func (t *Tavern) Required() *Tavern {
	r := t.lastRule()
	r.required = true
	return t
}

func (t *Tavern) Check() error {

}

func (t *Tavern) CheckAll() []error {

}

type rule struct {
	value       interface{}
	stringValue string
	intValue    int
	required    bool
	email       bool
	url         []string
	ip          []string
	in          []interface{}
	date        []string
	rangeList   [2]int
	length      [2]int
	min         number
	max         number
	typ         string
	equal       interface{}
	regexp      string
}

type number struct {
	value int
	has   bool
}
