# NextUI Theme Manager

A comprehensive theming solution for your NextUI device that lets you customize virtually every visual aspect of your handheld. No more manual file copying or tedious background swapping!

---

<p float="left">
  <img src="/documents/previews/Consolized.theme.png" width="32%" />
  <img src="/documents/previews/Deep-Space.theme.png" width="32%" />
  <img src="/documents/previews/Default.theme.png" width="32%" />
</p>

---

### NOTE: Currently only for the TrimUI Brick.

## Features

- **Complete Theme Management**: Import, export, and apply full theme packages that can customize every visual aspect of your device
- **Component-Level Customization**: Apply specific components like wallpapers, icons, accents, LEDs, fonts, and overlays
- **System-Specific Customization**: Apply custom wallpapers and icons for the main menu for each emulation system
- **Collection Theming**: Customize your collection folders with unique backgrounds and icons
- **Export Your Setup**: Save your current configuration as a shareable theme package
- **Theme Deconstruction**: Break down complex themes into individual components to mix and match your perfect setup

---

## Components You Can Customize

- **Wallpapers**: Change background images for the main menu, systems, tools, collections, and recently played
- **Icons**: Customize icons for all systems, tools, and collection folders
- **Accent Colors**: Modify the UI color scheme of your device
- **LEDs**: Configure custom lighting patterns and colors (TrimUI Brick only)
- **Fonts**: Replace system fonts with custom alternatives
- **Overlays**: Apply system-specific overlay images

---

## Installation

1. Download the latest `Theme-Manager.pak` from the releases page
2. Unzip the `Theme-Manager.pak` folder.
3. Copy it to your device's `Tools/tg5040` directory
4. Launch it from the Tools menu on your device

---

## Getting Started

### Browsing Themes
1. Launch Theme Manager from the Tools menu
2. Select `Sync Catalog` from the main menu to sync with the NextUI Themes repo, available here: https://github.com/Leviathanium/NextUI-Themes
3. Choose `Download Themes` to view the catalog of available themes to download
4. Confirm to download and apply the selected theme
5. You can view any downloaded/installed themes in `Installed Themes` and apply them there

### Managing Components
1. Select `Components` from the main menu
2. Choose the component type (Wallpapers, Icons, etc.)
3. Here, you can download components and apply installed components

### Exporting and Deconstructing
1. Selecting `Export` from the main menu will save your device's current configuration as a `.theme` package
2. Selecting `Export` from any component submenu will save that currently-applied component as its own package (`.bg`, `.icon`, `.over`, etc.)
3. Selecting `Deconstruct...` from the `Components` menu will allow you to deconstruct any installed `.theme` into any available component packages
4. Exported and deconstructed themes and components will be found in `Theme-Manager.pak/Exports` on your SD card.

---

## Documentation

- [Theme Documentation](documents/THEMES.md)
- [Theme Creation Guide](documents/THEME_BUILDING.md)
- [Component Documentation](documents/COMPONENTS.md)
- [Component Creation Guide](documents/COMPONENT_BUILDING.md)

## Credits

- Special thanks to the NextUI community for testing and feedback
- Original concept and development by @Leviathan
- @frysee for literally everything
- @kytz for the work on Noir-Minimal
- @GreenKraken22 for finding and suggesting arcade-dark
- @Fujykky for the work on Screens-Thematic
- Everyone else in the NextUI discord
- Epic Noir theme from https://github.com/c64-dev/es-theme-epicnoir
- All artwork and image source rights go to their respective owners.
