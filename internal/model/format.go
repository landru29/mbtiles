package model

// Format is the picture format.
type Format string

const (
	// FormatJPEG is the jpg format.
	FormatJPEG Format = "jpg"
	// FormatPNG is the png format.
	FormatPNG Format = "png"
	// FormatNoTransform keeps the original format.
	FormatNoTransform Format = "no-transform"
)

// String implements the pflag.Value interface.
func (f Format) String() string {
	return string(f)
}

// Set implements the pflag.Value interface.
func (f *Format) Set(data string) error {
	*f = Format(data)

	return nil
}

// Type implements the pflag.Value interface.
func (f *Format) Type() string {
	return "Image format"
}
