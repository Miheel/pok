package pok

type Delta interface {
	Undo(ed *Editor)
	Redo(ed *Editor)
}

const preAllocDelta = 16
var UndoStack = make([]Delta, 0, preAllocDelta)
var RedoStack = make([]Delta, 0, preAllocDelta)

var CurrentPencilDelta *PencilDelta = &PencilDelta{}
var CurrentEraserDelta *EraserDelta = &EraserDelta{}
var CurrentBucketDelta *BucketDelta = &BucketDelta{}
var CurrentObjectDelta *ObjectDelta = &ObjectDelta{}
var CurrentRemoveObjectDelta *RemoveObjectDelta = &RemoveObjectDelta{}
var CurrentLinkDelta *LinkDelta = &LinkDelta{}
var CurrentRemoveLinkDelta *RemoveLinkDelta = &RemoveLinkDelta{}
var CurrentAutotileDelta *AutotileDelta = &AutotileDelta{}
var CurrentNpcDelta *NpcDelta = &NpcDelta{}
var CurrentRemoveNpcDelta *RemoveNpcDelta = &RemoveNpcDelta{}

var CurrentResizeDelta *ResizeDelta = &ResizeDelta{}

func PerformUndo(ed *Editor) {
	if len(UndoStack) > 0 {
		top := UndoStack[len(UndoStack) - 1]
		top.Undo(ed)
		RedoStack = append(RedoStack, top)
		UndoStack = UndoStack[:len(UndoStack) - 1]
	}
}

func PerformRedo(ed *Editor) {
	if len(RedoStack) > 0 {
		top := RedoStack[len(RedoStack)-1]
		top.Redo(ed)
		UndoStack = append(UndoStack, top)
		RedoStack = RedoStack[:len(RedoStack) - 1]
	}
}

type PencilDelta struct {
	indicies []int
	oldTiles []int
	oldTextureIndicies []int
	z int
	tileMapIndex int
	newTile int
	newTextureIndex int
}

func (dp *PencilDelta) Undo(ed *Editor) {
	tm := ed.tileMaps[dp.tileMapIndex]
	for i, j := range dp.indicies {
		tm.Tiles[dp.z][j] = dp.oldTiles[i]
		tm.TextureIndicies[dp.z][j] = dp.oldTextureIndicies[i]
	}
}

func (dp *PencilDelta) Redo(ed *Editor) {
	tm := ed.tileMaps[dp.tileMapIndex]
	for _, j := range dp.indicies {
		tm.Tiles[dp.z][j] = dp.newTile
		tm.TextureIndicies[dp.z][j] = dp.newTextureIndex
	}
}

type EraserDelta struct {
	indicies []int
	oldTiles []int
	oldTextureIndicies []int
	z int
	tileMapIndex int
	newTextureIndex int
}

func (de *EraserDelta) Undo(ed *Editor) {
	tm := ed.tileMaps[de.tileMapIndex]
	for i, j := range de.indicies {
		tm.Tiles[de.z][j] = de.oldTiles[i]
		tm.TextureIndicies[de.z][j] = de.oldTextureIndicies[i]
	}
}

func (de *EraserDelta) Redo(ed *Editor) {
	tm := ed.tileMaps[de.tileMapIndex]
	for _, j := range de.indicies {
		tm.Tiles[de.z][j] = -1
		tm.TextureIndicies[de.z][j] = 0
	}
}

type BucketDelta struct {
	indicies []int
	oldTiles []int
	oldTextureIndicies[]int
	z int
	tileMapIndex int
	newTile int
	newTextureIndex int
}

func (bd *BucketDelta) Undo(ed *Editor) {
	tm := ed.tileMaps[bd.tileMapIndex]
	for i, j := range bd.indicies {
		tm.Tiles[bd.z][j] = bd.oldTiles[i]
		tm.TextureIndicies[bd.z][j] = bd.oldTextureIndicies[i]
	}
}

func (bd *BucketDelta) Redo(ed *Editor) {
	tm := ed.tileMaps[bd.tileMapIndex]
	for _, j := range bd.indicies {
		tm.Tiles[bd.z][j] = bd.newTile
		tm.TextureIndicies[bd.z][j] = bd.newTextureIndex
	}
}

