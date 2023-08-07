package converter

import (
	"regexp"
	"testing"
)

func TestIsFootnote(t *testing.T) {
	link := Link{ID: "^1"}
	if !link.IsFootnote() {
		t.Errorf("Expected IsFootnote() to return true, but got false")
	}
}

func TestIsReference(t *testing.T) {
	link := Link{ID: "abc"}
	if !link.IsReference() {
		t.Errorf("Expected IsReference() to return true, but got false")
	}

	link.ID = "^1"
	if link.IsReference() {
		t.Errorf("Expected IsReference() to return false, but got true")
	}

	link.ID = "1"
	if link.IsReference() {
		t.Errorf("Expected IsReference() to return false, but got true")
	}
}

func TestAsReference(t *testing.T) {
	link := Link{ID: "abc", URL: "http://example.com"}
	expected := "[abc]: http://example.com"
	if link.AsReference() != expected {
		t.Errorf("Expected AsReference() to return %q, but got %q", expected, link.AsReference())
	}

	link.ID = "1"
	expected = "[1]: http://example.com"
	if link.AsReference() != expected {
		t.Errorf("Expected AsReference() to return %q, but got %q", expected, link.AsReference())
	}
}

func TestAsMarkdownLink(t *testing.T) {
	link := Link{URL: "http://example.com"}
	expected := "<http://example.com>"
	if link.AsMarkdownLink() != expected {
		t.Errorf("Expected AsMarkdownLink() to return %q, but got %q", expected, link.AsMarkdownLink())
	}

	link.URL = ""
	expected = ""
	if link.AsMarkdownLink() != expected {
		t.Errorf("Expected AsMarkdownLink() to return %q, but got %q", expected, link.AsMarkdownLink())
	}
}

func TestIsReferenceRegex(t *testing.T) {
	link := Link{ID: "abc"}
	match, _ := regexp.MatchString(`^\D+$`, link.ID)
	if !match {
		t.Errorf("Expected ID %q to match reference regex, but it did not", link.ID)
	}

	link.ID = "^1"
	match, _ = regexp.MatchString(`^\D+$`, link.ID)
	if match {
		t.Errorf("Expected ID %q to not match reference regex, but it did", link.ID)
	}
}
