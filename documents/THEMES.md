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
2. Go to **Export** from the main menu and export your device's current configuration
3. Navigate to `Tools/tg5040/Theme-Manager.pak/Exports` on your device
4. Explore the newly created theme package to see how files are organized

---

## Detailed Theme Structure

### `manifest.json`

This file contains essential metadata about your theme. It will look something like this when exported:

```json5
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
### Important Manifest Notes

The manifest helps Theme Manager understand where files should be copied during import and contains useful metadata like your author name and theme version.

- When you first download/install a `.theme` and apply it, the `"content"` and `"path_mappings"` properties of the theme's `manifest.json` will **_automatically update_** to reflect your device. You should not need to tweak anything to get the manifest to work.
- The `"theme_info"`, `"accent_colors"`, and `"led_settings"` will always stay the same and will never update.
- If you're ever having trouble with a `.theme` pack not working correctly, you can always look at the `manifest.json` to see if the `"path_mappings"` are going to the right place. That's what it's there for!
- When you _**export**_ a `.theme` pack, the manifest will be heavily populated. Some of these properties should be kept, and some can be safely removed. For more details on how to export and create theme packs, check out the [Theme Creation Guide](documents/THEME_BUILDING.md) for best practices on how to do this.

---

## System Tags

System tags are crucial for correctly mapping files to systems. The are regex text expressions located **_at the end_** of folder and file names. Common system tags include:

- (GB) - Game Boy
- (NES) - Nintendo Entertainment System
- (MD) - Sega Genesis/Mega Drive
- (PS) - PlayStation
- (FBN) - FinalBurn Neo

Some systems support multiple system tags, like Game Boy Advance and Super Nintendo Entertainment system. This is to allow users to specify which emulator they'd like to use for those specific systems.

- (MGBA) Game Boy Advance, using the MGBA emulator
- (GBA) Game Boy Advance, using the alternative emulator
- (SUPA) Super Nintendo Entertainment System, using the Supafaust emulator
- (SFC) Super Nintendo Entertainment System, using the alternative emulator

Regardless of which emulator you choose, **_it is crucial to understand which systems have multiple tags like this_** because `.theme` components should encapsulate **_all possible Rom systems._** For example:

```
- Roms
    - Super Nintendo Entertainment System (SUPA)
    - Super Nintendo Entertainment System (SFC)
    - Game Boy Advance (GBA)
    - Game Boy Advance (MGBA)
    
Regardless of which Rom folder gets populated, ALL components should support ALL folders!
```
Read below for details on how to do this for each `.theme` component.

---

## Wallpapers



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

### Important Wallpaper Notes
When a `.theme` is applied, we:
1. Scan the `Roms` directory on your SD card
2. Match each Rom folder's system tag with the respective `.png` image
3. Move the `.png` image inside a `.media` folder inside that Rom directory
4. Rename the image to `bg.png`
5. Additionally, if a `.theme` does NOT have wallpapers, we delete any previously applied wallpapers.


Here's an example that shows what file names work and **do not work** for wallpapers:

```
In our Wallpaper Directory:

- Wallpapers
    - SystemWallpapers
        - (MGBA).png             <--- Recomended naming convention. Simple and works.
        - Sega Genesis (MD).png  <--- Also valid, contains system tag at the end.
        - Final Burn Neo.png     <--- This WILL NOT WORK. There is no system tag included at the end.
        - 01). Game Boy (GB).png <--- This ALSO WILL NOT WORK. You must REMOVE the "01)" from the .png name for it to work. Remember, we will AUTOMATICALLY find the system tag, so don't worry about adding order numbers to the .png images!


When we go to apply the wallpapers:


- Roms
    - Game Boy Advance (MGBA) <--- (MGBA).png would go inside here
    - Sega Genesis (MD)       <--- Megadrive (MD).png would also work, since we have the matching system tag "(MD)"
    - Arcade (FBN)            <--- Final Burn Neo.png would NOT WORK because there is no system tag in that .png name


Final result after wallpaper application:


- Roms
    - Game Boy Advance (MGBA)
        - .media
            - .bg.png <--- Previously (MGBA).png
    - Sega Genesis (MD)
        - .media
            - .bg.png <--- Previously Sega Genesis (MD).png
    - Arcade (FBN)    <--- Completely fails. No image.

```

Additionally, remember ***additional emulator tags for duplicate systems in NextUI***. If some wallpapers aren't working it might be because the `.png` image provided isn't getting duplicated to the correct Rom folder. For example:

```
In our Wallpaper Directory:

- Wallpapers
    - SystemWallpapers
        - (SUPA).png   <--- For the Supafaust Super Nintendo emulator
        - (SFC).png    <--- For the alternative Super Nintendo emulator
        - (MGBA).png   <--- For the MGBA Game Boy Advance emulator
        - (GBA).png    <--- For the alternative Game Boy Advance emulator

