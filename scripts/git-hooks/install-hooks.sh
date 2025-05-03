#!/bin/sh

# This script installs the git hooks for the project.

echo 'starting copy shell scripts...'

# copy shell script files to their desired hook destinations
cp ./scripts/git-hooks/commit-msg.sh .git/hooks/commit-msg

cp ./scripts/git-hooks/pre-commit.sh .git/hooks/pre-commit

cp ./scripts/git-hooks/pre-push.sh .git/hooks/pre-push

echo 'making hooks executables...'

# make hooks executable
chmod +x .git/hooks/commit-msg

chmod +x .git/hooks/pre-commit

chmod +x .git/hooks/pre-push