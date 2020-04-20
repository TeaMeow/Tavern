package tavern

import "context"

// Tavern 呈現了一個可驗證內容值的結構。
type Tavern struct {
	values []value
}

// New 會初始化一個可追加檢查值的結構體。
func New() Tavern {
	return Tavern{}
}

// value 呈現了一個值與其驗證器。
type value struct {
	name       string
	value      interface{}
	validators []Validator
}

// Error 是一個 Tavern 的驗證錯誤資料。
/*type Error struct {
	// Name 是欄位的名稱。
	Name string
	// Value 是欄位值。
	Value interface{}
	// Err 是錯誤訊息。
	Err error
}


func (e Error) Error() string {
	return e.Err.Error()
}*/

// Validator 是一個能夠驗證內容的驗證器，並且會將上下文結構體往下一個驗證器傳遞。
type Validator func(ctx context.Context, value interface{}) (context.Context, error)

// Add 會增加一個欲檢查的值與其檢查的驗證器。
func (t Tavern) Add(v interface{}, validators ...Validator) Tavern {
	return t.AddNamed("", v, validators...)
}

// AddNamed 會增加一個欲檢查的指名值與其檢查的驗證器，這能夠在錯誤發生時得知是哪個欄位出錯。
func (t Tavern) AddNamed(name string, v interface{}, validators ...Validator) Tavern {
	t.values = append(t.values, value{
		name:       name,
		value:      v,
		validators: validators,
	})
	return t
}

// CustomError 會接收一個驗證器，當該驗證器發生錯誤時會回傳自訂的錯誤而非 Tavern 原生錯誤。適合用在替每個欄位自訂錯誤的場景上。
func CustomError(validator Validator, err error) Validator {
	return func(ctx context.Context, v interface{}) (context.Context, error) {
		ctx, originalErr := validator(ctx, v)
		if originalErr != nil {
			return ctx, err
		}
		return ctx, nil
	}
}

// Validate 會開始校驗所有欲檢查的內容值。
func (t Tavern) Validate() (err error) {
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