Make sure you include BOTH .png images for wallpapers so that you cover users that might prefer one emulator over the other. Otherwise:


- Roms
    - Super Nintendo (SUPA)
        - Chrono Trigger
        - Super Mario RPG
        - etc...
        - .media
            - .bg  <--- Here, the (SUPA).png went to the directory with all the Roms in it, for users that prefer Supafaust.
    - Super Nintendo (SFC)
        - (EMPTY ROM DIRECTORY)
        - .media
            - .bg  <--- But we also place the (SFC).png image here if the directory exists. See how we cover users with multiple emulators like this?

In this case, the user is covered through BOTH Super Nintendo emulators.
```

---

## Icons

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

### Important Icon Notes

Icons follow a lot of the same rules as wallpapers, but with some exceptions. When a `.theme` is applied we:
1. Scan the Rom directory for all available systems by system tag
2. Copy the `.png` images over to their respective folders inside a `.media` folder.
3. We **_rename_** the icons according to the name of the appropriate Rom directory so that they match.
4. Additionally, if a `.theme` does NOT have any icons, we delete any previously applied icons.

Here are some examples of the process. It's very similar to how wallpapers work:

```
In our icon directory:

- Icons
    - SystemIcons
        - (MGBA).png             <--- Recomended naming convention. Simple and works.
        - Sega Genesis (MD).png  <--- Also valid, contains system tag at the end.
        - Arcade (FBN).png       <--- Valid name too.

When we go to apply the icons:


- Roms
    - Game Boy Advance (MGBA)
    - Megadrive (MD)
    - 01). Arcade (FBN)
    - .media
        - Game Boy Advance (MGBA).png  <--- RENAMED "(MGBA).png" to "Game Boy Advance (MGBA).png" to MATCH the Rom directory name. This icon should work.
        - Megadrive (MD).png           <--- RENAMED "Sega Genesis (MD).png" to "Megadrive (MD).png" to MATCH the Rom directory name. This icon should work.
        - 01). Arcade (FBN).png        <--- RENAMED "Arcade (FBN).png" to "01). Arcade (FBN).png" to MATCH the Rom directory name. This icon should work.

```
Additionally, Tools can have their own icons as well.

Just like wallpapers, make sure to remember ***additional emulator tags for duplicate systems in NextUI:*** 

```
In our Icon Directory:

- Icons
    - SystemIcons
        - (SUPA).png   <--- For the Supafaust Super Nintendo emulator
        - (SFC).png    <--- For the alternative Super Nintendo emulator
        - (MGBA).png   <--- For the MGBA Game Boy Advance emulator
        - (GBA).png    <--- For the alternative Game Boy Advance emulator

Make sure you include BOTH .png images for icons so that you cover users that might prefer one emulator over the other:


- Roms
    - .media
        - Super Nintendo (SUPA).png  <--- This icon goes to the SNES directory that uses the Supafaust emulator, which is populated with Roms in this example.
        - Super Nintendo (SFC).png   <--- This icon would go to SNES directory that uses the alternative emulator, with NO ROMS in it. This means that this icon will not be displayed.
    - Super Nintendo (SUPA)
        - Chrono Trigger
        - Super Mario RPG
        - Super Mario World
        - etc...
    - Super Nintendo (SFC)
        - (EMPTY ROM DIRECTORY)


In this case, the user is covered through BOTH Super Nintendo emulators.
```

## Overlays



```
Overlays/
├─ MGBA/                    # Overlays for GBA system
│  ├─ overlay1.png
│  └─ overlay2.png
└─ [other systems]/
   └─ [overlay files].png
```

### Important Overlay  Notes

When a `.theme` is applied, we:
1. Copy the `Overlays` directly onto the root of the SD card, where Overlays are stored.
2. That's it!
3. Additionally, if a `.theme` does NOT contain overlays, we delete any previously applied overlays.

Note that the directories containing systems **DO NOT** have parenthesis. That's just how NextUI overlays are stored:

```
~/SD
    - Overlays
        - MGBA               <--- This will work for the (MGBA) emulator. Note the LACK of parenthesis.
            - overlay1.png
            - overlay2.png
            - ...
        - (SUPA)             <--- This will NOT work at all. Parenthesis break overlays!
            - overlay1.png
            - overlay2.png
            - ...
