# Define log paths explicitly
LOGS_DIR="$PAK_DIR/Logs"
LAUNCH_LOG="$LOGS_DIR/launch.log"
ERROR_LOG="$LOGS_DIR/theme-manager-error.log"

# Create Logs directory if it doesn't exist
mkdir -p "$LOGS_DIR"

# Log environment
echo "PAK directory: $PAK_DIR" > "$LAUNCH_LOG"
echo "Environment variables:" >> "$LAUNCH_LOG"
env >> "$LAUNCH_LOG"

# Log binary permissions
echo "Binary permissions:" >> "$LAUNCH_LOG"
ls -la theme-manager minui-list minui-presenter >> "$LAUNCH_LOG" 2>&1

echo "Changed to directory: $(pwd)" >> "$LAUNCH_LOG"

# Make sure binaries are executable
chmod +x theme-manager minui-list minui-presenter 2>/dev/null
echo "Made binaries executable" >> "$LAUNCH_LOG"

# Log launch
echo "Launching theme-manager" >> "$LAUNCH_LOG"

# Launch with output redirection and better error handling
./theme-manager 2>>"$ERROR_LOG" || {
  echo "Application crashed with exit code: $?" >> "$LAUNCH_LOG"
  echo "Last 10 lines of error log:" >> "$LAUNCH_LOG"
  tail -10 "$ERROR_LOG" >> "$LAUNCH_LOG"
}

# Exit code
echo "Exit code: $?" >> "$LAUNCH_LOG"