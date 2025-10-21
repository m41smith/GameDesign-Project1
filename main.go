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
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/solarlune/resolv"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type gameState int

const (
	playState gameState = iota
	endState
)

type starsAndSpaceGame struct {
	player          playerShip
	background      *ebiten.Image
	backgroundXView int
	font            font.Face
	enemy           []*enemyShip
	score           int
	speed           int
	laser           []*laserBlast
	shotCount       int
	ammoCount       int
	space           *resolv.Space
	state           gameState
	audioContext    *audio.Context
	audioPlayer     *audio.Player
}

type playerShip struct {
	pic                 *ebiten.Image
	xPlayerPos          int
	yPlayerPos          int
	playerCollisionRect *resolv.ConvexPolygon
}

type enemyShip struct {
	pic                *ebiten.Image
	xEnemyPos          int
	yEnemyPos          int
	enemyCollisionRect *resolv.ConvexPolygon
}

type laserBlast struct {
	pic                *ebiten.Image
	xLaserPos          int
	yLaserPos          int
	laserSpeed         int
	laserCollisionRect *resolv.ConvexPolygon
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

// Trying to create collisions
func (spaceGame *starsAndSpaceGame) Init() {
	spaceGame.space = resolv.NewSpace(1000, 1000, 16, 16)
	spaceGame.player.playerCollisionRect = resolv.NewRectangle(float64(spaceGame.player.xPlayerPos)-27.5, float64(spaceGame.player.yPlayerPos)-12.5, 55, 25)
	spaceGame.space.Add(spaceGame.player.playerCollisionRect)

	for i := 0; i < len(spaceGame.enemy); i++ {
		spaceGame.enemy[i].enemyCollisionRect = resolv.NewRectangle(float64(spaceGame.enemy[i].xEnemyPos)-27.5, float64(spaceGame.enemy[i].yEnemyPos)-12.5, 55, 25)
		spaceGame.space.Add(spaceGame.enemy[i].enemyCollisionRect)
	}
	for i := 0; i < len(spaceGame.laser); i++ {
		spaceGame.laser[i].laserCollisionRect = resolv.NewRectangle(float64(spaceGame.laser[i].xLaserPos-21), float64(spaceGame.laser[i].yLaserPos-2), 42, 4)
		spaceGame.space.Add(spaceGame.laser[i].laserCollisionRect)
	}
}

func (spaceGame *starsAndSpaceGame) Update() error {
	if spaceGame.state == playState {
		backgroundWidth := spaceGame.background.Bounds().Dx()
		maxX := backgroundWidth * 2
		spaceGame.backgroundXView -= 4
		spaceGame.backgroundXView %= maxX
		enemyEnt, _, err := ebitenutil.NewImageFromFile("./Entities/UFO.png")
		if err != nil {
			fmt.Println("Unable to load Enemy entity:", err)
		}

		// Converts key inputs into movement for player unit
		if spaceGame.player.yPlayerPos <= -500 {
			if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
				spaceGame.player.yPlayerPos += 0
			} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
				spaceGame.player.yPlayerPos += 5
			}
		} else if spaceGame.player.yPlayerPos >= 475 {
			if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
				spaceGame.player.yPlayerPos += -5
			} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
				spaceGame.player.yPlayerPos += 0
			}
		} else {
			if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
				spaceGame.player.yPlayerPos += -5
			} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
				spaceGame.player.yPlayerPos += 5
			}
		}

		// Moves enemy units and subtracts a point if they get past you
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
				spaceGame.audioPlayer.Rewind()
				spaceGame.audioPlayer.Play()
				spaceGame.laser[spaceGame.shotCount].xLaserPos = spaceGame.player.xPlayerPos + spaceGame.player.pic.Bounds().Dx()
				spaceGame.laser[spaceGame.shotCount].yLaserPos = spaceGame.player.yPlayerPos + 512
				spaceGame.shotCount++
			}
		}

		// Attempt at collision
		//for _, enemies := range spaceGame.enemy {
		//	spaceGame.player.playerCollisionRect.IntersectionTest(resolv.IntersectionTestSettings{
		//		TestAgainst: enemies.enemyCollisionRect.SelectTouchingCells(1).FilterShapes(),
		//		OnIntersect: func(set resolv.IntersectionSet) bool {
		//			spaceGame.score -= 1000
		//			return true
		//		},
		//	})
		//}

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

		// Forces game over for testing until collisions work
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			spaceGame.state = endState
		}
		return nil
	} else {
		// Implement Game Over
		return nil
	}

}

