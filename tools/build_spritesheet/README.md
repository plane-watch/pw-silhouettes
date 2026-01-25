# ✈️ build_spritesheet

`build_spritesheet` is a small utility that builds a PNG sprite sheet from a directory of airframe definitions.

It reads aircraft definition JSON files, renders their associated SVG silhouettes using **Inkscape**, and packs them into a single spritesheet image. The tool is designed to extend an existing spritesheet as new airframes are added.

Future versions will also generate companion JavaScript metadata for sprite lookup.

---

## 📍 Important

**This tool must be run from the root of the repository.**  
Paths inside airframe definition files are resolved relative to the repo root.

---

## 🧠 What It Does

For each airframe definition:

1. Reads the airframe JSON metadata  
2. Locates the referenced SVG silhouette  
3. Uses **Inkscape v1+** to rasterise the SVG to PNG  
4. Places the rendered sprite into the correct position in the spritesheet grid  

The result is a single PNG containing all aircraft sprites, suitable for use in web or game UIs.

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
| `--output_png` | `-o` | ✅ | Path where the generated spritesheet PNG will be written |

---

## 📁 Airframe Definitions

Each airframe is defined by a JSON file, see the README.md at the root of this repo for details.

---

## 📦 Output

The tool currently outputs:

✔ A packed PNG spritesheet containing all airframes, and [original sprites](./cmd/build_spritesheet/original_sprites.png) at their original locations.  

Planned:

⬜ Companion JavaScript metadata describing sprite positions  
⬜ Optional sprite atlas JSON output  

---

## 🧩 Typical Workflow

1. Add a new airframe JSON + SVG  
2. Run `build_spritesheet` from the repository root  
3. Commit updated spritesheet to pw-ui repo
4. (Future) Commit updated JS atlas to pw-ui repo
