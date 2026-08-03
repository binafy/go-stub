package stub

import (
	"strings"
	"unicode"
)

// splitWords breaks s into lowercase word tokens. It splits on any character
// that is not a letter or digit (spaces, "_", "-", ".", etc.) and on camelCase
// boundaries, including acronym-to-word transitions such as "HTTPServer" ->
// ["http", "server"].
func splitWords(s string) []string {
	var words []string
	var cur []rune

	flush := func() {
		if len(cur) > 0 {
			words = append(words, strings.ToLower(string(cur)))
			cur = cur[:0]
		}
	}

	runes := []rune(s)
	for i, r := range runes {
		switch {
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			flush()
			continue
		case i > 0 && unicode.IsUpper(r) && unicode.IsLower(runes[i-1]):
			// lower -> Upper boundary: "userName" -> "user" | "Name"
			flush()
		case i > 0 && i+1 < len(runes) && unicode.IsUpper(r) &&
			unicode.IsUpper(runes[i-1]) && unicode.IsLower(runes[i+1]):
			// acronym -> word boundary: "HTTPServer" -> "HTTP" | "Server"
			flush()
		}
		cur = append(cur, r)
	}
	flush()

	return words
}

// title upper-cases the first rune of a lowercase word and leaves the rest.
func title(word string) string {
	if word == "" {
		return ""
	}
	r := []rune(word)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}

// ToSnake converts s to snake_case (e.g. "UserName" -> "user_name").
func ToSnake(s string) string {
	return strings.Join(splitWords(s), "_")
}

// ToScreamingSnake converts s to SCREAMING_SNAKE_CASE (e.g. "UserName" -> "USER_NAME").
func ToScreamingSnake(s string) string {
	return strings.ToUpper(ToSnake(s))
}

// ToKebab converts s to kebab-case (e.g. "UserName" -> "user-name").
func ToKebab(s string) string {
	return strings.Join(splitWords(s), "-")
}

// ToPascal converts s to PascalCase (e.g. "user_name" -> "UserName").
func ToPascal(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = title(w)
	}

	return strings.Join(words, "")
}

// ToCamel converts s to camelCase (e.g. "user_name" -> "userName").
func ToCamel(s string) string {
	words := splitWords(s)
	for i, w := range words {
		if i == 0 {
			words[i] = w
			continue
		}
		words[i] = title(w)
	}

	return strings.Join(words, "")
}
