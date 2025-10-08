// Marcus Smith
package main

import (
	"fmt"
	_ "image/png"
	"io"
	"log"
	"math/rand"
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
	enemy           []*enemyShip
	score           int
	speed           int
}

type enemyShip struct {
	pic       *ebiten.Image
	xEnemyPos float64
	yEnemyPos float64
}

func NewEnemy(xStart, yStart int, image *ebiten.Image) *enemyShip {
	return &enemyShip{
		pic:       image,
		xEnemyPos: float64(xStart),
		yEnemyPos: float64(rand.Intn(yStart)),
	}
}

func (spaceGame *starsAndSpaceGame) Update() error {
	backgroundWidth := spaceGame.background.Bounds().Dx()
	maxX := backgroundWidth * 2
	spaceGame.backgroundXView -= 4
	spaceGame.backgroundXView %= maxX
	if spaceGame.yPos <= -500 {
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			spaceGame.yPos += 0
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
			spaceGame.yPos += 5
		}
	} else if spaceGame.yPos >= 475 {
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
	//spaceGame.score = spaceGame.yPos

	//spaceGame.score = spaceGame.enemy[0].xEnemyPos

	for i := 0; i < len(spaceGame.enemy); i++ {
		if spaceGame.enemy[i].xEnemyPos > float64(0-spaceGame.enemy[i].pic.Bounds().Dx()) {
			spaceGame.enemy[i].xEnemyPos += -float64(spaceGame.speed)
		} else {
			spaceGame.enemy[i].xEnemyPos = spaceGame.enemy[i].xEnemyPos
		}
	}

	inpututil.IsKeyJustPressed(ebiten.KeyLeft)
	return nil
}

func (spaceGame *starsAndSpaceGame) Draw(screen *ebiten.Image) {

	// Draws scrolling background
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

	const x = 20
	drawFace := text.NewGoXFace(spaceGame.font)
	textOpts := &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions:    text.LayoutOptions{},
	}

	// Text explaining controls
	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 30)
	textOpts.ColorScale.ScaleWithColor(colornames.Red)
	text.Draw(screen, "How to Play: Controls: Arrow Keys to Move (or W and S).", drawFace, textOpts)
	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x+20, 60)
	text.Draw(screen, "Space to shoot.", drawFace, textOpts)

	// Text explaining game objective
	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 90)
	text.Draw(screen, "Objective: Shoot enemies to get points, any enemies that get past you take away points", drawFace, textOpts)

	// Text representing score value
	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 120)
	text.Draw(screen, "Score: "+strconv.Itoa(spaceGame.score), drawFace, textOpts)

	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 150)
	text.Draw(screen, "EnemyPos: "+strconv.FormatFloat(spaceGame.enemy[0].xEnemyPos, 'f', -1, 64), drawFace, textOpts)

	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 180)
	text.Draw(screen, "Enemy Count: "+strconv.Itoa(len(spaceGame.enemy)), drawFace, textOpts)

	// Draws enemy units on screen
	enemyDrawOpts := &ebiten.DrawImageOptions{}
	for _, swarm := range spaceGame.enemy {
		enemyDrawOpts.GeoM.Reset()
		enemyDrawOpts.GeoM.Translate(swarm.xEnemyPos, swarm.yEnemyPos)
		screen.DrawImage(swarm.pic, enemyDrawOpts)
	}

	// Draws player unit on screen
	playerDrawOpts := ebiten.DrawImageOptions{}
	playerDrawOpts.GeoM.Reset()
	playerDrawOpts.GeoM.Translate(x, 500)
	playerDrawOpts.GeoM.Translate(spaceGame.xPos, spaceGame.yPos)
	screen.DrawImage(spaceGame.player, &playerDrawOpts)
}

// Layout Sets window parameters
func (spaceGame starsAndSpaceGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetFullscreen(false)
	ebiten.SetWindowTitle("Scroller Example")
	//New image from file returns image as image.Image (_) and ebiten.Image
	backgroundPict, _, err := ebitenutil.NewImageFromFile("StarsAndSpace.png")
	if err != nil {
		fmt.Println("Unable to load background image:", err)
	}
	playerUnit, _, err := ebitenutil.NewImageFromFile("./Entities/UFO.png")
	if err != nil {
		fmt.Println("Unable to load Player unit:", err)
	}
	enemyUnits := make([]*enemyShip, 0)
	enemyEnt, _, err := ebitenutil.NewImageFromFile("./Entities/UFO.png")
	if err != nil {
		fmt.Println("Unable to load Enemy entity:", err)
	}
	for i := 0; i < 3; i++ {
		enemyUnits = append(enemyUnits, NewEnemy(1000, 950, enemyEnt))
	}
	enemyDelay()
	//enemyUnits = append(enemyUnits, NewEnemy(900, 950, enemyEnt))
	demo := starsAndSpaceGame{
		player:     playerUnit,
		enemy:      enemyUnits,
		background: backgroundPict,
		font:       LoadFont("Ubuntu-Regular.ttf", 18),
		speed:      2,
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

// Working on function to space out enemy spawns (WIP)
func enemyDelay() {
	enemyUnits := make([]*enemyShip, 0)
	enemyEnt, _, err := ebitenutil.NewImageFromFile("./Entities/UFO.png")
	if err != nil {
		fmt.Println("Unable to load Enemy entity:", err)
	}
	for i := 0; i < 3; i++ {
		enemyUnits = append(enemyUnits, NewEnemy(1000, 950, enemyEnt))
	}
}
