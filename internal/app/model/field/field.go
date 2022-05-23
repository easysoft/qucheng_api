package field

import (
	"github.com/go-playground/validator/v10"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service"
)

func New() map[string]validator.Func {
	return map[string]validator.Func{
		"namespace_exist": namespaceExist,
	}
}

func namespaceExist(fl validator.FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	// default reflect.String:
	return service.Namespaces(topField.String()).Has(field.String())
}
