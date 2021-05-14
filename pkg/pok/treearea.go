package pok

import (
	//"github.com/hajimehoshi/ebiten"
	"image/color"
)

type TreeAreaSelection struct {
	BeginX, BeginY *int
	EndX, EndY int
	TreeInfo *TreeAutoTileInfo
}

func (t *TreeAreaSelection) ClampToTileMap(tm *TileMap) {

}

func (t *TreeAreaSelection) Hold(x, y int, cam *Camera) {
	x, y = fitToCamera(x, y, cam)
	if t.BeginX == nil && t.BeginY == nil {
		t.BeginX = &x
		t.BeginY = &y
	} else {
		t.EndX = x
		t.EndY = y
	}
}

func (t *TreeAreaSelection) Draw(rend *Renderer) {
	if !t.IsHolding() {
		return
	}

	x0 := float64(*t.BeginX)
	y0 := float64(*t.BeginY)

	x1 := float64(t.EndX)
	y1 := float64(t.EndY)

	clr := color.RGBA{255, 0, 0, 255}

	line := DebugLine{}
	line.Clr = clr
	line.X1 = x0
	line.Y1 = y0
	line.X2 = x1
	line.Y2 = y0

	rend.DrawLine(line)

	line.X1 = x0
	line.Y1 = y1
	line.X2 = x1
	line.Y2 = y1

	rend.DrawLine(line)

	line.X1 = x0
	line.Y1 = y0
	line.X2 = x0
	line.Y2 = y1

	rend.DrawLine(line)

	line.X1 = x1
	line.Y1 = y0
	line.X2 = x1
	line.Y2 = y1

	rend.DrawLine(line)
}

func fitToCamera(x, y int, cam *Camera) (int, int) {
	x = int((float64(x) + cam.X) / cam.Scale)
	y = int((float64(y) + cam.Y) / cam.Scale)
	return x, y
}

func (t *TreeAreaSelection) IsHolding() bool {
	return t.BeginX != nil && t.BeginY != nil
}

func (t *TreeAreaSelection) Release() {
	if t.BeginX == nil || t.BeginY == nil {
		return
	}

	abs := func(i int) int {
		if i < 0 {
			return -i
		}
		return i
	}

	w := abs(*t.BeginX - t.EndX)
	h := abs(*t.BeginY - t.EndY)

	t.PopulateWithTrees(w, h)
	t.BeginX = nil
	t.BeginY = nil
}

func (t *TreeAreaSelection) PopulateWithTrees(w, h int) {

}
