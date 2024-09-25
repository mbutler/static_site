package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/russross/blackfriday/v2"
)

func main() {
	sourceDir := "./content"
	outputDir := "./public"

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			err := processFile(path, sourceDir, outputDir)
			if err != nil {
				log.Println(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Error walking the path %q: %v\n", sourceDir, err)
	}
}

func processFile(path, sourceDir, outputDir string) error {
	input, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}

	htmlContent := blackfriday.Run(input)

	relPath, err := filepath.Rel(sourceDir, path)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %v", err)
	}

	outputFilePath := filepath.Join(outputDir, relPath)
	outputFilePath = outputFilePath[:len(outputFilePath)-len(filepath.Ext(outputFilePath))] + ".html"

	err = os.MkdirAll(filepath.Dir(outputFilePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	f, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", outputFilePath, err)
	}
	defer f.Close()

	data := struct {
		Content string
	}{
		Content: string(htmlContent),
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}
