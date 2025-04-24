# Component Packages Documentation

The NextUI Theme Manager uses component packages to manage specific aspects of device customization. This modular approach allows you to mix and match elements from different themes, or build just one component, like overlays or wallpapers, and distribute them as an individual package.

Before reading about component packages it is **highly recommended** you start with `.theme` packages to fully understand how each component works. You can find that here in the [Theme Documentation](../documents/THEMES.md).

## What Are Component Packages?

Component packages are located in `Theme-Manager.pak/Components` and they are specialized theme elements that focus on a specific customization aspect. Unlike full themes, they only modify a single part of your device's appearance.

There are six types of component packages:

1. **Wallpaper** (`.bg`) - Background images
2. **Icon** (`.icon`) - System, tool, and collection icons
3. **Accent** (`.acc`) - UI color schemes
4. **LED** (`.led`) - LED lighting configurations
5. **Font** (`.font`) - System font replacements
6. **Overlay** (`.over`) - System-specific overlays

## Component Structure

Each component package is a directory with the appropriate extension (e.g., `.bg`, `.icon`) containing:

1. **manifest.json** - Component metadata and mappings
2. **preview.png** - Component preview image 
3. **Files specific to the component type** (wallpapers, icons, fonts, etc.)

Individual components are structured **identically** to the way they are stored in `.theme` packs. The only difference is that each component has **it's own `manifest.json` and `preview.png`** so that it can be installed and applied as its own package. 

## Detailed Component Structures

### Wallpaper Components (`.bg`)

```
component_name.bg/
├─ manifest.json
├─ preview.png
├─ SystemWallpapers/
│  ├─ Root.png
│  ├─ Recently Played.png
│  ├─ Tools.png
│  ├─ Collections.png
│  └─ Game Boy Advance (GBA).png
└─ CollectionWallpapers/
   └─ Handhelds.png
```

When installed, the manifest will contain:
```json5
{
  "component_info": {
    "name": "component_name",
    "type": "wallpaper",
    "version": "1.0.0",
    "author": "AuthorName",
    "creation_date": "2025-04-13T12:00:00Z",
    "exported_by": "Theme Manager v1.0"
  },
  "content": {
    "count": 6,
    "system_wallpapers": ["Root.png", "Recently Played.png", "Game Boy Advance (GBA).png"],
    "collection_wallpapers": ["Handhelds.png"]
  },
  "path_mappings": [
    {
      "theme_path": "SystemWallpapers/Root.png",
      "system_path": "/mnt/SDCARD/bg.png",
      "metadata": {
        "SystemName": "Root",
        "WallpaperType": "Main"
      }
    },
    /* Additional wallpaper mappings */
  ]
}
```

### Icon Components (`.icon`)

```
component_name.icon/
├─ manifest.json
├─ preview.png
├─ SystemIcons/
│  ├─ Collections.png
│  ├─ Recently Played.png
│  ├─ Tools.png
│  └─ Game Boy Advance (GBA).png
├─ ToolIcons/
│  └─ Battery.png
└─ CollectionIcons/
   └─ Favorites.png
```

When installed, the manifest will contain:
```json5
{
  "component_info": {
    "name": "component_name",
    "type": "icon",
    "version": "1.0.0",
    "author": "AuthorName",
    "creation_date": "2025-04-13T12:00:00Z",
    "exported_by": "Theme Manager v1.0"
  },
  "content": {
    "system_count": 4,
    "tool_count": 1,
    "collection_count": 1,
    "system_icons": ["Collections.png", "Recently Played.png", "Tools.png", "Game Boy Advance (GBA).png"],
    "tool_icons": ["Battery.png"],
    "collection_icons": ["Favorites.png"]
  },
  "path_mappings": [
    {
      "theme_path": "SystemIcons/Collections.png",
      "system_path": "/mnt/SDCARD/.media/Collections.png",
      "metadata": {
        "SystemName": "Collections",
        "IconType": "System"
      }
    },
    /* Additional icon mappings */
  ]
}
```

### Overlay Components (`.over`)

```
component_name.over/
├─ manifest.json
├─ preview.png
└─ Systems/
   ├─ MGBA/
   │  ├─ overlay1.png
   │  └─ overlay2.png
   └─ [other systems]/
      └─ [overlay files].png
```

When installed, the manifest will contain:
```json5
{
  "component_info": {
    "name": "component_name",
    "type": "overlay",
    "version": "1.0.0",
    "author": "AuthorName",
    "creation_date": "2025-04-13T12:00:00Z",
    "exported_by": "Theme Manager v1.0"
  },
  "content": {
    "systems": ["MGBA", "SFC", "MD"]
  },
  "path_mappings": [
    {
      "theme_path": "Systems/MGBA/overlay1.png",
      "system_path": "/mnt/SDCARD/Overlays/MGBA/overlay1.png",
      "metadata": {
        "SystemTag": "MGBA",
        "OverlayName": "overlay1.png"
      }
    },
    /* Additional overlay mappings */
  ]
}
```

### Accent Components (`.acc`)

```
component_name.acc/
├─ manifest.json
└─ preview.png
```

The manifest will always contain:
```json5
{
  "component_info": {
    "name": "component_name",
    "type": "accent",
    "version": "1.0.0",
    "author": "AuthorName",
    "creation_date": "2025-04-13T12:00:00Z",
    "exported_by": "Theme Manager v1.0"
  },
  "accent_colors": {
    "color1": "0xFFFFFF",
    "color2": "0x9B2257",
    "color3": "0x1E2329",
    "color4": "0xFFFFFF",
    "color5": "0x000000",
    "color6": "0xFFFFFF"
  }
}
```

