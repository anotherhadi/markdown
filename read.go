package markdown

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func (md *MarkdownFile) readFrontMatter(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	frontMatterLines := ""
	if line == "---" { // Front Matter
		for scanner.Scan() {
			line = scanner.Text()
			if line == "---" { // End of Front Matter
				break
			}
			frontMatterLines += line + "\n"
		}
	} else { // No Front Matter
		md.FrontMatter = nil
		return
	}

	// Parse Front Matter
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(frontMatterLines), &m)
	if err != nil {
		return err
	}

	md.FrontMatter = m
	return
}

// Read reads the markdown file and parses it into sections (also reads the Front Matter)
func (md *MarkdownFile) Read() (err error) {
	directory := filepath.Dir(md.Path)
	filename := filepath.Base(md.Path)
	currentDir := os.Getenv("PWD")
	err = os.Chdir(directory)
	if err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var currentLines []Line = []Line{}
	var currentSection Section = Section{SectionType: NullSection}
	var isFirstLine bool = true

	err = md.readFrontMatter(filename)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip the FrontMatter
		if line == "---" && isFirstLine {
			for {
				scanner.Scan()
				line = scanner.Text()
				if line == "---" {
					scanner.Scan()
					line = scanner.Text()
					break
				}
			}
		}
		if isFirstLine && line == "" {
			isFirstLine = false
			continue
		}
		isFirstLine = false

		trimmedLine := strings.TrimSpace(line)

		// if the line start with #, it's a new section
		if strings.HasPrefix(trimmedLine, "#") {
			sectionType := getSectionType(trimmedLine)
			if sectionType == H1 && md.Title == "" {
				md.Title = trimmedLine[strings.Index(trimmedLine, " ")+1:]
			}
			sectionText := trimmedLine[strings.Index(trimmedLine, " ")+1:] // The text after the first space
			if currentSection.SectionType != NullSection || len(currentLines) > 0 {
				md.Sections = append(md.Sections, currentSection)
			}
			currentSection = Section{SectionType: sectionType, Text: sectionText, originalText: line}
			currentLines = []Line{}
		} else { // Add to current section
			lineType := getLineType(trimmedLine)
			currentLines = append(currentLines, Line{Text: line, LineType: lineType, originalText: line})
			currentSection.Lines = currentLines
		}
	}
	if currentSection.SectionType != NullSection || len(currentLines) > 0 {
		md.Sections = append(md.Sections, currentSection)
	}

	// Go back to the original directory
	err = os.Chdir(currentDir)
	if err != nil {
		return err
	}
	return nil
}
