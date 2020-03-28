package tavern

import "context"

type Tavern struct {
	values []value
}

func New() *Tavern {
	return &Tavern{}
}

type value struct {
	value      interface{}
	validators []Validator
	//context context.Context
}

// type Validator interface {
// 	Validate(value interface{}) (err error)
// }

type Validator func(value interface{}, ctx context.Context) (error, context.Context)

func (t *Tavern) Add(v interface{}, validators ...Validator) *Tavern {
	t.values = append(t.values, value{
		value:      v,
		validators: validators,
	})
	return t
}

func (t *Tavern) Validate() (err error) {
	for _, v := range t.values {
		ctx := context.Background()
		for _, j := range v.validators {
			err, ctx = j(v.value, ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
