#!/bin/sh
# Debug launcher for theme manager

# Set to 1 to enable logging, 0 to disable logging
ENABLE_LOGGING=0

# Get the directory where this script is located
PAK_DIR="$(dirname "$0")"
cd "$PAK_DIR" || exit 1

if [ "$ENABLE_LOGGING" -eq 1 ]; then
    # Log environment
    echo "PAK directory: $PAK_DIR" > launch.log
    echo "Environment variables:" >> launch.log
    env >> launch.log

    # Log binary permissions
    echo "Binary permissions:" >> launch.log
    ls -la theme-manager minui-list minui-presenter >> launch.log 2>&1

    echo "Changed to directory: $(pwd)" >> launch.log

    # Log launch
    echo "Launching theme-manager" >> launch.log

    # Launch with output redirection
    ./theme-manager 2>>theme-manager-error.log

    # Exit code
    echo "Exit code: $?" >> launch.log
else
    # Make sure binaries are executable without logging
    chmod +x theme-manager minui-list minui-presenter 2>/dev/null

    # Launch without output redirection
    ./theme-manager 2>/dev/null
fi