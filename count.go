package readability

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	// Matches sentence-ending punctuation boundaries.
	reSentence = regexp.MustCompile(`\b[^.!?]+[.!?]*`)
	// Non-contraction apostrophes: apostrophes NOT between a letter and s/t/d/m/re/ll/ve.
	// Simple approach: remove leading/trailing single quotes from text chunks.
	reNonContractionApostrophe = regexp.MustCompile(`(?:^'|'$|'\s|\s')`)
	// Matches all non-word, non-space, non-apostrophe chars.
	// In Go regexp \w = [0-9A-Za-z_], matching Python's behavior.
	rePunct = regexp.MustCompile(`[^\w\s']`)
)

// listWords extracts words from text matching textstat's list_words(rm_punctuation=True).
// Removes non-contraction apostrophes, removes all punctuation except apostrophes
// (hyphens are removed, joining hyphenated words), splits on whitespace.
func listWords(text string) []string {
	// Remove non-contraction apostrophes (leading/trailing quotes).
	cleaned := reNonContractionApostrophe.ReplaceAllString(text, " ")
	// Remove all punctuation except apostrophes (hyphens, periods, &, etc. are removed).
	cleaned = rePunct.ReplaceAllString(cleaned, "")
	if strings.TrimSpace(cleaned) == "" {
		return nil
	}
	words := strings.Fields(cleaned)
	result := make([]string, 0, len(words))
	for _, w := range words {
		w = strings.Trim(w, "'")
		if w != "" {
			result = append(result, w)
		}
	}
	return result
}

// countWords returns the number of words in the text (punctuation removed).
func countWords(text string) int {
	return len(listWords(text))
}

// countRawWords returns the number of whitespace-separated tokens (punctuation kept).
// Used by charsPerWord for ARI formula.
func countRawWords(text string) int {
	return len(strings.Fields(text))
}

// countChars returns the character count excluding spaces.
func countChars(text string) int {
	count := 0
	for _, r := range text {
		if !unicode.IsSpace(r) {
			count++
		}
	}
	return count
}

// countLetters returns the count of word characters (letters, digits, underscores),
// matching textstat's count_letters which removes spaces then punctuation.
func countLetters(text string) int {
	// Remove spaces first.
	var sb strings.Builder
	for _, r := range text {
		if !unicode.IsSpace(r) {
			sb.WriteRune(r)
		}
	}
	// Remove all punctuation (including apostrophes): keep only \w chars.
	noSpace := sb.String()
	result := rePunct.ReplaceAllString(noSpace, "")
	// Also remove apostrophes.
	result = strings.ReplaceAll(result, "'", "")
	return len(result)
}

// countSentences splits text into sentences and returns the count.
// Fragments with <=2 words are ignored. Returns at least 1 for non-empty text.
func countSentences(text string) int {
	if strings.TrimSpace(text) == "" {
		return 0
	}
	matches := reSentence.FindAllString(text, -1)
	count := 0
	for _, m := range matches {
		words := strings.Fields(strings.TrimSpace(m))
		if len(words) > 2 {
			count++
		}
	}
	if count == 0 {
		return 1
	}
	return count
}

// countLongWords returns the number of words with more than 6 letters.
// Apostrophes are removed before measuring length (matching textstat's rm_apostrophe=True).
func countLongWords(text string) int {
	words := listWords(text)
	count := 0
	for _, w := range words {
		clean := strings.ReplaceAll(w, "'", "")
		if len(clean) > 6 {
			count++
		}
	}
	return count
}

// countPolysyllableWords returns the number of words with 3 or more syllables.
func countPolysyllableWords(text string) int {
	words := listWords(text)
	count := 0
	for _, w := range words {
		if countSyllablesWord(w) >= 3 {
			count++
		}
	}
	return count
}

// countMonosyllableWords returns the number of words with exactly 1 syllable.
func countMonosyllableWords(text string) int {
	words := listWords(text)
	count := 0
	for _, w := range words {
		if countSyllablesWord(w) == 1 {
			count++
		}
	}
	return count
}

// countDifficultWords returns the count of all occurrences of words
// not on the easy word list with at least syllableThreshold syllables.
// This matches textstat's count_difficult_words(unique=False).
func countDifficultWords(text string, syllableThreshold int) int {
	words := listWords(text)
	count := 0
	for _, w := range words {
		lower := strings.ToLower(w)
		if !isEasyWord(lower) {
			if syllableThreshold == 0 || countSyllablesWord(lower) >= syllableThreshold {
				count++
			}
		}
	}
	return count
}

// countMiniWords returns the number of words with max_size letters or fewer.
// Apostrophes are removed before measuring length (matching textstat's rm_apostrophe=True).
func countMiniWords(text string, maxSize int) int {
	words := listWords(text)
	count := 0
	for _, w := range words {
		clean := strings.ReplaceAll(w, "'", "")
		if len(clean) <= maxSize {
			count++
		}
	}
	return count
}

// Analysis holds all pre-computed values any formula needs, computed once per text.
// An Analysis is immutable after construction via NewAnalysis.
type Analysis struct {
	text              string
	words             []string
	charCount         int
	letterCount       int
	wordCount         int
	rawWordCount      int
	sentenceCount     int
	syllableCount     int
	longWordCount     int
	polysyllableCount int
	monosyllableCount int
	miniWordCount3    int
	difficultWords0   int     // syllable threshold=0 (Dale-Chall v1)
	difficultWords2   int     // threshold=2 (Dale-Chall v2, Spache)
	difficultWords3   int     // threshold=3 (Gunning Fog)
	wordsPerSentence  float64
	syllablesPerWord  float64
	charsPerWord      float64
	lettersPerWord    float64
	sentencesPerWord  float64
}

// NewAnalysis pre-computes all text statistics needed by readability formulas.
func NewAnalysis(text string) *Analysis {
	s := &Analysis{text: text}

	s.words = listWords(text)
	s.wordCount = len(s.words)
	s.charCount = countChars(text)
	s.letterCount = countLetters(text)
	s.sentenceCount = countSentences(text)
	s.rawWordCount = countRawWords(text)

	// Single loop over words to compute all word-level counts.
	for _, w := range s.words {
		syl := countSyllablesWord(w)
		s.syllableCount += syl

		clean := strings.ReplaceAll(w, "'", "")
		if len(clean) > 6 {
			s.longWordCount++
		}
		if len(clean) <= 3 {
			s.miniWordCount3++
		}

		if syl >= 3 {
			s.polysyllableCount++
		}
		if syl == 1 {
			s.monosyllableCount++
		}

		lower := strings.ToLower(w)
		if !isEasyWord(lower) {
			s.difficultWords0++
			if syl >= 2 {
				s.difficultWords2++
			}
			if syl >= 3 {
				s.difficultWords3++
			}
		}
	}

	// Derive ratio fields.
	if s.sentenceCount > 0 {
		s.wordsPerSentence = float64(s.wordCount) / float64(s.sentenceCount)
	}
	if s.wordCount > 0 {
		s.syllablesPerWord = float64(s.syllableCount) / float64(s.wordCount)
		s.lettersPerWord = float64(s.letterCount) / float64(s.wordCount) * 100
		s.sentencesPerWord = float64(s.sentenceCount) / float64(s.wordCount) * 100
	}
	if s.rawWordCount > 0 {
		s.charsPerWord = float64(s.charCount) / float64(s.rawWordCount)
	}

	return s
}

