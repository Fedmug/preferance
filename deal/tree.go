package deal

import (
	"fmt"
	"sync"
)

type TrickDistributionCounter struct {
	trickMap map[[NumberOfHands]int8]int
}

type NodeCounter int64

type NodeInfo struct {
	type_        MiniMaxType
	totalPaths   int64
	miniMaxPaths int64
	value        MiniMaxValue
	bestMove     Card
}

func ResetNodeInfo(node NodeInfo, dp *DealPlay) NodeInfo {
	value := MiniMaxInf
	if dp.MoverType() == MaxType {
		value = -value
	}
	return NodeInfo{dp.MoverType(), 0, 0, value, InvalidCard}
}

func (td TrickDistributionCounter) Update(other TrickDistributionCounter) TrickDistributionCounter {
	for tricks, nPaths := range other.trickMap {
		td.trickMap[tricks] += nPaths
	}
	return TrickDistributionCounter{td.trickMap}
}

func (td TrickDistributionCounter) TerminalValue(dp *DealPlay) TrickDistributionCounter {
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

func NewTree(play *DealPlay, node NodeInfo) Tree {
	deckSize := play.deckSize
	levels := make([]NodeMove, deckSize)
	var i int8
	for i = 0; i < deckSize; i++ {
		levels[i].moves = make([]Card, 0, 1+(deckSize-1-i)/NumberOfHands)
	}
	levels[0].NodeInfo = ResetNodeInfo(node, play)
	return Tree{levels, play}
}

func MultiEval(play *DealPlay, policy DensePolicy) {
	var wg sync.WaitGroup
	for _, move := range play.GetMoves(policy) {
		var newMatrix DealMatrix = *play.matrix
		newPlay := NewDealPlay(&newMatrix, play.trump, play.mover, play.goals)
		newPlay.DoMove(move)
		value := MiniMaxInf
		if newPlay.goals[newPlay.mover] == MaxType {
			value = -value
		}
		tree := NewTree(newPlay, NodeInfo{newPlay.goals[newPlay.mover], 0, 0, value, InvalidCard})
		wg.Add(1)
		// fmt.Printf("Move: %v, trick: %v, deck size: %d\n", move, newPlay.lastTrick(), newPlay.deckSize)
		go func() {
			defer wg.Done()
			tree.Eval(policy, 0)
			fmt.Println("Node type:", tree.levels[0].NodeInfo.type_)
			fmt.Println("Total paths:", tree.levels[0].NodeInfo.totalPaths)
			fmt.Println("MiniMax paths:", tree.levels[0].NodeInfo.miniMaxPaths)
			fmt.Println("Node value:", tree.levels[0].NodeInfo.value)
			fmt.Println("Best move:", tree.levels[0].NodeInfo.bestMove)
		}()
	}
	wg.Wait()
}

func (t Tree) TerminalValue() NodeInfo {
	var result MiniMaxValue
	dp := t.play
	if (dp.goals[0] == MaxType && dp.goals[1] == MinType && dp.goals[2] == MinType) ||
		(dp.goals[0] == MinType && dp.goals[1] == MaxType && dp.goals[2] == MaxType) {
		result = MiniMaxValue(dp.result[0])
	} else if (dp.goals[0] == MinType && dp.goals[1] == MaxType && dp.goals[2] == MinType) ||
		(dp.goals[0] == MaxType && dp.goals[1] == MinType && dp.goals[2] == MaxType) {
		result = MiniMaxValue(dp.result[1])
	} else if (dp.goals[0] == MinType && dp.goals[1] == MinType && dp.goals[2] == MaxType) ||
		(dp.goals[0] == MaxType && dp.goals[1] == MaxType && dp.goals[2] == MinType) {
		result = MiniMaxValue(dp.result[2])
	}
	return NodeInfo{UnknownType, 1, 1, result, InvalidCard}
}

func (t Tree) Update(level int8) NodeInfo {
	current := t.levels[level].NodeInfo
	var next NodeInfo
	if t.play.deckSize > 0 {
		next = t.levels[level+1].NodeInfo
	} else {
		next = t.TerminalValue()
	}
	if current.type_ == MinType {
		if next.value < current.value {
			current.miniMaxPaths = next.miniMaxPaths
			current.value = next.value
			current.bestMove = t.play.LastMove()
		} else if next.value == current.value {
			current.miniMaxPaths += next.miniMaxPaths
		}
	} else if current.type_ == MaxType {
		if next.value > current.value {
			current.miniMaxPaths = next.miniMaxPaths
			current.value = next.value
			current.bestMove = t.play.LastMove()
		} else if next.value == current.value {
			current.miniMaxPaths += next.miniMaxPaths
		}
	}
	current.totalPaths += next.totalPaths
	return current
}

func (t Tree) Eval(policy DensePolicy, level int8) {
	if t.play.deckSize == 0 {
		return
	}
	t.levels[level].moves = t.play.GetMoves(policy)
	for _, move := range t.levels[level].moves {
		// fmt.Println("Doing move", move)
		t.play.DoMove(move)
		// fmt.Println("Current trick:", t.play.lastTrick())
		if t.play.deckSize > 0 {
			t.levels[level+1].NodeInfo = ResetNodeInfo(t.levels[level].NodeInfo, t.play)
			// fmt.Printf("Level %d reset to %v\n", level+1, t.levels[level+1].NodeInfo)
			t.Eval(policy, level+1)
		}
		t.levels[level].NodeInfo = t.Update(level)
		// fmt.Println("Undoing move", move)
		t.play.UndoMove()
	}
}
