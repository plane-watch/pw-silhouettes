package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

const (
	svgNS      = "http://www.w3.org/2000/svg"
	inkscapeNS = "http://www.inkscape.org/namespaces/inkscape"

	wantSizePx = 70.0

	wantFill         = "#ffffff"
	wantStroke       = "#000000"
	wantStrokeWidth  = 0.26458333
	wantOpacity      = 1.0
	defaultStrokeTol = 0.0005 // "close enough" tolerance for stroke-width
	defaultSizeTolPx = 0.01
)

type Issue struct {
	File string
	Line int
	Msg  string
}

func runApp(_ context.Context, cmd *cli.Command) error {

	var issues []Issue

	is, err := ValidateSVG(cmd.String("svg"), defaultStrokeTol, defaultSizeTolPx)
	if err != nil {
		return fmt.Errorf("invalid svg file: %w", err)
	}
	issues = append(issues, is...)

	if len(issues) > 0 {
		for _, it := range issues {
			line := it.Line
			if line <= 0 {
				line = 1
			}
			log.Error().Int("line", line).Str("file", it.File).Msg(it.Msg)
		}
		return fmt.Errorf("%d issues", len(issues))
	}

	return nil
}

func ValidateSVG(path string, strokeWidthTol, sizeTol float64) ([]Issue, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := xml.NewDecoder(f)

	var issues []Issue

	// Visibility stack (effective hidden state), starting at "not hidden"
	hiddenStack := []bool{false}

	// Style stack: inherited properties (only meaningful for visible nodes)
	styleStack := []map[string]string{map[string]string{}}

	// If we enter a hidden subtree, we skip *all* checks and style processing until we exit it.
	skipDepth := 0

	seenRootSVG := false

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("xml parse error: %w", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			line := decoderLine(dec)

			// Determine parent hidden state
			parentHidden := hiddenStack[len(hiddenStack)-1]

			attrs := t.Attr
			thisHidden := elementHidden(attrs)
			effectiveHidden := parentHidden || thisHidden

			// Always push hidden state so EndElement pops stay aligned.
			hiddenStack = append(hiddenStack, effectiveHidden)

			// If we're already skipping, we do *nothing* (including style processing).
			// Just track depth so we know when we leave the skipped subtree.
			if skipDepth > 0 {
				skipDepth++
				// Keep styleStack aligned too (push a dummy entry). It will be popped on EndElement.
				styleStack = append(styleStack, styleStack[len(styleStack)-1])
				continue
			}

			// NEW: ignore everything under <defs>
			ignoreSubtree := t.Name.Local == "defs"

			// If this element is hidden, begin skipping its subtree entirely.
			// We also keep styleStack aligned with a dummy push.
			if effectiveHidden || ignoreSubtree {
				skipDepth = 1
				styleStack = append(styleStack, styleStack[len(styleStack)-1])
				continue
			}

			// ---- From here down: we are visible and not under any hidden ancestor ----

			// Compute style only for visible elements
			parentStyle := styleStack[len(styleStack)-1]
			thisStyle := mergeStyles(parentStyle, styleFrom(attrs))
			styleStack = append(styleStack, thisStyle)

			// Root <svg> width/height check (first svg element we see)
			if !seenRootSVG && t.Name.Local == "svg" {
				seenRootSVG = true
				w, okW := getAttr(attrs, "", "width")
				h, okH := getAttr(attrs, "", "height")
				if !okW || !okH {
					issues = append(issues, Issue{File: path, Line: line, Msg: "root <svg> missing width/height attributes"})
				} else {
					wpx, errW := parsePxLength(w)
					hpx, errH := parsePxLength(h)
					if errW != nil || errH != nil {
						msg := fmt.Sprintf("root <svg> width/height must be 70px/70px (got width=%q height=%q)", w, h)
						issues = append(issues, Issue{File: path, Line: line, Msg: msg})
					} else {
						if !closeEnough(wpx, wantSizePx, sizeTol) || !closeEnough(hpx, wantSizePx, sizeTol) {
							msg := fmt.Sprintf("root <svg> width/height must be 70px/70px (got width=%.6gpx height=%.6gpx)", wpx, hpx)
							issues = append(issues, Issue{File: path, Line: line, Msg: msg})
						}
					}
				}
			}

			// Element checks (only for visible elements)
			switch t.Name.Local {
			case "image":
				issues = append(issues, Issue{File: path, Line: line, Msg: "visible <image> found (reference artwork must be hidden)"})

			case "path", "rect", "circle", "ellipse", "polygon", "polyline", "line":
				issues = append(issues, validateDrawable(path, line, t.Name.Local, thisStyle, strokeWidthTol)...)
			}

		case xml.EndElement:
			// Pop hidden stack
			if len(hiddenStack) > 1 {
				hiddenStack = hiddenStack[:len(hiddenStack)-1]
			}

			// Pop style stack (kept aligned even when skipping)
			if len(styleStack) > 1 {
				styleStack = styleStack[:len(styleStack)-1]
			}

			// Manage skip depth
			if skipDepth > 0 {
				skipDepth--
			}
		}
	}

	// If the SVG never had a root <svg>, it’s malformed (but the parser would likely have errored)
	if !seenRootSVG {
		issues = append(issues, Issue{File: path, Line: 1, Msg: "no <svg> root element found"})
	}

	return issues, nil
}

