package blackjack

type ruleSet struct {
	blackjackPayout  float32
	numDecks         int
	insurance        bool
	surrender        bool
	hitSoft17        bool
	doubleAfterSplit bool
	resplitLimit     int
}
