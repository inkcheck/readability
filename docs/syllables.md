# Syllable Counting

Several readability formulas depend on syllable counts. This package uses a
rule-based heuristic to estimate syllables in English words.

## Algorithm

1. **Lookup table** - an irregular-word dictionary of 88 entries handles common
   words whose syllable counts are hard to derive from spelling (e.g.
   "beautiful" = 3, "business" = 2, "people" = 2).

2. **Numbers** - any word containing only digits counts as 1 syllable.

3. **Short words** - words of 2 or fewer characters count as 1 syllable.

4. **Vowel groups** - the base count is the number of contiguous vowel groups
   (a, e, i, o, u, y) in the lowercased, alpha-only form of the word.

5. **Silent endings** - the count is reduced for common silent patterns:
   - **Silent E**: a final `e` preceded by a consonant (but not `-le`).
   - **Silent ES**: a final `es` preceded by a consonant other than l, s, x,
     z, c, or g.
   - **Silent ED**: a final `ed` preceded by a consonant (excluding the words
     bed, fed, red, led, wed, shed).

6. **Hiatus additions** - the count is increased when adjacent vowels are
   pronounced separately:
   - `ia`, `io` patterns (minus `-iage`, `-tion`, `-sion`)
   - `iou`, `ien`, `ual`, `eo` (not `eou`), `ism`
   - `ua` not followed by `l`

7. **Overlap deductions** - suffixes like `-ian`, `-ium`, `-ial` are subtracted
   to avoid double-counting from the hiatus rules.

8. **Floor** - the result is clamped to a minimum of 1.

## Accuracy

The heuristic is designed for general English text. It matches the behavior of
the Python textstat library. Edge cases exist for borrowed words, technical
jargon, and proper nouns. The irregular-word table covers the most common
mispredictions.

## Extending the Irregular Words Table

The irregular words are defined in `syllable.go` in the `irregularWords` map.
To correct a mispredicted word, add it to the map with its known syllable
count:

```go
var irregularWords = map[string]int{
    // existing entries ...
    "segue": 2,
}
```
