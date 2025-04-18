# Creating Your Own Theme

This guide will walk you through the process of creating a custom theme for NextUI devices using the Theme Manager. We'll start by exporting your current setup and then modify it to create a completely custom theme.

## Getting Started: Export Your Current Setup

First, let's export your current configuration to use as a template:

1. Launch Theme Manager from the Tools menu
2. Select **Export** from the main menu
3. The app will create a new theme package in `Tools/tg5040/Theme-Manager.pak/Exports/` (typically named `theme_1.theme`, `theme_2.theme`, etc.)

This exported theme contains your current wallpapers, icons, accent colors, and other settings, providing an excellent starting point.

## Understanding Your Theme Directory

Navigate to the exported theme on your device or connect your SD card to a computer. You'll see a directory structure like this:

```
theme_1.theme/
├─ manifest.json
├─ preview.png
├─ Wallpapers/
│  ├─ SystemWallpapers/
│  └─ CollectionWallpapers/
├─ Icons/
│  ├─ SystemIcons/
│  ├─ ToolIcons/
│  └─ CollectionIcons/
├─ Fonts/
└─ Overlays/
```

## Creating Custom Wallpapers

Let's start by customizing the wallpapers:

1. Prepare your wallpaper images:
    - Create or obtain PNG images that match your device's resolution (1024x768 for TrimUI Brick)
    - Ensure your images are visually appealing and not too cluttered
    - Consider creating a cohesive set of wallpapers with a consistent style or theme

2. Replace the existing wallpapers:
    - For a global NextUI wallpaper: Add/replace `Wallpapers/SystemWallpapers/Root.png`
    - For system-specific wallpapers: Add files to `Wallpapers/SystemWallpapers/` following the naming convention `System Name (TAG).png`, e.g., `Game Boy Advance (GBA).png`
    - For collection wallpapers: Add files to `Wallpapers/CollectionWallpapers/` named after your collections, e.g., `Favorites.png`

3. Important wallpaper files to consider:
    - `Root.png` - Main menu background
    - `Recently Played.png` - Background for the recently played games list
    - `Tools.png` - Background for the tools menu
    - `Collections.png` - Background for the collections menu

## Designing Custom Icons

Next, let's create custom icons:

1. Prepare your icon images:
    - Create PNG images with transparency
    - Recommended size is at least 200x200 pixels (the system will resize as needed)
    - Use a consistent style for all icons to maintain visual harmony

2. Replace the existing icons:
    - For system icons: Add files to `Icons/SystemIcons/` with the same naming convention used for wallpapers
    - For tool icons: Add files to `Icons/ToolIcons/` named ***literally*** after the tools, e.g., `Settings.png` for the `Settings.pak` package.
    - For collection icons: Add files to `Icons/CollectionIcons/` named after your collections

3. Critical system icons:
    - `Collections.png` - Icon for the collections menu
    - `Recently Played.png` - Icon for the recently played list
    - `Tools.png` - Icon for the tools menu

## Creating Custom Accent Colors

To customize the UI color scheme:

1. Open `manifest.json` in a text editor
2. Locate the `accent_colors` section
3. Modify the color values (in hexadecimal format) for each element:
   ```json
   "accent_colors": {
     "color1": "0xFFFFFF", // Main UI color
     "color2": "0x9B2257", // Primary accent color
     "color3": "0x1E2329", // Secondary accent color
     "color4": "0xFFFFFF", // List text color
     "color5": "0x000000", // Selected list text color
     "color6": "0xFFFFFF"  // Hint text color
   }
   ```
4. Save the file

## Customizing LED Settings (TrimUI Brick)

If you're creating a theme for the TrimUI Brick, you can customize the LED settings:

1. Locate the `led_settings` section in `manifest.json`
2. Modify the settings for each LED zone:
   ```json
   "led_settings": {
     "f1_key": {
       "effect": 1,        // Effect type (1-7)
       "color1": "0xFFFFFF", // Primary color
       "color2": "0x000000", // Secondary color
       "speed": 1000,      // Animation speed
       "brightness": 100,  // Brightness (0-100)
       "trigger": 1,       // Trigger event
       "in_brightness": 100 // Info brightness
     },
     /* Repeat for f2_key, top_bar, and lr_triggers */
   }
   ```
3. Save the file

## Adding Custom Fonts

To replace the system fonts:

1. Prepare your font files in TTF format
2. Add them to the `Fonts/` directory with these names:
    - `OG.ttf` - Replacement for font2.ttf (original system font)
    - `Next.ttf` - Replacement for font1.ttf (alternative system font)
3. If you have the original font backups, include them as:
    - `OG.backup.ttf` - Backup of the original font2.ttf
    - `Next.backup.ttf` - Backup of the original font1.ttf

## Creating System Overlays

If you want to add custom overlays for specific systems:

1. Create PNG overlay images
2. Organize them by system tag **_without parenthesis_** in the `Overlays/` directory:
   ```
   Overlays/
   ├─ MGBA/
   │  ├─ overlay1.png
   │  └─ overlay2.png
   └─ SFC/
      └─ overlay1.png
   ```

## Creating a Preview Image

A good preview image is essential for your theme:

1. Create a `preview.png` image in the root of your theme directory
2. Recommended size is 640×480 pixels
3. The preview should represent the overall style of your theme

## Updating the Manifest

Now update the `manifest.json` file with your theme information:

