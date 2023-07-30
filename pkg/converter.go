package converter

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Link represents a link along with its reference number

// MarkdownConverter converts inline links to reference links in markdown files
type MarkdownConverter struct {
	originalContent []byte
	modifiedContent []byte
	Links           []Link
}

func (c *MarkdownConverter) extractFootnotesFromBuffer(content []byte) {
	footnoteRegex := regexp.MustCompile(`\[(\^\d+)\]:\s(.+)`)
	matches := footnoteRegex.FindAllSubmatch(content, -1)

	for _, match := range matches {
		c.addLink("", string(match[2]), string(match[1]))
	}
}
func (c *MarkdownConverter) extractMarkdownLinksFromBuffer(content []byte) {
	inlineLinkRegex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	matches := inlineLinkRegex.FindAllSubmatch(content, -1)

	for _, match := range matches {
		c.addLink(string(match[1]), string(match[2]), "")
	}

	refLinkRegex := regexp.MustCompile(`\[([^\]]*)?\]\[(\w+)\]`)
	matches = refLinkRegex.FindAllSubmatch(content, -1)

	for _, match := range matches {
		c.addLink(string(match[1]), "", string(match[2]))
	}
	c.extractReferenceLinksFromBuffer(content)
	c.extractFootnotesFromBuffer(content)
}
func (c *MarkdownConverter) extractReferenceLinksFromBuffer(content []byte) {
	refLinkRegex := regexp.MustCompile(`\[(.*?)\]:\s(.+)`)
	matches := refLinkRegex.FindAllSubmatch(content, -1)

	for _, match := range matches {
		matchID := string(match[1])
		matchURL := string(match[2])

		for i := range c.Links {
			if c.Links[i].ID == matchID {
				c.Links[i].URL = matchURL
				break
			}
		}
	}
}

func (c *MarkdownConverter) addLink(name string, url string, ID string) {
	if url != "" {
		for _, link := range c.Links {
			if link.URL == url {
				return
			}
		}
	}

	if ID != "" {
		for _, link := range c.Links {
			if link.ID == ID {
				return
			}
		}
	}

	link := Link{Name: name, URL: url, ID: ID}
	if !link.IsFootnote() {
		link.ReferenceNo = len(c.Links) + 1
	}
	c.Links = append(c.Links, link)
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

func (c *MarkdownConverter) cleanup() {
	for _, link := range c.Links {
		if link.IsFootnote() || link.IsReference() {
			c.modifiedContent = removeLineContainingString(c.modifiedContent, link.AsReference())
		} else {
			linkRef := fmt.Sprintf("[%d]", link.ReferenceNo)
			linkRegex := regexp.MustCompile(fmt.Sprintf(`\(%s\)`, link.URL))
			c.modifiedContent = linkRegex.ReplaceAll(c.modifiedContent, []byte(linkRef))
		}
	}
}

func (c *MarkdownConverter) addNewReferencesList() {

	if len(c.referencesList()) == 0 {
		return
	}
	c.modifiedContent = append(c.modifiedContent, []byte("\n\n")...)
	c.modifiedContent = append(c.modifiedContent, []byte(strings.Join(c.referencesList(), "\n"))...)
}

func (c *MarkdownConverter) referencesList() []string {
	var result []string
	for _, link := range c.Links {
		result = append(result, link.AsReference())
	}
	return result
}

func (c *MarkdownConverter) extractLinksFromReferences() {
	refLinkRegex := regexp.MustCompile(`\[(.*?)\]:\s(.+)`)
	matches := refLinkRegex.FindAllSubmatch(c.originalContent, -1)

	for _, match := range matches {
		c.addLink(string(""), string(match[2]), string(match[1]))
	}
}

func (c *MarkdownConverter) clearReferences() {
	refLinkRegex := regexp.MustCompile(`\[(.*?)\]:\s(.+)\n`)
	c.modifiedContent = refLinkRegex.ReplaceAll(c.originalContent, []byte(""))
}

func (c *MarkdownConverter) RunOnContent(content []byte) {
	c.originalContent = content
	c.extractLinksFromReferences()
	c.clearReferences()
	c.extractMarkdownLinksFromBuffer(c.modifiedContent)
	c.cleanup()
	c.addNewReferencesList()
}

func (c *MarkdownConverter) Run() {
	c.extractLinksFromReferences()
	c.clearReferences()
	c.extractMarkdownLinksFromBuffer(c.modifiedContent)
	c.cleanup()
	c.addNewReferencesList()
}

func RunOnContent(content []byte) (modifiedContent []byte) {
	converter := MarkdownConverter{originalContent: content}
	converter.Run()
	return converter.modifiedContent

}

func Run(filename string) {
	content, _ := os.ReadFile(filename)
	newContent := RunOnContent(content)

	err := os.WriteFile(filename, newContent, 0644)

	if err != nil {
		fmt.Printf("Error updating file: %v\n", err)
		return
	}

	fmt.Printf("File %s updated successfully!\n", filename)

}
