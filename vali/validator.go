package vali

import (
	"fmt"
	localen "github.com/go-playground/locales/en"
	localzh "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	transen "github.com/go-playground/validator/v10/translations/en"
	transzh "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

const defaultLocale = "en"

var defaultTranslateFunc = func(uniTrans ut.Translator, fe validator.FieldError) string {
	t, err := uniTrans.T(fe.Tag(), fe.Field())
	if err != nil {
		return fmt.Errorf("%w:field:%s,tag:%s", err, fe.Field(), fe.Tag()).Error()
	}
	return t
}

type Validate struct {
	origin         *validator.Validate
	locale         string
	translator     ut.Translator
	translateField string
	translateFunc  validator.TranslationFunc
}

func New(options ...Option) *Validate {
	v := &Validate{
		origin:        validator.New(),
		locale:        defaultLocale,
		translateFunc: defaultTranslateFunc,
	}
	for _, opt := range options {
		opt(v)
	}
	if v.locale != defaultLocale && v.translateField == "" {
		v.translateField = v.locale
	}
	v.init()
	return v
}

func (v *Validate) init() {
	en := localen.New()
	uni := ut.New(en, localzh.New(), en)
	translator, ok := uni.GetTranslator(v.locale)
	if !ok {
		panic(fmt.Sprintf("universalTranslator.GetTranslator(\"%s\") failed", v.locale))
	}
	v.translator = translator

	//翻译字段
	if v.locale != defaultLocale {
		v.origin.RegisterTagNameFunc(func(field reflect.StructField) string {
			value := field.Tag.Get(v.translateField)
			if value == "" {
				value = field.Name
			}
			return value
		})
	}

	//注册为默认翻译器
	var err error
	switch v.locale {
	case "zh":
		err = transzh.RegisterDefaultTranslations(v.origin, translator)
	default:
		err = transen.RegisterDefaultTranslations(v.origin, translator)
	}
	if err != nil {
		panic(err)
	}
}

func (v *Validate) Translator() ut.Translator {
	return v.translator
}

func (v *Validate) Validate() *validator.Validate {
	return v.origin
}

func (v *Validate) RegisterValidationTrans(tag string, fn validator.Func, registerFn validator.RegisterTranslationsFunc) (err error) {
	err = v.origin.RegisterValidation(tag, fn)
	if err != nil {
		return err
	}
	err = v.origin.RegisterTranslation(tag, v.translator, registerFn, v.translateFunc)
	return err
}
