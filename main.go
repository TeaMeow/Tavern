package tavern

import "context"

// Tavern 呈現了一個可驗證內容值的結構。
type Tavern struct {
	values []value
}

// New 會初始化一個可追加檢查值的結構體。
func New() *Tavern {
	return &Tavern{}
}

// value 呈現了一個值與其驗證器。
type value struct {
	value      interface{}
	validators []Validator
}

// Validator 是一個能夠驗證內容的驗證器，並且會將上下文結構體往下一個驗證器傳遞。
type Validator func(ctx context.Context, value interface{}) (context.Context, error)

// Add 會增加一個欲檢查的值，與其檢查的驗證器。
func (t *Tavern) Add(v interface{}, validators ...Validator) *Tavern {
	t.values = append(t.values, value{
		value:      v,
		validators: validators,
	})
	return t
}

// Validate 會開始校驗所有欲檢查的內容值。
func (t *Tavern) Validate() (err error) {
	for _, v := range t.values {
		ctx := context.Background()
		for _, j := range v.validators {
			ctx, err = j(ctx, v.value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
