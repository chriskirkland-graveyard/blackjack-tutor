package main

import (
	"./blackjack"
	"./blackjackui"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func winnerScreen() {
	fmt.Println("\nPlayer Wins!!!\n")
}

func loserScreen() {
	fmt.Println("\nDealer Wins.\n")
}

func pushScreen() {
	fmt.Println("\nPush.\n")
}

func randomize() {
	// seed RNG
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	randomize()

	gameCount := 0
	gamesBeforeShuffle := 5
	var myGame blackjack.Game
	for {
		// shuffle up and deal
		if gameCount%gamesBeforeShuffle == 0 {
			myGame = blackjack.NewGame()
		}
		myGame.NewHand()

		// player decisions
		var input string
		for myGame.PlayerCanHit() && input != "s" {
			clearScreen()
			fmt.Println(myGame)
			input := blackjackui.PromptUser("What do you want to do (h/s)?")
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
		myGame.GoDealer()
		clearScreen()
		fmt.Println(myGame)

		// decide winners
		switch myGame.GetWinner() {
		case blackjack.StatePlayerWins:
			winnerScreen()
		case blackjack.StateDealerWins:
			loserScreen()
		case blackjack.StatePush:
			pushScreen()
		}

		fmt.Print("Press ENTER to continue...")
		fmt.Scanln()

		gameCount++
	}
}
