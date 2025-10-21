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
	xPos            int
	yPos            int
	background      *ebiten.Image
	backgroundXView int
	font            font.Face
	enemy           []*enemyShip
	score           int
	speed           int
	laser           []*laserBlast
	shotCount       int
	ammoCount       int
}

type enemyShip struct {
	pic       *ebiten.Image
	xEnemyPos int
	yEnemyPos int
}

type laserBlast struct {
	pic        *ebiten.Image
	xLaserPos  int
	yLaserPos  int
	laserSpeed int
}

func NewEnemy(xStart, yStart int, image *ebiten.Image) *enemyShip {
	return &enemyShip{
		pic:       image,
		xEnemyPos: randRange(xStart, 2500),
		yEnemyPos: rand.Intn(yStart),
	}
}

func NewLaser(xStart, yStart int, image *ebiten.Image) *laserBlast {
	return &laserBlast{
		pic:       image,
		xLaserPos: xStart,
		yLaserPos: yStart,
	}
}

func (spaceGame *starsAndSpaceGame) Update() error {
	backgroundWidth := spaceGame.background.Bounds().Dx()
	maxX := backgroundWidth * 2
	//var shotCount int
	spaceGame.backgroundXView -= 4
	spaceGame.backgroundXView %= maxX
	enemyEnt, _, err := ebitenutil.NewImageFromFile("./Entities/UFO.png")
	if err != nil {
		fmt.Println("Unable to load Enemy entity:", err)
	}
	if spaceGame.yPos <= -500 {
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			spaceGame.yPos += 0
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			spaceGame.yPos += 5
		}
	} else if spaceGame.yPos >= 475 {
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			spaceGame.yPos += -5
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			spaceGame.yPos += 0
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			spaceGame.yPos += -5
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			spaceGame.yPos += 5
		}
	}

	for i := 0; i < len(spaceGame.enemy); i++ {
		if spaceGame.enemy[i].xEnemyPos > 0-spaceGame.enemy[i].pic.Bounds().Dx() {
			spaceGame.enemy[i].xEnemyPos += -spaceGame.speed
		} else {
			spaceGame.enemy[i] = NewEnemy(1000, 950, enemyEnt)
			spaceGame.score -= 1
		}
	}

	// Gets total number of lasers for ammoCount
	if spaceGame.ammoCount == 9000 {
		spaceGame.ammoCount = len(spaceGame.laser)
	}

	// Fires laser
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// Allows you to shoot only if you have ammo
		if spaceGame.ammoCount > 0 {
			spaceGame.ammoCount--
			spaceGame.laser[spaceGame.shotCount].xLaserPos = spaceGame.xPos + spaceGame.player.Bounds().Dx()
			spaceGame.laser[spaceGame.shotCount].yLaserPos = spaceGame.yPos + 512
			spaceGame.shotCount++
		}
	}
	// Prevents shotCount from going out of bounds
	if spaceGame.shotCount >= 10 {
		spaceGame.shotCount = 0
	}
	// Moves fired lasers
	for i := 0; i < len(spaceGame.laser); i++ {
		if spaceGame.laser[i].xLaserPos < 900 && spaceGame.laser[i].xLaserPos > 0 {
			spaceGame.laser[i].xLaserPos += spaceGame.speed
		} else if spaceGame.laser[i].xLaserPos >= 900 && spaceGame.laser[i].xLaserPos < 1000 {
			spaceGame.ammoCount++
			spaceGame.laser[i].xLaserPos = -200
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
	text.Draw(screen, "EnemyPos: "+strconv.Itoa(spaceGame.enemy[0].xEnemyPos), drawFace, textOpts)

	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 180)
	text.Draw(screen, "Enemy Count: "+strconv.Itoa(len(spaceGame.enemy)), drawFace, textOpts)

	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 210)
	text.Draw(screen, "LaserPos: "+strconv.Itoa(spaceGame.laser[0].xLaserPos), drawFace, textOpts)

	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 240)
	text.Draw(screen, "yPos: "+strconv.Itoa(spaceGame.yPos), drawFace, textOpts)

	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 270)
	text.Draw(screen, "Shot Count: "+strconv.Itoa(spaceGame.shotCount), drawFace, textOpts)

	textOpts.GeoM.Reset()
	textOpts.GeoM.Translate(x, 300)
	text.Draw(screen, "Ammo Count: "+strconv.Itoa(spaceGame.ammoCount), drawFace, textOpts)

	// Draws enemy units on screen
	enemyDrawOpts := &ebiten.DrawImageOptions{}
	for _, swarm := range spaceGame.enemy {
		enemyDrawOpts.GeoM.Reset()
		enemyDrawOpts.GeoM.Translate(float64(swarm.xEnemyPos), float64(swarm.yEnemyPos))
		screen.DrawImage(swarm.pic, enemyDrawOpts)
	}

	// Draws player unit on screen
	playerDrawOpts := ebiten.DrawImageOptions{}
	playerDrawOpts.GeoM.Reset()
	playerDrawOpts.GeoM.Translate(x, 500)
	playerDrawOpts.GeoM.Translate(float64(spaceGame.xPos), float64(spaceGame.yPos))
	screen.DrawImage(spaceGame.player, &playerDrawOpts)

	// Draws laser blast on screen
	laserDrawOpts := &ebiten.DrawImageOptions{}
	for _, shots := range spaceGame.laser {
		laserDrawOpts.GeoM.Reset()
		laserDrawOpts.GeoM.Translate(float64(shots.xLaserPos), float64(shots.yLaserPos))
		screen.DrawImage(shots.pic, laserDrawOpts)
	}

}

