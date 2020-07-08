package pok

import ( 
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"io/ioutil"
	"image"
	"image/color"
	"log"
)

type Game struct {
	Ows OverworldState
	As GameState
	Player Player
	Client Client
	Rend Renderer
	Audio Audio
}

func InitGame() {
	var err error
	playerImg, _, err = ebitenutil.NewImageFromFile("./resources/images/lucas.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	tileset, _, err = ebitenutil.NewImageFromFile("./resources/images/tileset1.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	selection, err = ebiten.NewImage(TileSize, TileSize, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	collisionMarker, err = ebiten.NewImage(TileSize, TileSize, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	exitMarker, err = ebiten.NewImage(TileSize, TileSize, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	selectionClr := color.RGBA{255, 0, 0, 255}

	for p := 0; p < selection.Bounds().Max.X; p++ {
		selection.Set(p, 0, selectionClr)
		selection.Set(p, selection.Bounds().Max.Y - 1, selectionClr)
	}

	for p := 1; p < selection.Bounds().Max.Y - 1; p++ {
		selection.Set(0, p, selectionClr)
		selection.Set(selection.Bounds().Max.Y - 1, p, selectionClr)
	}

	collisionClr := color.RGBA{255, 0, 255, 255}

	for p := 0; p < 4; p++ {
		for q := 0; q < 4; q++ {
			collisionMarker.Set(p, q, collisionClr)
		}
	}

	exitClr := color.RGBA{0, 0, 255, 255}

	for p:= 0; p < 4; p++ {
		for q := 0; q < 4; q++ {
			exitMarker.Set(p + 14, q, exitClr)
		}
	}
}

func (g *Game) TileIsOccupied(x int, y int, z int) bool {
	if x < 0 || x >= g.Ows.tileMap.Width || y < 0 ||  y >= g.Ows.tileMap.Height {
		return true
	}

	index := y * g.Ows.tileMap.Width + x

	// Out of bounds check
	if z < 0 || z >= len(g.Ows.tileMap.Tiles) {
		return true
	}

	if index >= len(g.Ows.tileMap.Tiles[z]) || index < 0 {
		return true
	}

	if g.Ows.tileMap.Collision[z][index] {
		return true
	}

	for _, p := range g.Client.playerMap.players {
		if p.X == x && p.Y == y {
			return true
		}
	}

	return false
}

func (g *Game) Update(screen *ebiten.Image) error {
	err := g.As.GetInputs(g)
	if err != nil {
		return err
	}
	err = g.As.Update(g)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.As.Draw(g, screen)
}

func (g *Game) Load(str string, entrypoint int) {
	data, err := ioutil.ReadFile(str)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &g.Ows.tileMap)
	if err != nil {
		panic(err)
	}
	currentLayer = 0
	selectedTile = 0
	g.Player.Location = str
	index := g.Ows.tileMap.GetEntryWithId(entrypoint)
	g.Player.X = g.Ows.tileMap.Entries[index].X
	g.Player.Y = g.Ows.tileMap.Entries[index].Y
	g.Player.Gx = float64(g.Player.X * TileSize)
	g.Player.Gy = float64(g.Player.Y * TileSize)
	g.Rend = NewRenderer(DisplaySizeX,
		DisplaySizeY,
		DisplaySizeX,
		DisplaySizeY,
	)
}

func (g *Game) Save() {
	bytes, err := json.Marshal(g.Ows.tileMap)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(g.Player.Location, bytes, 0644)
}

func (g *Game) DrawPlayer(player *Player) {
	playerOpt := &ebiten.DrawImageOptions{}
	playerOpt.GeoM.Scale(2,2)

	x := player.Gx + playerOffsetX
	y := player.Gy + playerOffsetY

	playerRect := image.Rect(
		player.Tx,
		player.Ty,
		player.Tx + TileSize,
		player.Ty + TileSize,
	)

	g.Rend.Draw(&RenderTarget{
		playerOpt,
		playerImg,
		&playerRect,
		x,
		y,
		3,
	})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return DisplaySizeX, DisplaySizeY
}

func (g *Game) PlayAudio() {
	g.Audio.audioPlayer.Play()
}
