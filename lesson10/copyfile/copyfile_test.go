package copyfile

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tests = []struct {
	fromFileName, toFileName string
	offset, limit            int
	err                      bool
	fromFileData, toFileData []byte
}{
	{
		offset: -1,
		err:    true,
	},
	{
		limit: -1,
		err:   true,
	},
	{
		fromFileName: "notfound.txt",
		err:          true,
	},
	{
		fromFileName: "found.txt",
		fromFileData: []byte("testbob"),
		toFileName:   "/root/1.txt",
		err:          true,
	},
	{
		fromFileName: "found.txt",
		fromFileData: []byte("testbob"),
		toFileName:   "/dev/null",
		offset:       100,
		err:          true,
	},
	{
		fromFileName: "found.txt",
		fromFileData: []byte("testbob"),
		toFileName:   "/dev/null",
		limit:        100,
		err:          true,
	},
	{
		fromFileName: "found.txt",
		fromFileData: []byte("testbob"),
		toFileName:   "/dev/null",
		offset:       5,
		limit:        5,
		err:          true,
	},
	{
		fromFileName: "found.txt",
		fromFileData: []byte("testbob"),
		toFileName:   "/dev/null",
		err:          false,
	},
	{
		fromFileName: "found.txt",
		fromFileData: []byte("testbob"),
		toFileName:   "/dev/null",
		offset:       1,
		limit:        1,
		err:          false,
	},
	{
		fromFileName: "found.txt",
		fromFileData: []byte("testbob"),
		toFileName:   "result.txt",
		toFileData:   []byte("testbob"),
		err:          false,
	},
	{
		fromFileName: "found.txt",
		fromFileData: []byte("testbobasd345"),
		toFileName:   "result.txt",
		toFileData:   []byte("stboba"),
		offset:       2,
		limit:        6,
		err:          false,
	},
}

func TestCopy(t *testing.T) {
	for _, data := range tests {
		t.Run("", func(t *testing.T) {
			if data.fromFileData != nil {
				err := ioutil.WriteFile(data.fromFileName, data.fromFileData, 0600)
				require.NoError(t, err)
			}
			err := Copy(data.fromFileName, data.toFileName, data.limit, data.offset)
			if data.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if data.toFileData != nil {
				buffer, err := ioutil.ReadFile(data.toFileName)
				require.NoError(t, err)
				assert.Equal(t, data.toFileData, buffer)
			}
		})
	}
}
