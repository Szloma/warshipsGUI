package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		g := NewGUI()
		if err := loop(w, g); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type GUI struct {
	theme                      *material.Theme
	startButton                *widget.Clickable
	exitButton                 *widget.Clickable
	personalizationButton      *widget.Clickable
	acceptButton               *widget.Clickable
	discardButton              *widget.Clickable
	nickname                   *widget.Editor
	profileDescription         *widget.Editor
	acceptShipPositions        *widget.Clickable
	backShipPositions          *widget.Clickable
	discardShipPositions       *widget.Clickable
	randomShipPositions        *widget.Clickable
	abandonButton              *widget.Clickable
	gameOverButton             *widget.Clickable
	showStatsButton            *widget.Clickable
	backFromStatsButton        *widget.Clickable
	botActiveButton            *widget.Clickable
	yourTurnIndicatorButton    *widget.Clickable
	refreshLobbyButton         *widget.Clickable
	lobbyButtons               []*widget.Clickable
	lobbyButtonsStates         [10]string
	youWinScreen               bool
	enemyDescriptionAvailable  bool
	sessionActive              bool
	lockRightTable             bool
	yourTurnIncidator          bool
	botActive                  bool
	lockLeftTable              bool
	youLoseScreen              bool
	displayPlayerAndEnemyBoard bool
	showShipSetUpMenu          bool
	showLeftTable              bool
	showTables                 bool
	showPersonalization        bool
	showStartMenu              bool
	inGame                     bool
	showLoadingMenu            bool
	showGameOver               bool
	showGameWon                bool
	showLeaderBoards           bool
	showStats                  bool
	selectionIincidatorState   [20]int
	leftShip                   int
	shipsPlaced                int
	selectionIndicatorButtons  []*widget.Clickable
	leftTableButtons           [][]*widget.Clickable
	leftTableLabels            [][]string
	leftTableStates            [][]int
	rightTableButtons          [][]*widget.Clickable
	rightTableLabels           [][]string
	rightTableStates           [][]int
	accuracy                   float64
	timeLeft                   string
	enemyName                  string
	enemyDescription           string
	timer                      *time.Ticker
}

func NewGUI() *GUI {
	gui := &GUI{
		theme:                      material.NewTheme(),
		startButton:                new(widget.Clickable),
		exitButton:                 new(widget.Clickable),
		personalizationButton:      new(widget.Clickable),
		acceptButton:               new(widget.Clickable),
		discardButton:              new(widget.Clickable),
		nickname:                   new(widget.Editor),
		profileDescription:         new(widget.Editor),
		acceptShipPositions:        new(widget.Clickable),
		discardShipPositions:       new(widget.Clickable),
		randomShipPositions:        new(widget.Clickable),
		backShipPositions:          new(widget.Clickable),
		abandonButton:              new(widget.Clickable),
		gameOverButton:             new(widget.Clickable),
		showStatsButton:            new(widget.Clickable),
		backFromStatsButton:        new(widget.Clickable),
		botActiveButton:            new(widget.Clickable),
		yourTurnIndicatorButton:    new(widget.Clickable),
		refreshLobbyButton:         new(widget.Clickable),
		enemyDescriptionAvailable:  false,
		sessionActive:              true,
		lockRightTable:             true,
		yourTurnIncidator:          false,
		botActive:                  true,
		leftShip:                   20,
		shipsPlaced:                0,
		youWinScreen:               false,
		lockLeftTable:              false,
		showStats:                  false,
		showLeaderBoards:           false,
		showGameWon:                false,
		showGameOver:               false,
		showLoadingMenu:            false,
		youLoseScreen:              false,
		inGame:                     false,
		displayPlayerAndEnemyBoard: false,
		showShipSetUpMenu:          false,
		showLeftTable:              false,
		showTables:                 false,
		showPersonalization:        false,
		showStartMenu:              true,
		timeLeft:                   "60",
		accuracy:                   0.0,
		enemyName:                  "Janusz",
		enemyDescription:           "aaa",
		timer:                      time.NewTicker(1 * time.Second),
	}
	gui.lobbyButtons, gui.lobbyButtonsStates = createLobbyButtons()
	gui.leftTableButtons, gui.leftTableLabels, gui.leftTableStates = createTable()
	gui.rightTableButtons, gui.rightTableLabels, gui.rightTableStates = createTable()
	gui.selectionIndicatorButtons = createButtonRow()
	gui.selectionIincidatorState = setSelectionIndidatorState(gui.leftShip)

	return gui
}

