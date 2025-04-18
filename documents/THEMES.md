# Theme Package Format Documentation

This document explains the structure and format of `.theme` packages used by the NextUI Theme Manager.

## What's in a Theme Package?

A theme package is a directory with a `.theme` extension containing:

1. **manifest.json** - Theme metadata and file mappings
2. **preview.png** - Theme `preview.png`
3. **Wallpapers/** - Background images for all screens
4. **Icons/** - System, tool, and collection icons
5. **Fonts/** - Font replacements and backups
6. **Overlays/** - System-specific overlays

## Quick Start

The fastest way to understand themes:

1. Open Theme Manager
2. Go to **Export** from the main menu
3. Navigate to `Tools/tg5040/Theme-Manager.pak/Exports` on your device
4. Explore the newly created theme package to see how files are organized.

## Detailed Theme Structure

### `manifest.json`

This file contains essential metadata about your theme:

```json
{
  "theme_info": {
    "name": "theme_1",
    "version": "1.0.0",
    "author": "AuthorName",
    "creation_date": "2025-04-13T12:00:00Z",
    "exported_by": "Theme Manager v1.0"
  },
  "content": {
    "wallpapers": {
      "present": true,
      "count": 15
    },
    "icons": {
      "present": true,
      "system_count": 12,
      "tool_count": 5,
      "collection_count": 3
    },
    "overlays": {
      "present": true,
      "systems": ["MGBA", "SFC", "MD"]
    },
    "fonts": {
      "present": true,
      "og_replaced": true,
      "next_replaced": true
    },
    "settings": {
      "accents_included": true,
      "leds_included": true
    }
  },
  "path_mappings": {
    /* File mappings between theme and system paths */
  },
  "accent_colors": {
    "color1": "#FFFFFF",
    "color2": "#9B2257",
    "color3": "#1E2329",
    "color4": "#FFFFFF", 
    "color5": "#000000",
    "color6": "#FFFFFF"
  },
  "led_settings": {
    "f1_key": {
      "effect": 1,
      "color1": "0xFFFFFF",
      "color2": "0x000000",
      "speed": 1000,
      "brightness": 100,
      "trigger": 1,
      "in_brightness": 100
    },
    /* Additional LED sections */
  }
}
```

The manifest helps Theme Manager understand where files should be copied during import and contains useful metadata like your author name and theme version.
 
### Wallpapers

Wallpapers directory structure:

```
Wallpapers/
├─ SystemWallpapers/          # All system and main section wallpapers
│  ├─ Root.png                # Default global NextUI background
│  ├─ Recently Played.png     # Recently played list
│  ├─ Tools.png               # Tools menu
│  ├─ Collections.png         # Main collections menu
│  └─ Game Boy Advance (GBA).png  # System wallpapers with tags
└─ CollectionWallpapers/      # Individual collection wallpapers
   └─ Handhelds.png           # Named after collection folders
```

**Important notes:**
- System wallpapers must can be named whatever you want, but they **MUST** contain their respective system tag in parentheses, like (PS), (MGBA), (MD), save for `Root`, `Recently Played`, `Tools`, and `Collections`, which are named literally as above.
- Resolution should ideally match your device's screen resolution (1024x768 for TrimUI Brick)
- PNG format is required

### Icons

Icons directory structure:

```
Icons/
├─ SystemIcons/             # Icons displayed in main menu
│  ├─ Collections.png
│  ├─ Recently Played.png
│  ├─ Tools.png
│  └─ Game Boy Advance (GBA).png  # System icons with tags
├─ ToolIcons/               # Icons for individual tools
│  └─ Battery.png            # Named as shown in Tools folder
└─ CollectionIcons/         # Icons for collections
   └─ Favorites.png         # Named as shown in Collections
```

**Important notes:**
- System icons should include the system tag in parentheses, save for `Tools`, `Recently Played`, and `Collections`, which are named literally.
- Any tool-specific icons also need to be named literally, for instance, `Battery.png` for the `Battery.pak` tool on your device.
- Icon format should be PNG with transparency
- Recommended size is a square of at least 200x200 px, as NextUI will resize the image.

### Fonts

Fonts directory structure:

```
Fonts/
├─ OG.ttf                   # Your replacement for font2.ttf
├─ Next.ttf                 # Your replacement for font1.ttf
├─ OG.backup.ttf            # Backup of original font2.ttf
└─ Next.backup.ttf          # Backup of original font1.ttf
```

**Why we save backups:**
Font backups are crucial for restoring the system to its original state if a custom font causes issues. When you replace a font in the Settings app, Theme Manager automatically creates a backup of the original font. These backups are included in theme packages to ensure a complete restoration is possible when importing.

To change the available themes in the font, simply replace the `OG.ttf` and/or `Next.ttf` font files with whatever you'd like, and the user will be able to cycle between either of them in their `Settings.pak`.

### Overlays

Overlays directory structure:

```
Overlays/
├─ MGBA/                    # Overlays for GBA system
│  ├─ overlay1.png
│  └─ overlay2.png
└─ [other systems]/
   └─ [overlay files].png
```

## System Tags

System tags are crucial for correctly mapping files to systems. Common system tags include:

- (GBA) - Game Boy Advance
- (SFC)/(SUPA) - Super Nintendo
- (NES) - Nintendo Entertainment System
- (MD) - Sega Genesis/Mega Drive
- (PS) - PlayStation
- (FBN) - FinalBurn Neo

## Accent Colors

Accent colors define the UI color scheme:

- **color1** - Main UI color
- **color2** - Primary accent color
- **color3** - Secondary accent color
- **color4** - List text color
- **color5** - Selected list text color
- **color6** - Hint/information text color

## LED Settings

LED settings control the behavior of the device's LED lights:

- **effect** - Lighting effect type (1-7)
- **color1/color2** - Primary and secondary colors (hex format)
- **speed** - Animation speed in milliseconds
- **brightness** - LED brightness (0-100)
- **trigger** - Button/event that activates the LED (1-14)
- **inbrightness** - Information LED brightness (0-100)

## Importing Themes

To import a theme:

1. Place your `.theme` package in `Tools/tg5040/Theme-Manager.pak/Themes/`
2. Open Theme Manager
3. Navigate to **Browse Themes**
4. Select your theme from the list
5. Confirm the import

During import, Theme Manager will:
- Clear any applied wallpapers and icons (if the theme includes these components)
- Read the manifest.json file
- Copy files to their appropriate system locations
- Create backups of fonts if necessary
- Apply accent and LED settings if included

## Exporting Themes

To export your current device setup:

1. Open Theme Manager
2. Navigate to **Export** from the main menu
3. Theme Manager will create a new theme package in `Tools/tg5040/Theme-Manager.pak/Exports/`

Exported themes are named sequentially (theme_1.theme, theme_2.theme, etc.)

## Troubleshooting

### Theme Import Issues

- **Missing Files**: Ensure all required directories and files exist
- **Permissions**: Make sure all files are readable
- **Manifest Errors**: Check manifest.json for syntax errors
- **System Tags**: Verify system tags in parentheses match your system
- **Logging**: Theme Manager comes with a detailed logger in `Logs/theme-manager.log`. You can always take a look here if there are any issues. Keep in mind this file fills up quickly, so make sure to clear it every once in awhile!

### Font Problems

If a custom font causes display issues:

1. You can restore the original fonts by applying the `Default.font` package in `Components -> Fonts -> Browse`
2. You can also manually place the backed up `.ttf` files in `.system/res` as `Next.ttf` and `OG.ttf` 