package envcfg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vporoshok/envcfg"
)

func TestSplitWords(t *testing.T) {
	cases := []struct {
		name string
		src  string
		res  []string
	}{
		{
			name: "camelCase",
			src:  "someCamelCase",
			res:  []string{"some", "Camel", "Case"},
		},
		{
			name: "snake_case",
			src:  "some_snake_case",
			res:  []string{"some", "snake", "case"},
		},
		{
			name: "abbr",
			src:  "JSONFileAndSOMEMore",
			res:  []string{"JSON", "File", "And", "SOME", "More"},
		},
		{
			name: "numbers",
			src:  "JSON42File42An42dSO42MEMore",
			res:  []string{"JSON42", "File42", "An42d", "SO42ME", "More"},
		},
		{
			name: "mix",
			src:  "JSON42_File42An42dSO42ME_More",
			res:  []string{"JSON42", "File42", "An42d", "SO42ME", "More"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := envcfg.SplitWords(c.src)
			assert.EqualValues(t, c.res, res)
		})
	}
}
