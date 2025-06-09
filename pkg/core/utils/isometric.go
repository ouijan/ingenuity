package utils

func ToIsoCoordinates(x, y, w, h int) (int, int) {
	// Isometric Transformation vars
	// ix := 1.0
	// iy := 0.5
	// jx := -1.0
	// jy := 0.5

	fx := float64(x)
	fy := float64(y)
	fw := float64(w)
	fh := float64(h)

	// isoX := fx*ix*.5*fw + fy*jx*.5*fw
	// isoY := fx*iy*.5*fh + fy*jy*.5*fh

	isoX := (fx - fy) * (fw / 2)
	isoY := (fx + fy) * (fh / 2)

	return int(isoX), int(isoY)
}
