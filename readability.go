// Package readability provides text readability analysis using standard formulas.
//
// It is a Go rewrite of the Python textstat library, targeting English text.
package readability

import (
	"fmt"
	"sort"
)

// Formula identifies a readability formula by name.
type Formula string

const (
	FleschReadingEase          Formula = "flesch_reading_ease"
	FleschKincaidGrade         Formula = "flesch_kincaid_grade"
	GunningFog                 Formula = "gunning_fog"
	SmogIndex                  Formula = "smog_index"
	ColemanLiauIndex           Formula = "coleman_liau_index"
	AutomatedReadabilityIndex  Formula = "automated_readability_index"
	DaleChallReadabilityScore  Formula = "dale_chall_readability_score"
	DaleChallReadabilityScoreV2 Formula = "dale_chall_readability_score_v2"
	LinsearWriteFormula        Formula = "linsear_write_formula"
	SpacheReadability          Formula = "spache_readability"
	Lix                        Formula = "lix"
	McalpineEFLAW              Formula = "mcalpine_eflaw"
	Rix                        Formula = "rix"
	TextStandard               Formula = "text_standard"
	ReadingTime                Formula = "reading_time"
)

// allFormulas is the canonical sorted list of every formula.
var allFormulas = []Formula{
	AutomatedReadabilityIndex,
	ColemanLiauIndex,
	DaleChallReadabilityScore,
	DaleChallReadabilityScoreV2,
	FleschKincaidGrade,
	FleschReadingEase,
	GunningFog,
	LinsearWriteFormula,
	Lix,
	McalpineEFLAW,
	ReadingTime,
	Rix,
	SmogIndex,
	SpacheReadability,
	TextStandard,
}

// formulaFuncs maps each Formula to its method on *Analysis, used for
// dynamic dispatch via Score/ScoreAll.
var formulaFuncs = map[Formula]func(*Analysis) float64{
	FleschReadingEase:           (*Analysis).FleschReadingEase,
	FleschKincaidGrade:          (*Analysis).FleschKincaidGrade,
	GunningFog:                  (*Analysis).GunningFog,
	SmogIndex:                   (*Analysis).SmogIndex,
	ColemanLiauIndex:            (*Analysis).ColemanLiauIndex,
	AutomatedReadabilityIndex:   (*Analysis).AutomatedReadabilityIndex,
	DaleChallReadabilityScore:   (*Analysis).DaleChallReadabilityScore,
	DaleChallReadabilityScoreV2: (*Analysis).DaleChallReadabilityScoreV2,
	LinsearWriteFormula:         (*Analysis).LinsearWriteFormula,
	SpacheReadability:           (*Analysis).SpacheReadability,
	Lix:                         (*Analysis).Lix,
	McalpineEFLAW:               (*Analysis).McalpineEFLAW,
	Rix:                         (*Analysis).Rix,
	TextStandard:                (*Analysis).TextStandard,
	ReadingTime:                 (*Analysis).ReadingTime,
}

// Valid reports whether f is a recognized formula.
func (f Formula) Valid() bool {
	_, ok := formulaFuncs[f]
	return ok
}

// AllFormulas returns a sorted list of all available formula names.
func AllFormulas() []Formula {
	out := make([]Formula, len(allFormulas))
	copy(out, allFormulas)
	return out
}

// Stats holds raw text counting statistics.
type Stats struct {
	Chars             int
	Letters           int
	Words             int
	Sentences         int
	Syllables         int
	LongWords         int
	PolysyllableWords int
	MonosyllableWords int
}

// Stats returns raw text counting statistics.
func (a *Analysis) Stats() Stats {
	return Stats{
		Chars:             a.charCount,
		Letters:           a.letterCount,
		Words:             a.wordCount,
		Sentences:         a.sentenceCount,
		Syllables:         a.syllableCount,
		LongWords:         a.longWordCount,
		PolysyllableWords: a.polysyllableCount,
		MonosyllableWords: a.monosyllableCount,
	}
}

// Score calculates the given readability formula.
func (a *Analysis) Score(formula Formula) (float64, error) {
	fn, ok := formulaFuncs[formula]
	if !ok {
		return 0, fmt.Errorf("readability: unknown formula %q", formula)
	}
	return fn(a), nil
}

// ScoreAll calculates all formulas and returns a map of formula → score.
func (a *Analysis) ScoreAll() map[Formula]float64 {
	results := make(map[Formula]float64, len(allFormulas))
	for _, f := range allFormulas {
		results[f] = formulaFuncs[f](a)
	}
	return results
}

// SpacheReadabilityInt returns the Spache score truncated to an integer.
func (a *Analysis) SpacheReadabilityInt() int {
	return int(a.SpacheReadability())
}

// Analyzer is a convenience type whose methods accept text and internally
// construct an Analysis for each call.
//
// Deprecated: prefer NewAnalysis(text) for direct access to methods and
// to reuse a single Analysis across multiple formula calls.
type Analyzer struct{}

// Score calculates the given readability formula for text.
func (az *Analyzer) Score(text string, formula Formula) (float64, error) {
	return NewAnalysis(text).Score(formula)
}

// ScoreAll calculates all formulas for text and returns a map of formula → score.
func (az *Analyzer) ScoreAll(text string) map[Formula]float64 {
	return NewAnalysis(text).ScoreAll()
}

// Analyze returns raw text counting statistics.
func (az *Analyzer) Analyze(text string) Stats {
	return NewAnalysis(text).Stats()
}

// ReadingTime returns estimated reading time in seconds using the given ms per character.
func (az *Analyzer) ReadingTime(text string, msPerChar float64) float64 {
	return NewAnalysis(text).ReadingTimeWithRate(msPerChar)
}

// LinsearWrite calculates the Linsear Write Formula with configurable bounds.
// strictLower: return 0 if text has fewer than 100 words.
// strictUpper: use only the first 100 words (default behavior).
func (az *Analyzer) LinsearWrite(text string, strictLower, strictUpper bool) float64 {
	return linsearWriteFormulaOpts(text, strictLower, strictUpper)
}

// TextStandardString returns grade level as a string like "6th and 7th grade".
func (az *Analyzer) TextStandardString(text string) string {
	return NewAnalysis(text).TextStandardString()
}

// SpacheReadabilityInt returns the Spache score truncated to an integer.
func (az *Analyzer) SpacheReadabilityInt(text string) int {
	return NewAnalysis(text).SpacheReadabilityInt()
}

// FormulaNames returns all formula names as strings, sorted.
func FormulaNames() []string {
	names := make([]string, len(allFormulas))
	for i, f := range allFormulas {
		names[i] = string(f)
	}
	sort.Strings(names)
	return names
}
