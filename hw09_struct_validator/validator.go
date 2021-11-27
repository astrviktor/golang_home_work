package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrProgramError             = errors.New("program error")
	ErrNotStruct                = errors.New("interface is not struct")
	ErrInvalidRegexpConst       = errors.New("invalid regexp const")
	ErrInvalidValidateTag       = errors.New("invalid validate tag")
	ErrInvalidValidateIntTag    = errors.New("invalid validate int tag")
	ErrInvalidValidateStringTag = errors.New("invalid validate string tag")

	ErrValidationError = errors.New("validation error")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	result := ""
	for _, value := range v {
		result += value.Field + ": " + value.Err.Error() + "\n"
	}
	return result
}

func Validate(v interface{}) error {
	// проверка что это структура
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("%q: %w", ErrNotStruct, ErrProgramError)
	}

	var valErrs ValidationErrors

	// смотрим поля структуры
	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := rv.Field(i)

		// если нет тэга validate - продолжаем
		if !existTagValidate(string(field.Tag)) {
			continue
		}

		// поверка тэга validate на валидность
		err := checkTagValidate(string(field.Tag))
		if err != nil {
			return err
		}

		// в тэге или ошибки, или получаем отдельные условия валидации
		exps, err := getValidateExpressions(field.Type.String(), string(field.Tag))
		if err != nil {
			return err
		}

		// проверка по условиям
		ok, errs := checkFieldValidate(field, value, exps)
		if !ok {
			valErrs = append(valErrs, errs...)
		}
	}

	if len(valErrs) > 0 {
		return fmt.Errorf("%q: %w", valErrs, ErrValidationError)
	}

	return nil
}

func checkFieldValidate(field reflect.StructField, value reflect.Value, exps []string) (bool, ValidationErrors) {
	switch field.Type.String() {
	case "int":
		if varInt, ok := value.Interface().(int); ok {
			if ok, errs := ValidateInt(field.Name, varInt, exps); !ok {
				return false, errs
			}
		}
	case "string":
		if varString, ok := value.Interface().(string); ok {
			if ok, errs := ValidateString(field.Name, varString, exps); !ok {
				return false, errs
			}
		}
	case "[]int":
		if sliceInt, ok := value.Interface().([]int); ok {
			if ok, errs := ValidateIntSlice(field.Name, sliceInt, exps); !ok {
				return false, errs
			}
		}
	case "[]string":
		if sliceString, ok := value.Interface().([]string); ok {
			if ok, errs := ValidateStringSlice(field.Name, sliceString, exps); !ok {
				return false, errs
			}
		}
	}

	return true, ValidationErrors{}
}

func existTagValidate(tag string) bool {
	// ищем tag validate
	return strings.Contains(tag, "validate:")
}

func checkTagValidate(tag string) error {
	// общая проверка тэга c validate
	r, err := regexp.Compile(regexpValidate)
	if err != nil {
		return ErrInvalidRegexpConst
	}
	if !r.MatchString(tag) {
		return ErrInvalidValidateTag
	}

	return nil
}

// https://regex101.com/
// для интов.
const regexpInt = "^min:[-]{0,1}[0-9]{1,}$|^max:[-]{0,1}[0-9]{1,}$|^in:[-]{0,1}[0-9]{1,}(,[-]{0,1}[0-9]{1,}){0,}$"

// для строк.
const regexpString = "^len:[0-9]{1,}$|^regexp:.{1,}$|^in:.{0,}(,.{0,}){0,}$"

// для validate.
const regexpValidate = "^.{1,}[ ]{1}validate:\".{1,}\"([ ]{1}.{1,}){0,}$|^validate:\".{1,}\"$"

// checkTagValidate получает тип поля и тэг, проверяет тэг на ошибки, возвращает ошибку или
// отдельные условия валидации.
func getValidateExpressions(fieldType string, tag string) ([]string, error) {
	// подготовка валидаторов
	ri, err := regexp.Compile(regexpInt)
	if err != nil {
		return []string{}, ErrInvalidRegexpConst
	}
	rs, err := regexp.Compile(regexpString)
	if err != nil {
		return []string{}, ErrInvalidRegexpConst
	}

	var expressions []string
	// разбиваем строку тэга на отдельные тэги - разделены пробелом. Ищем все validate
	words := strings.Split(tag, " ")
	for _, word := range words {
		if strings.HasPrefix(word, "validate:") && len(word) >= 11 {
			// убираем в строке validate:" и " в конце строки
			word = word[10 : len(word)-1]

			// разбиваем на отдельные выражения
			exps := strings.Split(word, "|")
			for _, exp := range exps {
				switch {
				case fieldType == "int" || fieldType == "[]int":
					if !ri.MatchString(exp) {
						return []string{}, ErrInvalidValidateIntTag
					}
					expressions = append(expressions, exp)
				case fieldType == "string" || fieldType == "[]string":
					if !rs.MatchString(exp) {
						return []string{}, ErrInvalidValidateStringTag
					}
					expressions = append(expressions, exp)
				}
			}
		}
	}

	return expressions, nil
}

