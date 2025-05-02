#!/bin/sh

# check file size, any file larger than 5 MB will cause
# the commit to fail
MAX_SIZE=5000000 # 5 MB

# Get the list of staged files
STAGED_FILES=$(git diff --cached --name-only)

# Exit if no files are staged
if [ -z "$STAGED_FILES" ]; then
  exit 0
fi

# Loop through staged files
for file in $STAGED_FILES; do
  # Check if the file exists in the Git index
  if [ -f "$file" ]; then
    # Get the size of the staged file content
    FILE_SIZE=$(git cat-file -s ":$file")
    if [ "$FILE_SIZE" -gt "$MAX_SIZE" ]; then
      echo "Error: $file exceeds the size limit of 5 MB." >&2
      exit 1
    fi
  fi
done

exit 0