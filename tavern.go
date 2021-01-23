package tavern

import "context"

// Rule has the validators to validate the value. The name is optional, a named rule is useful if you wanted to know which value is invalid.
type Rule struct {
	// Name of the value, optional.
	Name string
	// Value is the content will be validated by the validators.
	Value interface{}
	// Validators to validate the value.
	Validators []Validator
}

// Validator is used to validate the value, it accepts a context to pass between the validators.
type Validator func(ctx context.Context, value interface{}) (context.Context, error)

// Validate validates all the rules that passed in.
func Validate(rules []Rule) (err error) {
	for _, v := range rules {
		ctx := context.Background()
		for _, j := range v.Validators {
			ctx, err = j(ctx, v.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
