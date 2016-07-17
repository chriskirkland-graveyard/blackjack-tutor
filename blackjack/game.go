package blackjack

import "bytes"
import "fmt"
import "math"

const (
	StatePlayerWins = iota
	StateDealerWins
	StatePush
)

const (
	HAND_WIN = iota
	HAND_LOSE
	HAND_PUSH
)

type Record struct {
	handsPlayed int
	handStats   [3]int // [win, lose, push]
	chipCount   float32
}

func (r Record) String() string {
	safeHandsPlayed := math.Max(1, float64(r.handsPlayed))
	return fmt.Sprintf(
		"Hands: %-5d Win: %d (%1.2f%%)   Loss: %d (%1.2f%%)   Push: %-3d\nChip Count: %1.2f\n",
		r.handsPlayed,
		r.handStats[HAND_WIN],
		float64(r.handStats[HAND_WIN]-r.handStats[HAND_PUSH])/safeHandsPlayed*100,
		r.handStats[HAND_LOSE],
		float64(r.handStats[HAND_LOSE]-r.handStats[HAND_PUSH])/safeHandsPlayed*100,
		r.handStats[HAND_PUSH],
		r.chipCount,
	)
}

type Player struct {
	cards    []Card
	holeCard bool
	record   Record
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

func (p *Player) hasAce() bool {
	for _, card := range p.cards {
		if card.rank == 0 { // Ace
			return true
		}
	}
	return false
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
		// "Shoe: \n%v\nDealer: %v\nPlayer: %v\n",
		// g.shoe,
		"STATS:\n%v\n\nCURRENT HAND:\nDealer: %v\nPlayer: %v\n",
		g.player.record,
		g.dealer,
		g.player,
	)
}

func NewGame() Game {
	// initialize game
	myRuleset := newRuleSet()
	myShoe := NewShoe(myRuleset.numDecks)
	myPlayer := Player{record: Record{chipCount: 200}}
	return Game{
		shoe:    myShoe,
		ruleset: myRuleset,
		player:  myPlayer,
	}
}

func (g *Game) Shuffle() {
	g.shoe.Shuffle()
}

func (g *Game) NewHand() {
	// hide dealer hole card
	g.dealer.holeCard = true

	// reset cards
	g.player.cards = []Card{}
	g.dealer.cards = []Card{}

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

func (g *Game) QPlayerBust() bool {
	return g.player.Count() > 21
}

func (g *Game) QInsuranceAvailable() bool {
	return g.ruleset.insurance &&
		len(g.dealer.cards) == 2 &&
		g.dealer.cards[1].rank == 0 // Dealer showing Ace
}

func (g *Game) InsurancePays() bool {
	if g.dealer.hasBlackjack() {
		g.player.record.chipCount += 2
		return true
	} else {
		g.player.record.chipCount -= 1
		return false
	}
}

func (g *Game) DealPlayer() {
	g.player.addCard(g.shoe.DealCard())
}

func (g *Game) DealDealer() {
	g.dealer.addCard(g.shoe.DealCard())
}

func (g *Game) GoDealer() {
	if !g.player.hasBlackjack() {
		for g.dealer.Count() < 17 ||
			(g.ruleset.hitSoft17 && g.dealer.Count() == 17 && g.dealer.hasAce()) {
			g.DealDealer()
		}
	}
	g.dealer.holeCard = false // reveal card
}

func (g *Game) playerWins() int {
	if g.player.hasBlackjack() {
		g.player.record.chipCount += g.ruleset.blackjackPayout
	} else {
		g.player.record.chipCount += 1
	}
	g.player.record.handStats[HAND_WIN]++
	return StatePlayerWins
}

func (g *Game) dealerWins() int {
	g.player.record.handStats[HAND_LOSE]++
	g.player.record.chipCount -= 1
	return StateDealerWins
}

func (g *Game) push() int {
	g.player.record.handStats[HAND_PUSH]++
	return StatePush
}

func (g *Game) GetWinner() int {
	g.handsSinceShuffle++
	g.player.record.handsPlayed++

	playerCount := g.player.Count()
	dealerCount := g.dealer.Count()
	if g.player.isBust() {
		return g.dealerWins()
	} else if g.dealer.isBust() {
		return g.playerWins()
	} else if playerCount > dealerCount {
		return g.playerWins()
	} else if playerCount < dealerCount {
		return g.dealerWins()
	} else if playerCount == dealerCount {
		return g.push()
	} else {
		panic(fmt.Sprintf(
			"Invalid game state! (pc=%d,dc=%d)",
			playerCount,
			dealerCount,
		))
	}
}
