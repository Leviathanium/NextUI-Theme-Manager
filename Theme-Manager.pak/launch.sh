#!/bin/sh
# Debug launcher for theme manager

# Get the directory where this script is located
PAK_DIR="$(dirname "$0")"
cd "$PAK_DIR" || exit 1

# Create Logs directory if it doesn't exist
mkdir -p "$PAK_DIR/Logs"

# Log environment
echo "PAK directory: $PAK_DIR" > "$PAK_DIR/Logs/launch.log"
echo "Environment variables:" >> "$PAK_DIR/Logs/launch.log"
env >> "$PAK_DIR/Logs/launch.log"

# Log binary permissions
echo "Binary permissions:" >> "$PAK_DIR/Logs/launch.log"
ls -la theme-manager minui-list minui-presenter >> "$PAK_DIR/Logs/launch.log" 2>&1

echo "Changed to directory: $(pwd)" >> "$PAK_DIR/Logs/launch.log"

# Make sure binaries are executable
chmod +x theme-manager minui-list minui-presenter 2>/dev/null
echo "Made binaries executable" >> "$PAK_DIR/Logs/launch.log"

# Log launch
echo "Launching theme-manager" >> "$PAK_DIR/Logs/launch.log"

# Launch with output redirection
./theme-manager 2>>"$PAK_DIR/Logs/theme-manager-error.log"

# Exit code
echo "Exit code: $?" >> "$PAK_DIR/Logs/launch.log"