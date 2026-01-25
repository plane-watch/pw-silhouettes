package airframe

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	Airframe struct {
		Version int     `json:"version"`
		ICAO    ICAO    `json:"icao"`
		AliasOf *string `json:"aliasOf,omitempty"`
		Render  Render  `json:"render"`
		Art     Art     `json:"art"`
		Notes   string  `json:"notes"`
	}

	ICAO struct {
		Designator   string `json:"designator"`
		TypeCode     string `json:"typeCode"`
		WakeCategory string `json:"wakeCategory"`
	}

	Render struct {
		Scale    float64 `json:"scale"`
		Anchor   Anchor  `json:"anchor"`
		NoRotate bool    `json:"noRotate"`
	}

	Anchor struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	Art struct {
		Frames    []Frame `json:"frames"`
		FrameTime int     `json:"frameTime"`
	}

	Frame struct {
		Src      string `json:"src"`
		SpriteID int    `json:"spriteID,omitempty"`
	}
)

func FromFile(filename string) (*Airframe, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	af := new(Airframe)
	err = json.Unmarshal(b, af)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal airframe: %w", err)
	}
	return af, nil
}
