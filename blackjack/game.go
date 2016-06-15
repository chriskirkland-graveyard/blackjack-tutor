package blackjack

import "bytes"
import "fmt"

const (
	StatePlayerWins = iota
	StateDealerWins
	StatePush
)

type Player struct {
	cards    []Card
	holeCard bool
}

func (p Player) String() string {
	var buffer bytes.Buffer
	for ix, card := range p.cards {
		if ix == 0 && p.holeCard {
			// "Full Block" (█) == <Ctrl>+FB (in Insert mode)
			buffer.WriteString("██ ")
		} else {
			buffer.WriteString(fmt.Sprintf("%v ", card))
		}
	}
	if !p.holeCard {
		buffer.WriteString(fmt.Sprintf("(%d)", p.Count()))
	}
	return buffer.String()
}

func (p *Player) addCard(c Card) {
	p.cards = append(p.cards, c)
}

func (p Player) Count() int {
	var count int
	var aceCount int
	for _, card := range p.cards {
		count += rankValues[card.rank]
		if card.rank == 0 {
			aceCount++
		}
	}
	for count > 21 && aceCount > 0 {
		count -= 10
		aceCount--
	}
	return count
}

func (p *Player) hasBlackjack() bool {
	if len(p.cards) != 2 || p.Count() != 21 {
		return false
	} else {
		return true
	}
}

func (p *Player) isBust() bool {
	return p.Count() > 21
}

type Game struct {
	shoe              Shoe
	player            Player
	dealer            Player
	ruleset           ruleSet
	handsSinceShuffle int
}

func (g Game) String() string {
	return fmt.Sprintf(
		"Shoe: \n%v\nDealer: %v\nPlayer: %v\n",
		g.shoe,
		g.dealer,
		g.player,
	)
}

func NewGame() Game {
	// initialize game
	myRuleset := newRuleSet()
	myShoe := NewShoe(myRuleset.numDecks)
	return Game{shoe: myShoe, ruleset: myRuleset}
}

func (g *Game) NewHand() {
	g.player = Player{holeCard: false}
	g.dealer = Player{holeCard: true}

	// deal cards
	g.DealPlayer()
	g.DealDealer()
	g.DealPlayer()
	g.DealDealer()
}

func (g *Game) NeedsShuffle() bool {
	return g.handsSinceShuffle >= g.ruleset.handsBeforeShuffle
}

func (g *Game) PlayerCanHit() bool {
	return g.player.Count() < 21
}

func (g *Game) DealPlayer() {
	g.player.addCard(g.shoe.DealCard())
}

func (g *Game) DealDealer() {
	g.dealer.addCard(g.shoe.DealCard())
}

func (g *Game) GoDealer() {
	if !g.player.hasBlackjack() {
		for g.dealer.Count() < 17 {
			g.DealDealer()
		}
	}
	g.dealer.holeCard = false // reveal card
}

func (g *Game) GetWinner() int {
	g.handsSinceShuffle++

	playerCount := g.player.Count()
	dealerCount := g.dealer.Count()
	if g.player.hasBlackjack() {
		return StatePlayerWins
	} else if g.player.isBust() {
		return StateDealerWins
	} else if g.dealer.isBust() {
		return StatePlayerWins
	} else if playerCount > dealerCount {
		return StatePlayerWins
	} else if playerCount < dealerCount {
		return StateDealerWins
	} else if playerCount == dealerCount {
		return StatePush
	} else {
		panic(fmt.Sprintf(
			"Invalid game state! (pc=%d,dc=%d)",
			playerCount,
			dealerCount,
		))
	}
}
