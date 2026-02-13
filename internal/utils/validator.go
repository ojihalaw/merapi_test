package utils

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/spf13/viper"
)

type Validator struct {
	Validate   *validator.Validate
	Translator ut.Translator
}

func NewValidator(viper *viper.Viper) *Validator {
	validate := validator.New()

	locale := en.New()
	uni := ut.New(locale, locale)
	trans, _ := uni.GetTranslator("en")

	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	return &Validator{
		Validate:   validate,
		Translator: trans,
	}
}

func (v *Validator) TranslateError(err error) []string {
	var errs []string
	if err == nil {
		return errs
	}

	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, e.Translate(v.Translator))
	}
	return errs
}
