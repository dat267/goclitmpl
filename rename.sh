#!/bin/bash
set -e

# Checks if a target name is provided
if [ -z "$1" ]; then
    echo "Usage: ./rename.sh <new-module-path>"
    echo "Example: ./rename.sh github.com/username/mycli"
    exit 1
fi

NEW_MOD="$1"
OLD_MOD="github.com/dat267/goclitmpl"

echo "Renaming module from '$OLD_MOD' to '$NEW_MOD'..."

OS_NAME=$(uname)
if [ "$OS_NAME" = "Darwin" ]; then
    find . -type f -name "*.go" -exec sed -i '' "s|$OLD_MOD|$NEW_MOD|g" {} +
    find . -type f -name "go.mod" -exec sed -i '' "s|$OLD_MOD|$NEW_MOD|g" {} +
    find . -type f -name "go.sum" -exec sed -i '' "s|$OLD_MOD|$NEW_MOD|g" {} +
    find . -type f -name "Makefile" -exec sed -i '' "s|$OLD_MOD|$NEW_MOD|g" {} +
    find . -type f -name "README.md" -exec sed -i '' "s|$OLD_MOD|$NEW_MOD|g" {} +
    find . -type f -name "DEVELOPER.md" -exec sed -i '' "s|$OLD_MOD|$NEW_MOD|g" {} +
else
    find . -type f -name "*.go" -exec sed -i "s|$OLD_MOD|$NEW_MOD|g" {} +
    find . -type f -name "go.mod" -exec sed -i "s|$OLD_MOD|$NEW_MOD|g" {} +
    find . -type f -name "go.sum" -exec sed -i "s|$OLD_MOD|$NEW_MOD|g" {} +
    find . -type f -name "Makefile" -exec sed -i "s|$OLD_MOD|$NEW_MOD|g" {} +
    find . -type f -name "README.md" -exec sed -i "s|$OLD_MOD|$NEW_MOD|g" {} +
    find . -type f -name "DEVELOPER.md" -exec sed -i "s|$OLD_MOD|$NEW_MOD|g" {} +
fi

echo "Successfully renamed module to '$NEW_MOD'!"
echo "Running go mod tidy..."
go mod tidy
