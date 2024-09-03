package markdown

import (
	"regexp"
	"strings"
)

type LineType string

func getLineType(line string) LineType {
	trimmedLine := strings.TrimSpace(line)

	re := regexp.MustCompile(`^- \[.\] `)
	if re.MatchString(trimmedLine) {
		return Task
	}

	if strings.HasPrefix(trimmedLine, string(Code)) {
		return Code
	} else if strings.HasPrefix(trimmedLine, string(Image)) {
		return Image
	} else if strings.HasPrefix(trimmedLine, string(List)) {
		return List
	} else if strings.HasPrefix(trimmedLine, string(Quote)) {
		return Quote
	} else if strings.HasPrefix(trimmedLine, string(Table)) {
		return Table
	}

	re = regexp.MustCompile(`^\d+\.`)
	if re.MatchString(trimmedLine) {
		return NumberedList
	}

	return Normal
}

var (
	Normal       LineType = ""
	Code         LineType = "```"
	Image        LineType = "!["
	List         LineType = "- "
	Quote        LineType = "> "
	Table        LineType = "| "
	Task         LineType = "- [%] "
	NumberedList LineType = "%. "
)

type Line struct {
	Text         string // The original text, with prefix
	LineType     LineType
	originalText string
}

// Sections
type SectionType string

var (
	H1          SectionType = "#"
	H2          SectionType = "##"
	H3          SectionType = "###"
	H4          SectionType = "####"
	H5          SectionType = "#####"
	H6          SectionType = "######"
	NullSection SectionType = ""
)

func getSectionType(line string) SectionType {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, "#") {
		sectionType := SectionType(trimmedLine[:strings.Index(trimmedLine, " ")])
		return sectionType
	} else {
		return NullSection
	}
}

// A "section" is a part of the markdown file that starts with any of the H1-H6 headers
// A NullSection is a section that has no "section" (It's the first part of the file if no title at the beginning)
type Section struct {
	SectionType  SectionType
	Text         string // The text after the first space (Without the "#..")
	Lines        []Line
	originalText string
}

type MarkdownFile struct {
	Path        string
	Title       string // First "H1" section
	FrontMatter map[interface{}]interface{}
	Sections    []Section
}

func New(path string) MarkdownFile {
	return MarkdownFile{Path: path}
}
