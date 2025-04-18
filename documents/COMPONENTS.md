# Component Packages Documentation

The NextUI Theme Manager uses component packages to manage specific aspects of device customization. This modular approach allows you to mix and match elements from different themes.

## What Are Component Packages?

Component packages are specialized theme elements that focus on a specific customization aspect. Unlike full themes, they only modify a particular part of your device's appearance.

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

The manifest will contain:
```json
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

The manifest will contain:
```json
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

### Accent Components (`.acc`)

```
component_name.acc/
├─ manifest.json
└─ preview.png
```

The manifest will contain:
```json
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

The manifest will contain:
```json
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

The manifest will contain:
```json
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

The manifest will contain:
```json
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

## Creating Component Packages

Component packages can be created in three ways:

1. **Export from Theme Manager**: Use the Components menu to export any component type from your current setup
2. **Deconstruct an existing theme**: Select a theme in the Components menu and choose "Deconstruct..." to break it into components
3. **Manual creation**: Create the directory structure and manifest manually

**NOTE:** Exporting and Deconstructing will place all components in the `Exports` directory of the Theme Manager. They must be moved to re-import them. This is done to prevent unnecessary clutter.

## Using Component Packages

### Importing Components

1. Place the component package in `Tools/tg5040/Theme-Manager.pak/Components/[Type]` where `[Type]` is the component type folder (Wallpapers, Icons, etc.)
2. Launch Theme Manager
3. Navigate to **Components → [Component Type] → Browse**
4. Select the component to import
5. The component will be applied immediately

### Exporting Components

1. Launch Theme Manager
2. Navigate to **Components → [Component Type] → Export**
3. The current configuration for that component will be saved as a package in the `Exports` directory

### Mixing Components

The modular design allows you to mix and match:

1. Apply a full theme as a base
2. Override specific aspects with component packages
3. For instance, apply one theme's wallpapers with another theme's icons and a third theme's accent colors

## Best Practices

1. **Naming conventions**: Follow the naming conventions for files (especially for system wallpapers and icons)
2. **Preview images**: Include a representative preview.png
3. **Resolutions**: Match your device's screen resolution for wallpapers
4. **Font backups**: Always include font backups when creating font components
5. **Metadata**: Include accurate author and version information in your manifests

## Troubleshooting

If a component fails to apply:
1. Check the manifest.json for errors
2. Ensure files are in the correct directories
3. Verify system tags match your system
4. Make sure file permissions allow reading

## Known Issues

1. `preview.png` images are essential for your component to be browsed properly. If you're having trouble with it displaying:
    - Try a different color space. Sometimes the image space makes a difference. In Photoshop, try `Image -> Mode -> RGB Color`
    - If the background is white, you likely need to apply transparency. Check your prefered image editing software for how to properly do this!
    - Image permissions can sometimes be buggy. If you've been editing/renaming images directly on an SD card or over SSH, try doing so directly on your computer, then moving the complete `preview.png` to the correct directory.