package deal

import "strings"

type trickIndex int8

const invalidTrickIndex trickIndex = -1

type TrickMove struct {
	card                Card
	prevBeatenCardIndex trickIndex
}

type Trick struct {
	moves            []TrickMove
	highestCardIndex trickIndex
	// firstMover HandIndex
}

func NewTrick(cap int) *Trick {
	return &Trick{make([]TrickMove, 0, cap), invalidTrickIndex}
}

func (t *Trick) takerCard() Card {
	if t.highestCardIndex == invalidTrickIndex {
		return Card{InvalidCardCode, false}
	}
	return t.moves[t.highestCardIndex].card
}

func (t *Trick) takerHandIndex(firstMover HandIndex) HandIndex {
	if t.highestCardIndex == invalidTrickIndex {
		return InvalidHand
	}
	return (firstMover + HandIndex(t.highestCardIndex)) % NumberOfHands
}

func (t *Trick) append(card Card) {
	newMove := TrickMove{card, invalidTrickIndex}
	if card.Beats(t.takerCard()) {
		newMove.prevBeatenCardIndex = t.highestCardIndex
		t.highestCardIndex = trickIndex(len(t.moves))
	}
	t.moves = append(t.moves, newMove)
}

func (t *Trick) pop() Card {
	if len(t.moves) == 0 {
		panic("Trick must not be empty")
	}
	lastMove := t.moves[len(t.moves)-1]
	if len(t.moves) == 1 || lastMove.prevBeatenCardIndex != invalidTrickIndex {
		t.highestCardIndex = lastMove.prevBeatenCardIndex
	}
	t.moves = t.moves[:len(t.moves)-1]
	return lastMove.card
}

func (t *Trick) String() string {
	result := ""
	for _, move := range t.moves {
		result += move.card.String() + " "
	}
	return strings.TrimSpace(result)
}
