# pw_silhouettes

This repository contains aircraft silhouette artwork used by Plane Watch.

Silhouettes are created/curated by Plane Watch admins and are authored in Inkscape (.svg).

The primary goal is to keep silhouettes consistent, clean, and easy to update over time.

The eventual idea is to have a CI/CD workflow that uses ImageMagick or similar to iterate
through these files and rasterise into a spritesheet.

## File format

All silhouettes are stored as SVG files (Inkscape-native).

Each file must contain the following layers:

### Reference Artwork

Contains whatever reference artwork was used to create the line art (photo, diagram, manufacturer drawing, etc).

Notes:

- This layer may include raster images and/or imported vector artwork.
- This layer is intended for editing only.
- This layer should generally be hidden / not exported.

### Outline

Contains the traced silhouette.

Notes:

- This is the authoritative silhouette.
- Must be a clean outline suitable for export to other formats.
- Avoid unnecessary nodes.
- Keep paths closed where possible.

## Process

* Load Inkscape.
* Paste reference artwork into its own layer, lock the layer.
* Create a new layer. Use bezier curves & straight lines to trace the outline on the reference artwork. I find it helps if you set the stroke of the line to be a different colour and semi-transparent.
