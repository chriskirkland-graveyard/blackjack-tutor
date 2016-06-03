package blackjack

import "fmt"
import "math/rand"

// **** CONSTANTS ****
var ranks = "A23456789TJQK"
var suits = "HDCS"

// ***** STRUCTS *****
type Card struct {
	rank int
	suit int
}

type Shoe struct {
	cards          [13][4]int // 13 ranks * 4 suits
	remainingCards int
	numDecks       int
}

// **** STRUCT FUNCTIONS ***
func (c Card) String() string {
	return fmt.Sprintf("%c%c", ranks[c.rank], suits[c.suit])
}

func (s Shoe) String() string {
	var str string
	str += fmt.Sprintf("  %2c%2c%2c%2c\n", suits[0], suits[1], suits[2], suits[3])
	for row := 0; row < 13; row++ {
		str += fmt.Sprintf("%2c%2d%2d%2d%2d\n", ranks[row], s.cards[row][0], s.cards[row][1], s.cards[row][2], s.cards[row][3])
	}
	return str
}

func (s *Shoe) Shuffle() {
	for i := 0; i < 13; i++ {
		for j := 0; j < 4; j++ {
			s.cards[i][j] = s.numDecks
		}
	}
	s.remainingCards = 52 * s.numDecks
}

func (s *Shoe) DealCard() Card {
	for {
		rank := rand.Intn(13)
		suit := rand.Intn(4)
		if s.cards[rank][suit] > 0 {
			s.cards[rank][suit] -= 1
			s.remainingCards -= 1
			return Card{rank, suit}
		}
	}
}

func NewShoe(desiredNumDecks int) *Shoe {
	s := Shoe{numDecks: desiredNumDecks}
	s.Shuffle()
	return &s
}
