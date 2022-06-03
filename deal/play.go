package deal

const (
	trumpOffset            = 1
	contingencyTableOffset = 21
)

type DealResult [NumberOfHands]int8
type DealGoals [NumberOfHands]MiniMaxType

type DealPlay struct {
	tricks   []Trick
	matrix   *DealMatrix
	result   DealResult
	trump    Suit
	mover    HandIndex
	deckSize int8
	goals    DealGoals
}

func NewDealPlay(matrix *DealMatrix, trump Suit, firstMover HandIndex, goals DealGoals) *DealPlay {
	var dealResult DealResult
	deckSize := matrix.DeckSize()
	return &DealPlay{
		tricks:   make([]Trick, 0, deckSize/NumberOfHands),
		matrix:   matrix,
		result:   dealResult,
		trump:    trump,
		mover:    firstMover,
		deckSize: deckSize,
		goals:    goals,
	}
}

func (dp *DealPlay) MoverType() MiniMaxType {
	return dp.goals[dp.mover]
}

func (dp *DealPlay) LastMove() Card {
	if len(dp.lastTrick().moves) > 0 {
		return dp.lastTrick().moves[len(dp.lastTrick().moves)-1].card
	}
	if len(dp.tricks) > 1 {
		return dp.tricks[len(dp.tricks)-2].moves[len(dp.tricks[len(dp.tricks)-2].moves)-1].card
	}
	return InvalidCard
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
	dp.lastTrick().append(move)
	dp.matrix.Remove(move, dp.mover)
	dp.deckSize--
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
	dp.matrix.Add(move, dp.mover)
	dp.deckSize++
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
	result := make([]Card, 0, dp.deckSize/NumberOfHands)
	var i Suit
	for i = 0; i < NumberOfSuits; i++ {
		result = append(result, dp.matrix.GetMoves(i, dp.mover, policy, i == dp.trump)...)
	}
	return result
}

func (dp *DealPlay) Index(endian Endian) int64 {
	var result int64
	if dp.trump != InvalidSuit {
		result = 1
	}
	result += int64(contingencyTableMap[dp.matrix.ContingencyTable()]) << trumpOffset
	result += dp.matrix.Index(endian) << contingencyTableOffset
	return result
}

func NewDealPlayFromIndex(index int64, endian Endian, goals DealGoals) *DealPlay {
	var trump Suit = InvalidSuit
	if index%2 > 0 {
		trump = Spades
	}
	table := contingencyTables[index%(1<<contingencyTableOffset)>>trumpOffset]
	dealMatrix := DealMatrixFromIndex(table, index>>contingencyTableOffset, endian)
	deckSize := dealMatrix.DeckSize()
	return &DealPlay{
		tricks:   make([]Trick, 0, deckSize/NumberOfHands),
		matrix:   &dealMatrix,
		trump:    trump,
		mover:    FirstHand,
		deckSize: deckSize,
		goals:    goals,
	}
}
