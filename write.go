package markdown

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// AddSection adds a new section to the end of the markdown file
func (md *MarkdownFile) AddSection(line string) {
	sectionType := getSectionType(line)
	trimmedLine := strings.TrimSpace(line)
	text := trimmedLine[strings.Index(trimmedLine, " ")+1:]
	md.Sections = append(md.Sections, Section{SectionType: sectionType, Text: text, originalText: line})
}

// AddSectionAtIndex adds a new section at the specified index
func (md *MarkdownFile) AddSectionAtIndex(line string, index int) {
	sectionType := getSectionType(line)
	trimmedLine := strings.TrimSpace(line)
	text := trimmedLine[strings.Index(trimmedLine, " ")+1:]
	section := Section{SectionType: sectionType, Text: text, originalText: line}
	md.Sections = append(md.Sections[:index], append([]Section{section}, md.Sections[index:]...)...)
}

// AddLine adds a new line to the end of the specified section
func (s *Section) AddLine(text string) {
	lineType := getLineType(text)
	s.Lines = append(s.Lines, Line{Text: text, LineType: lineType, originalText: text})
}

// AddLineAtIndex adds a new line at the specified index in the specified section
func (s *Section) AddLineAtIndex(text string, index int) {
	lineType := getLineType(text)
	line := Line{Text: text, LineType: lineType, originalText: text}
	s.Lines = append(s.Lines[:index], append([]Line{line}, s.Lines[index:]...)...)
}

// If a string is provided, write to that file. Otherwise, write to the original file
func (md *MarkdownFile) Write(str ...string) (err error) {
	newPath := ""
	if len(str) > 0 {
		newPath = str[0]
	} else {
		newPath = md.Path
	}

	directory := filepath.Dir(newPath)
	filename := filepath.Base(newPath)
	currentDir := os.Getenv("PWD")
	err = os.Chdir(directory)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the file

	if md.FrontMatter != nil {
		d, err := yaml.Marshal(md.FrontMatter)
		if err != nil {
			return err
		}
		_, err = file.WriteString("---\n" + string(d) + "---\n\n")
		if err != nil {
			return err
		}
	}

	for _, section := range md.Sections {
		if section.SectionType != NullSection {
			_, err = file.WriteString(section.originalText + "\n")
		}
		if err != nil {
			return
		}
		for _, line := range section.Lines {
			_, err = file.WriteString(line.originalText + "\n")
			if err != nil {
				return
			}
		}
	}

	// Go back to the original directory
	err = os.Chdir(currentDir)
	if err != nil {
		return err
	}
	return nil
}
