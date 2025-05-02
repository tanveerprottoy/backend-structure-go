#!/bin/sh

# check file size, any file larger than 5 MB will cause
# the commit to fail
MAX_SIZE=5000000 # 5 MB

for file in $(git diff --cached --name-only); do
  if [ $(stat -c%s "$file") -gt $MAX_SIZE ]; then
    echo "Error: $file exceeds the size limit of 5 MB." >&2
    exit 1
  fi
done
