# Creating Your Own Component Packages

This guide will walk you through the process of creating custom component packages for NextUI devices using the Theme Manager. We'll cover each component type separately with step-by-step instructions.

Before creating your own component package, it's recommended to read how `.theme` packs work in the [Theme Documentation.](documents/THEMES.md)

## Wallpaper Components (.bg)

### 1. Export Your Current Wallpapers

Start by exporting your current wallpaper setup:

1. Launch Theme Manager from the Tools menu
2. Select **Components → Wallpapers → Export**
3. The app will create a new wallpaper package in `Tools/tg5040/Theme-Manager.pak/Exports/` (typically named with a timestamp like `wallpaper_20250424_153012.bg`)

### 2. Examine the Wallpaper Structure

Move the exported package to your computer and examine its structure:

```
wallpaper_20250424_153012.bg/
├─ manifest.json
├─ preview.png
├─ SystemWallpapers/
│  ├─ Root.png                # Default global NextUI background
│  ├─ Recently Played.png     # Recently played list
│  ├─ Tools.png               # Tools menu
│  ├─ Collections.png         # Main collections menu
│  └─ [System Name] (TAG).png # System wallpapers with tags
└─ CollectionWallpapers/
   └─ [Collection Name].png   # Named after collection folders
```

### 3. Replace Wallpaper Files

Now replace the wallpaper images with your custom designs:

1. Create your own wallpaper images (recommended resolution: 1024x768)
2. Place system wallpapers in the `SystemWallpapers` folder with proper naming
3. Place collection wallpapers in the `CollectionWallpapers` folder

**Important Naming Conventions:**
- System wallpapers must include the system tag in parentheses (e.g., `Game Boy Advance (GBA).png` or just `(GBA).png`)
- Include wallpapers for both emulator variants if applicable (e.g., both `(MGBA).png` and `(GBA).png`)
- Collection wallpapers must match the exact collection folder name (e.g., `Handhelds.png`)

### 4. Update the Preview Image

Create a compelling preview image:

1. Create an image that showcases your wallpaper design (1024x768px recommended)
2. Save it as `preview.png` in the root of your package
3. This image will be displayed in the Theme Manager when browsing components

### 5. Update the Manifest

Edit the `manifest.json` file to update basic information:

```json5
{
  "component_info": {
    "name": "MyAwesomeWallpapers",  // Update with your component name
    "type": "wallpaper",           // Don't change this
    "version": "1.0.0",            // Version number
    "author": "Your Name",         // Update with your name
    "creation_date": "2025-04-24T00:00:00Z",
    "exported_by": "Theme Manager v1.0.0"
  }
  // Don't worry about the other sections - they'll be automatically populated during import
}
```

### 6. Complete Wallpaper Package

Your final wallpaper package should look like this:

```
MyAwesomeWallpapers.bg/
├─ manifest.json
├─ preview.png
├─ SystemWallpapers/
│  ├─ Root.png
│  ├─ Recently Played.png
│  ├─ Tools.png
│  ├─ Collections.png
│  ├─ Game Boy Advance (GBA).png
│  ├─ Game Boy Advance (MGBA).png
│  └─ [Other system wallpapers with tags]
└─ CollectionWallpapers/
   └─ [Collection wallpapers]
```

## Overlay Components (.over)

### 1. Export Your Current Overlays

Start by exporting your current overlay setup:

1. Launch Theme Manager from the Tools menu
2. Select **Components → Overlays → Export**
3. The app will create a new overlay package in `Tools/tg5040/Theme-Manager.pak/Exports/` (typically named with a timestamp like `overlay_20250424_153012.over`)

### 2. Examine the Overlay Structure

Move the exported package to your computer and examine its structure:

```
overlay_20250424_153012.over/
├─ manifest.json
├─ preview.png
└─ Systems/
   ├─ MGBA/                   # Note: NO parentheses in folder names
   │  ├─ overlay1.png
   │  └─ overlay2.png
   └─ [Other system folders]/
      └─ [Overlay files].png
```

### 3. Create or Replace Overlay Files

Now customize your overlay files:

1. Create custom overlay PNG images for each system
2. Organize them in the appropriate system folders (without parentheses in folder names)
3. Each system can have multiple overlay options

**Important Notes:**
- System folders should NOT include parentheses (e.g., use `MGBA` not `(MGBA)`)
- Overlays are PNG files with transparency where needed
- Test overlays with actual games to ensure proper positioning

