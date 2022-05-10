package deal

type NodeInfo interface {
	Update(other NodeInfo)
	TerminalValue(dp *DealPlay) NodeInfo
}

func ResetNodeInfo(node NodeInfo) NodeInfo {
	return TrickDistributionCounter{make(map[[NumberOfHands]int8]int)}
}

type TrickDistributionCounter struct {
	trickMap map[[NumberOfHands]int8]int
}

func (td TrickDistributionCounter) Update(other NodeInfo) {
	for tricks, nPaths := range other.(TrickDistributionCounter).trickMap {
		td.trickMap[tricks] += nPaths
	}
}

func (td TrickDistributionCounter) TerminalValue(dp *DealPlay) NodeInfo {
	result := make(map[[NumberOfHands]int8]int)
	result[dp.result]++
	return TrickDistributionCounter{result}
}

type NodeMove struct {
	moves []Card
	NodeInfo
}

type Tree struct {
	levels []NodeMove
	play   *DealPlay
}

func NewTree(play *DealPlay) Tree {
	deckSize := play.deckSize
	levels := make([]NodeMove, deckSize)
	var i int8
	for i = 0; i < deckSize; i++ {
		levels[i].moves = make([]Card, 0, 1+(deckSize-1-i)/NumberOfHands)
	}
	levels[0].NodeInfo = ResetNodeInfo(levels[0].NodeInfo)
	return Tree{levels, play}
}

func (t Tree) Eval(policy DensePolicy, level int8) {
	if t.play.deckSize == 0 {
		return
	}
	if len(t.levels[level].moves) == 0 {
		t.levels[level].moves = t.play.GetMoves(policy)
	}
	for _, move := range t.levels[level].moves {
		t.play.DoMove(move)
		if t.play.deckSize > 0 {
			t.levels[level+1].NodeInfo = ResetNodeInfo(t.levels[level+1].NodeInfo)
			t.Eval(policy, level+1)
			t.levels[level].NodeInfo.Update(t.levels[level+1].NodeInfo)
		} else {
			t.levels[level].Update(t.levels[level].TerminalValue(t.play))
		}
		t.play.UndoMove()
	}
	t.levels[level].moves = t.levels[level].moves[:0]
}
