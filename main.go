package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// Templates for the files
const goModTemplate = `module {{.ModuleName}}

go 1.20
`

const mainGoTemplate = `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
`

const configYAMLTemplate = `server:
  port: 8080
database:
  host: localhost
  port: 5432
  user: user
  password: password
  dbname: mydb
`

const dockerfileTemplate = `# Use the official Golang image as a base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/app

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
`

const readmeTemplate = `# {{.ProjectName}}

This is the {{.ProjectName}} project.
`

const gitignoreTemplate = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Output of the go coverage tool, specifically when used with LiteIDE
*.out
`

const licenseTemplate = `MIT License

Copyright (c) {{.Year}} {{.Author}}

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`

// Main function to create the project scaffolding
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myscaffold <project-name>")
		return
	}

	projectName := os.Args[1]
	projectPath := filepath.Join(".", projectName)
	moduleName := fmt.Sprintf("github.com/yourusername/%s", projectName)

	// Create directories
	dirs := []string{
		"cmd/app",
		"internal/handler",
		"internal/service",
		"internal/repository",
		"pkg",
		"configs",
	}

	for _, dir := range dirs {
		path := filepath.Join(projectPath, dir)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", path, err)
			return
		}
		fmt.Println("Successfully created directory:", path)
	}

	// Create files
	files := map[string]string{
		"go.mod":              goModTemplate,
		"cmd/app/main.go":     mainGoTemplate,
		"configs/config.yaml": configYAMLTemplate,
		"Dockerfile":          dockerfileTemplate,
		"Makefile":            makefileTemplate,
		"README.md":           readmeTemplate,
		".gitignore":          gitignoreTemplate,
		"LICENSE":             licenseTemplate,
	}
	// Create the files from the templates
	for path, content := range files {
		createFile(filepath.Join(projectPath, path), content, map[string]string{
			"ModuleName":  moduleName,
			"ProjectName": projectName,
			"Year":        "2023",
			"Author":      "Raldes Krisnu",
		})
	}
	// Create println file
	fmt.Println("Project scaffold created successfully at", projectPath)

	// Initialize the Go module
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error initializing Go module:", err)
		return
	}
}

// Helper function to create a file from a template
func createFile(path string, templateContent string, data map[string]string) {
	tmpl, err := template.New(filepath.Base(path)).Parse(templateContent)
	if err != nil {
		fmt.Println("Error creating template:", err)
		return
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if data != nil {
		tmpl.Execute(file, data)
	} else {
		file.WriteString(templateContent)
	}
}
