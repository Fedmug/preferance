package deal

import (
	"fmt"
	"log"
)

func evalDealFromStrings(handStrings []string, delimiter rune,
	trump Suit, firstMover HandIndex, goals DealGoals) Tree {
	dealMatrix, err := DealMatrixFromStrings(handStrings, delimiter)
	if err != nil {
		log.Fatal(err)
	}
	dp := NewDealPlay(&dealMatrix, trump, firstMover, goals)
	// root := NewTree(dp, NodeCounter(0))
	//root := NewTree(dp, TrickDistributionCounter{})
	value := MiniMaxInf
	if goals[firstMover] == MaxType {
		value = -value
	}
	root := NewTree(dp, NodeInfo{goals[firstMover], 0, 0, value, InvalidCard})
	root.Eval(Max, 0)
	fmt.Println("Node type:", root.levels[0].NodeInfo.type_)
	fmt.Println("Total paths:", root.levels[0].NodeInfo.totalPaths)
	fmt.Println("MiniMax paths:", root.levels[0].NodeInfo.miniMaxPaths)
	fmt.Println("Node value:", root.levels[0].NodeInfo.value)
	fmt.Println("Best move:", root.levels[0].NodeInfo.bestMove)
	return root
	/*
		for k, v := range root.levels[0].NodeInfo.(TrickDistributionCounter).trickMap {
			fmt.Println(k, v)
		}*/
}

func ExampleEval_1() {
	handStrings := []string{"K", "Q", "A"}
	fmt.Println("First hand plays misere, second hand to move:")
	evalDealFromStrings(handStrings, '.', NT, SecondHand, FirstHandMisere)
	handStrings = []string{"K", "A", "Q"}
	fmt.Println("Second hand plays contract, first hand to move:")
	evalDealFromStrings(handStrings, '.', Spades, FirstHand, SecondHandContract)
	handStrings = []string{"K", "8", "J"}
	fmt.Println("Third hand plays misere, third hand to move:")
	evalDealFromStrings(handStrings, '.', NT, ThirdHand, ThirdHandMisere)
	// Output:
	// First hand plays misere, second hand to move:
	// Node type: 1
	// Total paths: 1
	// MiniMax paths: 1
	// Node value: 0
	// Best move: ♠Q
	// Second hand plays contract, first hand to move:
	// Node type: 0
	// Total paths: 1
	// MiniMax paths: 1
	// Node value: 1
	// Best move: ♤K
	// Third hand plays misere, third hand to move:
	// Node type: 0
	// Total paths: 1
	// MiniMax paths: 1
	// Node value: 0
	// Best move: ♠J
}

func ExampleEval_2() {
	handStrings := []string{"K.9", "AQ", "T.Q"}
	fmt.Println("Third hand plays contract, first hand to move:")
	evalDealFromStrings(handStrings, '.', NT, FirstHand, ThirdHandContract)
	fmt.Println("Second hand plays contract, second hand to move:")
	evalDealFromStrings(handStrings, '.', NT, SecondHand, SecondHandContract)
	fmt.Println("First hand plays misere, third hand to move:")
	evalDealFromStrings(handStrings, '.', NT, ThirdHand, FirstHandMisere)
	// Output:
	// Third hand plays contract, first hand to move:
	// Node type: 0
	// Total paths: 4
	// MiniMax paths: 1
	// Node value: 0
	// Best move: ♠K
	// Second hand plays contract, second hand to move:
	// Node type: 1
	// Total paths: 2
	// MiniMax paths: 1
	// Node value: 2
	// Best move: ♠A
	// First hand plays misere, third hand to move:
	// Node type: 1
	// Total paths: 4
	// MiniMax paths: 2
	// Node value: 1
	// Best move: ♠T
}

func ExampleEval_3() {
	handStrings := []string{"KJ", "Q9", "A8"}
	fmt.Println("First hand plays contract, first hand to move:")
	evalDealFromStrings(handStrings, '.', NT, FirstHand, FirstHandContract)
	fmt.Println("First hand plays contract, second hand to move:")
	evalDealFromStrings(handStrings, '.', NT, SecondHand, FirstHandContract)
	fmt.Println("First hand plays contract, third hand to move:")
	evalDealFromStrings(handStrings, '.', NT, ThirdHand, FirstHandContract)
	// Output:
	// First hand plays contract, first hand to move:
	// Node type: 1
	// Total paths: 8
	// MiniMax paths: 2
	// Node value: 0
	// Best move: ♠J
	// First hand plays contract, second hand to move:
	// Node type: 0
	// Total paths: 8
	// MiniMax paths: 6
	// Node value: 1
	// Best move: ♠9
	// First hand plays contract, third hand to move:
	// Node type: 0
	// Total paths: 8
	// MiniMax paths: 4
	// Node value: 1
	// Best move: ♠8
}

func ExampleEval_4() {
	handStrings := []string{"KJ.Q8.AQ", "Q9.AKJ.J", "AT.9.K98"}
	fmt.Println("First hand plays contract, first hand to move:")
	evalDealFromStrings(handStrings, '.', Diamonds, FirstHand, ThirdHandContract)
	// Output:
	// First hand plays contract, first hand to move:
	// Node type: 0
	// Total paths: 37288
	// MiniMax paths: 8598
	// Node value: 2
	// Best move: ♠J
}

func ExampleEval_5() {
	handStrings := []string{"♠109♣7♦QJT♥A", "♠QJ♣QJ109♦9", "♠AK♣A♦AK87"}
	// fmt.Println("First hand:", evalDealFromStrings(handStrings, NullDelimiter, Diamonds, FirstHand))
	// fmt.Println("Second hand:", evalDealFromStrings(handStrings, NullDelimiter, Diamonds, SecondHand))
	evalDealFromStrings(handStrings, NullDelimiter, Diamonds, FirstHand, ThirdHandContract)
	evalDealFromStrings(handStrings, NullDelimiter, Diamonds, SecondHand, ThirdHandContract)
	evalDealFromStrings(handStrings, NullDelimiter, Diamonds, ThirdHand, ThirdHandContract)
	// Output:
	// Node type: 0
	// Total paths: 5202
	// MiniMax paths: 885
	// Node value: 5
	// Best move: ♥A
	// Node type: 0
	// Total paths: 4317
	// MiniMax paths: 2061
	// Node value: 6
	// Best move: ♠Q
	// Node type: 1
	// Total paths: 4317
	// MiniMax paths: 2061
	// Node value: 6
	// Best move: ♠A
}

func ExampleMultiEval() {
	handStrings := []string{"♠109♣7♦QJT♥A", "♠QJ♣QJ109♦9", "♠AK♣A♦AK87"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	dp := NewDealPlay(&dealMatrix, Diamonds, FirstHand, ThirdHandContract)
	MultiEval(dp, Max)
	// Output:
	// Node type: 0
}
