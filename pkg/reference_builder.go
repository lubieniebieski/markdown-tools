package converter

import (
	"sort"
	"strconv"
	"strings"
)

func BuildReferenceLinks(links []Link) (output string) {

	if len(links) == 0 {
		return
	}

	var numberedRefs []string
	var otherRefs []string
	var footnotes []string

	// Sort links by ID
	sort.Slice(links, func(i, j int) bool {
		num1, err1 := strconv.Atoi(links[i].ID)
		num2, err2 := strconv.Atoi(links[j].ID)
		if err1 == nil && err2 == nil {
			return num1 < num2
		} else if err1 != nil && err2 != nil {
			return links[i].ID < links[j].ID
		} else {
			return err1 == nil
		}
	})

	// Separate links into numbered references, other references, and footnotes
	for _, link := range links {
		if link.IsFootnote() {
			footnotes = append(footnotes, link.AsReference())
		} else if link.IsReference() {
			otherRefs = append(otherRefs, link.AsReference())
		} else {
			numberedRefs = append(numberedRefs, link.AsReference())
		}
	}

	// Combine the three lists of links into a single list
	if len(numberedRefs) > 0 && len(otherRefs) > 0 {
		numberedRefs = append(numberedRefs, "")
	}
	allConvertedLinks := append(numberedRefs, otherRefs...)
	if len(allConvertedLinks) > 0 && len(footnotes) > 0 {
		allConvertedLinks = append(allConvertedLinks, "")
	}

	allConvertedLinks = append(allConvertedLinks, footnotes...)
	return strings.Join(allConvertedLinks, "\n")
}
