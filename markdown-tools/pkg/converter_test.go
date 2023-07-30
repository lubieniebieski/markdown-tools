package converter

import (
	"bytes"
	"os"
	"testing"
)

func TestAddLink(t *testing.T) {
	converter := MarkdownConverter{}
	converter.addLink("Google", "https://www.google.com", "ref")

	expectedLinks := []Link{
		{Name: "Google", URL: "https://www.google.com", ReferenceNo: 1, ID: "ref"},
	}

	if len(converter.Links) != len(expectedLinks) {
		t.Errorf("Expected %d links, but got %d", len(expectedLinks), len(converter.Links))
	}

	for i, link := range converter.Links {
		if link.Name != expectedLinks[i].Name {
			t.Errorf("Expected link name '%s', but got '%s'", expectedLinks[i].Name, link.Name)
		}
		if link.URL != expectedLinks[i].URL {
			t.Errorf("Expected link URL '%s', but got '%s'", expectedLinks[i].URL, link.URL)
		}
		if link.ReferenceNo != expectedLinks[i].ReferenceNo {
			t.Errorf("Expected link reference number '%d', but got '%d'", expectedLinks[i].ReferenceNo, link.ReferenceNo)
		}
		if link.ID != expectedLinks[i].ID {
			t.Errorf("Expected link ID '%s', but got '%s'", expectedLinks[i].ID, link.ID)
		}
	}
}

func TestExtractMarkdownLinksFromBuffer(t *testing.T) {
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
			{Name: "Google", URL: "https://www.google.com"},
			{Name: "GitHub", URL: "https://github.com", ID: "1"},
			{Name: "Wikipedia", URL: "https://www.wikipedia.org", ID: "ref"},
			{Name: "Example page", URL: "https://example.com", ID: "Example"},
		}

		converter := MarkdownConverter{}
		converter.extractMarkdownLinksFromBuffer(content)

		if len(converter.Links) != len(expectedLinks) {
			t.Errorf("Expected %d links, but got %d", len(expectedLinks), len(converter.Links))
		}

		for i, link := range converter.Links {
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
			{Name: "Google", URL: "https://www.google.com"},
			{Name: "GitHub", URL: "https://github.com", ID: "1"},
			{Name: "", URL: "some footnote", ID: "^1"},
		}

		converter := MarkdownConverter{}
		converter.extractMarkdownLinksFromBuffer(mixedContent)

		if len(converter.Links) != len(expectedLinks) {
			t.Errorf("Expected %d links, but got %d", len(expectedLinks), len(converter.Links))
		}

		for i, link := range converter.Links {
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

	if !bytes.Equal(newContent, expectedOutput) {
		t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectedOutput, newContent)
	}
}
func TestRunOnContent(t *testing.T) {
	content := []byte(`[Google](https://www.google.com) fdafd
[GitHub][1]
[Wikipedia][ref] fdsf ds
[Example page][Example]
[Invalid Link]
[1]: https://github.com
[ref]: https://www.wikipedia.org
[Example]: https://example.com`)

	expectedOutput := []byte(`[Google][4] fdafd
[GitHub][1]
[Wikipedia][ref] fdsf ds
[Example page][Example]
[Invalid Link]

[1]: https://github.com
[ref]: https://www.wikipedia.org
[Example]: https://example.com
[4]: https://www.google.com`)

	converter := MarkdownConverter{}
	converter.RunOnContent(content)

	if !bytes.Equal(converter.modifiedContent, expectedOutput) {
		t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectedOutput, converter.modifiedContent)
	}
}

func TestRun(t *testing.T) {
	content := []byte(`[Google](https://www.google.com) fdafd
[GitHub][1]
[Wikipedia][ref] fdsf ds
[Example page][Example]
[Invalid Link]
[1]: https://github.com
[ref]: https://www.wikipedia.org
[Example]: https://example.com`)

	expectedOutput := []byte(`[Google][4] fdafd
[GitHub][1]
[Wikipedia][ref] fdsf ds
[Example page][Example]
[Invalid Link]

[1]: https://github.com
[ref]: https://www.wikipedia.org
[Example]: https://example.com
[4]: https://www.google.com`)

	filename := "test.md"
	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	defer os.Remove(filename)

	Run(filename)

	newContent, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if !bytes.Equal(newContent, expectedOutput) {
		t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectedOutput, newContent)
	}
}
