#!/bin/bash
# Organize C# generated files into folder structure based on namespaces
# Run after: buf generate

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CSHARP_DIR="$SCRIPT_DIR/../gen/csharp"

echo -e "\nOrganizing C# files into folder structure..."
echo -e "Directory: $CSHARP_DIR"

if [ ! -d "$CSHARP_DIR" ]; then
    echo "Directory $CSHARP_DIR does not exist. Run buf generate first."
    exit 0
fi

MOVED_COUNT=0

# Find all C# files
find "$CSHARP_DIR" -maxdepth 1 -name "*.cs" -type f | while read -r file; do
    filename=$(basename "$file")
    
    # Extract namespace (first occurrence of namespace XXX;)
    namespace=$(grep -m 1 "^namespace " "$file" | sed -E 's/^namespace ([a-zA-Z0-9_\.]+)[;\{].*/\1/' | xargs)
    
    if [ -n "$namespace" ]; then
        # Convert namespace to folder path (replace . with /)
        folder_path="${namespace//./\/}"
        target_dir="$CSHARP_DIR/$folder_path"
        
        # Create directory
        mkdir -p "$target_dir"
        
        # Move file if the target is different from the current location
        target_file="$target_dir/$filename"
        if [ "$file" != "$target_file" ]; then
            mv "$file" "$target_file"
            MOVED_COUNT=$((MOVED_COUNT + 1))
            echo "  ✓ $filename -> $folder_path/"
        fi
    fi
done

echo -e "\nOrganized files into folder structure."
echo -e "Location: gen/csharp/"
