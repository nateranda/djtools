package db

import (
	"errors"
	"fmt"
)

func RgbToHex(r, g, b int) (string, error) {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return "", errors.New("RGB values must be between 0 and 255")
	}
	return fmt.Sprintf("#%02X%02X%02X", r, g, b), nil
}

func HexToRgb(hex string) (int, int, int, error) {
	var r, g, b int
	_, err := fmt.Sscanf(hex, "#%02X%02X%02X", &r, &g, &b)
	if err != nil {
		return 0, 0, 0, err
	}
	return r, g, b, nil
}
