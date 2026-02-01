package main

import "testing"

func TestCleanInput(t *testing.T) {
    cases := []struct {
		input    string
		expected []string
    } {
		{
			input: " hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "Test text  ",
			expected: []string{"test", "text"},
		},
		// add more cases here.
    }

    for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("ERR: Test failed!!\nlen(%v) != len(%v)", actual, c.expected)
			return
		}

		for i, _ := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("ERR: Test failed!!\n%v != %v", word, expectedWord)
				return
			}
		}
    }
}

func TestReplCommands(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
    } {
		{
			input: " hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "Test text  ",
			expected: []string{"test", "text"},
		},
		// add more cases here.
    }

    for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("ERR: Test failed!!\nlen(%v) != len(%v)", actual, c.expected)
			return
		}

		for i, _ := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("ERR: Test failed!!\n%v != %v", word, expectedWord)
				return
			}
		}
    }
}
