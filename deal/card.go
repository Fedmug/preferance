package deal

type CardCode int8

const InvalidCardCode CardCode = -1

// code = Rank + SuitLength * Suit
type Card struct {
	code  CardCode
	trump bool
}

var InvalidCard = Card{InvalidCardCode, false}

func NewCard(suit Suit, rank Rank, trump bool) Card {
	return Card{CardCode(rank) + SuitLength*CardCode(suit), trump}
}

func (c Card) Suit() Suit {
	return Suit(c.code / SuitLength)
}

func (c Card) Rank() Rank {
	return Rank(c.code % SuitLength)
}

func (c Card) Beats(other Card) bool {
	if c.code == InvalidCardCode {
		return false
	}
	return (other.code == InvalidCardCode) || (c.trump && !other.trump) ||
		(c.Suit() == other.Suit()) && (c.Rank() > other.Rank())
}

func (c Card) String() string {
	if c.Suit() == InvalidSuit || c.Rank() == InvalidRank {
		return ""
	}
	if c.trump {
		return string(TrumpSuitSymbols[c.Suit()]) + string(RankSymbols[c.Rank()])
	}
	return string(SuitSymbols[c.Suit()]) + string(RankSymbols[c.Rank()])
}
