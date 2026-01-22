# pw_silhouettes

This repository contains aircraft silhouette artwork used by Plane Watch.

Silhouettes are created/curated by Plane Watch admins and are authored in Inkscape (.svg).

The primary goal is to keep silhouettes consistent, clean, and easy to update over time.

The eventual idea is to have a CI/CD workflow that uses ImageMagick or similar to iterate
through these files and rasterise into a spritesheet.

## File format

All silhouettes are stored as SVG files (Inkscape-native).

It is suggested that each file contain the following layers:

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

* Document creation:
  * Load Inkscape.
  * Create a new document with size 70 x 70 px.
* Reference Artwork Layer:
  * Paste reference artwork into its own layer.
  * Rotate and scale down the reference artwork so that the top-down view of the aircraft faces up, and is within the 70 x 70 px document. I find it easier if you set the layer opacity to around 40%.
  * Lock the layer.
  * Rename the layer to "Reference Artwork" (exactly, as this is used by scripts).
* Outline Layer:
  * Create a new layer.
  * Use bezier curves & straight lines to trace the one half (down the middle) of the outline on the reference artwork. I find it helps if you set the stroke of the line to be a different colour and semi-transparent.
  * The path settings should be:
    * Fill: Flat colour, white.
    * Stroke paint: Flat colour, black.
    * Stroke style:
      * Width: 1px
      * Whichever "Join" and "Cap" settings look best for your outline.
  * Once done, copy/paste the path, choose Object > Flip Vertical, then line the other half of the outline up with the reference artwork.
  * Go to the node tool, and ensure the start and finish nodes are snapped to the start and finish nodes of the first half.
  * Select both halves, choose "Path > Union".
  * Use the node tool to perform any tidying up.
* Finalise:
  * Under "Layers and Objects", hide the "Reference Artwork" layer.
  * Select the outline, go to Align and Distribute, center the outline vertically and horizontally with reference to the page.
  * Save your work:
    * The filename should be the IATA code in all caps, followed by ".svg" (eg: A140.svg).
    * The file type should be "Inkscape SVG"
