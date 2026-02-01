# ✈️ build_spritesheet

`build_spritesheet` is a small utility that builds a PNG sprite sheet from a directory of airframe and generic airframe type definitions.

It reads aircraft and generic type definition JSON files, renders their associated SVG silhouettes using **Inkscape**, and packs them into a single spritesheet image with an accompanying JSON file that maps airframes and generic airframe types to sprite IDs in the sprite sheet.

The tool is designed to extend an existing spritesheet as new airframes are added.

---

## 📍 Important

**This tool must be run from the root of the repository.**  
Paths inside airframe definition files are resolved relative to the repo root.

---

## 🧠 What It Does

For each airframe definition:

1. Reads the airframe & generic type JSON metadata  
2. Locates the referenced SVG silhouette  
3. Uses **Inkscape v1+** to rasterise the SVG to PNG  
4. Places the rendered sprite into the correct position in the spritesheet grid  
5. Generates an accompanying JSON file that:
  - Defines the sprites. Sprites are defined by a map of one or more IDs within the sprite sheet, plus rotation and scaling information. 
  - Maps the airframe and generic types to sprites. 

The result is a single PNG containing all aircraft sprites, and accompanying metadata, suitable for use in UIs.

---

## ⚙️ Requirements

- **Go 1.25+**
- **Inkscape 1.0 or newer**  
  The tool calls Inkscape directly to convert SVG → PNG.

Check your version:

```bash
inkscape --version
```

---

## 🚀 Usage

```bash
go run ./cmd/build_spritesheet --inkscape_binary /usr/bin/inkscape --output_png ./spritesheet.png
```

### Flags

| Flag | Alias | Required | Description |
|------|-------|----------|-------------|
| `--inkscape_binary` | `--inkscape` | ✅ | Path to the Inkscape **v1+** binary |
| `--output_png` | `--op` | ✅ | Path where the generated sprite sheet PNG will be written |
| `—-output_json` | `—op` | ✅ | Path where the generated metadata JSON will be written |

---

## 📁 Airframe Definitions

Each airframe and also generic airframe types (for aircraft that we don’t have type designators for) is defined by a JSON file, see the README.md at the root of this repo for details.

---

## 📦 Output

The outputs:

✔ A packed PNG spritesheet containing all airframes, and [original sprites](./cmd/build_spritesheet/original_sprites.png) at their original locations.  

✔ A JSON file containing the sprite sheet metadata

---

## 🧩 Typical Workflow

This tool is run automatically on push to main when building release assets. 

You should not have to run manually. 
