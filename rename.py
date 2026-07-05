#!/usr/bin/env python3
import os
import sys

# Core definition keys to find and replace across the workspace
OLD_MODULE = "github.com/dat267/goclitmpl"
OLD_NAME = "goclitmpl"

# Folders completely excluded from string scanning operations
EXCLUDE_DIRS = {".git", ".github/assets", "node_modules", "__pycache__"}


def rename_project(new_module: str, new_name: str):
    print(f"[STATUS] Renaming project template to: {new_name}")
    print(f"[STATUS] Updating Go module target path to: {new_module}\n")

    for root, dirs, files in os.walk(".", topdown=True):
        # Filter out tracking folders in place to avoid infinite recursive walks
        dirs[:] = [
            d
            for d in dirs
            if os.path.join(root, d).lstrip("./") not in EXCLUDE_DIRS
            and d not in EXCLUDE_DIRS
        ]

        for file in files:
            if file == "rename.py":
                continue

            file_path = os.path.join(root, file)

            try:
                with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
                    content = f.read()

                # Process replacements if occurrences exist
                if OLD_MODULE in content or OLD_NAME in content:
                    updated_content = content.replace(OLD_MODULE, new_module)
                    updated_content = updated_content.replace(OLD_NAME, new_name)

                    with open(file_path, "w", encoding="utf-8") as f:
                        f.write(updated_content)
                    print(f"  [UPDATED] {file_path}")
            except Exception as e:
                print(f"  [SKIP] Error processing {file_path}: {e}")

    print("\n[SUCCESS] Workspace string transformations complete.")


if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python3 rename.py [new-module-path] [new-binary-name]")
        print("Example: python3 rename.py github.com/username/my-tool my-tool")
        sys.exit(1)

    rename_project(sys.argv[1], sys.argv[2])
