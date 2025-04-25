# Creating Your Own Theme

This guide will walk you through the process of creating a custom theme for NextUI devices using the Theme Manager. We'll start by exporting your current setup and then modify it to create a completely custom theme.

## 1. Export Your Current Setup

First, let's export your current configuration to use as a starting point:

1. Launch Theme Manager from the Tools menu
2. Select **Export** from the main menu
3. The app will create a new theme package in `Tools/tg5040/Theme-Manager.pak/Exports/` (typically named `theme_1.theme`, `theme_2.theme`, etc.)

This exported `.theme` contains your current wallpapers, icons, accent colors, and other settings, providing an excellent starting point.

## 2. Update Components

Once you've exported your `.theme`, move it to your computer and start working on fine tuning any components you want to add. For help, take a look at the detailed [Theme Documentation.](../documents/THEMES.md)

## 3. Update Metadata

Once you've updated the components, it's time to update the metadata:
- `preview.png` is essential for people to browse your theme. This image should ideally be `1024x768px` and can represent your theme however you'd like. Screenshots are a recommended starting point.
- `manifest.json` is slightly more involved. The best practice for doing this is to use the [template manifest.json here](https://github.com/Leviathanium/Template.theme/blob/main/manifest.json) _instead of_ the one created during export, since its heavily populated with your device's specific metadata.

After downloading the template, here are the recommended changes you may make to that `manifest.json`:

```
{
  "theme_info": {
    "name": "My Theme Template",               <-------- Update your .theme name here
    "version": "1.0.0",
    "author": "Your Name",                     <-------- Update your preferred author name
    "creation_date": "2025-04-22T00:00:00Z", 
    "exported_by": "Theme Manager v1.0.0"
  },
  "content": {                                 <-------- You may SAFELY IGNORE the "content"
    "wallpapers": {                                      section here, REGARDLESS of if you
      "present": false,                                  have any of these components. These
      "count": 0                                         will populate automatically when
    },                                                   the .theme pack is installed.
    "icons": {
      "present": false,
      "system_count": 0,
      "tool_count": 0,
      "collection_count": 0
    },
    "overlays": {
      "present": false,
      "systems": []
    },
    "fonts": {
      "present": true,
      "og_replaced": false,
      "next_replaced": false
    },
    "settings": {
      "accents_included": true,               
      "leds_included": true                   
    }
  },
  "path_mappings": {},                        <-------- You may also SAFELY IGNORE adding path mappings
  "accent_colors": {                          <-------- If you are adding accent colors, you SHOULD update
    "color1": "0xFFFFFF",                               them here. 
    "color2": "0x9B2257",                               
    "color3": "0x1E2329",                               If you are NOT adding accent colors, you may SAFELY
    "color4": "0xFFFFFF",                               REMOVE everything inside the squiggly brackets to look like:
    "color5": "0x000000",
    "color6": "0xFFFFFF"                                "accent_colors": {}
  },
  "led_settings": {                           <-------- Similar to accent_colors. You SHOULD update
    "f1_key": {                                         them here if you are adding them. Otherwise update them:
      "effect": 1,
      "color1": "0xFFFFFF",                             "led_settings": {}
      "color2": "0x000000",
      "speed": 1000,
      "brightness": 100,
      "trigger": 1,
      "in_brightness": 100
    },
    "f2_key": {
      "effect": 1,
      "color1": "0xFFFFFF",
      "color2": "0x000000",
      "speed": 1000,
      "brightness": 100,
      "trigger": 1,
      "in_brightness": 100
    },
    "top_bar": {
      "effect": 1,
      "color1": "0xFFFFFF",
      "color2": "0x000000",
      "speed": 1000,
      "brightness": 100,
      "trigger": 1,
      "in_brightness": 100
    },
    "lr_triggers": {
      "effect": 1,
      "color1": "0xFFFFFF",
      "color2": "0x000000",
      "speed": 1000,
      "brightness": 100,
      "trigger": 1,
      "in_brightness": 100
    }
  }
}

```
Then, place this new `manifest.json` inside your `.theme`. Additionally, create a **backup** of the above `manifest.json`. You'll see why in a moment.

---

## 4. Fine-Tuning

When creating `.theme` packs, importing is a one-way process because your device will **_populate the manifest.json_** when you import the `.theme` for the first time. And you don't want that fully-filled `manifest.json` to get sent with your theme because it can cause compatibility issues.

Here's the best practice for fine-tuning your `.theme`:
1. Once you've filled out the template `manifest.json`, create a backup of it.
2. Place that filled out `manifest.json` inside your `.theme` pack, along with the `preview.png`
3. Attempt a fresh import with the Theme Manager
4. Analyze how accurately the `.theme` imported
5. If something is missing or went wrong, you can take a look at the now-populated `manifest.json` to see what might have gone wrong
6. Repeat this process until you're confident all `.theme` components are working correctly!

**Bonus:** you may also use the [Default.theme](https://github.com/Leviathanium/NextUI-Themes/raw/main/Uploads/Themes/Default.theme.zip) if you'd like to reset your components as you're fine-tuning!

## 5. Sharing and Submitting

There are two ways to share and submit your finished theme:
1. By creating a `.zip` archive of your `.theme` pack and sharing it directly
2. By creating a Gihub repository hosting your `.theme` pack

Both work perfectly fine! If you just wanted the `.zip`, you're all set. This guide will continue for those that are familiar with GitHub repositories and would like to host their own `.theme` for use with the [NextUI Themes Repo.](https://github.com/Leviathanium/NextUI-Themes)

## 6. Create A New Github Repository

To make `.theme` submission/hosting easier, there is a `.theme` template available on GitHub that you can use to properly create/submit your themes:
1. Create a GitHub account if you don't have one yet
2. Navigate to the  [Theme Template Repo here](https://github.com/Leviathanium/Template.theme).
3. Select `Use This Template -> Create a New Repository`
4. Create the repository, ideally with a `.theme` extension like the template


## 7. Clone Your New Repository

Now we need to get your repository onto your computer to work on it:

```
git clone https://github.com/your.theme.git
cd your.theme
git pull
```
The template `.theme` repository should look very similar to the `.theme` pack you've already created. Carefully replace all of the files in your repository, commit, and push:

```
git add .
git commit -m "Initial commit"
git push -u origin main
```

## 8. Submitting Your Theme

You can submit your `.theme` through the various channels available on the NextUI official Discord using your GitHub repo link. We can take it from there!

Alternatively, if you're familiar with pull requests, you can create a special pull request using the `catalog.json` file located in the [Themes Repo.](https://github.com/Leviathanium/NextUI-Themes/blob/main/Catalog/catalog.json) The catalog looks like this:

```json5
  "last_updated": "2025-04-25T16:10:35.873049Z",
  "themes": {
    "Pop-Tarts.theme": {
      "author": "Shin",
      "repository": "https://github.com/KrutzOtrem/pop-tarts.theme",
      "commit": "eb3b3181ef1389e11e352de35531fcd70f991954",
      "preview_path": "Catalog/Themes/previews/Pop-Tarts.theme.png",
      "manifest_path": "Catalog/Themes/manifests/Pop-Tarts.theme.json",
      "description": "Pop-Tarts.theme",
      "URL": "https://github.com/Leviathanium/NextUI-Themes/raw/main/Uploads/Themes/Pop-Tarts.theme.zip"
    },
    "Consolized.theme": {
      "preview_path": "Catalog/Themes/previews/Consolized.theme.png",
      "manifest_path": "Catalog/Themes/manifests/Consolized.theme.json",
      "author": "Gamnrd",
      "description": "Consolized.theme",
      "URL": "https://github.com/Leviathanium/NextUI-Themes/raw/main/Uploads/Themes/Consolized.theme.zip"
    }
```

All you need to do is create a new JSON entry (preferably at the beginning if you want it to show up first in the catalog) that contains the following 3 properties, **exactly like this:**

```json5
"YourThemeName.theme": {
  "author": "YourPreferredName",
  "repository": "https://github.com/YourGitUsername/YourRepo.theme" // <-- Note that this ends in .theme!
  "commit": "a8sd6f89as6d9f8a7sd6" // <---- your full commit hash for whichever commit you want to submit
}, // <--- beginning of next theme
```

Then we can pull the `.theme` and update the catalog automatically!

---
## Index
- [README](../README.md)
- [Theme Documentation](../documents/THEMES.md)
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

- ---