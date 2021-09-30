package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

// wordFrequency - структура для сортировки slice по частоте и лексикографически.
type wordFrequency struct {
	word      string
	frequency int
}

// Top10 функция принимает на вход строку с текстом и
// возвращает слайс с 10-ю наиболее часто встречаемыми в тексте словами.
func Top10(text string) []string {
	// в тексте могут быть символы типа "?", "." и тд, перед разделением текста на слова их нужно заменить на пробелы
	// нельзя сделать такую константу, поэтому сразу в функцию
	replacer := strings.NewReplacer("...", " ", "\"-", " ", ",", " ", "!", " ", "?", " ", "\n", " ", "\"", " ",
		";", " ", "- ", " - ", ":", " ", ";", " ")
	text = replacer.Replace(text)

	// разделяем текст на слова
	splits := strings.Fields(text)

	// в этой мапе будут хранится слова и частота использования
	words := make(map[string]int)

	// заполняем мапу словами и частотой
	for _, word := range splits {
		words[word]++
	}

	// чтобы отсортировать, надо переложить в slice структур
	wordsFrequency := make([]wordFrequency, 0, len(words))
	for word, frequency := range words {
		wordsFrequency = append(wordsFrequency, wordFrequency{word, frequency})
	}

	// сортируем slice структур по частоте, если частота равна то лексикографически. но все в обратном порядке
	sort.Slice(wordsFrequency, func(i, j int) bool {
		if wordsFrequency[i].frequency == wordsFrequency[j].frequency {
			return wordsFrequency[i].word < wordsFrequency[j].word
		}
		return wordsFrequency[i].frequency > wordsFrequency[j].frequency
	})

	// в результат нужно вернуть 10, если слов меньше, то берем меньше
	maxLen := 10
	if len(wordsFrequency) < 10 {
		maxLen = len(wordsFrequency)
	}

	// формируем результат
	result := make([]string, 0, maxLen)
	for idx := 0; idx < maxLen; idx++ {
		result = append(result, wordsFrequency[idx].word)
	}

	return result
}
