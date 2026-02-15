# CLI Usage

## Installation

```bash
go install github.com/inkcheck/readability/cmd/readability@latest
```

## Synopsis

```
readability -f <formula> [file or directory ...]
readability -f <formula> < input.txt
```

## Options

- `-f <formula>` (required) - the readability formula to apply. Run
  `readability -f help` to see the full list.

## Input

The tool accepts one or more file paths or directory paths as arguments. When
given a directory it walks it recursively and processes files with these
extensions (case-insensitive):

- `.txt`
- `.md`
- `.rst`
- `.adoc`
- `.tex`

If no arguments are provided, text is read from standard input.

## Output

**Single file or stdin:** prints the numeric score on a line by itself.

```
$ echo "The cat sat on the mat." | readability -f flesch_reading_ease
116.15
```

**Multiple files:** prints one line per file in `path<tab>score` format.

```
$ readability -f flesch_kincaid_grade docs/
docs/intro.txt    6.20
docs/advanced.txt 11.40
```

## Examples

Score a single file with Flesch-Kincaid:

```bash
readability -f flesch_kincaid_grade essay.txt
```

Score all Markdown files in a directory:

```bash
readability -f gunning_fog ./content/
```

Pipe text from another command:

```bash
curl -s https://example.com/article.txt | readability -f smog_index
```

Get the consensus grade level:

```bash
readability -f text_standard report.txt
```

List available formulas:

```bash
readability -f help
```