### 4. Update the Preview Image

Create a preview image:

1. Create an image that represents your overlay designs (1024x768px recommended)
2. Save it as `preview.png` in the root of your package
3. This can be a collage showing different overlays or a single representative design

### 5. Update the Manifest

Edit the `manifest.json` file to update basic information:

```json5
{
  "component_info": {
    "name": "MyAwesomeOverlays",  // Update with your component name
    "type": "overlay",           // Don't change this
    "version": "1.0.0",          // Version number
    "author": "Your Name",       // Update with your name
    "creation_date": "2025-04-24T00:00:00Z",
    "exported_by": "Theme Manager v1.0.0"
  }
  // Don't worry about the other sections - they'll be automatically populated during import
}
```

### 6. Complete Overlay Package

Your final overlay package should look like this:

```
MyAwesomeOverlays.over/
├─ manifest.json
├─ preview.png
└─ Systems/
   ├─ MGBA/
   │  ├─ overlay1.png
   │  └─ overlay2.png
   ├─ SFC/
   │  └─ overlay1.png
   └─ [Other system folders]/
      └─ [Overlay files].png
```

## Icon Components (.icon)

### 1. Export Your Current Icons

Start by exporting your current icon setup:

1. Launch Theme Manager from the Tools menu
2. Select **Components → Icons → Export**
3. The app will create a new icon package in `Tools/tg5040/Theme-Manager.pak/Exports/` (typically named with a timestamp like `icon_20250424_153012.icon`)

### 2. Examine the Icon Structure

Move the exported package to your computer and examine its structure:

```
icon_20250424_153012.icon/
├─ manifest.json
├─ preview.png
├─ SystemIcons/             # Icons displayed in main menu
│  ├─ Collections.png
│  ├─ Recently Played.png
│  ├─ Tools.png
│  └─ [System Name] (TAG).png
├─ ToolIcons/               # Icons for individual tools
│  └─ [Tool Name].png       # Named as shown in Tools folder
└─ CollectionIcons/         # Icons for collections
   └─ [Collection Name].png # Named as shown in Collections
```

### 3. Replace Icon Files

Now replace the icon images with your custom designs:

1. Create your own icon images (recommended: square PNG files, 256x256px or larger)
2. Place system icons in the `SystemIcons` folder with proper naming
3. Place tool icons in the `ToolIcons` folder
4. Place collection icons in the `CollectionIcons` folder

**Important Naming Conventions:**
- System icons must include the system tag in parentheses (e.g., `Game Boy Advance (GBA).png` or just `(GBA).png`)
- Include icons for both emulator variants if applicable (e.g., both `(MGBA).png` and `(GBA).png`)
- Tool icons must match the exact tool name (e.g., `Battery.png`)
- Collection icons must match the exact collection folder name (e.g., `Favorites.png`)

### 4. Update the Preview Image

Create a compelling preview image:

1. Create an image that showcases your icon designs (1024x768px recommended)
2. Save it as `preview.png` in the root of your package
3. This can be a collage showing different icons or a representative selection

### 5. Update the Manifest

Edit the `manifest.json` file to update basic information:

```json5
{
  "component_info": {
    "name": "MyAwesomeIcons",    // Update with your component name
    "type": "icon",              // Don't change this
    "version": "1.0.0",          // Version number
    "author": "Your Name",       // Update with your name
    "creation_date": "2025-04-24T00:00:00Z",
    "exported_by": "Theme Manager v1.0.0"
  }
  // Don't worry about the other sections - they'll be automatically populated during import
}
```

### 6. Complete Icon Package

Your final icon package should look like this:

```
MyAwesomeIcons.icon/
├─ manifest.json
├─ preview.png
├─ SystemIcons/
│  ├─ Collections.png
│  ├─ Recently Played.png
│  ├─ Tools.png
│  ├─ Game Boy Advance (GBA).png
│  ├─ Game Boy Advance (MGBA).png
│  └─ [Other system icons with tags]
├─ ToolIcons/
│  └─ [Tool icons]
└─ CollectionIcons/
   └─ [Collection icons]
```

## Font Components (.font)

### 1. Export Your Current Fonts

Start by exporting your current font setup:

