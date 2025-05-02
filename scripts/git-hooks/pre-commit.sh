#!/bin/sh

# echo "Hello from Pre commit hook"

# Check code formatting
if ! gofmt -l . | grep -q '.'; then
  echo "Error: Code is not properly formatted. Run 'gofmt -w .' to fix." >&2
  exit 1
fi