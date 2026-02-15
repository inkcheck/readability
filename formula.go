package readability

import (
	"fmt"
	"math"
	"strings"
)

// FleschReadingEase calculates the Flesch Reading Ease score.
func (a *Analysis) FleschReadingEase() float64 {
	asl := a.wordsPerSentence
	asw := a.syllablesPerWord
	if asl == 0 || asw == 0 {
		return 0
	}
	return 206.835 - 1.015*asl - 84.6*asw
}

// FleschKincaidGrade calculates the Flesch-Kincaid Grade Level.
func (a *Analysis) FleschKincaidGrade() float64 {
	asl := a.wordsPerSentence
	asw := a.syllablesPerWord
	if asl == 0 || asw == 0 {
		return 0
	}
	return 0.39*asl + 11.8*asw - 15.59
}

// GunningFog calculates the Gunning Fog Index.
func (a *Analysis) GunningFog() float64 {
	if a.wordCount == 0 {
		return 0
	}
	asl := a.wordsPerSentence
	pctDifficult := float64(a.difficultWords3) / float64(a.wordCount) * 100
	return 0.4 * (asl + pctDifficult)
}

// SmogIndex calculates the SMOG Index.
func (a *Analysis) SmogIndex() float64 {
	if a.sentenceCount == 0 {
		return 0
	}
	return 1.043*math.Sqrt(float64(a.polysyllableCount)*30/float64(a.sentenceCount)) + 3.1291
}

// ColemanLiauIndex calculates the Coleman-Liau Index.
func (a *Analysis) ColemanLiauIndex() float64 {
	l := a.lettersPerWord
	sen := a.sentencesPerWord
	if l == 0 || sen == 0 {
		return 0
	}
	return 0.058*l - 0.296*sen - 15.8
}

// AutomatedReadabilityIndex calculates the Automated Readability Index.
func (a *Analysis) AutomatedReadabilityIndex() float64 {
	cpw := a.charsPerWord
	asl := a.wordsPerSentence
	if cpw == 0 || asl == 0 {
		return 0
	}
	return 4.71*cpw + 0.5*asl - 21.43
}

// DaleChallReadabilityScore calculates the Dale-Chall Readability Score.
func (a *Analysis) DaleChallReadabilityScore() float64 {
	if a.wordCount == 0 {
		return 0
	}
	asl := a.wordsPerSentence
	pctDifficult := float64(a.difficultWords0) / float64(a.wordCount) * 100
	score := 0.1579*pctDifficult + 0.0496*asl
	if pctDifficult > 5 {
		score += 3.6365
	}
	return score
}

// linsearWriteFormulaOpts calculates the Linsear Write Formula with configurable bounds.
// strictLower: return 0 if text has fewer than 100 words.
// strictUpper: use only the first 100 words (default behavior).
func linsearWriteFormulaOpts(text string, strictLower, strictUpper bool) float64 {
	// Get raw tokens (including punctuation) to preserve sentence boundaries.
	rawTokens := strings.Fields(text)
	if len(rawTokens) == 0 {
		return 0
	}

	// Collect cleaned words, optionally limiting to 100.
	var cleanedWords []string
	consumed := 0
	for _, token := range rawTokens {
		consumed++
		// Clean the token: remove non-contraction apostrophes and punctuation.
		w := rePunct.ReplaceAllString(token, "")
		w = strings.Trim(w, "'")
		if w != "" {
			cleanedWords = append(cleanedWords, w)
			if strictUpper && len(cleanedWords) >= 100 {
				break
			}
		}
	}

	if len(cleanedWords) == 0 {
		return 0
	}

	if strictLower && len(cleanedWords) < 100 {
		return 0
	}

	// Reconstruct text from original raw tokens for sentence counting.
	var truncatedText string
	if strictUpper {
		truncatedText = strings.Join(rawTokens[:consumed], " ")
	} else {
		truncatedText = text
	}
	sentences := countSentences(truncatedText)
	if sentences == 0 {
		sentences = 1
	}

	easyCount := 0
	hardCount := 0
	for _, w := range cleanedWords {
		if countSyllablesWord(w) >= 3 {
			hardCount++
		} else {
			easyCount++
		}
	}

	numerator := float64(easyCount*1 + hardCount*3)
	result := numerator / float64(sentences)
	if result <= 20 {
		result -= 2
	}
	return result / 2
}

// LinsearWriteFormula calculates the Linsear Write Formula with default options.
func (a *Analysis) LinsearWriteFormula() float64 {
	return linsearWriteFormulaOpts(a.text, false, true)
}

// LinsearWriteFormulaOpts calculates the Linsear Write Formula with configurable bounds.
// strictLower: return 0 if text has fewer than 100 words.
// strictUpper: use only the first 100 words.
func (a *Analysis) LinsearWriteFormulaOpts(strictLower, strictUpper bool) float64 {
	return linsearWriteFormulaOpts(a.text, strictLower, strictUpper)
}

