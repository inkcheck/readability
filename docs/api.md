# Library API

## Import

```go
import "github.com/inkcheck/readability"
```

## Types

### Formula

```go
type Formula string
```

A string identifying a readability formula. Use one of the 15 named constants
(e.g. `"flesch_reading_ease"`) or call `AllFormulas()` to list them.

**Methods:**

- `Valid() bool` - returns true if the formula name is recognized.

### Stats

```go
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
```

Raw text statistics returned by `Analyzer.Analyze`. Long words have more than
6 letters. Polysyllable words have 3 or more syllables. Monosyllable words
have exactly 1 syllable.

### Analyzer

```go
type Analyzer struct{}
```

The main entry point for scoring text. It caches statistics internally so
repeated calls with the same text avoid redundant computation.

## Analyzer Methods

### Score

```go
func (a *Analyzer) Score(text string, formula Formula) (float64, error)
```

Computes a single readability score. Returns an error if the formula name is
not recognized.

```go
a := readability.Analyzer{}
score, err := a.Score(text, "flesch_reading_ease")
```

### ScoreAll

```go
func (a *Analyzer) ScoreAll(text string) map[Formula]float64
```

Computes all 15 formulas and returns a map of formula name to score.

```go
scores := a.ScoreAll(text)
for formula, score := range scores {
    fmt.Printf("%s: %.2f\n", formula, score)
}
```

### Analyze

```go
func (a *Analyzer) Analyze(text string) Stats
```

Returns raw text statistics (word count, sentence count, syllable count, etc.)
without computing any formula.

### ReadingTime

```go
func (a *Analyzer) ReadingTime(text string, msPerChar float64) float64
```

Estimates reading time in seconds. Pass 0 for `msPerChar` to use the default
rate of 14.69 ms/character.

### LinsearWrite

```go
func (a *Analyzer) LinsearWrite(text string, strictLower, strictUpper bool) float64
```

Linsear Write formula with boundary control:

- `strictLower` - if true, returns 0 when the text has fewer than 100 words.
- `strictUpper` - if true (default behavior), truncates to the first 100 words.
  If false, uses the full text.

### TextStandardString

```go
func (a *Analyzer) TextStandardString(text string) string
```

Returns the consensus grade level as a human-readable string, e.g. `"6th and
7th grade"`.

### SpacheReadabilityInt

```go
func (a *Analyzer) SpacheReadabilityInt(text string) int
```

Returns the Spache readability score truncated to an integer.

## Package Functions

### AllFormulas

```go
func AllFormulas() []Formula
```

Returns all recognized formula values, sorted alphabetically.

### FormulaNames

```go
func FormulaNames() []string
```

Returns all formula names as strings, sorted alphabetically. Useful for
displaying help text or building CLI flags.
