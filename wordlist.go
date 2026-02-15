package readability

import (
	_ "embed"
	"strings"
)

//go:embed wordlist/dale_chall.txt
var daleChallData string

var easyWords map[string]struct{}

func init() {
	lines := strings.Split(strings.TrimSpace(daleChallData), "\n")
	easyWords = make(map[string]struct{}, len(lines))
	for _, line := range lines {
		word := strings.TrimSpace(line)
		if word != "" {
			easyWords[word] = struct{}{}
		}
	}
}

func isEasyWord(word string) bool {
	_, ok := easyWords[strings.ToLower(word)]
	return ok
}
