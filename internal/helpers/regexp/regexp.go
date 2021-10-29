package regexp_helpers

import (
	"ts-tg-bot/internal/types"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func IsHSEEmail(text string) (bool, error) {
	pattern := `^\w+@edu.hse.ru`

	return regexp.Match(pattern, []byte(text))
}

func GetLectureNumber(text string) (int64, error) {
	pattern := `^лекция \d+`

	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0, err
	}

	log.Printf("TEXT: %+v", text)
	found := strings.ReplaceAll(re.FindString(strings.ToLower(text)), "лекция ", "")
	if found == "" {
		return 0, types.ErrGetLectrureNumber
	}

	return strconv.ParseInt(found, 10, 64)
}

func StripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}
