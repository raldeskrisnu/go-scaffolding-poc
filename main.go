package main

import (
	"flag"
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

func main() {

	moduleName := flag.String("module", "", "The module name for go.mod (e.g., github.com/username/package)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: myscaffold --module=github.com/username/package <project-name>")
		return
	}

	projectName := flag.Args()[0]
	projectPath := filepath.Join(".", projectName)

	if *moduleName == "" {
		fmt.Println("Error: --module flag is required")
		return
	}

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
	}

	// Create files
	files := map[string]string{
		"go.mod":              goModTemplate,
		"cmd/app/main.go":     mainGoTemplate,
		"configs/config.yaml": configYAMLTemplate,
		"Dockerfile":          dockerfileTemplate,
	}

	for path, content := range files {
		createFile(filepath.Join(projectPath, path), content, map[string]string{"ModuleName": *moduleName})
	}

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

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	if data != nil {
		err = tmpl.Execute(file, data)
	} else {
		_, err = file.WriteString(templateContent)
	}

	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
