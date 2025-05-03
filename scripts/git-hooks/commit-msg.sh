#!/bin/sh

# This script is a commit-msg hook that checks the commit message for certain criteria.
# It ensures that the commit message is not empty and is longer than 5 characters.

# Get the commit message
COMMIT_MSG_FILE=$1
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Check if the commit message is empty
if [ -z "$COMMIT_MSG" ]; then
  echo "Error: Commit message cannot be empty" >&2
  exit 1
fi

# Check if the commit message length is less than or equal to 5 characters
if [ ${#COMMIT_MSG} -le 5 ]; then
  echo "Error: Commit message must be greater than 5 characters" >&2
  exit 1
fi
