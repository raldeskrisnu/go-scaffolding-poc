package main

const makefileTemplate = `# Build the application
build:
    go build -o bin/app ./cmd/app

# Run the application
run:
    go run ./cmd/app

# Clean the build files
clean:
    rm -rf bin/

# Build and run the Docker container
docker:
    docker build -t yourusername/yourproject .
    docker run -p 8080:8080 yourusername/yourproject
`
