package validate

type Func func(v interface{}) error

type Validator struct {
	validators []Func
}

func New(v []Func) *Validator {
	return &Validator{v}
}

func (v *Validator) Exec(m interface{}) error {
	for _, v := range v.validators {
		if err := v(m); err != nil {
			return err
		}
	}
	return nil
}