func ValidateIntSlice(field string, value []int, exps []string) (bool, ValidationErrors) {
	valFlag := true
	var valErrs ValidationErrors

	for _, v := range value {
		ok, errs := ValidateInt(field, v, exps)
		if !ok {
			valFlag = false
			valErrs = append(valErrs, errs...)
		}
	}

	return valFlag, valErrs
}

func ValidateStringSlice(field string, value []string, exps []string) (bool, ValidationErrors) {
	valFlag := true
	var valErrs ValidationErrors

	for _, v := range value {
		ok, errs := ValidateString(field, v, exps)
		if !ok {
			valFlag = false
			valErrs = append(valErrs, errs...)
		}
	}

	return valFlag, valErrs
}

func ValidateString(field string, value string, exps []string) (bool, ValidationErrors) {
	valFlag := true
	var valErrs ValidationErrors

	for _, exp := range exps {
		if strings.HasPrefix(exp, "len:") {
			ok, err := ValidateStringByLen(field, value, strings.TrimPrefix(exp, "len:"))
			if !ok {
				valFlag = false
				valErrs = append(valErrs, err)
			}
			continue
		}

		if strings.HasPrefix(exp, "in:") {
			ok, err := ValidateStringByIn(field, value, strings.TrimPrefix(exp, "in:"))
			if !ok {
				valFlag = false
				valErrs = append(valErrs, err)
			}
			continue
		}

		if strings.HasPrefix(exp, "regexp:") {
			ok, err := ValidateStringByRegexp(field, value, strings.TrimPrefix(exp, "regexp:"))
			if !ok {
				valFlag = false
				valErrs = append(valErrs, err)
			}
			continue
		}
	}

	return valFlag, valErrs
}

func ValidateInt(field string, value int, exps []string) (bool, ValidationErrors) {
	valFlag := true
	var valErrs ValidationErrors

	for _, exp := range exps {
		if strings.HasPrefix(exp, "min:") {
			ok, err := ValidateIntByMin(field, value, strings.TrimPrefix(exp, "min:"))
			if !ok {
				valFlag = false
				valErrs = append(valErrs, err)
			}
			continue
		}

		if strings.HasPrefix(exp, "max:") {
			ok, err := ValidateIntByMax(field, value, strings.TrimPrefix(exp, "max:"))
			if !ok {
				valFlag = false
				valErrs = append(valErrs, err)
			}
			continue
		}

		if strings.HasPrefix(exp, "in:") {
			ok, err := ValidateIntByIn(field, value, strings.TrimPrefix(exp, "in:"))
			if !ok {
				valFlag = false
				valErrs = append(valErrs, err)
			}
			continue
		}
	}

	return valFlag, valErrs
}

func ValidateIntByMin(field string, value int, exp string) (bool, ValidationError) {
	min, err := strconv.Atoi(exp)
	if err != nil {
		return false, ValidationError{field, fmt.Errorf("exp %s not convert to int", exp)}
	}
	if value < min {
		return false, ValidationError{field, fmt.Errorf("value %d less then min=%d", value, min)}
	}
	return true, ValidationError{}
}

func ValidateIntByMax(field string, value int, exp string) (bool, ValidationError) {
	max, err := strconv.Atoi(exp)
	if err != nil {
		return false, ValidationError{field, fmt.Errorf("exp %s not convert to int", exp)}
	}
	if value > max {
		return false, ValidationError{field, fmt.Errorf("value %d more then max=%d", value, max)}
	}
	return true, ValidationError{}
}

func ValidateIntByIn(field string, value int, exp string) (bool, ValidationError) {
	ins := strings.Split(exp, ",")
	for _, in := range ins {
		intIn, err := strconv.Atoi(in)
		if err != nil {
			return false, ValidationError{field, fmt.Errorf("exp %s not convert to int", in)}
		}

		if value == intIn {
			return true, ValidationError{}
		}
	}

	return false, ValidationError{field, fmt.Errorf("value %d not in %s", value, exp)}
}

func ValidateStringByLen(field string, value string, exp string) (bool, ValidationError) {
	length, err := strconv.Atoi(exp)
	if err != nil {
		return false, ValidationError{field, fmt.Errorf("exp %s not convert to int", exp)}
	}
	if len(value) != length {
		return false, ValidationError{field, fmt.Errorf("value %s has len=%d, expected len=%d", value, len(value), length)}
	}
	return true, ValidationError{}
}

func ValidateStringByRegexp(field string, value string, exp string) (bool, ValidationError) {
	r, err := regexp.Compile(exp)
	if err != nil {
		return false, ValidationError{field, fmt.Errorf("regex %s not compile", exp)}
	}
	if !r.MatchString(value) {
		return false, ValidationError{field, fmt.Errorf("value %s not match regexp %s", value, exp)}
	}
	return true, ValidationError{}
}

func ValidateStringByIn(field string, value string, exp string) (bool, ValidationError) {
	ins := strings.Split(exp, ",")
	for _, in := range ins {
		if value == in {
			return true, ValidationError{}
		}
	}
	return false, ValidationError{field, fmt.Errorf("value %s not in %s", value, exp)}
}
