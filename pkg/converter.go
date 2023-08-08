package converter

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
	refLinkRegex := regexp.MustCompile(`\[([^\]]*)?\]\[(\w+)\]`)
	refLinksMatches := refLinkRegex.FindAllSubmatch(content, -1)

	for _, match := range refLinksMatches {
		c.addLink(string(match[1]), "", string(match[2]))
	}

	inlineLinkRegex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	inlineLinksMatches := inlineLinkRegex.FindAllSubmatch(content, -1)

	for _, match := range inlineLinksMatches {
		c.addLink(string(match[1]), string(match[2]), "")
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
	if ID == "" {
		usedNumbers := make(map[int]bool)
		for _, l := range c.Links {
			if num, err := strconv.Atoi(l.ID); err == nil {
				usedNumbers[num] = true
			}
		}
		nextNumber := 1
		for usedNumbers[nextNumber] {
			nextNumber++
		}
		ID = strconv.Itoa(nextNumber)
	}

	link := Link{Name: name, URL: url, ID: ID}
	c.Links = append(c.Links, link)
}

func (c *MarkdownConverter) extractLinksFromReferences() {
	refLinkRegex := regexp.MustCompile(`\[(.*?)\]:\s(.+)`)
	matches := refLinkRegex.FindAllSubmatch(c.originalContent, -1)

	for _, match := range matches {
		c.addLink(string(""), string(match[2]), string(match[1]))
	}
}

func (c *MarkdownConverter) Run() {
	c.modifiedContent = c.originalContent
	c.extractLinksFromReferences()
	c.extractMarkdownLinksFromBuffer(c.modifiedContent)
	c.modifiedContent = cleanup(c.Links, c.modifiedContent)
	c.modifiedContent = append(c.modifiedContent, "\n"...)
	if len(c.Links) > 0 {
		c.modifiedContent = append(c.modifiedContent, "\n"...)
		c.modifiedContent = append(c.modifiedContent, []byte(BuildReferenceLinks(c.Links))...)
		c.modifiedContent = append(c.modifiedContent, "\n"...)
	}
}

func setupLogger(verbose bool) {
	log.SetOutput(io.Discard)
	if verbose {
		log.SetOutput(os.Stderr)
	}
}
func ConvertFilesInPath(path string, backup, verbose bool) {
	setupLogger(verbose)

	filepath.WalkDir(path, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing file %s: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			return err
		}
		mc := MarkdownConverter{originalContent: content}
		mc.Run()
		newContent := mc.modifiedContent

		if bytes.Equal(content, newContent) {
			log.Printf("%s: Nothing to update\n", path)
			return nil
		}
		if backup {
			backupFile(path)
		}
		err = os.WriteFile(path, newContent, 0644)

		if err != nil {
			fmt.Printf("Error updating file %s: %v\n", path, err)
			return err
		}
		log.Printf("%s updated successfully!\n", path)

		return nil
	})
	fmt.Printf("Completed!\n")
}

func backupFile(filename string) error {
	backupFilename := filename + ".bak"
	_, err := os.Stat(backupFilename)
	if err == nil {
		return fmt.Errorf("backup file already exists: %s", backupFilename)
	}
	if !os.IsNotExist(err) {
		return err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = os.WriteFile(backupFilename, data, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Backup created: %s\n", backupFilename)
	return nil
}