1. Locate the `theme_info` section
2. Update the metadata:
   ```json
   "theme_info": {
     "name": "My Awesome Theme",
     "version": "1.0.0",
     "author": "Your Name",
     "creation_date": "2025-04-13T12:00:00Z",
     "exported_by": "Theme Manager v1.0"
   }
   ```
3. The `content` section determines what your theme contains. We need to update it properly so that it will scan your `.theme` package and **_automatically populate_** with the correct settings:

```json
"content": {
    "wallpapers": {
      "present": true, <-- Update
      "count": 1
    },
    "icons": {
      "present": false, <-- Update
      "system_count": 0,
      "tool_count": 0,
      "collection_count": 0
    },
    "overlays": {
      "present": false, <-- Update
      "systems": []
    },
    "fonts": {
      "present": true, <-- Update
      "og_replaced": true,
      "next_replaced": true
    },
    "settings": {
      "accents_included": true, <-- Update
      "leds_included": true <-- Update
    }
```

4. Save the file

## Clearing/Updating Path Mappings

The `path_mappings` section of the manifest defines where each file in your theme should be copied when applied. In most cases, you won't need to modify these mappings if you're using the correct file names and directory structure because the `manifest.json` gets updated automatically when we apply the theme to match each user's device.

However, if you're creating a custom file organization, or you're sharing your `.theme` with others, it's wise to clear the current path mappings specifically for wallpapers, icons, and overlays so that they will properly populate on a new users' device. This can also be a way to troubleshoot any issues you may have with applying themes (Rom directories not working, etc.)

To do this, you may update the `manifest.json` to remove all currently-created mappings, as this will reset the paths and prepare it for another users' device, which would look like this:

```json
"path_mappings": {
  "wallpapers": [],
  "icons": [],
  "overlays": [],
  "fonts": {...
  
  // Fonts, Accents, and LEDs can be ignored. Font paths are required, and LEDs/Accents don't have paths.
```


## Testing Your Theme

Time to test your creation:

1. Copy your `.theme` directory to `Tools/tg5040/Theme-Manager.pak/Themes/` on your device
2. Launch Theme Manager
3. Select **Browse Themes** from the main menu
4. Find and select your theme
5. Confirm to apply it
6. Navigate through your device to see how your theme looks

## Iterative Refinement

After seeing your theme in action, you may want to make adjustments:

1. Make changes to the files in your theme directory
2. Re-apply the theme using Theme Manager
3. Keep an eye on the `manifest.json` to confirm everything is working correctly
3. Repeat until you're satisfied with the results

## Sharing Your Theme

Once your theme is complete, you can share it with others:

1. Update the `manifest.json` with your AuthorName, and clear the path mappings (recommended for maximum compatibility)
2. Confirm your `preview.png` works when browsing your theme and that the theme applies correctly from scratch.
3. Create a ZIP archive of your theme directory
4. Share the ZIP file
5. Recipients can extract it to their `Tools/tg5040/Theme-Manager.pak/Themes/` directory

## Theme Component Separation

For more flexibility, you can also break your theme into component packages:

1. Launch Theme Manager
2. Navigate to **Components → Deconstruction...**
3. Select your theme
4. The app will create separate component packages for wallpapers, icons, accents, etc.

This allows others to use specific parts of your theme while keeping their existing preferences for other aspects.

**NOTE:** For complete`.theme` packs, if you choose to deconstruct current theme packages into components to build your own theme, you may **_carefully discard_** each component's `preview.png` and `manifest.json` that gets created when you deconstruct a `.theme` package. Those are exclusively for those wanting to focus on one type of component.

For example, if you deconstructed a `.theme` because you just wanted the `wallpaper.bg` package containing all the wallpapers:

```
Deconstructed-Theme-Name.bg/
├─ manifest.json <-- Can safely disregard.
├─ preview.png <-- Can safely disregard.
├─ SystemWallpapers/
│  ├─ Root.png
│  ├─ Recently Played.png
│  ├─ Tools.png
│  ├─ Collections.png
│  └─ Game Boy Advance (GBA).png
└─ CollectionWallpapers/
   └─ Handhelds.png
```
## Tips and Best Practices

- **System compatibility**: Ensure your system wallpapers and icons use the correct system tags
- **Consistent styling**: Maintain a cohesive visual style across all elements
- **Font legibility**: Test custom fonts thoroughly to ensure they're readable at various sizes
- **Color harmony**: Choose accent colors that work well together
- **File naming**: Follow the naming conventions strictly to ensure proper application
- **Backups**: Always keep backups of your original system fonts
- **Testing**: Test your theme on the actual device, as some elements may look different on-screen
- **Optimization**: Keep file sizes reasonable for faster theme switching
- **Documentation**: Include a README.txt in your theme with any special instructions or credits

## Known Issues

1. `preview.png` images are essential for your `.theme` to be browsed properly. If you're having trouble with it displaying:
   - Try a different color space. Sometimes the image space makes a difference. In Photoshop, try `Image -> Mode -> RGB Color`
   - If the background is white, you likely need to apply transparency. Check your prefered image editing software for how to properly do this!
   - Image permissions can sometimes be buggy. If you've been editing/renaming images directly on an SD card or over SSH, try doing so directly on your computer, then moving the complete `preview.png` to the correct directory.

By following this guide, you'll be able to create beautiful custom themes that enhance your NextUI device experience. Happy theming!