func validateDrawable(file string, line int, name string, style map[string]string, strokeWidthTol float64) []Issue {
	var issues []Issue

	// Normalise colours to lowercase
	fill := strings.ToLower(strings.TrimSpace(style["fill"]))
	stroke := strings.ToLower(strings.TrimSpace(style["stroke"]))

	if fill != wantFill {
		issues = append(issues, Issue{File: file, Line: line, Msg: fmt.Sprintf("<%s> fill must be %s (got %q)", name, wantFill, fill)})
	}
	if stroke != wantStroke {
		issues = append(issues, Issue{File: file, Line: line, Msg: fmt.Sprintf("<%s> stroke must be %s (got %q)", name, wantStroke, stroke)})
	}

	// stroke-width: numeric, allow close enough
	swStr := strings.TrimSpace(style["stroke-width"])
	if swStr == "" {
		issues = append(issues, Issue{File: file, Line: line, Msg: fmt.Sprintf("<%s> missing stroke-width", name)})
	} else {
		sw, err := parseNumber(swStr)
		if err != nil {
			issues = append(issues, Issue{File: file, Line: line, Msg: fmt.Sprintf("<%s> invalid stroke-width %q", name, swStr)})
		} else if !closeEnough(sw, wantStrokeWidth, strokeWidthTol) {
			issues = append(issues, Issue{File: file, Line: line, Msg: fmt.Sprintf("<%s> stroke-width must be %.8f (got %.8f)", name, wantStrokeWidth, sw)})
		}
	}

	// Opacities
	if v := strings.TrimSpace(style["stroke-opacity"]); v == "" {
		issues = append(issues, Issue{File: file, Line: line, Msg: fmt.Sprintf("<%s> missing stroke-opacity", name)})
	} else if op, err := parseNumber(v); err != nil || !closeEnough(op, wantOpacity, 0.0001) {
		issues = append(issues, Issue{File: file, Line: line, Msg: fmt.Sprintf("<%s> stroke-opacity must be 1 (got %q)", name, v)})
	}

	if v := strings.TrimSpace(style["fill-opacity"]); v == "" {
		issues = append(issues, Issue{File: file, Line: line, Msg: fmt.Sprintf("<%s> missing fill-opacity", name)})
	} else if op, err := parseNumber(v); err != nil || !closeEnough(op, wantOpacity, 0.0001) {
		issues = append(issues, Issue{File: file, Line: line, Msg: fmt.Sprintf("<%s> fill-opacity must be 1 (got %q)", name, v)})
	}

	return issues
}

// elementHidden returns true if this element is hidden by its own attributes/style.
// Note: ancestral hidden state is handled by the stack in ValidateSVG.
func elementHidden(attrs []xml.Attr) bool {
	// direct attrs
	if v, ok := getAttr(attrs, "", "display"); ok && strings.TrimSpace(v) == "none" {
		return true
	}
	if v, ok := getAttr(attrs, "", "visibility"); ok && strings.TrimSpace(v) == "hidden" {
		return true
	}
	// style attr
	if style, ok := getAttr(attrs, "", "style"); ok {
		s := parseStyle(style)
		if strings.TrimSpace(s["display"]) == "none" || strings.TrimSpace(s["visibility"]) == "hidden" {
			return true
		}
	}
	return false
}

// styleFrom builds a property map from this element's attributes + style="...".
func styleFrom(attrs []xml.Attr) map[string]string {
	out := map[string]string{}

	// style="k:v; k2:v2"
	if s, ok := getAttr(attrs, "", "style"); ok {
		for k, v := range parseStyle(s) {
			out[k] = v
		}
	}

	// presentation attributes override style
	for _, key := range []string{
		"fill", "stroke", "stroke-width", "stroke-opacity", "fill-opacity",
	} {
		if v, ok := getAttr(attrs, "", key); ok {
			out[key] = v
		}
	}

	return out
}

func mergeStyles(parent, child map[string]string) map[string]string {
	// copy parent
	out := make(map[string]string, len(parent)+len(child))
	for k, v := range parent {
		out[k] = v
	}
	// apply child overrides
	for k, v := range child {
		out[k] = v
	}
	return out
}

func parseStyle(s string) map[string]string {
	out := map[string]string{}
	parts := strings.Split(s, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, ":", 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		out[k] = v
	}
	return out
}

func getAttr(attrs []xml.Attr, space, local string) (string, bool) {
	for _, a := range attrs {
		if a.Name.Local == local && (space == "" || a.Name.Space == space) {
			return a.Value, true
		}
	}
	return "", false
}

func parsePxLength(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(strings.ToLower(s), "px") {
		s = strings.TrimSpace(s[:len(s)-2])
	}
	// Accept plain numbers as px (common in SVG)
	return parseNumber(s)
}

func parseNumber(s string) (float64, error) {
	s = strings.TrimSpace(s)
	// Allow values like "0.26458333" or "0.26458333px"
	if strings.HasSuffix(strings.ToLower(s), "px") {
		s = strings.TrimSpace(s[:len(s)-2])
	}
	// Some SVGs might use scientific notation; ParseFloat handles it.
	return strconv.ParseFloat(s, 64)
}

func closeEnough(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func decoderLine(dec *xml.Decoder) int {
	// InputPos returns (line, column)
	line, _ := dec.InputPos()
	return line
}
