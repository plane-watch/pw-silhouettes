# Contributing Aircraft Silhouettes

Thank you for helping expand the Plane Watch aircraft silhouette library! ✈️

This project relies on consistent artwork and metadata so silhouettes render correctly in the UI. Please follow the guidelines below when adding or updating aircraft.

---

## 📁 What You Need to Provide

For each new airframe:

1. **An SVG silhouette** (or multiple if animated)
2. **An airframe definition JSON file**

Both must pass automated validation checks that run on pull requests.

---

## SVG Creation Process

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
  * Ensure all visible artwork adheres to the [styling rules for visible artwork](#-styling-rules-for-visible-artwork).
  * Select the outline, go to Align and Distribute, center the outline vertically and horizontally with reference to the page.
  * Save your work in the `silhouettes` directory:
    * The filename should be the ICAO code in all caps, followed by ".svg" (eg: A140.svg).
    * For animations, the filename should be appended with the frame number (eg: B06-1.svg for the first frame, B06-2.svg for the second frame...)
    * The file type should be "Inkscape SVG"
* Create airframe JSON file (see below)

---

## 📄 Airframe JSON Format

Each airframe has a JSON definition describing how it should appear in the spritesheet.

Schema file:

```
schemas/airframe.input.runtime.v1.schema.json
```

Pull requests automatically validate JSON files against this schema. If validation fails, the PR will show exactly which file and field is incorrect.

The best thing to do is duplicate an existing airframe JSON and edit it.

### File location and naming

- JSON files should be saved under the `airframes` directory.
- One JSON file per ICAO designator, e.g. `airframes/A306.json`, `airframes/A30B.json`, `airframes/B412.json`.
- SVG assets are referenced by path inside the JSON (e.g. `silhouettes/A306.svg`).

### Top-level structure

```json
{
  "version": 1,
  "icao": {
    "designator": "A306",
    "typeCode": "L2J",
    "wakeCategory": "H"
  },
  "aliasOf": null,
  "render": {
    "scale": 1,
    "anchor": { "x": 35, "y": 30 },
    "noRotate": false
  },
  "art": {
    "frames": [
      { "src": "silhouettes/A306.svg" }
    ],
    "frameTime": null
  },
  "notes": ""
}
```

### Field reference

#### `version` (integer, optional)

Schema version for the JSON file. If omitted, the effective version is `1`.

#### `icao` (object, required)

ICAO classification metadata.

- `designator` (string, required): ICAO aircraft type **designator** (e.g. `"A306"`, `"B412"`).
- `typeCode` (string, required): ICAO aircraft type **description code** (e.g. `"L2J"`, `"H2T"`).
  - Examples:
    - `L2J` = Landplane, 2 engines, Jet
    - `H2T` = Helicopter, 2 engines, Turbine (turboshaft)
- `wakeCategory` (string, required): Wake turbulence category:
  - `"L"` (Light), `"M"` (Medium), `"H"` (Heavy), `"J"` (Super)

#### `aliasOf` (string or null, optional)

If set to another designator, this airframe is treated as an **alias** of that designator.

- Use this when silhouettes are “close enough” at 70×70 and you don’t want duplicated artwork.
- Alias files may omit `render`, `noRotate`, and `art` entirely if the runtime should inherit from the canonical designator.

#### `render` (object, optional)

Runtime rendering hints.

- `scale` (number, optional, default `1`): Size multiplier applied at runtime.

  - Use to make a large aircraft (e.g. A225) appear larger than a small aircraft (e.g. SONX).
- `anchor` (object, optional, default `{x:35,y:35}`):

  - `x` (number): anchor X in pixels within the 70×70 cell
  - `y` (number): anchor Y in pixels within the 70×70 cell
  - For most aircraft, anchor should be near centre-of-mass rather than geometric centre.
- `noRotate` (boolean, optional, default `false`)

  - If `true`, the icon is **not rotated** by heading/track.
  - Intended for things like balloons, ground objects, towers, etc.
  - For aircraft/helicopters this is usually `false`.

#### `art` (object, optional)

Defines the artwork source(s) for this airframe.

- `frames` (array of objects, required if `art` is present):
  - Each frame object:
    - `src` (string, required): relative path to an SVG asset
- `frameTime` (integer or null, optional):
  - If `null` (or absent): the airframe is **static** (single frame).
  - If an integer: the airframe is animated, and the value is the per-frame time in milliseconds.

#### `notes` (string, optional)

Human-readable notes (why it aliases, source info, caveats, etc.).

### Defaults

If fields are omitted, the runtime should behave as though these defaults were applied:

- `version`: `1`
- `aliasOf`: `null`
- `render.scale`: `1`
- `render.anchor`: `{ "x": 35, "y": 35 }` (centre of 70×70 cell)
- `render.noRotate`: `false`
- `art.frameTime`: `null` (static)

### Examples

#### Canonical static airframe (A306)

```json
{
  "icao": { "designator": "A306", "typeCode": "L2J", "wakeCategory": "H" },
  "aliasOf": null,
  "render": { 
    "scale": 1,
    "anchor": { "x": 35, "y": 30 },
    "noRotate": false
  },
  "art": { "frames": [{ "src": "silhouettes/A306.svg" }], "frameTime": null },
  "notes": ""
}
```

#### Alias airframe (A30B → A306)

```json
{
  "icao": { "designator": "A30B", "typeCode": "L2J", "wakeCategory": "H" },
  "aliasOf": "A306",
  "notes": "At 70x70 pixels, the Airbus A300B2/A300B4/A300C4 look the same as the Airbus A300-600."
}
```

#### Animated airframe (B412)

```json
{
  "icao": { "designator": "B412", "typeCode": "H2T", "wakeCategory": "M" },
  "aliasOf": null,
  "render": { 
    "scale": 1,
    "anchor": { "x": 35, "y": 30 },
    "noRotate": false
  },
  "art": {
    "frames": [
      { "src": "silhouettes/B412-1.svg" },
      { "src": "silhouettes/B412-2.svg" },
      { "src": "silhouettes/B412-3.svg" }
    ],
    "frameTime": 50
  },
  "notes": ""
}
```

## 🎨 SVG Silhouette Requirements

SVGs must follow strict styling and structure rules to ensure visual consistency.

### 📐 Size

The root `<svg>` must be:

```
width="70px"
height="70px"
```

---

### 🖌 Styling Rules (for visible artwork)

All visible drawing elements (`path`, `rect`, `circle`, etc.) must have:


| Property         | Required Value                            |
| ---------------- | ----------------------------------------- |
| `fill`           | `#ffffff`                                 |
| `stroke`         | `#000000`                                 |
| `stroke-width`   | `1px` / `0.26458333` (± small tolerance) |
| `stroke-opacity` | `1`                                       |
| `fill-opacity`   | `1`                                       |

These may be set either as attributes or inside the `style` attribute.

- Must be a clean outline suitable for export to other formats.
- Avoid unnecessary nodes.
- Keep paths closed where possible.

---

### 🖼 Reference Artwork Layer

If you use reference artwork (e.g. a 3-view drawing):

- Place it in a layer named **“Reference Artwork”**
- The layer **must be hidden** (`display:none`)
- Any `<image>` elements must be inside hidden layers

Visible `<image>` elements will cause validation to fail.

---

### 🚫 Ignored SVG Content

The validator automatically ignores:

- Anything inside hidden layers/groups
- Anything inside `<defs>` sections

---

## 🔍 Automated Validation

When you open a pull request:


| Check         | What It Validates                              |
| ------------- | ---------------------------------------------- |
| JSON Schema   | Airframe definition files match the schema     |
| SVG Validator | SVG size, styles, hidden layers, and structure |

If a check fails, click into the failed job to see exactly what needs fixing.
