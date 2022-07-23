# ODphone (WIP)

ODphone is a phonetic algorithm for indexing Odia words by their pronunciation, like Metaphone for English. The algorithm generates three Romanized phonetic keys (hashes) of varying phonetic affinities for a given Odia word. This package implements the algorithm in Go.

The algorithm takes into account the context sensitivity of sounds, syntactic and phonetic gemination, compounding, modifiers, and other known exceptions to produce Romanized phonetic hashes of increasing phonetic affinity that are very faithful to the pronunciation of the original Odia word.

- `key0` = a broad phonetic hash comparable to a Metaphone key that doesn't account for hard sounds and phonetic modifiers
- `key1` = is a slightly more inclusive hash that accounts for hard sounds.
- `key2` = highly inclusive and narrow hash that accounts for hard sounds and phonetic modifiers.

### Examples

| Word       | Pronunciation | key0    | key1    | key2      |
| ---------- | ------------- | ------- | ------- | --------- |
| ପାକସ୍ଥଳୀ | pākasthali   | A3KS3KY | A3KS3KY | A3K6S3KY6 |
| ଅଭିସାର    | abhīsāra       | URSR    | UR1SR1  | UR14SR15  |
| ಈರಿತ       | īrita         | IR0     | IR0     | IR40      |
| ಒನಮಾಲೆ     | onamāle       | ONML    | ONML    | ONML6     |

### Go implementation

Install the package:
`go get -u github.com/soumendrak/odphone`

```go
package main

import (
	"fmt"

	"github.com/soumendrak/odphone"
)

func main() {
	k := odphone.New()
	fmt.Println(k.Encode("ପାକସ୍ଥଳୀ"))
	fmt.Println(k.Encode("ଅଭିସାର"))
}

```

License: GPLv3
