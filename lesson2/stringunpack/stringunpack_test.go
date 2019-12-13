package stringunpack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var successTests = []struct {
	input, output string
	err           bool
}{
	{
		input:  "a4bc2d5e",
		output: "aaaabccddddde",
		err:    false,
	},
	{
		input:  "abcd",
		output: "abcd",
		err:    false,
	},
	{
		input:  "45",
		output: "",
		err:    true,
	},
	{
		input:  "qwe\\4\\5",
		output: "qwe45",
		err:    false,
	},
	{
		input:  "qwe\\45",
		output: "qwe44444",
		err:    false,
	},
	{
		input:  "qwe\\\\5",
		output: "qwe\\\\\\\\\\",
		err:    false,
	},
}

func TestUnpack(t *testing.T) {
	for _, testData := range successTests {
		t.Run(testData.input, func(t *testing.T) {
			output, err := Unpack(testData.input)
			if testData.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testData.output, output)
			}
		})
	}
}
