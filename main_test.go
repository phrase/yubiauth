package yubioath

import (
	"testing"
)

func TestParseOutput(t *testing.T) {
	c := parseOutput([]byte(exampleOutput))
	tests := []struct{ Has, Want interface{} }{
		{len(c), 3},
		{c["Other Type"], "345678"},
		{c["Key 1"], "123456"},
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
