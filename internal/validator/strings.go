package validator

import (
	"fmt"
	"strings"
	"unicode"
)

func (v *Validator) FirstLetterUpper(s, key string) {
	firstRune := []rune(s)[0]
	v.Check(unicode.IsUpper(firstRune), key, "must start with an uppercase letter")
}

func (v *Validator) AllLettersLowercase(s, key string) {
	for _, r := range s {
		if unicode.IsUpper(r) {
			v.Check(false, key, "must only contain lowercase letters")
		}
	}
}

func (v *Validator) OnlyOneWord(s, key string) {
	trimmed := strings.TrimSpace(s)
	v.Check(!strings.Contains(trimmed, " "), key, "must only contain one word")
}

func (v *Validator) NoTrailingSpaces(s, key string) {
	trimmed := strings.TrimSpace(s)
	v.Check(trimmed == s, key, "must not have trailing spaces")
}

func (v *Validator) MaxLength(s, key string, length int) {
	v.Check(len(s) <= length, key, fmt.Sprintf("must be less than %d characters", length))
}

func (v *Validator) NotOnlySpaces(s, key string) {
	trimmed := strings.TrimSpace(s)
	v.Check(trimmed != "", key, "must not only contain spaces")
}
