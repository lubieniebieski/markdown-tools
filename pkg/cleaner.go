package converter

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode"
)

func cleanup(links []Link, content []byte) []byte {
	refLinkRegex := regexp.MustCompile(`[^\]]\[((\w+)|(\^\w+))\]:\s(.+)`)
	content = refLinkRegex.ReplaceAll(content, []byte(""))
	for _, link := range links {
		if link.IsFootnote() || link.IsReference() {
			content = removeLineContainingString(content, link.AsReference())
		} else {
			linkRef := fmt.Sprintf("[%s]", link.ID)
			linkRegex := regexp.MustCompile(fmt.Sprintf(`\(%s\)`, link.URL))
			content = linkRegex.ReplaceAll(content, []byte(linkRef))
		}
	}

	// Remove all empty lines if there is more than one in a row
	content = regexp.MustCompile(`(\n){3,}`).ReplaceAll(content, []byte("\n\n"))
	// // Remove the last line if it's empty
	content = bytes.TrimRightFunc(content, unicode.IsSpace)
	return content
}

func removeLineContainingString(buffer []byte, str string) []byte {
	lines := bytes.Split(buffer, []byte("\n"))
	var newLines [][]byte
	for _, line := range lines {
		if !bytes.Contains(line, []byte(str)) {
			newLines = append(newLines, line)
		}
	}
	return bytes.Join(newLines, []byte("\n"))
}
