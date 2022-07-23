// Package odphone (Odia Phone) is a phonetic algorithm for indexing
// unicode Odia words by their pronounciation, like Metaphone for English.
// The algorithm generates three Romanized phonetic keys (hashes) of varying
// phonetic proximity for a given Odia word.
//
// The algorithm takes into account the context sensitivity of sounds, syntactic
// and phonetic gemination, compounding, modifiers, and other known exceptions
// to produce Romanized phonetic hashes of increasing phonetic affinity that are
// faithful to the pronunciation of the original Odia word.
//
// `key0` = a broad phonetic hash comparable to a Metaphone key that doesn't account
// for hard sounds or phonetic modifiers
//
// `key1` = is a slightly more inclusive hash that accounts for hard sounds
//
// `key2` = highly inclusive and narrow hash that accounts for hard sounds
// and phonetic modifiers
//
// odphone was created to aid spelling tolerant Odia word search, but may
// be useful in tasks like spell checking, word suggestion etc.
//
// This is based on MLphone (https://github.com/knadh/mlphone/) for Malayalam.
//
// Soumendra Kumar Sahoo (c) 2022. https://www.soumendrak.com | License: GPLv3
package odphone

import (
	"regexp"
	"strings"
)

var vowels = map[string]string{
	"ଅ": "A",
    "ଆ": "A",
    "ଇ": "E",
    "ଈ": "E",
    "ଉ": "U",
    "ଊ": "U",
    "ଋ": "R",
	"ୠ": "R",
    "ଏ": "E",
    "ଐ": "AI",
    "ଓ": "O",
    "ଔ": "O",
}

var consonants = map[string]string{
	"କ": "K",
    "ଖ": "K",
    "ଗ": "G",
    "ଘ": "G",
    "ଙ": "WN",
    "ଚ": "C",
    "ଛ": "C",
    "ଜ": "J",
    "ଝ": "J",
    "ଞ": "N",
    "ଟ": "T",
    "ଠ": "T",
    "ଡ": "D",
    "ଢ": "D",
    "ଣ": "N",
    "ତ": "T",
    "ଥ": "T",
    "ଦ": "D",
    "ଧ": "D",
    "ନ": "N",
    "ପ": "P",
    "ଫ": "F",
    "ବ": "B",
    "ଭ": "V",
    "ମ": "M",
    "ଯ": "J",
    "ର": "R",
    "ଲ": "L",
    "ଳ": "L",
    "ଵ": "B",
    "ଶ": "S",
    "ଷ": "S",
    "ସ": "S",
    "ହ": "H",
    "ୟ": "Y",
    "ୱ": "W",
}

var compounds = map[string]string{
	// TODO: Tobe done for Odia
	"ಕ್ಕ": "K2", "ಗ್ಗಾ": "K", "ಙ್ಙ": "NG",
	"ಚ್ಚ": "C2", "ಜ್ಜ": "J", "ಞ್ಞ": "NJ",
	"ಟ್ಟ": "T2", "ಣ್ಣ": "N2",
	"ತ್ತ": "0", "ದ್ದ": "D", "ದ್ಧ": "D", "ನ್ನ": "NN",
	"ಬ್ಬ": "B",
	"ಪ್ಪ": "P2", "ಮ್ಮ": "M2",
	"ಯ್ಯ": "Y", "ಲ್ಲ": "L2", "ವ್ವ": "V", "ಶ್ಶ": "S1", "ಸ್ಸ": "S",
	"ಳ್ಳ": "L12",
	"ಕ್ಷ": "KS1",
}

var modifiers = map[string]string{
	"ଁ": "1",
	"ଂ": "1",
	"ଃ": "3",
	"଼": "4",
	"ଽ": "",
	"ା": "6",
	"ି": "7",
	"ୀ": "7",
	"ୁ": "8",
	"ୂ": "8",
	"ୃ": "9",
	"ୄ": "9",
	"େ": "2",
	"ୈ": "2",
	"ୋ": "5",
	"ୌ": "5",
	"୍": "4",
	"ୖ": "2",
	"ୗ": "",
}

