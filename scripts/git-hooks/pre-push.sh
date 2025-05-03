#!/bin/sh

# This script is a pre-push hook that checks for common issues before pushing code.

# Check code formatting with gofmt
echo "Running gofmt..."
UNFORMATTED_FILES=$(gofmt -l .)
if [ -n "$UNFORMATTED_FILES" ]; then
  echo "Error: The following files are not properly formatted:"
  echo "$UNFORMATTED_FILES"
  echo "Run 'gofmt -w .' to format the code."
  
  exit 1
fi

# Run go vet for static analysis
echo "Running go vet..."
if ! go vet ./...; then
  echo "Error: go vet found issues in the code."
  exit 1
fi

echo "Pre-push checks passed."
exit 0

# lint check with golangci
# golangci-lint run --tests=0 ./...