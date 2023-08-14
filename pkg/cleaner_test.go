package converter

import (
	"bytes"
	"testing"
)

func TestCleanup(t *testing.T) {
	t.Run("removes reference links", func(t *testing.T) {
		links := []Link{
			{ID: "1", URL: "https://www.example1.com"},
			{ID: "2", URL: "https://www.example2.com"},
			{ID: "3", URL: "https://www.example3.com"},
		}
		content := []byte(`This is some text with a reference [link][1].
    [1]: https://www.example1.com`)

		expectedOutput := []byte(`This is some text with a reference [link][1].`)

		output := cleanup(links, content)

		compareResults(output, expectedOutput, t)
	})

	t.Run("replaces inline links with reference links", func(t *testing.T) {
		links := []Link{
			{ID: "1", URL: "https://www.google.com"},
		}
		content := []byte(`This is some text with an inline link to [Google](https://www.google.com).`)

		expectedOutput := []byte(`This is some text with an inline link to [Google][1].`)

		output := cleanup(links, content)

		compareResults(output, expectedOutput, t)

	})

	t.Run("removes duplicated empty lines", func(t *testing.T) {
		links := []Link{}
		content := []byte(`This is some text with an empty line.


This is some more text.`)

		expectedOutput := []byte(`This is some text with an empty line.

This is some more text.`)

		output := cleanup(links, content)

		compareResults(output, expectedOutput, t)

	})

	t.Run("removes trailing whitespace", func(t *testing.T) {
		links := []Link{}
		content := []byte(`This is some text with trailing whitespace.

This is some more text.`)

		expectedOutput := []byte(`This is some text with trailing whitespace.

This is some more text.`)

		output := cleanup(links, content)

		compareResults(output, expectedOutput, t)

	})

	t.Run("removes footnote links", func(t *testing.T) {
		links := []Link{
			{ID: "1", URL: "https://www.example1.com"},
			{ID: "2", URL: "https://www.example2.com"},
			{ID: "3", URL: "https://www.example3.com"},
		}
		content := []byte(`This is some text with a footnote link[^1].
[^1]: https://www.example1.com`)

		expectedOutput := []byte(`This is some text with a footnote link[^1].`)

		output := cleanup(links, content)

		if !bytes.Equal(output, expectedOutput) {
			t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectedOutput, output)
		}
	})
}

func TestRemoveLineContainingString(t *testing.T) {
	content := []byte(`
		This is a test file.
		It has multiple lines.
		Some lines contain the word "test".
		This line should be removed because of test.
		This line should also be removed because... test.
		This line should stay.
	`)

	expectedOutput := []byte(`
		It has multiple lines.
		This line should stay.
	`)

	newContent := removeLineContainingString(content, "test")
	compareResults(newContent, expectedOutput, t)

}

func compareResults(output []byte, expectedOutput []byte, t *testing.T) {
	t.Helper()
	if !bytes.Equal(output, expectedOutput) {
		t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectedOutput, output)
	}
}
