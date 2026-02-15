package readability

import (
	"math"
	"strings"
	"testing"
)

// Test texts from textstat: tests/backend/resources.py
const shortText = "Cool dogs wear da sunglasses."

const punctText = `I said: 'This is a test sentence to test the remove_punctuation function.
It's short and not the work of a singer-songwriter. But it'll suffice.'
Your answer was: "I don't know. If I were you I'd write a test; just to make
sure, you're really just removing the characters you want to remove!" Didn't`

const longText = "Playing ... games has always been thought to be " +
	"important to the development of well-balanced and " +
	"creative children; however, what part, if any, " +
	"they should play in the lives of adults has never " +
	"been researched that deeply. I believe that " +
	"playing games is every bit as important for adults " +
	"as for children. Not only is taking time out to " +
	"play games with our children and other adults " +
	"valuable to building interpersonal relationships " +
	"but is also a wonderful way to release built up " +
	"tension.\n" +
	"There's nothing my husband enjoys more after a " +
	"hard day of work than to come home and play a game " +
	"of Chess with someone. This enables him to unwind " +
	"from the day's activities and to discuss the highs " +
	"and lows of the day in a non-threatening, kick back " +
	"environment. One of my most memorable wedding " +
	"gifts, a Backgammon set, was received by a close " +
	"friend. I asked him why in the world he had given " +
	"us such a gift. He replied that he felt that an " +
	"important aspect of marriage was for a couple to " +
	"never quit playing games together. Over the years, " +
	"as I have come to purchase and play, with other " +
	"couples & coworkers, many games like: Monopoly, " +
	"Chutes & Ladders, Mastermind, Dweebs, Geeks, & " +
	"Weirdos, etc. I can reflect on the integral part " +
	"they have played in our weekends and our " +
	"\"shut-off the T.V. and do something more " +
	"stimulating\" weeks. They have enriched my life and " +
	"made it more interesting. Sadly, many adults " +
	"forget that games even exist and have put them " +
	"away in the cupboards, forgotten until the " +
	"grandchildren come over.\n" +
	"All too often, adults get so caught up in working " +
	"to pay the bills and keeping up with the " +
	"\"Joneses'\" that they neglect to harness the fun " +
	"in life; the fun that can be the reward of " +
	"enjoying a relaxing game with another person. It " +
	"has been said that \"man is that he might have " +
	"joy\" but all too often we skate through life " +
	"without much of it. Playing games allows us to: " +
	"relax, learn something new and stimulating, " +
	"interact with people on a different more " +
	"comfortable level, and to enjoy non-threatening " +
	"competition. For these reasons, adults should " +
	"place a higher priority on playing games in their " +
	"lives"

const easyText = "Anna and her family love doing puzzles. Anna is best at " +
	"little puzzles. Anna and her brother work on medium size " +
	"puzzles together. Anna's Brother likes puzzles with cars " +
	"in them. When the whole family does a puzzle, they do really " +
	"big puzzles. It can take them a week to finish a really " +
	"big puzzle. Last year they did a puzzle with 500 pieces! " +
	"Anna tries to finish one small puzzle a day by her. " +
	"Her puzzles have about 50 pieces. They all glue their " +
	"favorite puzzles together and frame them. The puzzles look " +
	"so nice on the wall."

func round3(v float64) float64 {
	return math.Round(v*1000) / 1000
}

