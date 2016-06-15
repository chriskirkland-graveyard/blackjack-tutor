package blackjack

type ruleSet struct {
	blackjackPayout    float32
	numDecks           int
	handsBeforeShuffle int
	insurance          bool
	surrender          bool
	hitSoft17          bool
	doubleAfterSplit   bool
	resplitLimit       int
}

func newRuleSet() ruleSet {
	return ruleSet{
		blackjackPayout:    1.5,
		numDecks:           2,
		handsBeforeShuffle: 8,
		insurance:          false,
		surrender:          false,
		hitSoft17:          false,
		doubleAfterSplit:   false,
		resplitLimit:       1,
	}
}
