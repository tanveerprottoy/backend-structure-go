#!/bin/sh

cp ./scripts/git-hooks/commit-msg.sh .git/hooks/commit-msg
cp ./scripts/git-hooks/pre-commit.sh .git/hooks/pre-commit
cp ./scripts/git-hooks/pre-push.sh .git/hooks/pre-push

chmod +x .git/hooks/commit-msg
chmod +x .git/hooks/pre-commit
chmod +x .git/hooks/pre-push