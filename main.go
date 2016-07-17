package main

import (
	"fmt"
	"math/rand"
	"time"

	"./blackjack"
	"./blackjackui"
)

func randomize() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	// seed RNG
	randomize()

	// setup UI
	ui := new(blackjackui.ShellUI)

	gameCount := 0
	myGame := blackjack.NewGame()
	for {
		// shuffle up and deal
		if myGame.NeedsShuffle() {
			myGame.Shuffle()
		}
		myGame.NewHand()

		// player decisions
		var input string
		for myGame.PlayerCanHit() && input != "s" {
			ui.Redraw(myGame)

			// insurance check
			if myGame.QInsuranceAvailable() {
				input := ui.PromptUser("Would you like insurances (y/n)?")
				switch input {
				case "y":
					if myGame.InsurancePays() {
						ui.InsuranceWin()
					}
				case "n":
				default:
					panic(fmt.Sprintf("Invalid input: expected (y/n) but found \"%s\""))
				}
			}
			input := ui.PromptUser("What do you want to do (h/s)?")
			if input == "h" {
				// player gets and card and loop
				myGame.DealPlayer()
			} else if input == "s" {
				break
			} else {
				panic("Invalid input for player action!")
			}
		}

		// dealer does stuff
		if !myGame.QPlayerBust() {
			myGame.GoDealer()
		}
		ui.Redraw(myGame)

		// decide winners
		switch myGame.GetWinner() {
		case blackjack.StatePlayerWins:
			ui.WinnerScreen()
		case blackjack.StateDealerWins:
			ui.LoserScreen()
		case blackjack.StatePush:
			ui.PushScreen()
		}

		// continue...?
		ui.QContinue()
		gameCount++
	}
}
