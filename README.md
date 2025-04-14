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

### 1. What is a theme?

A theme consists of a combination of many UI assets, along with lots of metadata. We store themes in folders that end with the extension `.theme` for simplicity. Kind of like `.paks` for apps.

If you open Theme Manager for the first time, navigate to `Themes -> Export` and export your device exactly as it is. Then, navigate to `Tools/tg5040/Theme-Manager.pak/Themes/Exports` to find it. This will make the next part easier to understand.

A `.theme` contains:

1. `manifest.json` file with necessary metadata
2. `preview.png` image to preview the theme
3. `Fonts` directory containing font replacements for Next and OG fonts, as well as backups
4. `Icons` directory containing all system icons in their own directories
5. `Overlays` directory containing any system overlays
6. `Settings` directory containing backups of `ledsettings_brick.txt` and `minuisettings.txt`
7. `Wallpapers` directory containing wallpapers for every directory

You can import themes by placing the `.theme` directory inside `Theme-Manager.pak/Themes/Imports`.

---
### 2. What is `manifest.json?`

`manifest.json` is used by the Theme Manager to identify a lot of the metadata for applying the theme, but it also contains an `author` section:

```
{
  "theme_info": {
    "name": "theme_1",
    "version": "1.0.0",
    "author": "AuthorName",
    "creation_date": "2025-04-13T17:02:44.839046216-05:00",
    "exported_by": "Theme Manager v1.0.0"
  },
```

---




## Sources

- @frysee for literally everything
- @kytz for the work on Noir-Minimal
- @GreenKraken22 for finding and suggesting arcade-dark
- @Fujykky for the work on Screens-Thematic
- Everyone else in the NextUI discord
- Epic Noir theme from https://github.com/c64-dev/es-theme-epicnoir
- All artwork and image source rights go to their respective owners.
