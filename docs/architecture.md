# Architecture

## Package Structure

```
readability/
├── cmd/
│   └── readability/
│       └── main.go          # CLI entry point
├── docs/                    # Documentation
├── wordlist/
│   └── dale_chall.txt       # Embedded word list (2,940 words)
├── count.go                 # Text counting functions
├── formula.go               # Formula implementations
├── readability.go           # Public API (Analyzer, Formula, Stats)
├── syllable.go              # Syllable counting heuristic
├── wordlist.go              # Word list loading via go:embed
├── readability_test.go      # Tests and benchmarks
└── go.mod
```

## Source Files

| File              | Purpose                                              |
|-------------------|------------------------------------------------------|
| `readability.go`  | Public types (`Analyzer`, `Formula`, `Stats`) and method dispatch |
| `formula.go`      | One function per formula, all unexported              |
| `count.go`        | Word, sentence, character, and syllable counting      |
| `syllable.go`     | Syllable estimation heuristic and irregular-word table |
| `wordlist.go`     | Embeds `dale_chall.txt` and provides `isEasyWord()`  |
| `cmd/readability/main.go` | CLI flag parsing, file walking, stdin support |

## Design Decisions

### Zero Dependencies

The package uses only the Go standard library. The Dale-Chall word list is
embedded at compile time with `//go:embed`, so the resulting binary is fully
self-contained.

### Statistic Caching

`Analyzer` precomputes all text statistics (counts, ratios, difficult-word
counts at multiple thresholds) into an internal `textStats` struct on first
use. This avoids redundant work when multiple formulas are applied to the same
text, as in `ScoreAll`.

### Python Compatibility

The formulas are ported from the Python textstat library. Tests compare output
against known textstat results. The package uses banker's rounding
(`math.RoundToEven`) where needed to match Python's `round()` behavior.

### Unexported Internals

All counting functions and formula implementations are unexported. The public
surface is intentionally small: `Analyzer`, `Formula`, `Stats`, and a few
package-level helpers.

## Data Flow

```
text
  │
  ├─► countWords, countSentences, countChars, ...  (count.go)
  │     └─► countSyllables  (syllable.go)
  │     └─► isEasyWord      (wordlist.go)
  │
  ├─► textStats (precomputed ratios and counts)
  │
  └─► formula functions  (formula.go)
        └─► float64 score
```

## Testing

Tests live in `readability_test.go` and use four sample texts of varying
complexity (6 words to ~470 words). Expected values are taken from the Python
textstat library. The test suite includes:

- Unit tests for each formula
- Counting function tests
- Analyzer integration tests (single score, all scores, empty text)
- Benchmarks for syllable counting and full scoring
- Fuzz tests for syllable counting robustness
