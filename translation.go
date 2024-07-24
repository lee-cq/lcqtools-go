package main

import (
	"unicode"
	"unicode/utf8"
)

func HasEnglish(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && (unicode.Is(unicode.Latin, r)) {
			return true
		}
	}
	return false
}

func isChinese(r rune) bool {
	return r >= '\u4e00' && r <= '\u9fff'
}

func MostChinese(s string) bool {
	var totalRuneCount, chineseRuneCount int

	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			s = s[size:]
			continue
		}
		totalRuneCount++
		if isChinese(r) {
			chineseRuneCount++
		}
	}

	return float64(chineseRuneCount)/float64(totalRuneCount) > 0.5
}

func ToLang(s string) string {
	if MostChinese(s) {
		return "en"
	} else {
		return "zh"
	}
}
