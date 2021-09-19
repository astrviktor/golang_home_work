package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString - ошибка в формате строки для распаковки c повтором символов.
var ErrInvalidString = errors.New("invalid string")

// Unpack - Распаковка строки c повторами символов.
func Unpack(str string) (string, error) {
	// преобразование в слайс рун
	runes := []rune(str)

	idx := 0
	var builder strings.Builder

	for idx < len(runes) {
		// если тут число значит строка некорректная
		if unicode.IsDigit(runes[idx]) {
			return "", ErrInvalidString
		}

		// тут последнияя руна строки, и это не число (предыдущая проверка)
		if idx == len(runes)-1 {
			builder.WriteString(string(runes[idx]))
			idx++
			continue
		}

		// оталось 2 руны или больше
		if idx < len(runes)-1 {
			// если руна и число - строим строку
			if !unicode.IsDigit(runes[idx]) && unicode.IsDigit(runes[idx+1]) {
				i, _ := strconv.Atoi(string(runes[idx+1]))
				builder.WriteString(strings.Repeat(string(runes[idx]), i))
				idx += 2
				continue
			}

			// если руна и руна, добавляем перкую руну и идем дальше
			if !unicode.IsDigit(runes[idx]) && !unicode.IsDigit(runes[idx+1]) {
				builder.WriteString(string(runes[idx]))
				idx++
				continue
			}
		}
	}

	// финальный результат
	return builder.String(), nil
}
