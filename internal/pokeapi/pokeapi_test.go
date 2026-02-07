package pokeapi

import "testing"

func TestGetLocationAreas(t *testing.T) {
	NewConfig()

	goodCases := []struct {
		name string
		input int
		expected []LocationArea
	} {
		{
			name: "valid FORWARD direction",
			input: FORWARD,
			expected: []LocationArea{},
		},
		{
			name: "valid FORWARD direction",
			input: FORWARD,
			expected: []LocationArea{},
		},
		{
			name: "valid FORWARD direction",
			input: FORWARD,
			expected: []LocationArea{},
		},
		{
			name: "valid BACK direction",
			input: BACK,
			expected: []LocationArea{},
		},
		{
			name: "valid BACK direction",
			input: BACK,
			expected: []LocationArea{},
		},
	}

	errCase := []struct {
		name          string
		input         int
		expectedError string
	}{
		{
			name:          "invalid direction returns error",
			input:         42,
			expectedError: "you didn't think this through, did you?",
		},
		{
			name:          "first page BACK returns error",
			input:         BACK,
			expectedError: "you're on the first page",
		},
		// add more cases here.
	}

	for _, c := range goodCases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := GetLocationAreas(c.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(actual) != 20 {
				t.Errorf("ERR: Test failed!!\nlen(actual) != len(expected) [20]")
				return
			}
			// further checks can be added here to compare actual and expected contents.
		})
	}

	for _, c := range errCase {
		t.Run(c.name, func(t *testing.T) {
			_, err := GetLocationAreas(c.input)
			if c.expectedError != "" {
				if err == nil {
					t.Errorf("expected error %q but got nil", c.expectedError)
					return
				}
				if err.Error() != c.expectedError {
					t.Errorf("expected error %q but got %q", c.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}
