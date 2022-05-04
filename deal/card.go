package deal

const (
	NumberOfSuits = 4
	SuitLength    = 8
	NumberOfHands = 3
)

var SuitSymbols = [NumberOfSuits]rune{'\u2660', '\u2663', '\u2666', '\u2665'}
var TrumpSuitSymbols = [NumberOfSuits]rune{'\u2664', '\u2667', '\u2662', '\u2661'}
var RankSymbols = [SuitLength]rune{'7', '8', '9', 'T', 'J', 'Q', 'K', 'A'}

var SuitSymbolsMap = map[rune]Suit{
	'\u2660': Spades,
	'\u2663': Clubs,
	'\u2666': Diamonds,
	'\u2665': Hearts,
}

var RankSymbolsMap = map[rune]Suit{
	'7': 0,
	'8': 1,
	'9': 2,
	'T': 3,
	'J': 4,
	'Q': 5,
	'K': 6,
	'A': 7,
}

type Suit int8
type Rank int8
type HandIndex int8

const (
	Spades Suit = iota
	Clubs
	Diamonds
	Hearts
	NT
	InvalidSuit = -1
)

const (
	Seven Rank = iota
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
	InvalidRank = -1
)

const (
	FirstHand HandIndex = iota
	SecondHand
	ThirdHand
	FourthHand
	InvalidHand = -1
)

// code = Rank + SuitLength * Suit
type Card struct {
	code  int8
	trump bool
}

func NewCard(suit Suit, rank Rank, trump bool) Card {
	return Card{int8(rank) + SuitLength*int8(suit), trump}
}

func (c Card) Suit() Suit {
	return Suit(c.code / SuitLength)
}

func (c Card) Rank() Rank {
	return Rank(c.code % SuitLength)
}

func (c Card) Beats(other Card) bool {
	return (c.trump && !other.trump) || (c.Suit() == other.Suit()) && (c.Rank() > other.Rank())
}

func (c Card) String() string {
	if c.trump {
		return string(TrumpSuitSymbols[c.Suit()]) + string(RankSymbols[c.Rank()])
	}
	return string(SuitSymbols[c.Suit()]) + string(RankSymbols[c.Rank()])
}
