package main

type (

	// Output is the schema used to generate the output JSON
	Output struct {

		// Version represents the schema version
		Version int `json:"version"`

		Metadata Metadata `json:"metadata"`

		// AirframeToSprite maps an airframe ICAO (key) to a Sprite (value)
		AirframeToSprite map[string]string `json:"airframeToSprite"`

		// Sprites represents the artwork in the spritesheet. It is named after an airframe ICAO (the key).
		Sprites map[string]Sprite `json:"sprites"`
	}

	// Sprite represents sprite details in the output JSON
	Sprite struct {
		IDs       []int   `json:"ids"`
		Scale     float64 `json:"scale"`
		Anchor    Anchor  `json:"anchor"`
		NoRotate  bool    `json:"noRotate,omitempty"`
		FrameTime *int    `json:"frameTime,omitempty"`
	}

	Metadata struct {
		PNG          string `json:"png"`
		SpriteWidth  int    `json:"spriteWidth"`
		SpriteHeight int    `json:"spriteHeight"`
	}

	// Airframe represents the input JSON airframe schema defined at the root of this repo
	Airframe struct {
		Version int     `json:"version"`
		ICAO    ICAO    `json:"icao"`
		AliasOf *string `json:"aliasOf,omitempty"`
		Render  Render  `json:"render"`
		Art     Art     `json:"art"`
		Notes   string  `json:"notes"`
	}

	// ICAO represents the ICAO information from the input JSON airframe schema
	ICAO struct {
		Designator   string `json:"designator"`
		TypeCode     string `json:"typeCode"`
		WakeCategory string `json:"wakeCategory"`
	}

	// Render represents the sprite rendering information from the input JSON airframe schema
	Render struct {
		Scale    float64 `json:"scale"`
		Anchor   Anchor  `json:"anchor"`
		NoRotate bool    `json:"noRotate"`
	}

	// Anchor defines an x,y point from the top-left of the sprite, that shall be the
	// anchor point of the sprite. Any rotation should be done about this point.
	// The sprite should be drawn with this point on the pixel the sprite is indended to be drawn at.
	Anchor struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	// Art represents the sprite artwork from the input JSON airframe schema
	Art struct {
		Frames    []Frame `json:"frames"`
		FrameTime int     `json:"frameTime"`
	}

	// Frame represents the sprite artwork from the input JSON airframe schema
	Frame struct {
		Src string `json:"src"`
	}
)
