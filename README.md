# NextUI Theme Manager

This Pak allows user to customize their NextUI devices to their hearts' content. No more dragging and dropping `bg.png`!

_NOTE: THERE WILL BE BUGS! Also, this is currently only for the TrimUI Brick._

## Features

- Theme Management: Import/Export your device's current customization options to share them with the world!
- Global Customization: Change global wallpapers, icons, fonts, LEDs, and accents.
- System Customization: Change wallpapers and icons for specific systems.
- Exporting: Export Custom LED and Accent settings, and import/apply them as well.
- Reset: Clear all background/icon settings to start fresh.

With more features to come...

---
## Gallery

![Retro-Programmer](Theme-Manager.pak/Wallpapers/Retro-Programmer/bg.png)
![Blackstreets](Theme-Manager.pak/Wallpapers/Blackstreets/bg.png)
![Cozy](Theme-Manager.pak/Wallpapers/Cozy/bg.png)
![Firewatch](Theme-Manager.pak/Wallpapers/Firewatch/bg.png)
![Retro-Mario-Chill](Theme-Manager.pak/Wallpapers/Retro-Mario-Chill/bg.png)
![Sunset-Forest](Theme-Manager.pak/Wallpapers/Sunset-Forest/bg.png)

---
## Guide


_Theme Management_
1. **Themes (working alpha):** Import or export themes to your device. See documentation for more details on what constitutes a theme.
2. **Customization:** Here, you can change specific wallpaper/icon/font elements for your device.
3. **Reset:** This offers the ability to clear all background or icon files across NextUI if you want to start over fresh.
---
_Customization_
---
1. **Global Options:** Apply global icons/wallpapers to the entire device from `Theme-Manager.pak/Themes/Global` and `Theme-Manager.pak/Icons`. There are a few presets to mess around with, or if you create your own, you can apply them from here.
2. **System Options:** Apply system-specific main menu icons/wallpapers if you want to get really technical. This also includes `Recently Played`, `Tools`, and `Collections` as well!
3. **Accents:** Choose from several preset options, or export your custom-built options from Settings.pak! They are exported as `.txt` files in `Theme-Manager.pak/Accents/Custom`.
4. **LEDs:** Choose from several preset options, or export your custom-built options from LedControl.pak! They are exported as `.txt` files in `Theme-Manager.pak/LEDs/Custom`.
5. **Fonts:** You can _replace_ and _restore_ the `Next` and `OG` fonts in Settings.pak with any fonts you choose. Fonts are located in `Theme-Manager.pak/Fonts`.
---
## Installation
1. Clone or download `Theme-Manager.zip`.
2. Move the `Theme-Manager.pak` folder into your `Tools/tg5040` directory on your SD card.
3. Launch it on your Brick and start changing your theme!
---
# Documentation

## What's in a Theme Pack?

A theme pack is a directory with a `.theme` extension containing:

1. **manifest.json** - Theme metadata and file mappings
2. **preview.png** - Theme preview image
3. **Wallpapers** - Background images for all screens
4. **Icons** - System, tool, and collection icons
5. **Fonts** - Font replacements and backups
6. **Settings** - Accent and LED configuration
7. **Overlays** - System-specific overlays

## Quick Start

For the fastest way to understand themes:

1. Open Theme Manager
2. Go to **Themes → Export Current Settings**
3. Navigate to `Tools/tg5040/Theme-Manager.pak/Themes/Exports` on your device
4. Explore the newly created theme package to see how files are organized

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
  }
}
```

The manifest helps Theme Manager understand where files should be copied during import and contains useful metadata like your author name.

### Wallpapers

Wallpapers directory structure:

```
Wallpapers/
├─ SystemWallpapers/          # All system and main section wallpapers
│  ├─ Root.png                # Main menu background
│  ├─ Recently Played.png     # Recently played list
│  ├─ Tools.png               # Tools menu
│  ├─ Collections.png         # Main collections menu
│  └─ Game Boy Advance (GBA).png  # System wallpapers with tags
└─ CollectionWallpapers/      # Individual collection wallpapers
   └─ Handhelds.png           # Named after collection folders
```

**Important notes:**
- System wallpapers must be named with their respective system tag, like (PS), (MGBA), (MD), save for `Root`, `Recently Played`, `Tools`, and `Collections`, which are named literally as above.
- Resolution should ideally match your TrimUI Brick's screen resolution (1024x768)
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
│  └─ Tetris.png            # Named as shown in Tools folder
└─ CollectionIcons/         # Icons for collections
   └─ Favorites.png         # Named as shown in Collections
```

**Important notes:**
- System icons should include the system tag in parentheses, save for `Tools`, `Recently Played`, and `Collections`, which are literal, as above.
- Icon format should be PNG with transparency
- Recommended size is at a square of at least 200x200 px, the image will be resized by NextUI.

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

