package converter

import (
	"fmt"
	"regexp"
	"strings"
)

type Link struct {
	Name        string
	URL         string
	ReferenceNo int
	ID          string
}

func (l *Link) IsFootnote() bool {
	return strings.HasPrefix(l.ID, "^")
}

func (l *Link) IsReference() bool {
	match, _ := regexp.MatchString(`^\D+$`, l.ID)
	return match && !l.IsFootnote()
}

func (l *Link) AsReference() string {
	var ref string
	if l.ID != "" {
		ref = l.ID
	} else {
		ref = fmt.Sprint(l.ReferenceNo)
	}
	return fmt.Sprintf("[%s]: %s", ref, l.URL)
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
