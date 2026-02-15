# Formulas Reference

This document describes each readability formula, its mathematical definition,
and how to interpret its output.

## Flesch Reading Ease

**Formula name:** `flesch_reading_ease`

```
206.835 - 1.015 * (words / sentences) - 84.6 * (syllables / words)
```

Returns a score from roughly 0 to 100 (can exceed 100 for very simple text).
Higher scores indicate easier text.

| Score     | Difficulty       |
|-----------|------------------|
| 90-100    | Very easy        |
| 80-89     | Easy             |
| 70-79     | Fairly easy      |
| 60-69     | Standard         |
| 50-59     | Fairly difficult |
| 30-49     | Difficult        |
| 0-29      | Very confusing   |

## Flesch-Kincaid Grade Level

**Formula name:** `flesch_kincaid_grade`

```
0.39 * (words / sentences) + 11.8 * (syllables / words) - 15.59
```

Returns a US school grade level. A score of 8.0 means an eighth-grader should
be able to understand the text.

## Gunning Fog Index

**Formula name:** `gunning_fog`

```
0.4 * ((words / sentences) + 100 * (complex_words / words))
```

Complex words are those with three or more syllables, excluding proper nouns,
familiar jargon, and compound words. Returns a US grade level.

## SMOG Index

**Formula name:** `smog_index`

```
1.043 * sqrt(polysyllable_count * (30 / sentences)) + 3.1291
```

SMOG (Simple Measure of Gobbledygook) estimates the years of education needed
to understand a piece of writing. Polysyllable words have three or more
syllables.

## Coleman-Liau Index

**Formula name:** `coleman_liau_index`

```
0.0588 * L - 0.296 * S - 15.8
```

Where L is the average number of letters per 100 words and S is the average
number of sentences per 100 words. Unlike most formulas, Coleman-Liau uses
character counts instead of syllable counts.

## Automated Readability Index

**Formula name:** `automated_readability_index`

```
4.71 * (characters / words) + 0.5 * (words / sentences) - 21.43
```

Character-based formula that returns a US grade level. Like Coleman-Liau, it
avoids syllable counting.

## Dale-Chall Readability Score

**Formula name:** `dale_chall_readability_score`

```
0.1579 * (difficult_words / words * 100) + 0.0496 * (words / sentences)
```

If more than 5% of words are "difficult" (not on the Dale-Chall list of ~3,000
familiar words), add 3.6365 to the raw score.

| Score       | Grade Level |
|-------------|-------------|
| 4.9 or less | Grade 4     |
| 5.0 - 5.9  | Grades 5-6  |
| 6.0 - 6.9  | Grades 7-8  |
| 7.0 - 7.9  | Grades 9-10 |
| 8.0 - 8.9  | Grades 11-12|
| 9.0 - 9.9  | Grades 13-15|

## Dale-Chall Readability Score V2

**Formula name:** `dale_chall_readability_score_v2`

A variant that uses a syllable threshold of 2 when identifying difficult words,
and triggers the adjustment when the raw score exceeds 0.05 rather than when
the difficult-word percentage exceeds 5%.

## Linsear Write Formula

**Formula name:** `linsear_write_formula`

Operates on the first 100 words of the text. Classifies each word as easy
(1-2 syllables, weighted 1) or hard (3+ syllables, weighted 3).

```
raw = (easy_count * 1 + hard_count * 3) / sentences
if raw > 20: grade = raw / 2
else:         grade = (raw - 1) / 2
```

The `Analyzer.LinsearWrite` method accepts `strictLower` and `strictUpper`
flags to control the 100-word boundary behavior.

## Spache Readability

**Formula name:** `spache_readability`

```
0.141 * (words / sentences) + 0.086 * (difficult_words / words * 100) + 0.839
```

Designed for primary-grade text (grades 1-3). Difficult words are those not on
the Dale-Chall list and with 2 or more syllables.

## LIX (Läsbarhetsindex)

**Formula name:** `lix`

```
(words / sentences) + 100 * (long_words / words)
```

A Scandinavian readability index. Long words have more than 6 characters.

| Score  | Difficulty |
|--------|------------|
| < 25   | Very easy  |
| 25-35  | Easy       |
| 35-45  | Standard   |
| 45-55  | Difficult  |
| > 55   | Very hard  |

## RIX

**Formula name:** `rix`

```
long_words / sentences
```

A simplified variant of LIX. Long words have more than 6 characters.

## McAlpine EFLAW

**Formula name:** `mcalpine_eflaw`

```
(words + mini_words) / sentences
```

Mini-words are words with 3 or fewer characters. Higher scores indicate more
difficult text.

## Text Standard

**Formula name:** `text_standard`

Computes the grade level from multiple formulas (Flesch-Kincaid, Flesch Reading
Ease, SMOG, Coleman-Liau, ARI, Dale-Chall, Linsear Write, and Gunning Fog),
then returns the most frequent (mode) grade level as a consensus estimate.

The `Analyzer.TextStandardString` method returns a human-readable string like
"6th and 7th grade".

## Reading Time

**Formula name:** `reading_time`

```
characters * ms_per_char / 1000
```

Estimates reading time in seconds. The default rate is 14.69 milliseconds per
character. A custom rate can be set via `Analyzer.ReadingTime`.
