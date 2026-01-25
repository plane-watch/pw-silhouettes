package main

import (
	"fmt"
	"os"
	"os/exec"
)

func convertSVGtoPNG(inkscapeBinary, src, dst string) error {
	cmd := exec.Command(
		inkscapeBinary,
		src,
		"--export-type=png",
		"--export-overwrite",
		"--export-filename="+dst,
	)

	out, err := cmd.CombinedOutput()

	// Inkscape sometimes logs useful info even on "success".
	if err != nil {
		return fmt.Errorf("inkscape failed: %w\noutput:\n%s", err, out)
	}

	// Don’t assume success: ensure the file exists and is non-zero.
	st, statErr := os.Stat(dst)
	if statErr != nil {
		return fmt.Errorf("inkscape produced no output file: %v\noutput:\n%s", statErr, out)
	}
	if st.Size() == 0 {
		return fmt.Errorf("inkscape produced empty output file (%s)\noutput:\n%s", dst, out)
	}

	return nil
}
