package request

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/stubborn-gaga-0805/prepare2go/api/ecode"
	"github.com/stubborn-gaga-0805/prepare2go/api/response"
	"go.uber.org/zap"
	"strings"
)

type Validator struct {
	v     *validator.Validate
	uni   *ut.UniversalTranslator
	trans ut.Translator
}

func V() *Validator {
	var (
		enTrans = en.New()
		zhTrans = zh.New()
		uni     = ut.New(enTrans, zhTrans)
	)
	trans, _ := uni.GetTranslator("zh")
	return &Validator{
		v:     validator.New(),
		uni:   uni,
		trans: trans,
	}
}

func (v *Validator) Check(ptr interface{}) (err error) {
	if err := zhtrans.RegisterDefaultTranslations(v.v, v.trans); err != nil {
		return err
	}
	if err = v.v.Struct(ptr); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			zap.S().Errorf("Validator Error. err: %v", err)
			return err
		}
	}
	return v.parseError(err)
}

func (v *Validator) parseError(vErr error) (err error) {
	var str = ""
	for k, val := range v.removeStructName(vErr.(validator.ValidationErrors).Translate(v.trans)) {
		str += fmt.Sprintf("[%s] %s\n", k, val)
	}
	return response.ThrowErr(ecode.ParseParamsError)
}

func (v *Validator) removeStructName(fields map[string]string) map[string]string {
	result := map[string]string{}
	for field, err := range fields {
		result[field[strings.Index(field, ".")+1:]] = err
	}
	return result
}