// TestFleschReadingEase tests against textstat expected values.
func TestFleschReadingEase(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 78.918},
		{"short", shortText, 83.32},
		{"punct", punctText, 81.148},
		{"long", longText, 59.771},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).FleschReadingEase())
			if got != tt.want {
				t.Errorf("fleschReadingEase = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFleschKincaidGrade tests against textstat expected values.
func TestFleschKincaidGrade(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 4.488},
		{"short", shortText, 2.88},
		{"punct", punctText, 4.574},
		{"long", longText, 10.359},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).FleschKincaidGrade())
			if got != tt.want {
				t.Errorf("fleschKincaidGrade = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGunningFog tests against textstat expected values.
func TestGunningFog(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 4.004},
		{"short", shortText, 10.0},
		{"punct", punctText, 7.259},
		{"long", longText, 11.441},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).GunningFog())
			if got != tt.want {
				t.Errorf("gunningFog = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSmogIndex tests against textstat expected values.
func TestSmogIndex(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 7.348},
		{"short", shortText, 8.842},
		{"punct", punctText, 8.239},
		{"long", longText, 11.67},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).SmogIndex())
			if got != tt.want {
				t.Errorf("smogIndex = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestColemanLiauIndex tests against textstat expected values.
func TestColemanLiauIndex(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 5.4},
		{"short", shortText, 6.12},
		{"punct", punctText, 6.249},
		{"long", longText, 9.134},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).ColemanLiauIndex())
			if got != tt.want {
				t.Errorf("colemanLiauIndex = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAutomatedReadabilityIndex tests against textstat expected values.
func TestAutomatedReadabilityIndex(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 3.575},
		{"short", shortText, 4.62},
		{"punct", punctText, 5.82},
		{"long", longText, 11.408},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).AutomatedReadabilityIndex())
			if got != tt.want {
				t.Errorf("automatedReadabilityIndex = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDaleChallReadabilityScore tests against textstat expected values.
func TestDaleChallReadabilityScore(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 7.592},
		{"short", shortText, 13.359}, // Go rounds 0.5 away from zero; Python rounds to even (13.358)
		{"punct", punctText, 6.248},
		{"long", longText, 8.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).DaleChallReadabilityScore())
			if got != tt.want {
				t.Errorf("daleChallReadabilityScore = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLinsearWriteFormula tests against textstat expected values.
// Our implementation uses strict_upper=True (first 100 words).
func TestLinsearWriteFormula(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"short", shortText, 2.5},
		{"punct", punctText, 5.1},
		{"easy", easyText, 4.045},
		{"long", longText, 15.25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).LinsearWriteFormula())
			if got != tt.want {
				t.Errorf("linsearWriteFormula = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSpacheReadability tests against textstat expected values.
func TestSpacheReadability(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 3.585},
		{"short", shortText, 3.264},
		{"punct", punctText, 3.469},
		{"long", longText, 5.473},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).SpacheReadability())
			if got != tt.want {
				t.Errorf("spacheReadability = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLix tests against textstat expected values.
func TestLix(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 22.131},
		{"short", shortText, 25.0},
		{"punct", punctText, 23.808},
		{"long", longText, 42.581},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).Lix())
			if got != tt.want {
				t.Errorf("lix = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRix tests against textstat expected values.
func TestRix(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 1.182},
		{"short", shortText, 1.0},
		{"punct", punctText, 1.4},
		{"long", longText, 4.529},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).Rix())
			if got != tt.want {
				t.Errorf("rix = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestReadingTime tests against textstat expected values.
// textstat uses configurable ms_per_char; we hardcode 14.69.
// We verify using ms_per_char=1.0 character counts from textstat.
func TestReadingTime(t *testing.T) {
	// textstat reading_time at ms_per_char=1.0 gives char count / 1000.
	// EASY_TEXT: 0.431 → 431 chars; SHORT_TEXT: 0.025 → 25 chars
	// Our formula: chars * 14.69 / 1000
	tests := []struct {
		name string
		text string
	}{
		{"empty", ""},
		{"easy", easyText},
		{"short", shortText},
		{"long", longText},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAnalysis(tt.text).ReadingTime()
			expected := float64(countChars(tt.text)) * 14.69 / 1000
			if round3(got) != round3(expected) {
				t.Errorf("readingTime = %v, want %v", got, expected)
			}
		})
	}
}

// TestTextStandard tests against textstat expected values.
func TestTextStandard(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 4.0},
		{"short", shortText, 2.0},
		{"punct", punctText, 6.0},
		{"long", longText, 11.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAnalysis(tt.text).TextStandard()
			if got != tt.want {
				t.Errorf("textStandard = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCountWords verifies basic word counting.
func TestCountWords(t *testing.T) {
	tests := []struct {
		text string
		want int
	}{
		{"", 0},
		{"hello world", 2},
		{"The cat sat on the mat.", 6},
		{"don't won't can't", 3},
	}
	for _, tt := range tests {
		got := countWords(tt.text)
		if got != tt.want {
			t.Errorf("countWords(%q) = %d, want %d", tt.text, got, tt.want)
		}
	}
}

// TestCountSyllablesWord verifies syllable counting for individual words.
func TestCountSyllablesWord(t *testing.T) {
	tests := []struct {
		word string
		want int
	}{
		{"the", 1},
		{"hello", 2},
		{"beautiful", 3},
		{"immediately", 5},
		{"a", 1},
		{"people", 2},
		{"eye", 1},
	}
	for _, tt := range tests {
		got := countSyllablesWord(tt.word)
		if got != tt.want {
			t.Errorf("countSyllablesWord(%q) = %d, want %d", tt.word, got, tt.want)
		}
	}
}

// TestIsEasyWord verifies the Dale-Chall word list lookup.
func TestIsEasyWord(t *testing.T) {
	if !isEasyWord("the") {
		t.Error("expected 'the' to be easy word")
	}
	if !isEasyWord("about") {
		t.Error("expected 'about' to be easy word")
	}
	if isEasyWord("antidisestablishmentarianism") {
		t.Error("expected 'antidisestablishmentarianism' to not be easy word")
	}
}

// TestDaleChallReadabilityScoreV2 tests against textstat expected values.
func TestDaleChallReadabilityScoreV2(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 6.794},
		{"short", shortText, 7.043},
		{"punct", punctText, 6.248},
		{"long", longText, 7.566},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).DaleChallReadabilityScoreV2())
			if got != tt.want {
				t.Errorf("daleChallReadabilityScoreV2 = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMcalpineEFLAW tests against textstat expected values.
func TestMcalpineEFLAW(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{"empty", "", 0.0},
		{"easy", easyText, 12.182},
		{"short", shortText, 6.0},
		{"punct", punctText, 15.2},
		{"long", longText, 30.765},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round3(NewAnalysis(tt.text).McalpineEFLAW())
			if got != tt.want {
				t.Errorf("mcalpineEFLAW = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestReadingTimeCustomRate verifies ReadingTime with a custom ms_per_char.
func TestReadingTimeCustomRate(t *testing.T) {
	a := &Analyzer{}
	// With ms_per_char=1.0, result should be chars/1000.
	tests := []struct {
		name string
		text string
	}{
		{"empty", ""},
		{"easy", easyText},
		{"short", shortText},
		{"long", longText},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.ReadingTime(tt.text, 1.0)
			want := float64(countChars(tt.text)) / 1000
			if round3(got) != round3(want) {
				t.Errorf("ReadingTime(msPerChar=1.0) = %v, want %v", got, want)
			}
		})
	}
}

// TestLinsearWriteStrictLower verifies returns 0 for short text when strictLower=true.
func TestLinsearWriteStrictLower(t *testing.T) {
	a := &Analyzer{}
	// shortText and easyText have fewer than 100 words, should return 0.
	for _, tt := range []struct {
		name string
		text string
	}{
		{"short", shortText},
		{"easy", easyText},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := a.LinsearWrite(tt.text, true, true)
			if got != 0 {
				t.Errorf("LinsearWrite(strictLower=true) = %v, want 0", got)
			}
		})
	}
	// longText has more than 100 words, should return non-zero.
	got := a.LinsearWrite(longText, true, true)
	if got == 0 {
		t.Error("LinsearWrite(longText, strictLower=true) = 0, want non-zero")
	}
}

// TestLinsearWriteNoStrictUpper verifies uses all words when strictUpper=false.
func TestLinsearWriteNoStrictUpper(t *testing.T) {
	a := &Analyzer{}
	// For longText (>100 words), strictUpper=false should give a different result
	// than strictUpper=true since it uses all words.
	withUpper := a.LinsearWrite(longText, false, true)
	withoutUpper := a.LinsearWrite(longText, false, false)
	if withUpper == withoutUpper {
		t.Errorf("expected different results: strictUpper=true gave %v, strictUpper=false gave %v", withUpper, withoutUpper)
	}
	// For short text (<100 words), both should give the same result.
	withUpper = a.LinsearWrite(easyText, false, true)
	withoutUpper = a.LinsearWrite(easyText, false, false)
	if withUpper != withoutUpper {
		t.Errorf("for short text expected same results: strictUpper=true gave %v, strictUpper=false gave %v", withUpper, withoutUpper)
	}
}

// TestTextStandardString verifies string output for each test text.
func TestTextStandardString(t *testing.T) {
	a := &Analyzer{}
	tests := []struct {
		name string
		text string
		want string
	}{
		{"empty", "", "-1st and 0th grade"},
		{"easy", easyText, "3rd and 4th grade"},
		{"short", shortText, "1st and 2nd grade"},
		{"punct", punctText, "5th and 6th grade"},
		{"long", longText, "10th and 11th grade"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.TextStandardString(tt.text)
			if got != tt.want {
				t.Errorf("TextStandardString = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestSpacheReadabilityInt verifies truncated int output.
func TestSpacheReadabilityInt(t *testing.T) {
	a := &Analyzer{}
	tests := []struct {
		name string
		text string
		want int
	}{
		{"empty", "", 0},
		{"easy", easyText, 3},
		{"short", shortText, 3},
		{"punct", punctText, 3},
		{"long", longText, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.SpacheReadabilityInt(tt.text)
			if got != tt.want {
				t.Errorf("SpacheReadabilityInt = %d, want %d", got, tt.want)
			}
		})
	}
}

// BenchmarkCountSyllablesWord benchmarks syllable counting for individual words.
func BenchmarkCountSyllablesWord(b *testing.B) {
	words := []string{
		"the", "hello", "beautiful", "immediately", "communication",
		"environment", "comfortable", "interesting", "relationship",
		"playing", "stimulating", "nonthreatening", "interpersonal",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range words {
			countSyllablesWord(w)
		}
	}
}

// BenchmarkCountSyllables benchmarks syllable counting for full text.
func BenchmarkCountSyllables(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		countSyllables(longText)
	}
}

// BenchmarkScoreAll benchmarks computing all formulas for a text.
func BenchmarkScoreAll(b *testing.B) {
	a := &Analyzer{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.ScoreAll(longText)
	}
}

// FuzzCountSyllablesWord fuzzes syllable counting to catch panics on unusual input.
func FuzzCountSyllablesWord(f *testing.F) {
	f.Add("hello")
	f.Add("beautiful")
	f.Add("immediately")
	f.Add("")
	f.Add("123")
	f.Add("don't")
	f.Add("🎉")
	f.Add("café")
	f.Add(strings.Repeat("a", 10000))
	f.Fuzz(func(t *testing.T, word string) {
		n := countSyllablesWord(word)
		if n < 0 {
			t.Errorf("countSyllablesWord(%q) = %d, want >= 0", word, n)
		}
	})
}

// TestAllFormulas verifies formula count.
func TestAllFormulas(t *testing.T) {
	formulas := AllFormulas()
	if len(formulas) != 15 {
		t.Errorf("AllFormulas() returned %d formulas, want 15", len(formulas))
	}
}

// TestAnalyzerScore verifies Score dispatches without error.
func TestAnalyzerScore(t *testing.T) {
	a := &Analyzer{}
	for _, f := range AllFormulas() {
		score, err := a.Score(longText, f)
		if err != nil {
			t.Errorf("Score(%q) error: %v", f, err)
		}
		if math.IsNaN(score) || math.IsInf(score, 0) {
			t.Errorf("Score(%q) non-finite: %f", f, score)
		}
	}
	_, err := a.Score(longText, "nonexistent")
	if err == nil {
		t.Error("expected error for unknown formula")
	}
}

// TestAnalyzerScoreAll verifies ScoreAll returns all formulas.
func TestAnalyzerScoreAll(t *testing.T) {
	a := &Analyzer{}
	results := a.ScoreAll(longText)
	if len(results) != 15 {
		t.Errorf("ScoreAll returned %d results, want 15", len(results))
	}
}

// TestEmptyText verifies all formulas handle empty input gracefully.
func TestEmptyText(t *testing.T) {
	a := &Analyzer{}
	for _, f := range AllFormulas() {
		score, err := a.Score("", f)
		if err != nil {
			t.Errorf("Score on empty for %q error: %v", f, err)
		}
		if math.IsNaN(score) || math.IsInf(score, 0) {
			t.Errorf("Score on empty for %q non-finite: %f", f, score)
		}
	}
}
