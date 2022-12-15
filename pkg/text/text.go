package text

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/dchest/uniuri"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// RemoveNonASCII removes all non-ascii characters
func RemoveNonASCII(s string) string {
	return strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return -1
		}
		return r
	}, s)
}

type mnSet struct{}

func (s mnSet) Contains(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func ToASCII(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(mnSet{}))
	r, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}
	return r, nil
}

func SanitizeCVFileName(s string) string {
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return s
	}
	fileName := RemoveNonASCII(strings.Join(parts[0:len(parts)-1], " "))
	fileType := parts[len(parts)-1]
	return fmt.Sprintf("%s--%s.%s", fileName, uniuri.NewLen(5), fileType)
}
