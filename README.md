# Readability

A Go library and CLI tool for measuring text readability. It implements 15
standard readability formulas commonly used in education, publishing, and
content analysis. This is a Go rewrite of the Python
[textstat](https://github.com/textstat/textstat) library.

## Features

- **15 readability formulas** including Flesch-Kincaid, Gunning Fog, SMOG,
  Coleman-Liau, Dale-Chall, and more
- **Zero dependencies** beyond the Go standard library
- **Library and CLI** for use in Go programs or from the command line
- **Tested against textstat** to ensure compatible results

## Installation

### Library

```bash
go get github.com/inkcheck/readability
```

### CLI

Using Homebrew:

```bash
brew install inkcheck/tap/readability
```

Build from source:

```bash
go build -o readability ./cmd/readability
```

Or install to your `$GOPATH/bin` (so it's available globally):

```bash
go install ./cmd/readability
```

## Quick Start

### As a Library

```go
package main

import (
    "fmt"
    "github.com/inkcheck/readability"
)

func main() {
    text := "The quick brown fox jumps over the lazy dog."
    a := readability.Analyzer{}

    score, err := a.Score(text, "flesch_reading_ease")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Flesch Reading Ease: %.2f\n", score)

    // Compute all formulas at once
    scores := a.ScoreAll(text)
    for formula, s := range scores {
        fmt.Printf("%s: %.2f\n", formula, s)
    }

    // Get raw text statistics
    stats := a.Analyze(text)
    fmt.Printf("Words: %d, Sentences: %d, Syllables: %d\n",
        stats.Words, stats.Sentences, stats.Syllables)
}
```

### As a CLI Tool

```bash
# Score a single file
readability -f flesch_kincaid_grade essay.txt

# Score all text files in a directory
readability -f gunning_fog ./documents/

# Read from stdin
echo "The cat sat on the mat." | readability -f flesch_reading_ease

# List available formulas
readability -f help
```

## Available Formulas

| Formula                          | Output          | Description                                |
|----------------------------------|-----------------|--------------------------------------------|
| `flesch_reading_ease`            | 0-100+ score    | Higher = easier to read                    |
| `flesch_kincaid_grade`           | US grade level  | Grade level needed to understand text      |
| `gunning_fog`                    | US grade level  | Estimates years of education needed        |
| `smog_index`                     | US grade level  | Based on polysyllable word count           |
| `coleman_liau_index`             | US grade level  | Character-based, no syllable counting      |
| `automated_readability_index`    | US grade level  | Character-based formula                    |
| `dale_chall_readability_score`   | Raw score       | Uses familiar-word list                    |
| `dale_chall_readability_score_v2`| Raw score       | Variant with syllable threshold            |
| `linsear_write_formula`          | US grade level  | Easy/hard word ratio over first 100 words  |
| `spache_readability`             | US grade level  | Designed for primary-grade text            |
| `lix`                            | Index           | Scandinavian readability index             |
| `rix`                            | Index           | Simplified LIX variant                     |
| `mcalpine_eflaw`                 | Index           | Counts mini-words (3 letters or fewer)     |
| `text_standard`                  | US grade level  | Consensus grade from multiple formulas     |
| `reading_time`                   | Seconds         | Estimated reading time                     |

See the [formulas reference](docs/formulas.md) for details on each formula.

## Documentation

- [Formulas Reference](docs/formulas.md) - detailed descriptions and
  mathematical definitions
- [Library API](docs/api.md) - exported types and methods
- [CLI Usage](docs/cli.md) - command-line tool reference
- [Syllable Counting](docs/syllables.md) - how syllable counting works
- [Architecture](docs/architecture.md) - package structure and design decisions

## License

This project is a Go port of the [textstat](https://github.com/textstat/textstat)
library created by Shivam Bansal, built using
[Claude Code](https://claude.ai/claude-code).
See [LICENSE](LICENSE) for details.

