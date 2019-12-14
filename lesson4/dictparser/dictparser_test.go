package dictparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var successTests = []struct {
	input  string
	output []string
}{
	{
		input:  "",
		output: []string{},
	},
	{
		input:  "aa aa aa bb bb d d d d c",
		output: []string{"d", "aa", "bb", "c"},
	},
	{
		input:  "a b b c c c d d d d e e e e e g g g g g g h h h h h h h i i i i i i i i j j j j j j j j j j k k k k k k k k k k k k  ",
		output: []string{"k", "j", "i", "h", "g", "e", "d", "c", "b", "a"},
	},
}

func TestTop10(t *testing.T) {
	for _, testData := range successTests {
		t.Run(testData.input, func(t *testing.T) {
			output := Top10(testData.input)
			assert.Equal(t, testData.output, output)
		})
	}
}