1. Launch Theme Manager from the Tools menu
2. Select **Components → Fonts → Export**
3. The app will create a new font package in `Tools/tg5040/Theme-Manager.pak/Exports/` (typically named with a timestamp like `font_20250424_153012.font`)

### 2. Examine the Font Structure

Move the exported package to your computer and examine its structure:

```
font_20250424_153012.font/
├─ manifest.json
├─ preview.png
├─ OG.ttf                   # Font for "OG" theme selection
├─ Next.ttf                 # Font for "Next" theme selection
├─ OG.backup.ttf            # Backup of original font2.ttf
└─ Next.backup.ttf          # Backup of original font1.ttf
```

### 3. Replace Font Files

Now replace the font files with your custom fonts:

1. Obtain or create TTF font files for your package
2. Replace `OG.ttf` and/or `Next.ttf` with your custom fonts
3. Keep the backup font files if possible, as they provide a way to restore original fonts

**Important Notes:**
- Use high-quality, readable fonts
- Test fonts at small sizes to ensure readability in menus
- Not all applications will use these system fonts; some apps use their own embedded fonts

### 4. Update the Preview Image

Create a preview image:

1. Create an image that demonstrates your fonts in use (1024x768px recommended)
2. Save it as `preview.png` in the root of your package
3. This can show text samples in different sizes and contexts

### 5. Update the Manifest

Edit the `manifest.json` file to update basic information:

```json5
{
  "component_info": {
    "name": "MyAwesomeFonts",    // Update with your component name
    "type": "font",              // Don't change this
    "version": "1.0.0",          // Version number
    "author": "Your Name",       // Update with your name
    "creation_date": "2025-04-24T00:00:00Z",
    "exported_by": "Theme Manager v1.0.0"
  }
  // Don't worry about the other sections - they'll be automatically populated during import
}
```

### 6. Complete Font Package

Your final font package should look like this:

```
MyAwesomeFonts.font/
├─ manifest.json
├─ preview.png
├─ OG.ttf
├─ Next.ttf
├─ OG.backup.ttf
└─ Next.backup.ttf
```

## Accent Components (.acc)

### 1. Export Your Current Accents

Start by exporting your current accent colors:

1. Launch Theme Manager from the Tools menu
2. Select **Components → Accents → Export**
3. The app will create a new accent package in `Tools/tg5040/Theme-Manager.pak/Exports/` (typically named with a timestamp like `accent_20250424_153012.acc`)

### 2. Examine the Accent Structure

Move the exported package to your computer and examine its structure:

```
accent_20250424_153012.acc/
├─ manifest.json
└─ preview.png
```

Unlike other components, accent packages don't contain additional files beyond the manifest and preview image. All the color information is stored in the manifest.

### 3. Update the Accent Colors

Edit the `manifest.json` file to customize your accent colors:

```json5
{
  "component_info": {
    "name": "MyAwesomeAccents",   // Update with your component name
    "type": "accent",             // Don't change this
    "version": "1.0.0",           // Version number
    "author": "Your Name",        // Update with your name
    "creation_date": "2025-04-24T00:00:00Z",
    "exported_by": "Theme Manager v1.0.0"
  },
  "accent_colors": {
    "color1": "0xFFFFFF",         // Main UI color
    "color2": "0x9B2257",         // Primary accent color
    "color3": "0x1E2329",         // Secondary accent color
    "color4": "0xFFFFFF",         // List text color
    "color5": "0x000000",         // Selected list text color
    "color6": "0xFFFFFF"          // Hint/information text color
  }
}
```

Customize the hex color values to create your desired color scheme. All colors use the format `0xRRGGBB`.

### 4. Update the Preview Image

Create a preview image:

1. Create an image that showcases your accent color scheme (1024x768px recommended)
2. Save it as `preview.png` in the root of your package
3. This can show UI elements with your color scheme applied or a color palette

### 5. Complete Accent Package

Your final accent package should look like this:

```
MyAwesomeAccents.acc/
├─ manifest.json
└─ preview.png
```

## LED Components (.led)

### 1. Export Your Current LED Settings

Start by exporting your current LED configuration:

1. Launch Theme Manager from the Tools menu
2. Select **Components → LEDs → Export**
3. The app will create a new LED package in `Tools/tg5040/Theme-Manager.pak/Exports/` (typically named with a timestamp like `led_20250424_153012.led`)

### 2. Examine the LED Structure

Move the exported package to your computer and examine its structure:

