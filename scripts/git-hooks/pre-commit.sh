#!/bin/sh

# Configuration
MAX_FILE_SIZE_MB=5  # Set your maximum file size in megabytes

MAX_FILE_SIZE_BYTES=$((MAX_FILE_SIZE_MB * 1024 * 1024))

# Function to check if a file exceeds the size limit
check_file_size() {
  local filepath="$1"
  local size=$(git cat-file -s "$filepath")
  if [ "$size" -gt "$MAX_FILE_SIZE_BYTES" ]; then
    echo "Error: File '$filepath' exceeds the maximum allowed size of $MAX_FILE_SIZE_MB MB."
    return 1 # Indicate failure
  fi
  return 0 # Indicate success
}

# Get a list of all staged files
staged_files=$(git diff --cached --name-only)

# Check files in the root directory
echo "Checking root directory files..."
for file in $staged_files; do
  # Check if the file is directly in the root
  if [[ "$file" != */* ]]; then
    if check_file_size "$file"; then
      exit 1
    fi
  fi
done

# Check files within all subdirectories
echo "Checking files in all subdirectories..."
for file in $staged_files; do
  # Check if the file is in a subdirectory
  if [[ "$file" == */* ]]; then
    if check_file_size "$file"; then
      exit 1
    fi
  fi
done

echo "Pre commit check passed"
exit 0