// Layout Sets window parameters
func (spaceGame starsAndSpaceGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetFullscreen(false)
	ebiten.SetWindowTitle("Scroller Example")
	backgroundPict, _, err := ebitenutil.NewImageFromFile("StarsAndSpace.png")
	if err != nil {
		fmt.Println("Unable to load background image:", err)
	}
	playerUnit := makePlayer()
	enemyUnits := makeEnemy()
	lasers := makeLaser()
	spaceScrollerGame := starsAndSpaceGame{
		player:     playerUnit,
		enemy:      enemyUnits,
		laser:      lasers,
		background: backgroundPict,
		font:       LoadFont("Ubuntu-Regular.ttf", 18),
		speed:      2,
		ammoCount:  9000,
	}
	err = ebiten.RunGame(&spaceScrollerGame)
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

// Function to select a random integer between a minimum and maximum value
func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

// Function that creates player entity
func makePlayer() *ebiten.Image {
	playerUnit, _, err := ebitenutil.NewImageFromFile("./Entities/FriendlyUFO.png")
	if err != nil {
		fmt.Println("Unable to load Player unit:", err)
	}
	return playerUnit
}

// Function that creates enemy units
func makeEnemy() []*enemyShip {
	enemyUnits := make([]*enemyShip, 0)
	enemyEnt, _, err := ebitenutil.NewImageFromFile("./Entities/UFO.png")
	if err != nil {
		fmt.Println("Unable to load Enemy entity:", err)
	}
	for i := 0; i < 3; i++ {
		enemyUnits = append(enemyUnits, NewEnemy(1000, 950, enemyEnt))
	}
	return enemyUnits
}

// Function that creates laser blasts for the game
func makeLaser() []*laserBlast {
	laserShots := make([]*laserBlast, 0)
	laserEnt, _, err := ebitenutil.NewImageFromFile("./Entities/laserBlast.png")
	if err != nil {
		fmt.Println("Unable to load laser entity:", err)
	}
	for i := 0; i < 10; i++ {
		laserShots = append(laserShots, NewLaser(-200, 0, laserEnt))
		laserShots[i].laserSpeed = 2
	}
	return laserShots
}