// SpacheReadability calculates the Spache Readability Formula.
func (a *Analysis) SpacheReadability() float64 {
	if a.wordCount == 0 {
		return 0
	}
	asl := a.wordsPerSentence
	pctDifficult := float64(a.difficultWords2) / float64(a.wordCount) * 100
	return 0.141*asl + 0.086*pctDifficult + 0.839
}

// Lix calculates the LIX readability index.
func (a *Analysis) Lix() float64 {
	if a.wordCount == 0 {
		return 0
	}
	asl := a.wordsPerSentence
	return asl + 100*float64(a.longWordCount)/float64(a.wordCount)
}

// Rix calculates the RIX readability index.
// Uses the same long word count as LIX (> 6 letters, apostrophes removed).
func (a *Analysis) Rix() float64 {
	if a.sentenceCount == 0 {
		return 0
	}
	return float64(a.longWordCount) / float64(a.sentenceCount)
}

// DaleChallReadabilityScoreV2 calculates the New Dale-Chall Readability Score.
// Unlike v1, uses syllable_threshold=2 and checks raw_score > 0.05 for adjustment.
func (a *Analysis) DaleChallReadabilityScoreV2() float64 {
	if a.wordCount == 0 {
		return 0
	}
	asl := a.wordsPerSentence
	pctDifficult := float64(a.difficultWords2) / float64(a.wordCount) * 100
	score := 0.1579*pctDifficult + 0.0496*asl
	if score > 0.05 {
		score += 3.6365
	}
	return score
}

// McalpineEFLAW calculates the McAlpine EFLAW readability score.
func (a *Analysis) McalpineEFLAW() float64 {
	if a.sentenceCount == 0 {
		return 0
	}
	return float64(a.wordCount+a.miniWordCount3) / float64(a.sentenceCount)
}

// ReadingTimeWithRate returns estimated reading time in seconds using the given ms per character.
func (a *Analysis) ReadingTimeWithRate(msPerChar float64) float64 {
	return float64(a.charCount) * msPerChar / 1000
}

// ReadingTime returns estimated reading time in seconds using the default 14.69 ms per character.
func (a *Analysis) ReadingTime() float64 {
	return a.ReadingTimeWithRate(14.69)
}

// TextStandard estimates grade level as the mode of grade estimates from multiple formulas.
func (a *Analysis) TextStandard() float64 {
	if a.wordCount == 0 {
		return 0
	}

	grades := make([]int, 0, 24)

	addGrades := func(v float64) {
		// Uses RoundToEven (banker's rounding) to match Python's round().
		grades = append(grades, int(math.Floor(v)), int(math.Ceil(v)), int(math.RoundToEven(v)))
	}

	// Flesch-Kincaid Grade
	addGrades(a.FleschKincaidGrade())

	// Flesch Reading Ease → grade mapping
	fre := a.FleschReadingEase()
	switch {
	case fre >= 90:
		grades = append(grades, 5)
	case fre >= 80:
		grades = append(grades, 6)
	case fre >= 70:
		grades = append(grades, 7)
	case fre >= 60:
		grades = append(grades, 8, 9)
	case fre >= 50:
		grades = append(grades, 10)
	case fre >= 40:
		grades = append(grades, 11)
	case fre >= 30:
		grades = append(grades, 12)
	default:
		grades = append(grades, 13)
	}

	// SMOG
	addGrades(a.SmogIndex())

	// Coleman-Liau
	addGrades(a.ColemanLiauIndex())

	// ARI
	addGrades(a.AutomatedReadabilityIndex())

	// Dale-Chall
	addGrades(a.DaleChallReadabilityScore())

	// Linsear Write
	addGrades(a.LinsearWriteFormula())

	// Gunning Fog
	addGrades(a.GunningFog())

	// Find mode.
	if len(grades) == 0 {
		return 0
	}
	freq := make(map[int]int)
	for _, g := range grades {
		freq[g]++
	}
	modeVal := grades[0]
	modeCount := 0
	for val, cnt := range freq {
		if cnt > modeCount || (cnt == modeCount && val < modeVal) {
			modeVal = val
			modeCount = cnt
		}
	}
	return float64(modeVal)
}

// gradeSuffix returns the English ordinal suffix for an integer.
func gradeSuffix(n int) string {
	if n < 0 {
		n = -n
	}
	mod100 := n % 100
	if mod100 >= 11 && mod100 <= 13 {
		return "th"
	}
	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

// TextStandardString returns the text standard grade level as a string like "6th and 7th grade".
func (a *Analysis) TextStandardString() string {
	grade := int(a.TextStandard())
	lower := grade - 1
	return fmt.Sprintf("%d%s and %d%s grade", lower, gradeSuffix(lower), grade, gradeSuffix(grade))
}
