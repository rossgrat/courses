package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

func syncGuide(c *course, check bool) error {
	source := filepath.Join(c.dir, c.Source)
	target := filepath.Join(c.dir, c.Guide)
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	file, reader, err := pdf.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()
	var output strings.Builder
	fmt.Fprintf(&output, "Source: %s\nPages: %d\nSHA256: %x\nText and link destinations extracted without changing the source.\n", c.Source, reader.NumPage(), sha256.Sum256(data))
	for index := 1; index <= reader.NumPage(); index++ {
		page := reader.Page(index)
		content, err := page.GetPlainText(nil)
		if err != nil {
			return fmt.Errorf("page %d: %w", index, err)
		}
		fmt.Fprintf(&output, "\n=== PAGE %d ===\n\n", index)
		content = strings.ReplaceAll(content, "\x00", "\n")
		for _, line := range strings.Split(content, "\n") {
			if line = strings.Join(strings.Fields(line), " "); line != "" {
				fmt.Fprintln(&output, line)
			}
		}
		annotations := page.V.Key("Annots")
		for index := 0; index < annotations.Len(); index++ {
			if uri := annotations.Index(index).Key("A").Key("URI").Text(); uri != "" {
				fmt.Fprintln(&output, "LINK:", uri)
			}
		}
	}
	derived := []byte(output.String())
	if check {
		existing, err := os.ReadFile(target)
		if err != nil {
			return err
		}
		if !bytes.Equal(existing, derived) {
			return fmt.Errorf("the derived guide differs from the PDF. Run sync-guide --course %s", c.ID)
		}
		fmt.Println("The derived course guide matches the PDF.")
		return nil
	}
	if err := writeAtomic(target, derived); err != nil {
		return err
	}
	fmt.Println("Updated", target)
	return nil
}
