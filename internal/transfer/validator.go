package transfer

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var currencyRegexp = regexp.MustCompile(`^[A-Z]{3}$`)

// RegisterValidations wires custom validators required by the transfer payloads.
func RegisterValidations(v *validator.Validate) {
	_ = v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		s, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return currencyRegexp.MatchString(s)
	})

	_ = v.RegisterValidation("fee_lt_amount", func(fl validator.FieldLevel) bool {
		fee, ok := fl.Field().Interface().(float64)
		if !ok {
			return false
		}
		amount := fl.Parent().FieldByName("SourceAmount").Float()
		return fee < amount
	})
}