```

Any applied overlays can be tweaked by going to `Settings -> Frontend` while in your game of choice and selecting the preferred overlay.

---

## Fonts



```
Fonts/
├─ OG.ttf                   # Your replacement for font2.ttf
├─ Next.ttf                 # Your replacement for font1.ttf
├─ OG.backup.ttf            # Backup of original font2.ttf
└─ Next.backup.ttf          # Backup of original font1.ttf
```

### Important Fonts Notes

When a `.theme` is applied, we:

1. Copy all four of the above files to `.system/res` on the SD card
2. That's it!

NextUI currently supports 2 fonts, tweakable in the `Settings.pak` tool. 
To change the available themes in the font, simply replace the `OG.ttf` and/or `Next.ttf` font files with whatever you'd like, and the user will be able to cycle between either of them in their `Settings.pak`.

It's important to consider that not all `.pak` apps with NextUI actually _use this system font._ Many use their own hard-coded font. So if some apps don't update the font, this is the reason!

To revert back to the original fonts, just rename the backups in the `.system/res` directory.

---

## Other Settings

In the `manifest.json`, you there are other optional settings that can be stored in a `.theme` pack:

### Accent Colors

```json5
  "accent_colors": {
    "color1": "#FFFFFF",
    "color2": "#9B2257",
    "color3": "#1E2329",
    "color4": "#FFFFFF", 
    "color5": "#000000",
    "color6": "#FFFFFF"
  }
```

- **color1** - Main UI color
- **color2** - Primary accent color
- **color3** - Secondary accent color
- **color4** - List text color
- **color5** - Selected list text color
- **color6** - Hint/information text color

### LED Settings

```json5
  "led_settings": {
    "f1_key": {
      "effect": 1,
      "color1": "0xFFFFFF",
      "color2": "0x000000",
      "speed": 1000,
      "brightness": 100,
      "trigger": 1,
      "in_brightness": 100
    }
```
- **effect** - Lighting effect type (1-7)
- **color1/color2** - Primary and secondary colors (hex format)
- **speed** - Animation speed in milliseconds
- **brightness** - LED brightness (0-100)
- **trigger** - Button/event that activates the LED (1-14)
- **inbrightness** - Information LED brightness (0-100)

**NOTE:** Themes can include settings for **ALL FOUR LEDs.**

---

## Exporting Themes

_Theme Exporting_ involves the full process of locating all six components on your device (if they exist), and copying them to a `.theme` folder inside `Theme-Manager.pak/Exports`. When exporting, we do the following:
1. For wallpapers, we traverse the filesystem in search of all possible `bg.png` files, place them inside `.theme/Wallpapers`, and then _rename them_ according to the directory they came from, like (MGBA), Root, Recently Played, etc.
2. For icons, we traverse the filesystem similarly to wallpapers, looking for all relevant `.png` images that are either literal like `Tools.png`, `Recently Played.png`, or by system tag, like `(MGBA).png`, `Super Nintendo Entertainment System (SUPA).png`, etc, and place them inside `.theme/Icons`.
3. For overlays, we simply pull the `~Overlays` directory and place it inside `.theme/Overlays/Systems`.
4. For fonts, we pull the font files in `./system/res` and place them in `.theme/Fonts`.
5. For accents and LEDs, we pull the settings directly from `.userdata/shared/minuisettings.txt` and `.userdata/shared/ledsettings_brick.txt` and place them inside the `.theme/manifest.json`.

You may then rename the `.theme` folder whatever you'd like, and you can re-import the `.theme` by placing it inside `Theme-Manager.pak/Themes` and applying it via `Installed Themes` in the Theme Manager.

For more details on how to submit/share `.theme` packs, take a look at the [Theme Creation Guide](../documents/THEME_BUILDING.md).

---

## Applying Themes

To apply a theme:

1. Place your `.theme` package in `Tools/tg5040/Theme-Manager.pak/Themes/`
2. Open Theme Manager
3. Navigate to **Installed Themes**
4. Select your theme from the list
5. Confirm the import

During import, Theme Manager will:
- Clear any applied wallpapers and icons (if the theme includes these components)
- Read the manifest.json file
- Copy files to their appropriate system locations
- Create backups of fonts if necessary
- Apply accent and LED settings if included

---
## Index
- [README](../README.md)
- [Theme Creation Guide](../documents/THEME_BUILDING.md)
- [Component Documentation](../documents/COMPONENTS.md)
- [Component Creation Guide](../documents/COMPONENT_BUILDING.md)

## Templates
- [Theme Template](https://github.com/Leviathanium/Template.theme)
- [Wallpaper Template](https://github.com/Leviathanium/Template.bg)
- [Icon Template](https://github.com/Leviathanium/Template.icon)
- [Overlay Template](https://github.com/Leviathanium/Template.over)
- [Font Template](https://github.com/Leviathanium/Template.font)
- [Accent Template](https://github.com/Leviathanium/Template.acc)
- [LED Template](https://github.com/Leviathanium/Template.led)

---

## Troubleshooting

### Theme Import Issues

- **Missing Files**: Ensure all required directories and files exist
- **Permissions**: Make sure all files are readable
- **Manifest Errors**: Check manifest.json for syntax errors
- **System Tags**: Verify system tags in parentheses match your system
- **Logging**: Theme Manager comes with a detailed logger in `Logs/theme-manager.log`. You can always take a look here if there are any issues. Keep in mind this file fills up quickly, so make sure to clear it every once in awhile!
