package vali

import "github.com/go-playground/validator/v10"

type Option func(v *Validate)

func Origin(origin *validator.Validate) Option {
	return func(v *Validate) {
		v.origin = origin
	}
}

func Locale(locale string) Option {
	return func(v *Validate) {
		v.locale = locale
	}
}

func TranslateField(translateField string) Option {
	return func(v *Validate) {
		v.translateField = translateField
	}
}

func TranslateFunc(translateFunc validator.TranslationFunc) Option {
	return func(v *Validate) {
		v.translateFunc = translateFunc
	}
}
