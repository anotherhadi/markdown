package markdown

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	md := New("./testfiles/empty.md")
	if md.Path != "./testfiles/empty.md" {
		t.Errorf("Expected './testfiles/empty.md', got %s", md.Path)
	}
}

func TestRead(t *testing.T) {
	md := New("./testfiles/empty.md")
	err := md.Read()
	if err != nil {
		t.Errorf("Error reading file 1: %s", err)
	}
	if len(md.Sections) != 0 {
		t.Errorf("Expected 0 sections, got %d", len(md.Sections))
	}
	if md.Title != "" {
		t.Errorf("Expected empty title, got %s", md.Title)
	}

	md = New("./testfiles/empty_section.md")
	err = md.Read()
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}
	if len(md.Sections[1].Lines) != 0 {
		t.Errorf("Expected 0 lines in section 1, got %d", len(md.Sections[1].Lines))
	}
	printMarkdown(md)

	md = New("./testfiles/lorem.md")
	err = md.Read()
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}
	if md.Title != "Lorem Ipsum" {
		t.Errorf("Expected 'Lorem Ipsum', got %s", md.Title)
	}
	if len(md.Sections) != 5 {
		t.Errorf("Expected 5 sections, got %d", len(md.Sections))
	}
	if md.Sections[1].Text != "Etiam" {
		t.Errorf("Expected 'Etiam', got %s", md.Sections[1].Text)
	}
	if md.Sections[2].Lines[3].LineType != List {
		t.Errorf("Expected List, got %s", md.Sections[2].Lines[3].LineType)
	}
	if md.Sections[2].Lines[6].LineType != List {
		t.Errorf("Expected List, got %s", md.Sections[2].Lines[6].LineType)
	}
	if md.Sections[2].Lines[8].LineType != Task {
		t.Errorf("Expected Task, got %s", md.Sections[2].Lines[8].LineType)
	}
	if md.Sections[2].Lines[13].LineType != NumberedList {
		t.Errorf("Expected Task, got %s", md.Sections[2].Lines[13].LineType)
	}

	md = New("./testfiles/start_with_no_section.md")
	err = md.Read()
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}
	if md.Title != "My title and first section" {
		t.Errorf("Expected 'My title and first section', got %s", md.Title)
	}
	if len(md.Sections) != 2 {
		t.Errorf("Expected 2 sections, got %d", len(md.Sections))
	}
	if md.Sections[0].Text != "" {
		t.Errorf("Expected '', got %s", md.Sections[1].Text)
	}
	if md.Sections[0].Lines[0].Text != "This is a NullSection SectionType because it has no 'section'." {
		t.Errorf("Expected text, got %s", md.Sections[1].Text)
	}

	md = New("./testfiles/metadata.md")
	err = md.Read()
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}
	if md.FrontMatter["author"] != "Another Hadi" {
		t.Errorf("Expected 'Another Hadi', got %s", md.FrontMatter["author"])
	}
	if md.FrontMatter["unknown"] != nil {
		t.Errorf("Expected nil, got %s", md.FrontMatter["unknown"])
	}
	if len(md.Sections) != 1 {
		t.Errorf("Expected 1s sections, got %d", len(md.Sections))
	}
}

func TestReadAndWrite(t *testing.T) {
	filenames := []string{
		"./testfiles/empty.md",
		"./testfiles/lorem.md",
		"./testfiles/empty_section.md",
		"./testfiles/metadata.md",
		"./testfiles/start_with_no_section.md",
		"./testfiles/hyprland.md",
	}

	for _, filename := range filenames {
		md := New(filename)
		err := md.Read()
		if err != nil {
			t.Errorf("Error reading file: %s", err)
		}

		err = md.Write(filename + ".tmp")
		if err != nil {
			t.Errorf("Error writing file: %s", err)
		}

		if !deepCompare(filename, filename+".tmp") {
			t.Errorf("Files are not the same: %s, %s", filename, filename+".tmp")
		} else {
			os.Remove(filename + ".tmp")
		}
	}
}

func TestAddSection(t *testing.T) {
	filename := "./testfiles/empty.md"
	md := New(filename)
	err := md.Read()
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}
	md.AddSection("# New Section")
	md.AddSection("## New Sub Section")
	if len(md.Sections) != 2 {
		t.Errorf("Expected 2 section, got %d", len(md.Sections))
	}
	if md.Sections[0].Text != "New Section" {
		t.Errorf("Expected 'New Section', got %s", md.Sections[0].Text)
	}
}

func TestAddSectionAtIndex(t *testing.T) {
	filename := "./testfiles/lorem.md"
	md := New(filename)
	err := md.Read()
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}

	md.AddSectionAtIndex("## New Section", 1)
	if len(md.Sections) != 6 {
		t.Errorf("Expected 6 sections, got %d", len(md.Sections))
	}
	if md.Sections[1].Text != "New Section" {
		t.Errorf("Expected 'New Section', got %s", md.Sections[1].Text)
	}

	err = md.Write(filename + ".tmp")
	if err != nil {
		t.Errorf("Error writing file: %s", err)
	}
}

func TestAddLine(t *testing.T) {
	filename := "./testfiles/empty.md"
	md := New(filename)
	err := md.Read()
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}
	md.AddSection("# Test 1")
	md.Sections[0].AddLine("New Line")
	if len(md.Sections[0].Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(md.Sections[0].Lines))
	}
	if md.Sections[0].Lines[0].Text != "New Line" {
		t.Errorf("Expected 'New Line', got %s", md.Sections[0].Lines[0].Text)
	}

	filename = "./testfiles/empty.md"
	md = New(filename)
	err = md.Read()
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}
	md.AddSection("")
	md.Sections[0].AddLine("New Line")
	if len(md.Sections[0].Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(md.Sections[0].Lines))
	}
	if md.Sections[0].Lines[0].Text != "New Line" {
		t.Errorf("Expected 'New Line', got %s", md.Sections[0].Lines[0].Text)
	}

	filename = "./testfiles/lorem.md"
	md = New(filename)
	err = md.Read()
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}
	md.Sections[0].AddLine("New Line")
	if md.Sections[0].Lines[3].Text != "New Line" {
		t.Errorf("Expected 'New Line', got %s", md.Sections[0].Lines[3].Text)
	}
}

func TestAddLineAtIndex(t *testing.T) {
	filename := "./testfiles/lorem.md"
	md := New(filename)
	err := md.Read()
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}

	md.Sections[0].AddLineAtIndex("New Line", 1)
	if md.Sections[0].Lines[1].Text != "New Line" {
		t.Errorf("Expected 'New Line', got %s", md.Sections[0].Lines[1].Text)
	}

	err = md.Write(filename + ".tmp")
	if err != nil {
		t.Errorf("Error writing file: %s", err)
	}
}
