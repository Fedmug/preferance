package deal

import "strings"

type TrickMove struct {
	card  Card
	taker HandIndex
}

type Trick struct {
	moves      []TrickMove
	firstMover HandIndex
}

func (t *Trick) TakerMove() TrickMove {
	if len(t.moves) == 0 {
		return TrickMove{Card{-1, false}, InvalidHand}
	}
	return t.moves[len(t.moves)-1]
}

func NewTrick(firstMover HandIndex, cap int) *Trick {
	return &Trick{make([]TrickMove, 0, cap), firstMover}
}

func (t *Trick) Append(card Card) {
	var nextTaker HandIndex
	if len(t.moves) == 0 {
		nextTaker = t.firstMover
	} else {
		highestCardIndex := (t.moves[len(t.moves)-1].taker + NumberOfHands - t.firstMover) % NumberOfHands
		if card.Beats(t.moves[highestCardIndex].card) {
			nextTaker = (t.firstMover + HandIndex(len(t.moves))) % NumberOfHands
		} else {
			nextTaker = t.moves[len(t.moves)-1].taker
		}
	}
	t.moves = append(t.moves, TrickMove{card, nextTaker})
}

func (t *Trick) Pop() Card {
	if len(t.moves) == 0 {
		panic("Trick must not be empty")
	}
	result := t.moves[len(t.moves)-1].card
	t.moves = t.moves[:len(t.moves)-1]
	return result
}

func (t *Trick) String() string {
	result := ""
	for _, move := range t.moves {
		result += move.card.String() + " "
	}
	return strings.TrimSpace(result)
}
