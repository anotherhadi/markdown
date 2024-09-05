package markdown

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// Returns the value of the specified front matter key or the default value if the key does not exist
func (md MarkdownFile) GetFrontMatter(key string, defaultValue interface{}) (value interface{}) {
	if md.FrontMatter == nil {
		return defaultValue
	}
	if val, ok := md.FrontMatter[key]; ok {
		return val
	}
	return defaultValue
}

// GetSection returns the first section of the specified type with the specified text
func (md *MarkdownFile) GetSection(sectionType SectionType, text string) (section *Section) {
	for _, s := range md.Sections {
		if s.SectionType == sectionType && s.Text == text {
			return &s
		}
	}
	return nil
}

// SearchSection returns a list of sections that contain the search string (with fuzzy search)
func (md *MarkdownFile) SearchSection(searchString string) (sections *[]Section) {
	var foundSections []Section
	for _, s := range md.Sections {
		if fuzzy.Match(searchString, s.Text) {
			foundSections = append(foundSections, s)
		}
	}
	return &foundSections
}

// SearchSection returns a list of sections with specified section type that contain the search string (with fuzzy search)
func (md *MarkdownFile) SearchSectionWithType(searchString string, sectionType SectionType) (sections *[]Section) {
	var foundSections []Section
	for _, s := range md.Sections {
		if s.SectionType == sectionType {
			if fuzzy.Match(searchString, s.Text) {
				foundSections = append(foundSections, s)
			}
		}
	}
	return &foundSections
}