type ObjectDelta struct {
	placedObjectIndex int
	objectIndex int
	tileMapIndex int
	origin int
	z int
}

func (do *ObjectDelta) Undo(ed *Editor) {
	tm := ed.tileMaps[do.tileMapIndex]
	pobj := placedObjects[do.tileMapIndex][do.placedObjectIndex]
	obj := &ed.objectGrid.objs[do.objectIndex]
	tm.EraseObject(pobj, obj)
	copy(placedObjects[do.tileMapIndex][do.placedObjectIndex:], placedObjects[do.tileMapIndex][do.placedObjectIndex+1:])
	placedObjects[do.tileMapIndex] = placedObjects[do.tileMapIndex][:len(placedObjects[do.tileMapIndex])-1]
}

func (do *ObjectDelta) Redo(ed *Editor) {
	tm := ed.tileMaps[do.tileMapIndex]
	obj := &ed.objectGrid.objs[do.objectIndex]
	tm.InsertObject(obj, do.placedObjectIndex, do.origin, do.z, &placedObjects[do.tileMapIndex])
}

type RemoveObjectDelta struct {
	objectDelta *ObjectDelta
}

func (dor *RemoveObjectDelta) Undo(ed *Editor) {
	dor.objectDelta.Redo(ed)
}

func (dor *RemoveObjectDelta) Redo(ed *Editor) {
	dor.objectDelta.Undo(ed)
}

type LinkDelta struct {
	linkBegin *LinkData
	linkEnd *LinkData
	linkIdA int
	linkIdB int
}

func (dl *LinkDelta) Undo(ed *Editor) {
	tmA := ed.tileMaps[dl.linkBegin.TileMapIndex]
	tmB := ed.tileMaps[dl.linkEnd.TileMapIndex]

	for i := range tmA.Exits {
		if tmA.Exits[i].Id == dl.linkIdA {
			ed.removeLink(dl.linkBegin.TileMapIndex, i)
			break
		}
	}

	for i := range tmB.Exits {
		if tmB.Exits[i].Id == dl.linkIdB {
			ed.removeLink(dl.linkEnd.TileMapIndex, i)
			break
		}
	}

}

func (dl *LinkDelta) Redo(ed *Editor) {
	ed.tryConnectTileMaps(dl.linkBegin, dl.linkEnd)
}

type RemoveLinkDelta struct {
	entry *Entry
	exit *Exit
	tileMapIndex int
}

func (drl *RemoveLinkDelta) Undo(ed *Editor) {
	tm := ed.tileMaps[drl.tileMapIndex]
	if drl.entry != nil {
		tm.Entries = append(tm.Entries, *drl.entry)
	}
	if drl.exit != nil {
		tm.Exits = append(tm.Exits, *drl.exit)
	}
}

func (drl *RemoveLinkDelta) Redo(ed *Editor) {
	tm := ed.tileMaps[drl.tileMapIndex]
	if drl.entry != nil {
		entryIndex := -1

		for i := range tm.Entries {
			if tm.Exits[i].X == drl.entry.X && tm.Exits[i].Y == drl.entry.Y {
				entryIndex = i
				break
			}
		}

		if entryIndex != -1 {
			tm.Entries[entryIndex] = tm.Entries[len(tm.Entries)-1]
			tm.Entries = tm.Entries[:len(tm.Entries)-1]
		}
	}
	if drl.exit != nil {
		exitIndex := -1

		for i := range tm.Exits {
			if tm.Exits[i].X == drl.entry.X && tm.Exits[i].Y == drl.entry.Y {
				exitIndex = i
				break
			}
		}

		if exitIndex != -1 {
			tm.Exits[exitIndex] = tm.Exits[len(tm.Exits)-1]
			tm.Exits = tm.Exits[:len(tm.Exits)-1]
		}
	}
}

type AutotileDelta struct {
	oldValues map[int]ModifiedTile
	newValues map[int]ModifiedTile
	tileMapIndex int
	z int
}

type ModifiedTile struct {
	tile int
	textureIndex int
}

