#!/bin/sh

# This script is a pre-commit hook that checks the size of staged files before committing.

# Maximum file size limit in bytes
# You can set this limit by exporting the GIT_FILE_SIZE_LIMIT environment variable
limit=${GIT_FILE_SIZE_LIMIT:-5000000} # Default 5MB
limit_in_mb=$(( $limit / 1000000 ))

# Move to the repo root so git files paths make sense
repo_root=$( git rev-parse --show-toplevel )
cd $repo_root

empty_tree=$( git hash-object -t tree /dev/null )

if git rev-parse --verify HEAD > /dev/null 2>&1
then
	against=HEAD
else
	against="$empty_tree"
fi

# Set split so that for loop below can handle spaces in file names by splitting on line breaks
IFS='
'

echo "Checking staged file sizes"
for file in $(git diff-index --cached --name-only "$against"); do
    # Check if the file exists
    if [ -f "$file" ]; then
        # Get the file size using stat
        file_size=$(stat -c%s "$file")

        # Check if the file size exceeds the limit
        if [ "$file_size" -gt "$limit" ]; then
            echo "File $file is $((file_size / 1000000)) MB, which is larger than our configured limit of $limit_in_mb MB"
            echo "If you really need to commit this file, you can override the size limit by setting the GIT_FILE_SIZE_LIMIT environment variable, e.g. GIT_FILE_SIZE_LIMIT=20000000 for 20MB. Or, commit with the --no-verify switch to skip the check entirely."
            echo "Commit aborted"

            exit 1
        fi
    fi
done