const (
	Empty = iota
	Ship
	Hit
	Miss
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func handlePlayerBoard(g *GUI) {
	for i := range g.leftTableButtons {
		for y := range g.leftTableButtons[i] {
			//fmt.Printf("%s: %d\n", g.leftTableLabels[i][y], g.leftTableStates[i][y])
			for a := range gameProperties.opp_shots {

				if g.leftTableLabels[i][y] == gameProperties.opp_shots[a] {
					if g.leftTableStates[i][y] == Empty {
						g.leftTableStates[i][y] = Miss
					}
					if g.leftTableStates[i][y] == Ship {
						g.leftTableStates[i][y] = Hit
					}
				}

			}

		}
	}
}

func placeShipsOnLeftTable(gtx layout.Context, g *GUI) {
	for i := range g.leftTableButtons {
		for y := range g.leftTableButtons[i] {
			//fmt.Printf("%s: %d\n", g.leftTableLabels[i][y], g.leftTableStates[i][y])
			for a := range gameProperties.Board {
				var tmpValue = g.leftTableLabels[i][y]
				if gameProperties.Board[a] == tmpValue {
					g.leftTableStates[i][y] = Ship
				}
			}

		}
	}
}

func getShipsFromTable(g *GUI) [20]string {
	var board [20]string
	index := 0
	for i := range g.leftTableLabels {
		for y := range g.leftTableLabels[i] {
			if g.leftTableStates[i][y] == 1 {
				board[index] = g.leftTableLabels[i][y]
				index += 1
			}
		}
	}
	return board

}
func putShipsOnLeftTable(g *GUI) {
	for i := range g.leftTableLabels {
		for y := range g.leftTableLabels[i] {
			for a := range gameProperties.Board {
				if g.leftTableLabels[i][y] == gameProperties.Board[a] {
					g.leftTableStates[i][y] = Ship
				}
			}
		}
	}
}

func loop(w *app.Window, g *GUI) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	var ops op.Ops

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if g.inGame {
				if g.sessionActive {

					gameProperties.gameStatus, _ = Status()

					var previousTime = g.timeLeft
					if len(fmt.Sprintf("%.0f", gameProperties.gameStatus.Body["timer"])) > 2 {
						g.timeLeft = previousTime
					} else {
						g.timeLeft = fmt.Sprintf("%.0f", gameProperties.gameStatus.Body["timer"])
					}

					if gameProperties.gameStatus.Body["should_fire"] == true {
						g.yourTurnIncidator = true

						if !g.enemyDescriptionAvailable {
							gameProperties.gameDescription, _ = GameDescription()

							g.enemyDescription = fmt.Sprintf("%s", gameProperties.gameDescription.Body["opp_desc"])
							g.enemyName = fmt.Sprintf("%s", gameProperties.gameDescription.Body["opponent"])
							fmt.Println("gamedescription")
							for key, value := range gameProperties.gameStatus.Body {
								fmt.Printf("gamedesc: %s: %v\n", key, value)
							}

						}
						g.enemyDescriptionAvailable = true

					} else {
						g.yourTurnIncidator = false
					}
					var tmpOppShots = fmt.Sprintf("%s", gameProperties.gameStatus.Body["opp_shots"])
					gameProperties.opp_shots = stringToSlice(tmpOppShots)

					fmt.Println(" opp shots:")

					var sessionTerminate = fmt.Sprintf("%s", gameProperties.gameStatus.Body["message"])
					if sessionTerminate == "session not found" {
						g.sessionActive = false
					}
					fmt.Println(" player shots:")

					for i := range gameProperties.PlayerShoots {
						fmt.Sprintf("%s", gameProperties.PlayerShoots[i])
					}

					var gameStatus = fmt.Sprintf("%s", gameProperties.gameStatus.Body["game_status"])
					if gameStatus == "lose" {

						g.sessionActive = false
						g.youLoseScreen = true
						g.displayPlayerAndEnemyBoard = false
						g.inGame = false
					}
					if gameStatus == "win" {
						g.sessionActive = false
						g.youLoseScreen = true
						g.displayPlayerAndEnemyBoard = false
						g.inGame = false
					}

					for i := range gameProperties.opp_shots {
						fmt.Sprintf("%s,", gameProperties.opp_shots[i])
					}
					handlePlayerBoard(g)
					for key, value := range gameProperties.gameStatus.Body {

						fmt.Printf("%s: %v\n", key, value)
					}
				}

			}
		}
	}()
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if g.startButton.Clicked(gtx) {
				g.showStartMenu = false
				g.showShipSetUpMenu = true

			}
			if g.personalizationButton.Clicked(gtx) {
				g.showStartMenu = false
				g.showPersonalization = true
			}
			if g.acceptButton.Clicked(gtx) {
				g.showStartMenu = true
				g.showPersonalization = false

				gameProperties.Nick = g.nickname.Text()
				gameProperties.Description = g.nickname.Text()

				fmt.Printf("Nickname: %s\n", g.nickname.Text())
				fmt.Printf("Description: %s\n", g.profileDescription.Text())
			}
			if g.discardButton.Clicked(gtx) {
				g.showStartMenu = true
				g.showPersonalization = false
				gameProperties.Nick = ""
				gameProperties.Description = ""
				fmt.Printf("Nickname: %s\n", g.nickname.Text())
				fmt.Printf("Description: %s\n", g.profileDescription.Text())

			}
			if g.backFromStatsButton.Clicked(gtx) {
				g.showStartMenu = true
				g.showStats = false
			}
			if g.showStatsButton.Clicked(gtx) {
				g.showStartMenu = false
				g.showStats = true
			}

			if g.exitButton.Clicked(gtx) {
				os.Exit(0)
			}
			if g.botActiveButton.Clicked(gtx) {
				if g.botActive {
					g.botActive = false
				} else {
					g.botActive = true
				}
			}

			if g.acceptShipPositions.Clicked(gtx) {
				g.lockLeftTable = true
				g.lockRightTable = false
				g.showShipSetUpMenu = false
				g.displayPlayerAndEnemyBoard = true
				g.inGame = true

				if g.nickname.Text() == "" {
					gameProperties.Nick = "BetonoJanusz"
				}

				gameProperties.Nick = g.nickname.Text()
				gameProperties.Description = g.profileDescription.Text()
				gameProperties.Board = getShipsFromTable(g)
				if g.botActive {
					gameProperties.wpBot = true
				} else {
					gameProperties.wpBot = false
				}

				fmt.Println("resultingtable")
				for e := range gameProperties.Board {
					fmt.Println(gameProperties.Board[e])
				}

				err := InitGame()
				if err != nil {
					panic(err)
				}
				gameProperties.gameStatus, err = Status()
				if err != nil {
					panic(err)
				}

				fmt.Println("gamestatus")
				for key, value := range gameProperties.gameStatus.Body {
					fmt.Printf("%s: %v\n", key, value)
				}

				fmt.Print("gameproperties board")
				fmt.Print(gameProperties.Board)
				placeShipsOnLeftTable(gtx, g)
				fmt.Printf("accepted ship positions")
			}
			if g.discardShipPositions.Clicked(gtx) {
				g.leftTableStates = createEmptyState(10, 10)
				fmt.Printf("discarded ship positions")
			}
			if g.randomShipPositions.Clicked(gtx) {
				gameProperties.Board, _ = customBoard()
				g.leftTableStates = createEmptyState(10, 10)
				putShipsOnLeftTable(g)
				//g.showShipSetUpMenu = false
				//tmp for testing
				//g.showLoadingMenu = true
				fmt.Printf("random ship positions")
			}
			if g.backShipPositions.Clicked(gtx) {
				g.showStartMenu = true
				g.displayPlayerAndEnemyBoard = false
			}

			if g.abandonButton.Clicked(gtx) {
				DeleteGame()
				g.displayPlayerAndEnemyBoard = false
				g.showStartMenu = true
			}
			if g.sessionActive == false {
				g.displayPlayerAndEnemyBoard = false
				g.showStartMenu = true
			}
			if g.refreshLobbyButton.Clicked(gtx) {

				playerlist, err := PlayerList()
				if err != nil {
					panic(err)
				}
				for key, value := range playerlist.Body {

					fmt.Printf("%s: %v\n", key, value)
				}
			}

			Layout(gtx, g)

			e.Frame(gtx.Ops)
		}

	}
}

