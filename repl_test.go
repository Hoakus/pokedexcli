package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  clefairy gardevoir",
			expected: []string{"clefairy", "gardevoir"},
		},
		{
			input:    "piKAchu        BULBAsaur SQUIRTLE",
			expected: []string{"pikachu", "bulbasaur", "squirtle"},
		},
		{
			input:    "  BLASTOISE jirachi",
			expected: []string{"blastoise", "jirachi"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("result: %s |expected: %s", word, expectedWord)
			}
		}
	}

}

func TestCommandMapB(t *testing.T) {
	previous := "https://pokeapi.co/api/v2/location-area/"
	c := config{Next: nil, Previous: &previous}

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	commandMap(&c)

	w.Close()

	os.Stdout = originalStdout

	outBytes, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	r.Close()

	outputStr := string(outBytes)

	expected := "you're on the last page"
	if !strings.Contains(outputStr, expected) {
		t.Errorf("Expected stdout to contain %q, got %q", expected, outputStr)
	}
	t.Log(outputStr)
}