### LED Components (`.led`)

```
component_name.led/
└─ manifest.json
```

The manifest will always contain:
```json5
{
  "component_info": {
    "name": "component_name",
    "type": "led",
    "version": "1.0.0",
    "author": "AuthorName",
    "creation_date": "2025-04-13T12:00:00Z",
    "exported_by": "Theme Manager v1.0"
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
    "f2_key": {
      /* F2 key settings */
    },
    "top_bar": {
      /* Top bar settings */
    },
    "lr_triggers": {
      /* L&R triggers settings */
    }
  }
}
```

### Font Components (`.font`)

```
component_name.font/
├─ manifest.json
├─ preview.png
├─ OG.ttf
├─ Next.ttf
├─ OG.backup.ttf
└─ Next.backup.ttf
```

When installed, the manifest will contain:
```json5
{
  "component_info": {
    "name": "component_name",
    "type": "font",
    "version": "1.0.0",
    "author": "AuthorName",
    "creation_date": "2025-04-13T12:00:00Z",
    "exported_by": "Theme Manager v1.0"
  },
  "content": {
    "og_replaced": true,
    "next_replaced": true
  },
  "path_mappings": {
    "OG": {
      "theme_path": "OG.ttf",
      "system_path": "/mnt/SDCARD/.userdata/shared/font2.ttf"
    },
    "Next": {
      "theme_path": "Next.ttf",
      "system_path": "/mnt/SDCARD/.userdata/shared/font1.ttf"
    },
    "OG.backup": {
      "theme_path": "OG.backup.ttf",
      "system_path": "/mnt/SDCARD/.userdata/shared/font2.ttf.bak"
    },
    "Next.backup": {
      "theme_path": "Next.backup.ttf",
      "system_path": "/mnt/SDCARD/.userdata/shared/font1.ttf.bak"
    }
  }
}
```

---

## Exporting

_Exporting_ is the process of copying the **currently applied components** on your device and generating a component package with those resources, located in `Theme-Manager.pak/Exports`. When exporting, the package is given a hash and shows up in this directory. You can export any component you'd like in the Theme Manager under `Components -> (component) -> Export`.

Let's say you wanted to export the current wallpapers active on your device:

```
- Theme-Manager.pak
   - Exports
      - wallpaper_90283742982734.bg <--- Newly exported wallpaper pack, you can name this whatever you'd like!
```

To then re-import it on its own, you may **move this pack to its respective folder** inside `Theme-Manager.pak/Components`, then browse to it inside the Theme Manager and apply it:

```
- Theme-Manager.pak
   - Components
      - Wallpapers
         - MyWallpaperPack.bg <--- Theme Manager will be able to find this pack!
```

This approach allows you to mix and match the components on your device to make it however you'd like!

For details on how to submit and share your component packages, read the [Component Creation Guide](../documents/COMPONENT_BUILDING.md).

---

## Deconstruction

_Deconstruction_
 is the process of breaking a full `.theme` package into supported component packages (`.bg`, `.icon`, `.overlay`, etc.)
You can find this option in `Components -> Deconstruct...` in the Theme Manager. It will allow you to select a currently installed `.theme` package inside your `Theme-Manager.pak/Themes` directory and deconstruct it. Keep in mind, the `.theme` pack **does not need to be active for this to work**. You just need the `.theme` package to be located on your device.

When deconstructing, ALL components created will be found inside your `Theme-Manager.pak/Exports` directory upon completion:

```
- Theme-Manager.pak
   - Themes
      - Consolized.theme <--- Let's say we wanted to deconstruct @Gamrnd's excellent Consolized .theme pack into component packages.
   - Exports
      - Consolized.acc   <--- Here is the accent package.
      - Consolized.bg    <--- Here is the wallpaper package.
      - Consolized.icon  <--- Here is the icon package.
      - ...
```
It's important to consider that deconstruction only creates these packages for components inside the deconstructed `.theme`. For example, there won't be a `.icon` pack if a `.theme` doesn't contain icons to begin with!

If you'd like to re-import these components packages or work with them, you **must move them** to their respective directory inside `Theme-Manager.pak/Components` to be able to re-import them. You can't just re-import them from the `Exports` directory!

---

## Important Component Package Notes

Component packages can be created in three ways:

1. **Export from Theme Manager**: Use the Components menu to export any component type from your current setup
2. **Deconstruct an existing theme**: Select a theme in the Components menu and choose "Deconstruct..." to break it into components
3. **Manual creation**: Create the directory structure and manifest manually (not recommended though!)

**NOTE:** Exporting and Deconstructing will place all components in the `Exports` directory of the Theme Manager. They must be moved to re-import them. This is done to prevent unnecessary clutter.



### Applying Components

1. Place the component package in `Tools/tg5040/Theme-Manager.pak/Components/[Type]` where `[Type]` is the component type folder (Wallpapers, Icons, etc.)
2. Launch Theme Manager
3. Navigate to **Components → [Component Type] → Installed**
4. Select the component to import
5. The component will be applied immediately

---

## Index
- [README](../README.md)
- [Theme Documentation](../documents/THEMES.md)
- [Theme Creation Guide](../documents/THEME_BUILDING.md)
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