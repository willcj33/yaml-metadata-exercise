package services

import (
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"

	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

const dnsName = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
const ip = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
const urlSchema = `((ftp|tcp|udp|wss?|https?):\/\/)`
const urlUserName = `(\S+(:\S*)?@)`
const urlPath = `((\/|\?|#)[^\s]*)`
const urlPort = `(:(\d{1,5}))`
const urlIP = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
const urlSubdomain = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`

var urlExp = `^` + urlSchema + `?` + urlUserName + `?` + `((` + urlIP + `|(\[` + ip + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + urlSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + urlPort + `?` + urlPath + `?$`

//Validator contains storage utilities for the metadata
//TODO: potentially move this to a real DB in the future
type Validator struct {
	Service    *validator.Validate
	translator ut.Translator
}

var validatorInstance *Validator
var validatorOnce sync.Once
var urlRegex *regexp.Regexp

//GetValidator gets the singleton validator instance
func GetValidator() *Validator {
	validatorOnce.Do(func() {
		en := en.New()
		uni := ut.New(en, en)
		trans, _ := uni.GetTranslator("en")
		validatorInstance = &Validator{
			Service:    validator.New(),
			translator: trans,
		}
		urlRegex = regexp.MustCompile(urlExp)
		validatorInstance.Service.RegisterValidation("isURL", ValidateURL)
		en_translations.RegisterDefaultTranslations(validatorInstance.Service, trans)
	})
	return validatorInstance
}

func ValidateURL(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return urlRegex.MatchString(str)
}

func (v *Validator) Translate(err error) validator.ValidationErrorsTranslations {
	validationErrors := err.(validator.ValidationErrors)
	return validationErrors.Translate(v.translator)
}
