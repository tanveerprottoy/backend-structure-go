#!/bin/sh

cp ./scripts/git-hooks/pre-commit.sh .git/hooks/pre-commit
cp ./scripts/git-hooks/pre-push.sh .git/hooks/pre-push
chmod 755 .git/hooks/pre-commit
chmod 755 .git/hooks/pre-push