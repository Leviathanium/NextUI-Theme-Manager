#!/bin/sh
# Debug launcher for theme manager

# Get the directory where this script is located
PAK_DIR="$(dirname "$0")"
cd "$PAK_DIR" || exit 1

# Log environment
echo "PAK directory: $PAK_DIR" > launch.log
echo "Environment variables:" >> launch.log
env >> launch.log

# Log binary permissions
echo "Binary permissions:" >> launch.log
ls -la theme-manager minui-list minui-presenter >> launch.log 2>&1

echo "Changed to directory: $(pwd)" >> launch.log

# Make sure binaries are executable
chmod +x theme-manager minui-list minui-presenter 2>/dev/null
echo "Made binaries executable" >> launch.log

# Log launch
echo "Launching theme-manager" >> launch.log

# Launch with output redirection
./theme-manager 2>>theme-manager-error.log

# Exit code
echo "Exit code: $?" >> launch.log