# Changes

## UI

The flow for the user interface will be changing significantly. Here is the outline of the new UI we need to implement:

```
- Themes
    - *Gallery of theme available theme packs from catalog* <--- If theme is installed, label as [Installed] like we do now
        IF NOT DOWNLOADED: - Download Theme?
            - Downloading theme..
                - Theme Downloaded!
                    - Apply Theme? <--- *Create up to 3 backups if turned on, like skyrim autosave*
                        - Applying theme...
                            - Theme Applied!

- Overlays
    - *Gallery of theme available overlay packs from catalog* <--- If overlay pack is installed, label as [Installed] like we do now
        IF NOT DOWNLOADED: - Download Overlays?
            - Downloading overlays...
                - Overlays Downloaded!
                    - Apply Overlays? <--- *Create up to 3 backups if turned on, like skyrim autosave*
                        - Applying overlays...
                            - Overlays Applied!

- Sync Catalog
    - Syncing...
        - Sync complete!
        
- Backup
    - Backup Theme
        - Are you sure? Confirmation
            - Creating theme backup...
                - Theme backup created!
    - Backup Overlays
        - Are you sure? Confirmation
            - Creating overlays backup
                - Overlays backup created!
    - Auto-Backup
        - True/False <-- Have the header message for this screen say "Enable Auto-Backup?"
        
- Revert
    - Revert Theme
        - *gallery of backups*
            - Are you sure? Confirmation
                - Reverting from backup...
                    - Reverted!
    - Revert Overlays
        - *gallery of backups*
            - Are you sure? Confirmation
                - Reverting from backup...
                    - Reverted!

- Purge
    - WARNING: Erase everything?
        - Purging...
            Purge complete!
```

## Packages
1. So there are currently 7 package types: `.theme`, `.icon`, `.acc`, `.font`, `.bg` and `.over`. This will need to change significantly. We will only support _2 package types for now:_
- `.theme` packs, which change **Backgrounds, Icons, and Accents ONLY**
- `.over` packs, which change **OVERLAYS ONLY**

There will no longer be any other pack of any kind.

Deconstruction is also not necessary at all.

## `.theme` packs

Theme packs will be structured like this:

```
- Backgrounds <-- renamed from wallpapers
    - MainMenuBackgrounds <--- OLD "SystemWallpapers", contains your `bg.png` images
    - SystemBackgrounds <--- NEW option, replaces "ListWallpapers". This will contain the `bglist.png` images, named in some format
- Icons
    - MainMenuIcons <--- OLD "SystemIcons". This demonstrates that we are working with MAIN MENU ICONS FOR SYSTEMS, as well as Root, Collections, Tools, Recently Played.
    - CollectionIcons
    - ToolIcons
- Fonts
    - Next.ttf
    - Next.backup.ttf
    - OG.ttf
    - OG.backup.ttf
- Accents
    - accents.txt <--- New file containing all 6 accent colors. Done in favor of moving hard metadata OUT of the manifest.
- preview.png
- manifest.yml  <--- NEW replacement of manifest.json
```

Overlay packs will stay exactly the way they are for now, but the overlay names can change. So they don't have to be `overlay1.png`, etc. They can be any png image.

## `manifest.yml`

Manifest files will work very differently. There will be a section that includes the actual pack data near the top. 