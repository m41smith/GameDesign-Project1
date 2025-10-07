// Marcus Smith
package main

import (
	"fmt"
	_ "image/png"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type starsAndSpaceGame struct {
	player          *ebiten.Image
	xPos            float64
	yPos            float64
	background      *ebiten.Image
	backgroundXView int
	font            font.Face
	enemy           []*enemyUnit
	score           float64
}

type enemyUnit struct {
	pic  *ebiten.Image
	xPos float64
	yPos float64
}

func (spaceGame *starsAndSpaceGame) Update() error {
	const x = 350
	backgroundWidth := spaceGame.background.Bounds().Dx()
	maxX := backgroundWidth * 2
	spaceGame.backgroundXView -= 4
	spaceGame.backgroundXView %= maxX
	maxY := spaceGame.background.Bounds().Dy() - 625
	if spaceGame.yPos <= -500 {
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			spaceGame.yPos += 0
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
			spaceGame.yPos += 5
		}
	} else if spaceGame.yPos > float64(maxY) {
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			spaceGame.yPos += -5
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
			spaceGame.yPos += 0
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			spaceGame.yPos += -5
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
			spaceGame.yPos += 5
		}
	}
	spaceGame.score = spaceGame.yPos
	inpututil.IsKeyJustPressed(ebiten.KeyLeft)
	return nil
}

func (spaceGame *starsAndSpaceGame) Draw(screen *ebiten.Image) {

	drawOps := ebiten.DrawImageOptions{}
	const repeat = 3
	backgroundWidth := spaceGame.background.Bounds().Dx()
	for count := 0; count < repeat; count += 1 {
		drawOps.GeoM.Reset()
		drawOps.GeoM.Translate(float64(backgroundWidth*count),
			float64(0))
		drawOps.GeoM.Translate(float64(spaceGame.backgroundXView), 0)
		screen.DrawImage(spaceGame.background, &drawOps)
	}

	const x = 350
	drawFace := text.NewGoXFace(spaceGame.font)
	textOpts := &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions:    text.LayoutOptions{},
	}
	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 90)
	textOpts.ColorScale.ScaleWithColor(colornames.Red)
	text.Draw(screen, "How to Play: Controls: Arrow Keys to Move (or W and S).", drawFace, textOpts)

	textOpts = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions:    text.LayoutOptions{},
	}
	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x+20, 120)
	textOpts.ColorScale.ScaleWithColor(colornames.Red)
	text.Draw(screen, "Space to shoot.", drawFace, textOpts)

	textOpts = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions:    text.LayoutOptions{},
	}
	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 150)
	textOpts.ColorScale.ScaleWithColor(colornames.Red)
	text.Draw(screen, "Objective: Shoot enemies to get points, any enemies that get past you take away points", drawFace, textOpts)

	//textOpts = &text.DrawOptions{
	//	DrawImageOptions: ebiten.DrawImageOptions{},
	//	LayoutOptions:    text.LayoutOptions{},
	//}
	//textOpts.GeoM.Reset()
	//textOpts.GeoM.Translate(x, 250)
	//textOpts.ColorScale.ScaleWithColor(colornames.Red)
	//text.Draw(screen, strconv.Itoa(spaceGame.background.Bounds().Dx()), drawFace, textOpts)

	textOpts = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions:    text.LayoutOptions{},
	}
	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 250)
	textOpts.ColorScale.ScaleWithColor(colornames.Red)
	text.Draw(screen, strconv.FormatFloat(spaceGame.score, 'f', -1, 64), drawFace, textOpts)

	//playerUnit := spaceGame.player.Bounds().Dx()
	playerOpts := ebiten.DrawImageOptions{}
	playerOpts.GeoM.Reset()
	playerOpts.GeoM.Translate(350, 500)
	playerOpts.GeoM.Translate(spaceGame.xPos, spaceGame.yPos)
	screen.DrawImage(spaceGame.player, &playerOpts)
}

// Layout Sets window parameters
func (spaceGame starsAndSpaceGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	//	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Scroller Example")
	//New image from file returns image as image.Image (_) and ebiten.Image
	backgroundPict, _, err := ebitenutil.NewImageFromFile("StarsAndSpace.png")
	if err != nil {
		fmt.Println("Unable to load background image:", err)
	}

	playerUnit, _, errP := ebitenutil.NewImageFromFile("UFO.png")
	if errP != nil {
		fmt.Println("Unable to load Player unit:", err)
	}

	demo := starsAndSpaceGame{
		player:     playerUnit,
		background: backgroundPict,
		font:       LoadFont("Ubuntu-Regular.ttf", 18),
	}
	err = ebiten.RunGame(&demo)
	if err != nil {
		fmt.Println("Failed to run game", err)
	}
}

// LoadFont implements the font used in the game
func LoadFont(fontFile string, size float64) font.Face {
	fileHandle, err := os.Open(fontFile)
	if err != nil {
		log.Fatal(err)
	}
	fontData, err := io.ReadAll(fileHandle)
	if err != nil {
		log.Fatal(err)
	}
	ttFont, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}
	fontFace, err := opentype.NewFace(ttFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	return fontFace
}
