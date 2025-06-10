package math

// IsoGridToPixel converts screen coordinates (x, y) to isometric coordinates based on the tile width (w) and height (h).
// Reference: https://clintbellanger.net/articles/isometric_math
func IsoGridToPixel(x, y, w, h int) (int, int) {
	halfWidth := w / 2
	halfHeight := h / 2
	isoX := (x - y) * halfWidth
	isoY := (x + y) * halfHeight
	return isoX, isoY
}

// IsoPixelToGrid converts isometric coordinates (x, y) back to screen coordinates based on the tile width (w) and height (h).
// Reference: https://clintbellanger.net/articles/isometric_math
func IsoPixelToGrid(x, y, w, h int) (int, int) {
	halfWidth := w / 2
	halfHeight := h / 2
	screenX := (x/halfWidth + y/halfHeight) / 2
	screenY := (y/halfHeight - (x / halfWidth)) / 2
	return screenX, screenY
}
