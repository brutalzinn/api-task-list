package converter_util

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// thahnks to https://gist.github.com/micheltlutz/8398f801abf3ac61ff002dd5be7caeb4
func Normalize(oldText string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	newText, _, err := transform.String(t, oldText)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(newText))
	return strings.ToLower(newText)
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.IsSpace(r)
}