### Settings

Settings directory structure:

```
Settings/
├─ minuisettings.txt        # Accent color settings
└─ ledsettings_brick.txt    # LED settings
```

**Why these are text files:**
Both accent and LED settings are stored as simple text files on the device, making them easy to back up and restore. These files contain key-value pairs that configure the UI colors and LED behaviors.

**minuisettings.txt example:**
```
color1=0xFFFFFF
color2=0x9B2257
color3=0x1E2329
color4=0xFFFFFF
color5=0x000000
color6=0xFFFFFF
font_id=0
show_gamearts=1
show_recents=1
...
```

**ledsettings_brick.txt example:**
```
[F1 key]
effect=4
color1=0xFFFFFF
color2=0x000000
speed=1000
brightness=100
trigger=1
filename=
inbrightness=100

[F2 key]
...
```

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

## Creating Your Own Theme Package

To create a new theme package manually (without exporting):

1. Create a new directory with a `.theme` extension (e.g., `my_awesome_theme.theme`)
2. Create the directory structure outlined above
3. Add your custom files to the appropriate directories
4. Create a basic `manifest.json` file with at least the theme_info section:
   ```json
   {
     "theme_info": {
       "name": "My Awesome Theme",
       "version": "1.0.0",
       "author": "Your Name",
       "creation_date": "2025-04-14T00:00:00Z",
       "exported_by": "Manual Creation"
     }
   }
   ```
5. Create a `preview.png` image (640×480) to showcase your theme
6. Place the completed theme in `Tools/tg5040/Theme-Manager.pak/Themes/Imports/`

## Tips and Best Practices

### System Tags

System tags are crucial for correctly mapping files to systems. Common system tags include:

- (GBA) - Game Boy Advance
- (SFC)/(SUPA) - Super Nintendo
- (NES) - Nintendo Entertainment System
- (MD) - Sega Genesis/Mega Drive
- (PS) - PlayStation
- (FBN) - FinalBurn Neo

### LED Settings

LED settings control the behavior of the device's LED lights:

- **effect** - Lighting effect type (1-7)
- **color1/color2** - Primary and secondary colors (hex format)
- **speed** - Animation speed in milliseconds
- **brightness** - LED brightness (0-100)
- **trigger** - Button/event that activates the LED (1-14)
- **inbrightness** - Information LED brightness (0-100)

### Accent Colors

Accent colors define the UI color scheme:

- **color1** - Main UI color
- **color2** - Primary accent color
- **color3** - Secondary accent color
- **color4** - List text color
- **color5** - Selected list text color
- **color6** - Hint/information text color

## Importing Themes

To import a theme:

1. Place your theme package in `Tools/tg5040/Theme-Manager.pak/Themes/Imports/`
2. Open Theme Manager
3. Navigate to **Themes → Import Theme**
4. Select your theme from the list
5. Confirm the import

During import, Theme Manager will:
- Clear any applied icons and backgrounds
- Read the manifest.json file
- Copy files to their appropriate system locations
- Create backups of fonts if necessary
- Apply accent and LED settings if included

## Exporting Themes

To export your current device setup:

1. Open Theme Manager
2. Navigate to **Themes → Export Current Settings**
3. Theme Manager will create a new theme package in `Tools/tg5040/Theme-Manager.pak/Themes/Exports/`

Exported themes are named sequentially (theme_1.theme, theme_2.theme, etc.)

## Troubleshooting

### Theme Import Issues

- **Missing Files**: Ensure all required directories and files exist
- **Permissions**: Make sure all files are readable
- **Manifest Errors**: Check manifest.json for syntax errors
- **System Tags**: Verify system tags in parentheses match your system

### Font Problems

If a custom font causes display issues:

1. You can restore the original fonts using Theme Manager's font restoration option
2. Or manually restore from the backup fonts included in the theme

## Advanced: Understanding Path Mappings

The manifest.json file contains path mappings that tell Theme Manager exactly where each file should be copied:

```json
"path_mappings": {
  "wallpapers": [
    {
      "theme_path": "Wallpapers/Root/bg.png",
      "system_path": "/mnt/SDCARD/bg.png",
      "metadata": {
        "SystemName": "Root",
        "SystemTag": "ROOT"
      }
    },
    // More wallpaper mappings...
  ],
  // Icons, fonts, settings mappings...
}
```

These mappings support special metadata to help Theme Manager correctly place files even if your system structure differs from the one used when creating the theme.



## Sources

- @frysee for literally everything
- @kytz for the work on Noir-Minimal
- @GreenKraken22 for finding and suggesting arcade-dark
- @Fujykky for the work on Screens-Thematic
- Everyone else in the NextUI discord
- Epic Noir theme from https://github.com/c64-dev/es-theme-epicnoir
- All artwork and image source rights go to their respective owners.