func (spaceGame *starsAndSpaceGame) Draw(screen *ebiten.Image) {
	if spaceGame.state == playState {
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
		text.Draw(screen, "Space to shoot. When out of ammo, wait for reload", drawFace, textOpts)

		// Text explaining game objective
		textOpts.GeoM.Reset()
		textOpts.GeoM.Translate(x, 90)
		text.Draw(screen, "Objective: Shoot enemies to get points, any enemies that get past you take away points", drawFace, textOpts)

		// Text representing score value
		textOpts.GeoM.Reset()
		textOpts.GeoM.Translate(x, 120)
		text.Draw(screen, "Score: "+strconv.Itoa(spaceGame.score), drawFace, textOpts)

		// Text displaying remaining ammunition
		textOpts.GeoM.Reset()
		textOpts.GeoM.Translate(x, 150)
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
		playerDrawOpts.GeoM.Translate(float64(spaceGame.player.xPlayerPos), float64(spaceGame.player.yPlayerPos))
		screen.DrawImage(spaceGame.player.pic, &playerDrawOpts)

		// Draws laser blast on screen
		laserDrawOpts := &ebiten.DrawImageOptions{}
		for _, shots := range spaceGame.laser {
			laserDrawOpts.GeoM.Reset()
			laserDrawOpts.GeoM.Translate(float64(shots.xLaserPos), float64(shots.yLaserPos))
			screen.DrawImage(shots.pic, laserDrawOpts)
		}
	} else {
		drawFace := text.NewGoXFace(LoadFont("Ubuntu-Regular.ttf", 30))
		screen.Fill(colornames.Black)
		textOpts := &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{},
			LayoutOptions:    text.LayoutOptions{},
		}
		textOpts.GeoM.Reset()
		textOpts.GeoM.Translate(400, 470)
		textOpts.ColorScale.ScaleWithColor(colornames.Red)
		text.Draw(screen, "Game Over", drawFace, textOpts)
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
	soundContext := audio.NewContext(48000)

	spaceScrollerGame := starsAndSpaceGame{
		player:       playerUnit,
		enemy:        enemyUnits,
		laser:        lasers,
		background:   backgroundPict,
		font:         LoadFont("Ubuntu-Regular.ttf", 18),
		speed:        2,
		ammoCount:    9000,
		audioContext: soundContext,
		audioPlayer:  LoadWav("LaserBlastSound.wav", soundContext),
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

// randRange returns a random integer between a minimum and maximum value
func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

// makePlayer creates the player entity
func makePlayer() playerShip {
	playerEnt, _, err := ebitenutil.NewImageFromFile("./Entities/FriendlyUFO.png")
	if err != nil {
		fmt.Println("Unable to load Player unit:", err)
	}
	playerUnit := playerShip{
		pic: playerEnt,
	}
	return playerUnit
}

// makeEnemy creates the enemy entities
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

// Collision code borrowed from slides for testing
func playerCollisionCheck(spaceGame *starsAndSpaceGame) {
	for _, enemies := range spaceGame.enemy {
		if hit := spaceGame.player.playerCollisionRect.Intersection(enemies.enemyCollisionRect); !hit.IsEmpty() {
			//spaceGame.state = endState
			spaceGame.score += 1000
		}
	}
}

func LoadWav(name string, context *audio.Context) *audio.Player {
	laserFile, err := os.Open(name)
	if err != nil {
		fmt.Println("Error Loading sound: ", err)
	}
	laserSound, err := wav.DecodeWithoutResampling(laserFile)
	if err != nil {
		fmt.Println("Error interpreting sound file: ", err)
	}
	soundPlayer, err := context.NewPlayer(laserSound)
	if err != nil {
		fmt.Println("Couldn't create sound player: ", err)
	}
	return soundPlayer
}
