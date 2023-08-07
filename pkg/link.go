package converter

import (
	"fmt"
	"regexp"
	"strings"
)

type Link struct {
	Name string
	URL  string
	ID   string
}

func (l *Link) IsFootnote() bool {
	return strings.HasPrefix(l.ID, "^")
}

func (l *Link) IsReference() bool {
	match, _ := regexp.MatchString(`^\D+$`, l.ID)
	return match && !l.IsFootnote()
}

func (l *Link) AsReference() string {
	return fmt.Sprintf("[%s]: %s", l.ID, l.URL)
}

func (l *Link) AsMarkdownLink() string {
	if l.URL == "" {
		return l.Name
	}
	if l.Name == "" {
		return fmt.Sprintf("<%s>", l.URL)
	}
	return fmt.Sprintf("[%s](%s)", l.Name, l.URL)
}
