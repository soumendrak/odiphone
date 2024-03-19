# ODIphone (WIP)

ODIphone is a phonetic algorithm for indexing Odia words by their pronunciation, like Metaphone for English. The algorithm generates three Romanized phonetic keys (hashes) of varying phonetic affinities for a given Odia word. This package implements the algorithm in Go.

The algorithm takes into account the context sensitivity of sounds, syntactic and phonetic gemination, compounding, modifiers, and other known exceptions to produce Romanized phonetic hashes of increasing phonetic affinity that are very faithful to the pronunciation of the original Odia word.

- `key0` = a broad phonetic hash comparable to a Metaphone key that doesn't account for hard sounds and phonetic modifiers
- `key1` = is a slightly more inclusive hash that accounts for hard sounds.
- `key2` = highly inclusive and narrow hash that accounts for hard sounds and phonetic modifiers.

### Examples

| Word   | Pronunciation | key0   | key1    | key2    |
| ------ | ------------- | ------ | ------- | ------- |
| ଅଂଶ    | ansha         | ASH    | ASH     | A7SH    |
| ଭ୍ରମର  | vramara       | BHRMR  | BH2RMR  | BH2RMR  |
| ଭ୍ରମରେ | vramarè       | BHRMR  | BH2RMR3 | BH2RMR3 |
| ଭ୍ରମଣ  | vramańa       | BHRMNH | BH2RMNH | BH2RMNH |


### Go implementation

Install the package:
`go get -u github.com/soumendrak/odiphone`

```go
package main

import (
	"fmt"

	"github.com/soumendrak/odiphone"
)

func main() {
	od := odiphone.New()
	fmt.Println(od.Encode("ଭ୍ରମର"))
	fmt.Println(od.Encode("ଭ୍ରମରେ"))
	fmt.Println(od.Encode("ଭ୍ରମଣ"))
}

```

License: GPLv3
