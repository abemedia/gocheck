package condense

import "testing"

func TestCollapseLine(t *testing.T) {
	tests := []struct {
		input  string
		result string
		ok     bool
	}{
		{
			input: `myStruct{
	A: 1,
	B: 2,
}`,
			result: "myStruct{A: 1, B: 2}",
			ok:     true,
		},
		{
			input: `	{
		A: 1,
		B: 2,
	},`,
			result: "	{A: 1, B: 2},",
			ok:     true,
		},
		{
			input:  `myStruct{A: 1, B: 2}`,
			result: "",
			ok:     false,
		},
	}
	for _, test := range tests {
		result, ok := collapseLine([]byte(test.input))
		if string(result) != test.result {
			t.Errorf("collapseLine(`%s`) = `%s`, want `%s`", test.input, result, test.result)
		}
		if ok != test.ok {
			t.Errorf("collapseLine(`%s`) ok = %v, want %v", test.input, ok, test.ok)
		}
	}
}
