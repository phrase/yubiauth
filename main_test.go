package yubiauth

import (
	"context"
	"fmt"
	"testing"
)

func TestParseOutput(t *testing.T) {
	c := parseOutput([]byte(exampleOutput))
	tests := []struct{ Has, Want interface{} }{
		{len(c), 3},
		{c["Other Type"], 345678},
		{c["Key 1"], 123456},
	}
	for i, tc := range tests {
		if tc.Has != tc.Want {
			t.Errorf("%d: want=%#v has=%#v", i+1, tc.Want, tc.Has)
		}
	}
}

type testCommanderResult struct {
	Output string
	Error  error
}

type testCommanderResults []*testCommanderResult

func testCommander(results testCommanderResults) func(cmd string, args ...string) ([]byte, error) {
	i := 0
	return func(cmd string, args ...string) ([]byte, error) {
		if len(results) > i {
			r := results[i]
			i++
			return []byte(r.Output), r.Error
		}
		return nil, fmt.Errorf("no output configured for run %d", i)
	}
}

func TestWaitFor(t *testing.T) {
	cmder := testCommander(testCommanderResults{
		{Output: string(msgYubiKeyNotFound), Error: fmt.Errorf("no yubikey")},
		{Output: exampleOutput},
	})
	ctx := context.Background()
	v, err := waitForKeysWithCommander(ctx, cmder)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct{ Has, Want interface{} }{
		{len(v), 3},
		{v["Key 1"], 123456},
		{v["Key 1 Private"], 234567},
	}
	for i, tc := range tests {
		if tc.Has != tc.Want {
			t.Errorf("%d: want=%#v has=%#v", i+1, tc.Want, tc.Has)
		}
	}
}

const exampleOutput = `Key 1                                           123456
Key 1 Private                                   234567
Other Type                                      345678
`
