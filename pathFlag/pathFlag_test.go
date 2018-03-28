package pathFlag

import "testing"

func TestPathFlagInterface(t *testing.T) {
	tests := []struct {
		name           string
		flagValue      []string
		requiredOutput map[string]string
		errors         bool
	}{
		{
			name:      "good 1 value",
			flagValue: []string{"/_status:GET /_status"},
			requiredOutput: map[string]string{
				"/_status": "GET /_status",
			},
			errors: false,
		},
		{
			name: "good 3 values",
			flagValue: []string{
				"/_status:GET /_status",
				"/_test1:GET /_test1",
				"/_test2:GET /_test2",
			},
			requiredOutput: map[string]string{
				"/_status": "GET /_status",
				"/_test1":  "GET /_test1",
				"/_test2":  "GET /_test2",
			},
			errors: false,
		},
		{
			name:      "bad 1 value",
			flagValue: []string{"/_status:GET:/_status"},
			requiredOutput: map[string]string{
				"nothing": "nothing",
			},
			errors: true,
		},
	}
	for _, test := range tests {
		p := new(PF)
		for _, fv := range test.flagValue {
			p.Set(fv)
		}
		m, err := p.Split()
		if err != nil {
			// If we get an error and expect it return ok
			if test.errors {
				t.Log("Test:", test.name, "Errored:", err)
				return
			}
			// Else fail the run
			t.Errorf("Test: %s, got err: %s", test.name, err)
			return
		}

		// Run through the table and check the values.
		// They should all match
		for k, v := range m {
			if test.requiredOutput[k] != v {
				t.Errorf("Test: %s, value for %s is missing or incorrect from returned map. Got '%s', Want '%s'",
					test.name,
					k,
					v,
					test.requiredOutput[k],
				)
			}
		}
	}
}

func TestFlagStringer(t *testing.T) {
	tests := []struct {
		name           string
		flagValue      []string
		requiredOutput string
		errors         bool
	}{
		{
			name:           "good 1 value",
			flagValue:      []string{"/_status:GET /_status"},
			requiredOutput: "/_status:GET /_status",
		},
		{
			name: "good 3 values",
			flagValue: []string{
				"/_status:GET /_status",
				"/_test1:GET /_test1",
				"/_test2:GET /_test2",
			},
			requiredOutput: "/_status:GET /_status,/_test1:GET /_test1,/_test2:GET /_test2",
		},
		{
			name:           "bad 1 value",
			flagValue:      []string{"/_status:GET:/_status"},
			requiredOutput: "/_status:GET:/_status",
		},
	}

	for _, test := range tests {
		p := new(PF)
		for _, fv := range test.flagValue {
			p.Set(fv)
		}
		testString := p.String()
		if testString != test.requiredOutput {
			t.Errorf("Test: %s, did not match expected output.\nGot: %s\nWant: %s", test.name, testString, test.requiredOutput)
		}
	}
}
