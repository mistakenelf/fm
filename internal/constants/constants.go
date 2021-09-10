package constants

// DirectoryTypes contains the different types of directories.
type DirectoryTypes struct {
	CurrentDirectory  string
	PreviousDirectory string
	HomeDirectory     string
}

// Directories contains the different kinds of directories and their values.
var Directories = DirectoryTypes{
	CurrentDirectory:  ".",
	PreviousDirectory: "..",
	HomeDirectory:     "~",
}

// VersionTypes contains the different types of versions.
type VersionTypes struct {
	AppVersion string
}

// Versions contains the different kinds of versions and their values.
var Versions = VersionTypes{
	AppVersion: "0.2.0",
}

// DimensionTypes contains the different types of dimensions.
type DimensionTypes struct {
	StatusBarHeight int
}

// Dimensions contains the different kinds of dimensions and their values.
var Dimensions = DimensionTypes{
	StatusBarHeight: 1,
}

// ColorTypes contains the different types of colors.
type ColorTypes struct {
	White       string
	Pink        string
	LightPurple string
	DarkPurple  string
	DarkGray    string
	Blue        string
	Red         string
}

// Colors contains the different kinds of colors and their values.
var Colors = ColorTypes{
	White:       "#FFFDF5",
	Pink:        "#F25D94",
	LightPurple: "#A550DF",
	DarkPurple:  "#6124DF",
	DarkGray:    "#353533",
	Blue:        "#1D4ED8",
	Red:         "#DC2626",
}