func emptyLayoutDebug(gtx layout.Context, g *GUI) layout.Dimensions {
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx C) D {
				margins := layout.Inset{
					Top:    unit.Dp(25),
					Bottom: unit.Dp(25),
					Right:  unit.Dp(35),
					Left:   unit.Dp(35),
				}
				return margins.Layout(gtx,
					func(gtx C) D {
						btn := material.Button(g.theme, g.startButton, "Start")
						return btn.Layout(gtx)
					},
				)
			},
		))
}

func displayEnemyNameAndDescription(gtx layout.Context, g *GUI) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceEvenly}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			timerText := fmt.Sprintf("Enemy Name: %s", g.enemyName)
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.H4(g.theme, timerText).Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			timerText := fmt.Sprintf("Enemy Description: %s", g.enemyDescription)
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.H6(g.theme, timerText).Layout(gtx)
				})
			})
		}),
	)
}

func Layout(gtx layout.Context, g *GUI) layout.Dimensions {
	if g.showStartMenu {
		return startMenu(gtx, g)
	}
	if g.showPersonalization {
		return personalizationMenu(gtx, g)
	}
	if g.showShipSetUpMenu {
		return boardSelectMenu(gtx, g)
	}
	if g.displayPlayerAndEnemyBoard {
		return displayPlayerAndEnemyBoard(gtx, g)
	}
	if g.showLoadingMenu {
		return loadingMenu(gtx, g)
	}
	if g.showGameOver {
		return gameOver(gtx, g)
	}
	if g.showGameWon {
		return gameWon(gtx, g)
	}
	if g.showStats {
		return showStatsMenu(gtx, g)
	}
	if !g.sessionActive {
		return sessionTerminated(gtx, g)
	}
	if g.youWinScreen {
		return gameWon(gtx, g)
	}
	if g.youLoseScreen {
		return gameOver(gtx, g)
	}
	return emptyLayoutDebug(gtx, g)
}