func (a *AutotileDelta) Join(other *AutotileDelta) {
	if len(a.oldValues) == 0 {
		a.oldValues = other.oldValues
	}

	if len(a.newValues) == 0 {
		a.newValues = other.newValues
	}

	for i, j := range other.oldValues {
		if _, ok := a.oldValues[i]; !ok {
			a.oldValues[i] = j
		}
	}

	for i, j := range other.newValues {
		a.newValues[i] = j
	}
}

func (dat *AutotileDelta) Undo(ed *Editor) {
	tm := ed.tileMaps[dat.tileMapIndex]

	for i, j := range dat.oldValues {
		tm.Tiles[dat.z][i] = j.tile
		tm.TextureIndicies[dat.z][i] = j.textureIndex
	}
}

func (dat *AutotileDelta) Redo(ed *Editor) {
	tm := ed.tileMaps[dat.tileMapIndex]

	for i, j := range dat.newValues {
		tm.Tiles[dat.z][i] = j.tile
		tm.TextureIndicies[dat.z][i] = j.textureIndex
	}
}

type ResizeDelta struct {
	oldExits map[int]Exit
	oldEntries map[int]Entry
	exitIndicies []int
	entryIndicies []int

	offsetDeltaX, offsetDeltaY float64
	dx, dy, origin int
	tileMapIndex int
}

func (dr *ResizeDelta) InsertExitsAndEntries(tm *TileMap) {

	findIndexLessThan := func(value int, span []int) *int {
		if len(span) == 0 {
			return nil
		}

		var min *int = nil

		for i := range span {
			if span[i] > value {
				min = &span[i]
			}
		}

		if min == nil {
			return nil
		}

		for i := range span {
			if span[i] > value && span[i] < *min {
				min = &span[i]
			}
		}

		return min
	}

	for i := findIndexLessThan(-1, dr.exitIndicies[:]); i != nil; i = findIndexLessThan(*i, dr.exitIndicies[:]){
		tm.Exits = append(tm.Exits[:*i+1], tm.Exits[*i:]...)
		tm.Exits[*i] = dr.oldExits[*i]
	}

	for i := findIndexLessThan(-1, dr.entryIndicies[:]); i != nil; i = findIndexLessThan(*i, dr.entryIndicies[:]){
		tm.Entries = append(tm.Entries[:*i+1], tm.Entries[*i:]...)
		tm.Entries[*i] = dr.oldEntries[*i]
	}
}

func (dr *ResizeDelta) Undo(ed *Editor) {
	tm := ed.tileMaps[dr.tileMapIndex]

	tm.Resize(-dr.dx, -dr.dy, dr.origin)
	dr.InsertExitsAndEntries(tm)

	ed.tileMapOffsets[dr.tileMapIndex].X -= dr.offsetDeltaX
	ed.tileMapOffsets[dr.tileMapIndex].Y -= dr.offsetDeltaY
}

func (dr *ResizeDelta) Redo(ed *Editor) {
	tm := ed.tileMaps[dr.tileMapIndex]

	tm.Resize(dr.dx, dr.dy, dr.origin)
	_, _, _, _ = ed.removeInvalidLinks(dr.tileMapIndex)

	ed.tileMapOffsets[dr.tileMapIndex].X += dr.offsetDeltaX
	ed.tileMapOffsets[dr.tileMapIndex].Y += dr.offsetDeltaY
}

type NpcDelta struct {
	npcInfo *NpcInfo
	tileMapIndex int
	npcIndex int
}

func (dn *NpcDelta) Undo(ed *Editor) {
	tm := ed.tileMaps[dn.tileMapIndex]
	tm.RemoveNpc(dn.npcIndex)
}

func (dn *NpcDelta) Redo(ed *Editor) {
	tm := ed.tileMaps[dn.tileMapIndex]
	dn.npcIndex = len(tm.npcs)
	tm.PlaceNpc(dn.npcInfo)
}

type RemoveNpcDelta struct {
	npcDelta *NpcDelta
}

func (drn *RemoveNpcDelta) Undo(ed *Editor) {
	drn.npcDelta.Redo(ed)
}

func (drn *RemoveNpcDelta) Redo(ed *Editor) {
	drn.npcDelta.Undo(ed)
}
