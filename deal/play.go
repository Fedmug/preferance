package deal

type DealResult [NumberOfHands]int8

type DealPlay struct {
	tricks []Trick
	matrix *DealMatrix
	result DealResult
	trump  Suit
	mover  HandIndex
}

func NewDealPlay(matrix *DealMatrix, trump Suit, firstMover HandIndex) *DealPlay {
	var dealResult DealResult
	return &DealPlay{
		tricks: make([]Trick, 0, matrix.DeckSize()/NumberOfHands),
		matrix: matrix,
		result: dealResult,
		trump:  trump,
		mover:  firstMover,
	}
}

func (dp *DealPlay) addNewTrick(firstMover HandIndex) {
	newTrick := NewTrick(firstMover, NumberOfHands)
	dp.tricks = append(dp.tricks, *newTrick)
}

func (dp *DealPlay) lastTrick() *Trick {
	return &dp.tricks[len(dp.tricks)-1]
}

func (dp *DealPlay) taker() HandIndex {
	return dp.lastTrick().takerHandIndex()
}

func (dp *DealPlay) toNextTrick() {
	if len(dp.tricks) > 0 {
		dp.mover = dp.taker()
	}
	dp.addNewTrick(dp.mover)
}

func (dp *DealPlay) toPreviousTrick() {
	dp.tricks = dp.tricks[:len(dp.tricks)-1]
	dp.mover = InvalidHand
}

func (dp *DealPlay) DoMove(move Card) {
	if len(dp.tricks) == 0 || len(dp.lastTrick().moves) == NumberOfHands {
		dp.toNextTrick()
	}
	suit := move.Suit()
	rank := move.Rank()
	dp.lastTrick().append(move)
	dp.matrix[suit][dp.mover] -= 1 << rank
	if len(dp.lastTrick().moves) < NumberOfHands {
		dp.mover = (dp.mover + 1) % NumberOfHands
	} else {
		dp.mover = dp.taker()
		dp.result[dp.mover]++
	}

}

func (dp *DealPlay) UndoMove() {
	if len(dp.lastTrick().moves) == 0 {
		dp.toPreviousTrick()
	}
	if len(dp.lastTrick().moves) == NumberOfHands {
		dp.result[dp.taker()]--
	}
	move := dp.lastTrick().pop()
	dp.mover = (dp.lastTrick().firstMover + HandIndex(len(dp.lastTrick().moves))) % NumberOfHands
	dp.matrix[move.Suit()][dp.mover] += 1 << move.Rank()
}

func (dp *DealPlay) GetMoves(policy DensePolicy) []Card {
	lenOfLastTrick := 0
	if len(dp.tricks) > 0 {
		lenOfLastTrick = len(dp.lastTrick().moves)
	}
	if lenOfLastTrick > 0 && lenOfLastTrick < NumberOfHands {
		firstTrickCardSuit := dp.lastTrick().moves[0].card.Suit()
		if dp.matrix[firstTrickCardSuit][dp.mover] > 0 {
			return dp.matrix.GetMoves(firstTrickCardSuit, dp.mover, policy, firstTrickCardSuit == dp.trump)
		}
		if dp.trump < NT && dp.matrix[dp.trump][dp.mover] > 0 {
			return dp.matrix.GetMoves(dp.trump, dp.mover, policy, true)
		}
	}
	result := make([]Card, 0, dp.matrix.DeckSize()/NumberOfHands)
	var i Suit
	for i = 0; i < NumberOfSuits; i++ {
		result = append(result, dp.matrix.GetMoves(i, dp.mover, policy, i == dp.trump)...)
	}
	return result
}
