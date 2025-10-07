// Marcus Smith
package main

import (
	"fmt"
	_ "image/png"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type gameState int

const (
	gameStateStart gameState = iota
	gameStatePlay
)

type starsAndSpaceGame struct {
	player          *ebiten.Image
	xPos            int
	yPos            int
	background      *ebiten.Image
	backgroundXView int
	state           gameState
	font            font.Face
	enemy           []*enemyUnit
}

type enemyUnit struct {
	pic  *ebiten.Image
	xPos float64
	yPos float64
}

func (spaceGame *starsAndSpaceGame) Update() error {
	if spaceGame.state == gameStateStart {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			spaceGame.state = gameStatePlay
		}
		return nil
	} else {
		backgroundWidth := spaceGame.background.Bounds().Dx()
		maxX := backgroundWidth * 2
		spaceGame.backgroundXView -= 4
		spaceGame.backgroundXView %= maxX
		inpututil.IsKeyJustPressed(ebiten.KeyLeft)
		return nil
	}
}

func (spaceGame *starsAndSpaceGame) Draw(screen *ebiten.Image) {
	// Draws start screen with "How to Play" text
	if spaceGame.state == gameStateStart {
		const x = 350
		screen.Fill(colornames.Khaki)
		drawFace := text.NewGoXFace(spaceGame.font)
		textOpts := &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{},
			LayoutOptions:    text.LayoutOptions{},
		}
		textOpts.GeoM.Reset()
		textOpts.GeoM.Translate(x, 390)
		textOpts.ColorScale.ScaleWithColor(colornames.Red)
		text.Draw(screen, "How to Play: ", drawFace, textOpts)

		textOpts = &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{},
			LayoutOptions:    text.LayoutOptions{},
		}
		textOpts.GeoM.Reset()
		textOpts.GeoM.Translate(x+20, 430)
		textOpts.ColorScale.ScaleWithColor(colornames.Red)
		text.Draw(screen, "Controls:", drawFace, textOpts)

		textOpts = &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{},
			LayoutOptions:    text.LayoutOptions{},
		}
		textOpts.GeoM.Reset()
		textOpts.GeoM.Translate(x+40, 470)
		textOpts.ColorScale.ScaleWithColor(colornames.Red)
		text.Draw(screen, "Arrow Keys to Move (or A and D).", drawFace, textOpts)

		textOpts = &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{},
			LayoutOptions:    text.LayoutOptions{},
		}
		textOpts.GeoM.Reset()
		textOpts.GeoM.Translate(x+40, 510)
		textOpts.ColorScale.ScaleWithColor(colornames.Red)
		text.Draw(screen, "Space to shoot.", drawFace, textOpts)

		textOpts = &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{},
			LayoutOptions:    text.LayoutOptions{},
		}
		textOpts.GeoM.Reset()
		textOpts.GeoM.Translate(x+20, 550)
		textOpts.ColorScale.ScaleWithColor(colornames.Red)
		text.Draw(screen, "Objective:", drawFace, textOpts)

		textOpts = &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{},
			LayoutOptions:    text.LayoutOptions{},
		}
		textOpts.GeoM.Reset()
		textOpts.GeoM.Translate(x+40, 590)
		textOpts.ColorScale.ScaleWithColor(colornames.Red)
		text.Draw(screen, "Shoot enemies to get points, any enemies that get past you take away points", drawFace, textOpts)
	} else {
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
	}
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

	demo := starsAndSpaceGame{
		player:     nil,
		background: backgroundPict,
		font:       LoadFont("Ubuntu-Regular.ttf", 24),
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
