package converter

import "testing"

func TestBuildReferenceLinks(t *testing.T) {
	t.Run("returns empty string for empty input", func(t *testing.T) {
		links := []Link{}
		expectedOutput := ""

		output := BuildReferenceLinks(links)

		if output != expectedOutput {
			t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectedOutput, output)
		}
	})

	t.Run("sorts links by ID", func(t *testing.T) {
		links := []Link{
			{ID: "3", URL: "https://www.example3.com"},
			{ID: "1", URL: "https://www.example1.com"},
			{ID: "2", URL: "https://www.example2.com"},
		}
		expectedOutput := "[1]: https://www.example1.com\n[2]: https://www.example2.com\n[3]: https://www.example3.com"

		output := BuildReferenceLinks(links)

		if output != expectedOutput {
			t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectedOutput, output)
		}
	})

	t.Run("separates links into numbered references, other references, and footnotes", func(t *testing.T) {
		links := []Link{
			{ID: "1", URL: "https://www.example1.com"},
			{ID: "^2", URL: "https://www.example2.com"},
			{ID: "3", URL: "https://www.example3.com"},
			{ID: "ref", URL: "https://www.example4.com"},
			{ID: "5", URL: "https://www.example5.com"},
		}
		expectedOutput := `[1]: https://www.example1.com
[3]: https://www.example3.com
[5]: https://www.example5.com
[ref]: https://www.example4.com
[^2]: https://www.example2.com`

		output := BuildReferenceLinks(links)

		if output != expectedOutput {
			t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectedOutput, output)
		}
	})

	t.Run("combines numbered references, other references, and footnotes into a single list", func(t *testing.T) {
		links := []Link{
			{ID: "1", URL: "https://www.example1.com"},
			{ID: "2", URL: "https://www.example2.com"},
			{ID: "3", URL: "https://www.example3.com"},
			{ID: "4", URL: "https://www.example4.com"},
			{ID: "^1", URL: "https://www.example5.com"},
		}
		expectedOutput := `[1]: https://www.example1.com
[2]: https://www.example2.com
[3]: https://www.example3.com
[4]: https://www.example4.com
[^1]: https://www.example5.com`

		output := BuildReferenceLinks(links)

		if output != expectedOutput {
			t.Errorf("Expected output:\n%s\n\nBut got:\n%s", expectedOutput, output)
		}
	})
}
