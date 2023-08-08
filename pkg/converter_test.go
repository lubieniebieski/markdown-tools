package converter

import (
	"bytes"
	"os"
	"testing"
)

func TestAddLink(t *testing.T) {
	converter := MarkdownConverter{}
	converter.addLink("Google", "https://www.google.com", "ref")
	converter.addLink("Github", "https://www.github.com", "")

	expectedLinks := []Link{
		{Name: "Google", URL: "https://www.google.com", ID: "ref"},
		{Name: "Github", URL: "https://www.github.com", ID: "1"},
	}
	assertLinksEqual(t, converter.Links, expectedLinks)
}

func TestExtractMarkdownLinksFromBuffer(t *testing.T) {
	t.Run("extracts proper reference", func(t *testing.T) {
		content := []byte(`[Google][2]`)

		expectedLinks := []Link{
			{Name: "Google", URL: "", ID: "2"},
		}

		converter := MarkdownConverter{}
		converter.extractMarkdownLinksFromBuffer(content)

		assertLinksEqual(t, converter.Links, expectedLinks)

	})
	t.Run("works with inline links", func(t *testing.T) {
		content := []byte(`
		[Google](https://www.google.com) fdafd
		[GitHub][1]
		[Wikipedia][ref] fdsf ds
		[Example page][Example]
		[Invalid Link]
		[1]: https://github.com
		[ref]: https://www.wikipedia.org
		[Example]: https://example.com
	`)

		expectedLinks := []Link{
			{Name: "GitHub", URL: "https://github.com", ID: "1"},
			{Name: "Wikipedia", URL: "https://www.wikipedia.org", ID: "ref"},
			{Name: "Example page", URL: "https://example.com", ID: "Example"},
			{Name: "Google", URL: "https://www.google.com", ID: "2"},
		}

		converter := MarkdownConverter{}
		converter.extractMarkdownLinksFromBuffer(content)

		assertLinksEqual(t, converter.Links, expectedLinks)

	})
	t.Run("works with footnotes too", func(t *testing.T) {
		mixedContent := []byte(`
		[Google](https://www.google.com)
		[GitHub][1]
		footnote example[^1]
		[1]: https://github.com
		[^1]: some footnote
	`)
		expectedLinks := []Link{
			{Name: "GitHub", URL: "https://github.com", ID: "1"},
			{Name: "Google", URL: "https://www.google.com", ID: "2"},
			{Name: "", URL: "some footnote", ID: "^1"},
		}

		converter := MarkdownConverter{}
		converter.extractMarkdownLinksFromBuffer(mixedContent)

		assertLinksEqual(t, converter.Links, expectedLinks)

	})
}

func TestRunOnContent(t *testing.T) {

	t.Run("works with inline links", func(t *testing.T) {
		content := []byte(`[Google](https://www.google.com) fdafd
[GitHub][1]
[Wikipedia][ref] fdsf ds
[Third link](https://www.example3.com)
[Fourth link](https://www.example4.com)
[Example page][Example]
[Invalid Link]
[1]: https://github.com
[ref]: https://www.wikipedia.org
[Example]: https://example.com`)

		expectedOutput := []byte(`[Google][2] fdafd
[GitHub][1]
[Wikipedia][ref] fdsf ds
[Third link][3]
[Fourth link][4]
[Example page][Example]
[Invalid Link]

[1]: https://github.com
[2]: https://www.google.com
[3]: https://www.example3.com
[4]: https://www.example4.com
[Example]: https://example.com
[ref]: https://www.wikipedia.org
`)

		compareConvertResults(t, content, expectedOutput)

	})
	t.Run("doesn't add any new lines if there are some links already defined", func(t *testing.T) {
		content := []byte(`first line
	second line
[1]: https://github.com
`)

		expectedOutput := []byte(`first line
	second line

[1]: https://github.com
`)
		compareConvertResults(t, content, expectedOutput)
	})

	t.Run("doesn't remove empty lines between paragraphs", func(t *testing.T) {
		content := []byte(`first line

last line
`)

		expectedOutput := []byte(`first line

last line
`)
		compareConvertResults(t, content, expectedOutput)

	})

	t.Run("references are listed in order", func(t *testing.T) {
		content := []byte(`first line
	second line
[1]: https://github.com
`)

		expectedOutput := []byte(`first line
	second line

[1]: https://github.com
`)
		compareConvertResults(t, content, expectedOutput)

	})

	t.Run("doesn't remove last line in the file", func(t *testing.T) {
		content := []byte(`first line
second line
`)

		expectedOutput := []byte(`first line
second line
`)
		compareConvertResults(t, content, expectedOutput)
	})
}

func TestRun(t *testing.T) {
	content := []byte(`[Google](https://www.google.com) fdafd
[GitHub][1]
[Wikipedia][ref] fdsf ds
[Example page][Example]
[Invalid Link]
[1]: https://github.com
[Example]: https://example.com
[ref]: https://www.wikipedia.org`)

	expectedOutput := []byte(`[Google][2] fdafd
[GitHub][1]
[Wikipedia][ref] fdsf ds
[Example page][Example]
[Invalid Link]

[1]: https://github.com
[2]: https://www.google.com
[Example]: https://example.com
[ref]: https://www.wikipedia.org
`)

	filename := "test.md"
	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	defer os.Remove(filename)

	ConvertFilesInPath(filename, false)

	newContent, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if !bytes.Equal(newContent, expectedOutput) {
		t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectedOutput, newContent)
	}
}

func assertLinksEqual(t *testing.T, links []Link, expectedLinks []Link) {
	if len(links) != len(expectedLinks) {
		t.Errorf("Expected %d links, but got %d", len(expectedLinks), len(links))
	}

	for i, link := range links {
		if link.Name != expectedLinks[i].Name {
			t.Errorf("Expected link name '%s', but got '%s'", expectedLinks[i].Name, link.Name)
		}
		if link.URL != expectedLinks[i].URL {
			t.Errorf("Expected link URL '%s', but got '%s'", expectedLinks[i].URL, link.URL)
		}

		if link.ID != expectedLinks[i].ID {
			t.Errorf("Expected link ID '%s', but got '%s'", expectedLinks[i].ID, link.ID)
		}
	}
}

func compareConvertResults(t *testing.T, input []byte, expectations []byte) {
	converter := MarkdownConverter{originalContent: input}
	converter.Run()

	if !bytes.Equal(converter.modifiedContent, expectations) {
		t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectations, converter.modifiedContent)
	}
}