```
led_20250424_153012.led/
└─ manifest.json
```

LED packages are the simplest component type, containing only the manifest file with all settings stored inside it.

### 3. Update the LED Settings

Edit the `manifest.json` file to customize your LED configuration:

```json5
{
  "component_info": {
    "name": "MyAwesomeLEDs",    // Update with your component name
    "type": "led",              // Don't change this
    "version": "1.0.0",         // Version number
    "author": "Your Name",      // Update with your name
    "creation_date": "2025-04-24T00:00:00Z",
    "exported_by": "Theme Manager v1.0.0"
  },
  "led_settings": {
    "f1_key": {
      "effect": 1,              // Lighting effect type (1-7)
      "color1": "0x8833FF",     // Primary color
      "color2": "0x000000",     // Secondary color
      "speed": 1000,            // Animation speed in milliseconds
      "brightness": 100,        // LED brightness (0-100)
      "trigger": 1,             // Button/event trigger (1-14)
      "in_brightness": 100      // Information LED brightness
    },
    "f2_key": {
      // Same structure as f1_key
    },
    "top_bar": {
      // Same structure as f1_key
    },
    "lr_triggers": {
      // Same structure as f1_key
    }
  }
}
```

Customize the LED effects, colors, speeds, and other parameters to create your desired lighting effects.

**Effect Types:**
1. Static color
2. Breathing effect
3. Color cycle
4. Rainbow effect
5. Strobe effect
6. Wave effect
7. Random patterns

### 4. Complete LED Package

Your final LED package should look like this:

```
MyAwesomeLEDs.led/
└─ manifest.json
```

Unlike other components, LED packages typically don't include a preview image.

---

# Sharing and Submitting Your Components

Once you've created your component package, you can share it with the community. Here's how:

## 1. Testing Your Component

Before sharing, test your component:

1. Move your component package to the appropriate folder:
   - `Tools/tg5040/Theme-Manager.pak/Components/[Type]/`
   - For example: `Components/Wallpapers/` for a `.bg` package
2. Launch Theme Manager
3. Navigate to **Components → [Component Type] → Installed**
4. Apply your component and verify everything works correctly

## 2. Package for Sharing

There are two ways to share your components:

### Option 1: Direct Zip Sharing

1. Create a ZIP archive of your component folder (keeping the extension in the filename):
   - For example: `MyAwesomeOverlays.over.zip` (not just `MyAwesomeOverlays.zip`)
2. Share this ZIP file with others, who can extract it to their Components directory

### Option 2: GitHub Template Repository (Recommended)

We've provided six different templates, one for each component type, so you have an easy place to start when creating your own Git repository:
- [Wallpaper Template](https://github.com/Leviathanium/Template.bg)
- [Icon Template](https://github.com/Leviathanium/Template.icon)
- [Overlay Template](https://github.com/Leviathanium/Template.over)
- [Font Template](https://github.com/Leviathanium/Template.font)
- [Accent Template](https://github.com/Leviathanium/Template.acc)
- [LED Template](https://github.com/Leviathanium/Template.led)






### Example: Creating an Overlay Component Repository

Here's how to create a repository for an overlay component:

1. Navigate to the `Template.over` repository
2. Select "Use This Template → Create a New Repository"
3. Name your repository (e.g., `MyAwesomeOverlays.over`)
4. Clone your new repository to your computer:
   ```
   git clone https://github.com/yourusername/MyAwesomeOverlays.over.git
   cd MyAwesomeOverlays.over
   git pull
   ```
5. Replace all files with your component files
6. Update the README.md with information about your component
7. Commit and push your changes:
   ```
   git add .
   git commit -m "Initial commit"
   git push -u origin main
   ```
8. Share your repository URL with others

The Template.over repository is structured identically to a regular exported `.over` package, making it easy to transition from your exported component to a hosted repository.

## 3. Submitting to the Community

You can submit your component through various channels:

1. Share your ZIP file or GitHub repository URL in the NextUI Discord community
2. Submit a pull request to have your component added to the official theme catalog
3. Include screenshots or a GIF showing your component in action to attract interest

---

## Index
- [README](../README.md)
- [Theme Documentation](../documents/THEMES.md)
- [Theme Creation Guide](../documents/THEME_BUILDING.md)
- [Component Documentation](../documents/COMPONENTS.md)
- ---