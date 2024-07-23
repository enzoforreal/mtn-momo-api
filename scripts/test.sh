#!/bin/bash
echo "Running tests with detailed output and coverage..."
go test -v -coverprofile=coverage.out ./...
echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html
echo "Opening HTML coverage report..."

# Adjust the path to your browser executable
BROWSER="/mnt/c/Program Files/Google/Chrome/Application/chrome.exe"
if [ -f "$BROWSER" ]; then
    "$BROWSER" coverage.html
else
    echo "Browser not found. Please open 'coverage.html' manually."
fi
