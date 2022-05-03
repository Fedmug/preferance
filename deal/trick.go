package deal

import "strings"

type Trick struct {
	cards         []Card
	partialTakers []HandIndex
}

func (t *Trick) Len() int {
	return len(t.cards)
}

func (t *Trick) TakerCard() Card {
	if len(t.partialTakers) == 0 {
		return Card{-1, false}
	}
	return t.cards[t.partialTakers[len(t.partialTakers)-1]]
}

func NewTrick() *Trick {
	cards := make([]Card, 0, NumberOfHands)
	partialTakers := make([]HandIndex, 0, NumberOfHands)
	return &Trick{cards, partialTakers}
}

func (t *Trick) Append(card Card) {
	t.cards = append(t.cards, card)
	var nextTaker HandIndex
	if t.Len() == 1 || card.Beats(t.TakerCard()) {
		nextTaker = HandIndex(t.Len() - 1)
	} else {
		nextTaker = t.partialTakers[len(t.partialTakers)-1]
	}
	t.partialTakers = append(t.partialTakers, nextTaker)
}

func (t *Trick) Pop() Card {
	if t.Len() == 0 {
		panic("Trick must not be empty")
	}
	result := t.cards[t.Len()-1]
	t.partialTakers = t.partialTakers[:len(t.partialTakers)-1]
	t.cards = t.cards[:t.Len()-1]
	return result
}

func (t *Trick) String() string {
	result := ""
	for _, card := range t.cards {
		result += card.String() + " "
	}
	return strings.TrimSpace(result)
}
