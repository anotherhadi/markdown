package markdown

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

const chunkSize = 64000

func deepCompare(file1, file2 string) bool {
	// Check file size ...

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func printMarkdown(md MarkdownFile) {
	fmt.Println("Title:", md.Title)
	fmt.Println()

	for s, section := range md.Sections {
		fmt.Println("Section n", s, " :", section.SectionType, section.Text)
		for l, line := range section.Lines {
			fmt.Println("Line n", l, " (", line.LineType, ") :", line.Text)
		}
		fmt.Println()
	}
}
