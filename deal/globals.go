package deal

import "math"

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
type HandPosition int8
type MiniMaxType int8
type MiniMaxValue int8

const (
	MinType MiniMaxType = iota
	MaxType
	UnknownType = -1
)

var FirstHandContract = DealGoals{MaxType, MinType, MinType}
var SecondHandContract = DealGoals{MinType, MaxType, MinType}
var ThirdHandContract = DealGoals{MinType, MinType, MaxType}

var FirstHandMisere = DealGoals{MinType, MaxType, MaxType}
var SecondHandMisere = DealGoals{MaxType, MinType, MaxType}
var ThirdHandMisere = DealGoals{MaxType, MaxType, MinType}

const MiniMaxInf MiniMaxValue = math.MaxInt8

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

const (
	North HandPosition = iota
	East
	South
	West
	Drop
	InvalidPosition = -1
)
