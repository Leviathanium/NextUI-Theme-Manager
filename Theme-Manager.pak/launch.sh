#!/bin/sh

# Get the PAK directory path
PAK_DIR="$(dirname "$0")"
echo "PAK directory: $PAK_DIR" > "$PAK_DIR/launch.log"

# Print environment info for debugging
echo "Environment variables:" >> "$PAK_DIR/launch.log"
env >> "$PAK_DIR/launch.log"

# Check file permissions
echo "Binary permissions:" >> "$PAK_DIR/launch.log"
ls -la "$PAK_DIR/theme-manager" >> "$PAK_DIR/launch.log"
ls -la "$PAK_DIR/minui-list" >> "$PAK_DIR/launch.log"
ls -la "$PAK_DIR/minui-presenter" >> "$PAK_DIR/launch.log"

# Change to the PAK directory
cd "$PAK_DIR" || exit 1
echo "Changed to directory: $(pwd)" >> "$PAK_DIR/launch.log"

# Make sure binaries are executable
chmod +x "$PAK_DIR/theme-manager"
chmod +x "$PAK_DIR/minui-list"
chmod +x "$PAK_DIR/minui-presenter"
echo "Made binaries executable" >> "$PAK_DIR/launch.log"

# Launch the theme manager
echo "Launching theme-manager" >> "$PAK_DIR/launch.log"
./theme-manager
echo "theme-manager exited with code: $?" >> "$PAK_DIR/launch.log"