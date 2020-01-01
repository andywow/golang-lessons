package funcrunner

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func testSuccess() error {
	fmt.Println("Func ended - success")
	return nil
}

func testFailed() error {
	fmt.Println("Func ended - error")
	return fmt.Errorf("Test func failed")
}

var testData = []struct {
	testName  string
	functions []func() error
	N, M      int
	err       bool
}{
	{
		testName:  "successTest - tasks > N",
		functions: []func() error{testSuccess, testFailed, testSuccess, testSuccess, testSuccess},
		N:         2,
		M:         2,
		err:       false,
	},
	{
		testName:  "successTest - tasks < N",
		functions: []func() error{testSuccess, testFailed, testSuccess, testSuccess, testSuccess},
		N:         10,
		M:         2,
		err:       false,
	},
	{
		testName:  "failedTest - tasks > N",
		functions: []func() error{testSuccess, testFailed, testFailed, testSuccess, testSuccess},
		N:         2,
		M:         2,
		err:       true,
	},
	{
		testName: "failedTest - task < N",
		functions: []func() error{
			testSuccess, testFailed, testSuccess, testSuccess, testFailed,
			testSuccess, testFailed, testSuccess, testSuccess, testFailed,
			testSuccess, testFailed, testSuccess, testSuccess, testFailed,
			testSuccess, testFailed, testSuccess, testSuccess, testFailed},
		N:   10,
		M:   2,
		err: true,
	},
}

func TestRun(t *testing.T) {
	for _, data := range testData {
		t.Run(data.testName, func(t *testing.T) {
			fmt.Println("Test started: ", data.testName)
			startGoroutines := runtime.NumGoroutine()
			err := Run(data.functions, data.N, data.M)
			if data.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, startGoroutines, runtime.NumGoroutine())
		})
	}

}
