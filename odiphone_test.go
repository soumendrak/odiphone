package odiphone

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testVal struct {
	word     string
	expected expected
}

type expected struct {
	val1, val2, val3 string
}

func TestODIPhone(t *testing.T) {
	phone := New()
	testStrings := []testVal{
		{
			word: "ଅଂଶ",
			expected: expected{"VRMR",
				"V2RMR",
				"V2RMR",
			},
		},
		{
			word: "ଭ୍ରମରେ",
			expected: expected{
				"VRMR",
				"V2RMR3",
				"V2RMR3",
			},
		},
		{
			word: "ଭ୍ରମଣ",
			expected: expected{
				"VRMN",
				"V2RMN",
				"V2RMN",
			},
		},
	}
	for _, v := range testStrings {
		out1, out2, out3 := phone.Encode(v.word)
		require.Equal(t, v.expected.val1, out1)
		require.Equal(t, v.expected.val2, out2)
		require.Equal(t, v.expected.val3, out3)
	}
}
