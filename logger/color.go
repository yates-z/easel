package logger

import (
	"fmt"
	"strconv"
)

// Color represents a text color.
type Color uint8

// Foreground colors.
const (
	DefaultColor Color = 0
	Black        Color = iota + 29
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func (c Color) Paint(s string) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, 0, 0, uint8(c), s, 0x1B)
}

func (c Color) PaintWith(color Color, s string) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, 0, color, uint8(c), s, 0x1B)
}

func (c Color) String() string {
	return strconv.Itoa(int(c))
}

func (c Color) IsDefault() bool {
	return c == DefaultColor
}

func (c Color) ToBackground() Color {
	if c == DefaultColor {
		return c
	}
	return c + 10
}
