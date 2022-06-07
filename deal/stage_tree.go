package deal

import "log"

const nDealGoals = 6

type StagePlay struct {
	matrix DealMatrix
	trick  Trick
	trump  bool
}

func (sp *StagePlay) GetMoves(policy DensePolicy) []Card {
	moverIndex := HandIndex(len(sp.trick.moves))
	if moverIndex > 0 {
		firstTrickCardSuit := sp.trick.moves[0].card.Suit()
		if sp.matrix[firstTrickCardSuit][moverIndex] > 0 {
			return sp.matrix.GetMoves(firstTrickCardSuit, moverIndex, policy, sp.trump && firstTrickCardSuit == 0)
		}
		if sp.trump && sp.matrix[0][moverIndex] > 0 {
			return sp.matrix.GetMoves(0, moverIndex, policy, true)
		}
	}
	result := make([]Card, 0, sp.matrix.DeckSize()/NumberOfHands)
	var i Suit
	for i = 0; i < NumberOfSuits; i++ {
		result = append(result, sp.matrix.GetMoves(i, moverIndex, policy, sp.trump && i == 0)...)
	}
	return result
}

func (sp *StagePlay) DoMove(move Card) {
	moverIndex := HandIndex(len(sp.trick.moves))
	if moverIndex == NumberOfHands {
		log.Fatalf("trick %v is full", sp.trick)
	}
	sp.trick.append(move)
	sp.matrix.Remove(move, moverIndex)
}

func (sp *StagePlay) UndoMove() {
	if len(sp.trick.moves) == 0 {
		log.Fatalf("trick %v is empty", sp.trick)
	}
	move := sp.trick.pop()
	sp.matrix.Add(move, HandIndex(len(sp.trick.moves)))
}

func NewStagePlay(matrix DealMatrix, trump bool) *StagePlay {
	return &StagePlay{matrix, *NewTrick(FirstHand, NumberOfHands), trump}
}

type NodeResult struct {
	nTricks int8
	nPaths  int64
}

type MiniMaxNodeInfo struct {
	totalPaths int64
	results    [nDealGoals]NodeResult
	move       Card
}

type StageTree struct {
	rootInfo MiniMaxNodeInfo
	levels   [NumberOfHands][]MiniMaxNodeInfo
	play     *StagePlay
}

func NewStageTree(sp *StagePlay, stage int8) StageTree {
	var levels [NumberOfHands][]MiniMaxNodeInfo
	for i := 0; i < NumberOfHands; i++ {
		levels[i] = make([]MiniMaxNodeInfo, 0, stage)
	}
	return StageTree{
		rootInfo: MiniMaxNodeInfo{move: InvalidCard},
		levels:   levels,
		play:     sp,
	}
}

/*
type DealGoal int8

const (
	MaxMinMin DealGoal = iota
	MinMaxMin
	MinMinMax
	MinMaxMax
	MaxMinMax
	MaxMaxMin
)

func (st StageTree) UpdateTaker(node *MiniMaxNode) {
	switch st.play.trick.takerHandIndex() {
	case FirstHand:
		node.results[MaxMinMin].nTricks++
		node.results[MinMaxMax].nTricks++
	case SecondHand:
		node.results[MinMaxMin].nTricks++
		node.results[MaxMinMax].nTricks++
	case ThirdHand:
		node.results[MinMinMax].nTricks++
		node.results[MaxMaxMin].nTricks++
	}
}

func (st StageTree) UpdateLastTrick(node *MiniMaxNode) {
	takerIndex := st.play.trick.takerHandIndex()
	lastTrick := *NewTrick(takerIndex, NumberOfHands)
	var handIndex HandIndex
	for ; handIndex < NumberOfHands; handIndex++ {
		for i := 0; i < NumberOfSuits; i++ {
			index := (handIndex + takerIndex) % NumberOfHands
			if st.play.matrix[i][index] > 0 {
				lastTrick.append(st.play.matrix.GetMoves(Suit(i), index, All, st.play.trump)[0])
				break
			}
		}
	}
	node.MiniMaxNodeInfo.totalPaths++
	switch (lastTrick.takerHandIndex() + takerIndex) % NumberOfHands {
	case FirstHand:
		node.results[MaxMinMin].nTricks++
		node.results[MinMaxMax].nTricks++
		node.results[MaxMinMin].nPaths++
		node.results[MinMaxMax].nPaths++
	case SecondHand:
		node.results[MinMaxMin].nTricks++
		node.results[MaxMinMax].nTricks++
		node.results[MinMaxMin].nPaths++
		node.results[MaxMinMax].nPaths++
	case ThirdHand:
		node.results[MinMinMax].nTricks++
		node.results[MaxMaxMin].nTricks++
		node.results[MinMinMax].nPaths++
		node.results[MaxMaxMin].nPaths++
	}
}

func (st StageTree) Eval(node *MiniMaxNode, policy DensePolicy) {
	deckSize := st.play.matrix.DeckSize()
	if deckSize%NumberOfHands == 0 {
		st.UpdateTaker(node)
	}
	if deckSize == NumberOfHands {
		st.UpdateLastTrick(node)
		return
	}
	moves := st.play.GetMoves(policy)
	node.children = make([]*MiniMaxNode, 0, len(moves))
	for _, move := range moves {
		node.children = append(node.children, &MiniMaxNode{move: move, parent: node})
		st.play.DoMove(move)
		st.Eval(node.children[len(node.children)-1], policy)
		// Update
		st.play.UndoMove()
	}
}
*/
