#!/bin/sh
# Script to create a sample dynamic theme directory structure

if [ -z "$1" ]; then
  echo "Please provide a theme name."
  echo "Usage: ./setup_dynamic_theme.sh MyTheme"
  exit 1
fi

THEME_NAME="$1"
BASE_DIR="Themes/Dynamic/$THEME_NAME"

# Create main theme directory
mkdir -p "$BASE_DIR"
echo "Creating theme: $THEME_NAME"

# Create system directories
mkdir -p "$BASE_DIR/Root"
mkdir -p "$BASE_DIR/Recently Played"
mkdir -p "$BASE_DIR/Tools"
mkdir -p "$BASE_DIR/Roms"

# Create common system tag directories
# Add more as needed
mkdir -p "$BASE_DIR/Roms/GBA"
mkdir -p "$BASE_DIR/Roms/SNES"
mkdir -p "$BASE_DIR/Roms/NES"
mkdir -p "$BASE_DIR/Roms/GBC"
mkdir -p "$BASE_DIR/Roms/GB"
mkdir -p "$BASE_DIR/Roms/MD"
mkdir -p "$BASE_DIR/Roms/PCE"

# Create placeholder files (touch only creates empty files)
# Users should replace these with actual 320x240 PNG images
touch "$BASE_DIR/Root/bg.png"
touch "$BASE_DIR/Recently Played/bg.png"
touch "$BASE_DIR/Tools/bg.png"
touch "$BASE_DIR/Roms/default.png"
touch "$BASE_DIR/Roms/GBA/bg.png"
touch "$BASE_DIR/Roms/SNES/bg.png"
touch "$BASE_DIR/Roms/NES/bg.png"
touch "$BASE_DIR/Roms/GBC/bg.png"
touch "$BASE_DIR/Roms/GB/bg.png"
touch "$BASE_DIR/Roms/MD/bg.png"
touch "$BASE_DIR/Roms/PCE/bg.png"

echo "Dynamic theme structure created at $BASE_DIR"
echo "Please replace the empty files with actual 320x240 PNG images."
echo ""
echo "Note: The following system tags are included by default:"
echo "- GBA  (Game Boy Advance)"
echo "- SNES (Super Nintendo)"
echo "- NES  (Nintendo Entertainment System)"
echo "- GBC  (Game Boy Color)"
echo "- GB   (Game Boy)"
echo "- MD   (Mega Drive/Genesis)"
echo "- PCE  (PC Engine/TurboGrafx-16)"
echo ""
echo "You can add more system tags as needed based on what's installed on your device."