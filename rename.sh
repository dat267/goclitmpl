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

OLD_APP="goclitmpl"
NEW_APP=$(basename "$NEW_MOD")

OLD_ENV="GOCLITMPL"
NEW_ENV=$(echo "$NEW_APP" | tr '[:lower:]' '[:upper:]')

echo "Renaming module from '$OLD_MOD' to '$NEW_MOD'..."
echo "Renaming application name from '$OLD_APP' to '$NEW_APP'..."
echo "Renaming environment prefix from '$OLD_ENV' to '$NEW_ENV'..."

OS_NAME=$(uname)
if [ "$OS_NAME" = "Darwin" ]; then
    # Replace module namespace
    find . -type f \( -name "*.go" -o -name "go.mod" -o -name "go.sum" -o -name "Makefile" -o -name "*.md" -o -name "*.yml" -o -name "*.yaml" -o -name "*.ps1" -o -name "*.sh" \) -exec sed -i '' "s|$OLD_MOD|$NEW_MOD|g" {} +
    # Replace AppName/binary name
    find . -type f \( -name "*.go" -o -name "go.mod" -o -name "go.sum" -o -name "Makefile" -o -name "*.md" -o -name "*.yml" -o -name "*.yaml" -o -name "*.ps1" -o -name "*.sh" \) -exec sed -i '' "s|$OLD_APP|$NEW_APP|g" {} +
    # Replace EnvPrefix
    find . -type f \( -name "*.go" -o -name "Makefile" -o -name "*.ps1" -o -name "*.sh" \) -exec sed -i '' "s|$OLD_ENV|$NEW_ENV|g" {} +
else
    # Replace module namespace
    find . -type f \( -name "*.go" -o -name "go.mod" -o -name "go.sum" -o -name "Makefile" -o -name "*.md" -o -name "*.yml" -o -name "*.yaml" -o -name "*.ps1" -o -name "*.sh" \) -exec sed -i "s|$OLD_MOD|$NEW_MOD|g" {} +
    # Replace AppName/binary name
    find . -type f \( -name "*.go" -o -name "go.mod" -o -name "go.sum" -o -name "Makefile" -o -name "*.md" -o -name "*.yml" -o -name "*.yaml" -o -name "*.ps1" -o -name "*.sh" \) -exec sed -i "s|$OLD_APP|$NEW_APP|g" {} +
    # Replace EnvPrefix
    find . -type f \( -name "*.go" -o -name "Makefile" -o -name "*.ps1" -o -name "*.sh" \) -exec sed -i "s|$OLD_ENV|$NEW_ENV|g" {} +
fi

# Rename the main entrypoint cmd directory
if [ -d "cmd/$OLD_APP" ]; then
    echo "Renaming cmd/$OLD_APP directory to cmd/$NEW_APP..."
    mv "cmd/$OLD_APP" "cmd/$NEW_APP"
fi

echo "Successfully renamed module to '$NEW_MOD'!"
echo "Running go mod tidy..."
go mod tidy
