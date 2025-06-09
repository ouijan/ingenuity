package utils_test

import (
	"testing"

	"github.com/ouijan/ingenuity/pkg/core/utils"
)

func TestToIsoCoordinates(t *testing.T) {
	tests := []struct {
		x, y, w, h int
		wantX      int
		wantY      int
	}{
		{x: 0, y: 0, w: 128, h: 64, wantX: 0, wantY: 0},
		// {x: 1, y: 0, w: 128, h: 64, wantX: -32, wantY: 16},
		// {x: 0, y: 1, w: 128, h: 64, wantX: -64, wantY: 96},
		{x: 2, y: 1, w: 128, h: 64, wantX: 64, wantY: 96},
		// {x: 2, y: 0, w: 64, h: 32, wantX: 64, wantY: 32},
		// {x: 0, y: 2, w: 64, h: 32, wantX: -64, wantY: 32},
	}

	for _, tt := range tests {
		gotX, gotY := utils.ToIsoCoordinates(tt.x, tt.y, tt.w, tt.h)
		if gotX != tt.wantX || gotY != tt.wantY {
			t.Errorf("ToIsoCoordinates(%d, %d, %d, %d) = (%d, %d), want (%d, %d)",
				tt.x, tt.y, tt.w, tt.h, gotX, gotY, tt.wantX, tt.wantY)
		}
	}
}
