package readability

import (
	_ "embed"
	"log"
	"strconv"
	"strings"
	"unicode"
)

//go:embed wordlist/syllables.txt
var syllablesData string

// irregularSyllables maps words with known syllable counts (matching CMU Dict).
var irregularSyllables map[string]int

func init() {
	lines := strings.Split(strings.TrimSpace(syllablesData), "\n")
	irregularSyllables = make(map[string]int, len(lines))
	for i, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			log.Fatalf("readability: malformed syllables.txt at line %d", i+1)
		}
		count, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatalf("readability: invalid count in syllables.txt at line %d: %v", i+1, err)
		}
		irregularSyllables[parts[0]] = count
	}
}

// isVowel checks if a byte is a vowel (including y).
// Used for vowel-group counting where [aeiouy]+ is the pattern.
func isVowel(b byte) bool {
	switch b {
	case 'a', 'e', 'i', 'o', 'u', 'y':
		return true
	}
	return false
}

// isVowelNoY checks if a byte is a vowel (excluding y).
// Used for silent-ending checks where [^aeiou] is the pattern.
func isVowelNoY(b byte) bool {
	switch b {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

// stripNonAlpha removes all non-lowercase-alpha characters from word.
func stripNonAlpha(word string) string {
	var b strings.Builder
	b.Grow(len(word))
	for i := 0; i < len(word); i++ {
		c := word[i]
		if c >= 'a' && c <= 'z' {
			b.WriteByte(c)
		}
	}
	return b.String()
}

// countVowelGroups counts contiguous vowel groups (aeiouy) in word.
func countVowelGroups(word string) int {
	count := 0
	inVowel := false
	for i := 0; i < len(word); i++ {
		if isVowel(word[i]) {
			if !inVowel {
				count++
				inVowel = true
			}
		} else {
			inVowel = false
		}
	}
	return count
}

// hasSilentE checks for [^aeiou]e$ pattern.
func hasSilentE(word string) bool {
	n := len(word)
	return n >= 2 && word[n-1] == 'e' && !isVowelNoY(word[n-2])
}

// hasSilentES checks for [^aeioulsxzcg]es$ pattern.
func hasSilentES(word string) bool {
	n := len(word)
	if n < 3 || word[n-1] != 's' || word[n-2] != 'e' {
		return false
	}
	c := word[n-3]
	switch c {
	case 'a', 'e', 'i', 'o', 'u', 'l', 's', 'x', 'z', 'c', 'g':
		return false
	}
	return true
}

// hasSilentED checks for [^aeiou]ed$ pattern.
func hasSilentED(word string) bool {
	n := len(word)
	return n >= 3 && word[n-1] == 'd' && word[n-2] == 'e' && !isVowelNoY(word[n-3])
}

// countUANotL counts occurrences of ua[^l] in word.
func countUANotL(word string) int {
	count := 0
	n := len(word)
	for i := 0; i+2 < n; i++ {
		if word[i] == 'u' && word[i+1] == 'a' && word[i+2] != 'l' {
			count++
		}
	}
	return count
}

// countEONotU counts occurrences of eo[^u] in word.
func countEONotU(word string) int {
	count := 0
	n := len(word)
	for i := 0; i+2 < n; i++ {
		if word[i] == 'e' && word[i+1] == 'o' && word[i+2] != 'u' {
			count++
		}
	}
	return count
}

// hasDigit checks if a string contains any digit.
func hasDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// countSyllablesWord counts syllables in a single word using rule-based heuristics.
func countSyllablesWord(word string) int {
	word = strings.ToLower(strings.TrimSpace(word))

	if word == "" {
		return 0
	}

	// Numbers count as 1 syllable.
	if hasDigit(word) {
		return 1
	}

	// Check irregulars before stripping non-alpha (to handle contractions).
	if count, ok := irregularSyllables[word]; ok {
		return count
	}

	word = stripNonAlpha(word)

	if word == "" {
		return 0
	}
	if len(word) <= 2 {
		return 1
	}

	// Check irregulars again after stripping (e.g., word with trailing punctuation).
	if count, ok := irregularSyllables[word]; ok {
		return count
	}

	// Count vowel groups as base syllable count.
	count := countVowelGroups(word)

	// Subtract for silent endings.
	if hasSilentE(word) && !strings.HasSuffix(word, "le") {
		count--
	}
	if hasSilentES(word) {
		count--
	}
	if hasSilentED(word) && word != "bed" && word != "fed" && word != "red" && word != "led" && word != "wed" && word != "shed" {
		count--
	}

	// Add back for hiatus patterns (two vowels pronounced separately).
	iaCount := strings.Count(word, "ia")
	ioCount := strings.Count(word, "io")

	// Subtract tion/sion matches from IO count (these are not hiatus).
	ioCount -= strings.Count(word, "tion")
	ioCount -= strings.Count(word, "sion")
	if ioCount < 0 {
		ioCount = 0
	}

	// Subtract iage matches from IA count (these are not hiatus).
	iaCount -= strings.Count(word, "iage")
	if iaCount < 0 {
		iaCount = 0
	}

	count += iaCount
	count += ioCount

	if strings.HasSuffix(word, "ian") && !strings.Contains(word, "tion") {
		count-- // already counted by ia pattern
	}
	if strings.HasSuffix(word, "ium") {
		count-- // already counted by io/ia pattern
	}
	count += strings.Count(word, "iou")
	if strings.HasSuffix(word, "ial") {
		count-- // already counted by ia
	}
	count += strings.Count(word, "ien")
	if strings.HasSuffix(word, "ual") {
		count++
	}
	count += countUANotL(word)
	if strings.Contains(word, "iet") {
		count++
	}
	count += countEONotU(word)
	if strings.HasSuffix(word, "ism") {
		count++
	}

	if count < 1 {
		count = 1
	}
	return count
}

// countSyllables counts total syllables across all words in text.
func countSyllables(text string) int {
	words := listWords(text)
	total := 0
	for _, w := range words {
		total += countSyllablesWord(w)
	}
	return total
}