var (
	regexKey0, _       = regexp.Compile(`[1,2,4-9]`)
	regexKey1, _       = regexp.Compile(`[2,4-9]`)
	regexNonOdia, _ = regexp.Compile(`[\P{Odia}]`)
	regexAlphaNum, _   = regexp.Compile(`[^0-9A-Z]`)
)

// ODphone is the Odia-phone tokenizer.
type ODphone struct {
	modCompounds  *regexp.Regexp
	modConsonants *regexp.Regexp
	modVowels     *regexp.Regexp
}

// New returns a new instance of the ODPhone tokenizer.
func New() *ODphone {
	var (
		glyphs []string
		mods   []string
		od     = &ODphone{}
	)

	// modifiers.
	for k := range modifiers {
		mods = append(mods, k)
	}

	// compounds.
	for k := range compounds {
		glyphs = append(glyphs, k)
	}
	od.modCompounds, _ = regexp.Compile(`((` + strings.Join(glyphs, "|") + `)(` + strings.Join(mods, "|") + `))`)

	// consonants.
	glyphs = []string{}
	for k := range consonants {
		glyphs = append(glyphs, k)
	}
	od.modConsonants, _ = regexp.Compile(`((` + strings.Join(glyphs, "|") + `)(` + strings.Join(mods, "|") + `))`)

	// vowels.
	glyphs = []string{}
	for k := range vowels {
		glyphs = append(glyphs, k)
	}
	od.modVowels, _ = regexp.Compile(`((` + strings.Join(glyphs, "|") + `)(` + strings.Join(mods, "|") + `))`)

	return od
}

// Encode encodes a unicode Odia string to its Roman ODPhone hash.
// Ideally, words should be encoded one at a time, and not as phrases
// or sentences.
func (k *ODphone) Encode(input string) (string, string, string) {
	// key2 accounts for hard and modified sounds.
	key2 := k.process(input)

	// key1 loses numeric modifiers that denote phonetic modifiers.
	key1 := regexKey1.ReplaceAllString(key2, "")

	// key0 loses numeric modifiers that denote hard sounds, doubled sounds,
	// and phonetic modifiers.
	key0 := regexKey0.ReplaceAllString(key2, "")

	return key0, key1, key2
}

func (k *ODphone) process(input string) string {
	// Remove all non-malayalam characters.
	input = regexNonOdia.ReplaceAllString(strings.Trim(input, ""), "")

	// All character replacements are grouped between { and } to maintain
	// separatability till the final step.

	// Replace and group modified compounds.
	input = k.replaceModifiedGlyphs(input, compounds, k.modCompounds)

	// Replace and group unmodified compounds.
	for k, v := range compounds {
		input = strings.ReplaceAll(input, k, `{`+v+`}`)
	}

	// Replace and group modified consonants and vowels.
	input = k.replaceModifiedGlyphs(input, consonants, k.modConsonants)
	input = k.replaceModifiedGlyphs(input, vowels, k.modVowels)

	// Replace and group unmodified consonants.
	for k, v := range consonants {
		input = strings.ReplaceAll(input, k, `{`+v+`}`)
	}

	// Replace and group unmodified vowels.
	for k, v := range vowels {
		input = strings.ReplaceAll(input, k, `{`+v+`}`)
	}

	// Replace all modifiers.
	for k, v := range modifiers {
		input = strings.ReplaceAll(input, k, v)
	}

	// Remove non alpha numeric characters (losing the bracket grouping).
	return regexAlphaNum.ReplaceAllString(input, "")
}

func (k *ODphone) replaceModifiedGlyphs(input string, glyphs map[string]string, r *regexp.Regexp) string {
	for _, matches := range r.FindAllStringSubmatch(input, -1) {
		for _, m := range matches {
			if rep, ok := glyphs[m]; ok {
				input = strings.ReplaceAll(input, m, rep)
			}
		}
	}
	return input
}