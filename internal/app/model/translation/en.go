package translation

import (
	"fmt"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"k8s.io/klog/v2"
)

func RegisterCustomENTranslations(v *validator.Validate, trans ut.Translator) (err error) {
	translations := []struct {
		tag             string
		translation     string
		override        bool
		customRegisFunc validator.RegisterTranslationsFunc
		customTransFunc validator.TranslationFunc
	}{
		{
			tag: "namespace_exist",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("namespace_exist_string", "{0} {1} not exist", false); err != nil {
					return
				}
				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				t, err = ut.T("namespace_exist_string", strings.ToLower(fe.Field()), fe.Value().(string))
				if err != nil {
					goto END
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
	}

	for _, t := range translations {

		if t.customTransFunc != nil && t.customRegisFunc != nil {

			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)

		} else if t.customTransFunc != nil && t.customRegisFunc == nil {

			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)

		} else if t.customTransFunc == nil && t.customRegisFunc != nil {

			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFunc)

		} else {
			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		}

		if err != nil {
			return
		}
	}
	return
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {

	return func(ut ut.Translator) (err error) {

		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return

	}

}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {

	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		klog.Errorf("warning: error translating FieldError: %#v", fe)
		return fe.(error).Error()
	}

	return t
}
