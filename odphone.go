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
	"fmt"
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
	"ଐ": "EI",
	"ଓ": "O",
	"ଔ": "OU",
}

var consonants = map[string]string{
	"କ": "K",
	"ଖ": "KH",
	"ଗ": "G",
	"ଘ": "GH",
	"ଙ": "WN",
	"ଚ": "CH",
	"ଛ": "CHH",
	"ଜ": "J",
	"ଝ": "JH",
	"ଞ": "NY",
	"ଟ": "T",
	"ଠ": "TH",
	"ଡ": "D",
	"ଢ": "DH",
	"ଣ": "N",
	"ତ": "T",
	"ଥ": "TH",
	"ଦ": "D",
	"ଧ": "DH",
	"ନ": "N",
	"ପ": "P",
	"ଫ": "F",
	"ବ": "B",
	"ଭ": "V",
	"ମ": "M",
	"ଯ": "J",
	"ର": "R",
	"ଲ": "L",
	"ଳ": "LH",
	"ଵ": "B",
	"ଶ": "SH",
	"ଷ": "SH",
	"ସ": "S",
	"ହ": "H",
	"ୟ": "Y",
	"ୱ": "W",
}

var compounds = map[string]string{
	// TODO: Tobe done for Odia
	"କ୍ତ": "K2",
	"ଙ୍କ": "K3",
	"ଙ୍ଗ": "NG",
	"ଙ୍ଘ": "NG2",
	"ଞ୍ଜ": "NJ",
}

var modifiers = map[string]string{
	"ା": "1",
	"଼": "2",
	"୍": "2",
	"ୖ": "3",
	"େ": "3",
	"ୈ": "3",
	"ୗ": "3",
	"ୋ": "4",
	"ୌ": "4",
	"ି": "5",
	"ୀ": "5",
	"ୁ": "6",
	"ୂ": "6",
	"ୃ": "6",
	"ଃ": "7",
	"ଁ": "7",
	"ଂ": "7",
	"ୄ": "8",
	"ଽ": "8",
}

var (
	regexKey0, _     = regexp.Compile(`[1-8]`)
	regexKey1, _     = regexp.Compile(`[7-8]`)
	regexNonOdia, _  = regexp.Compile(`\P{Oriya}`)
	regexAlphaNum, _ = regexp.Compile(`[^\dA-Z]`)
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
	// Remove all non-odia characters.
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

	// Remove non-alphanumeric characters (losing the bracket grouping).
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

// This is for testing purpose
func main() {
	od := New()
	fmt.Println(od.Encode("ଭ୍ରମର"))
	fmt.Println(od.Encode("ଭ୍ରମରେ"))
	fmt.Println(od.Encode("ଭ୍ରମଣ"))
}
