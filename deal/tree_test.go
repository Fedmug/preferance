package deal

import (
	"fmt"
	"log"
	"time"
)

func evalDealFromStrings(handStrings []string, delimiter rune, trump Suit, firstMover HandIndex) {
	dealMatrix, err := DealMatrixFromStrings(handStrings, delimiter)
	if err != nil {
		log.Fatal(err)
	}
	dp := NewDealPlay(&dealMatrix, trump, firstMover)
	root := NewTree(dp)
	root.Eval(Min, 0)
	for k, v := range root.levels[0].NodeInfo.(TrickDistributionCounter).trickMap {
		fmt.Println(k, v)
	}
}

func ExampleEval_1() {
	handStrings := []string{"K", "Q", "A"}
	evalDealFromStrings(handStrings, '.', Clubs, SecondHand)
	handStrings = []string{"K", "A", "Q"}
	evalDealFromStrings(handStrings, '.', NT, FirstHand)
	handStrings = []string{"K", "8", "J"}
	evalDealFromStrings(handStrings, '.', Spades, ThirdHand)
	// Output:
	// [0 0 1] 1
	// [0 1 0] 1
	// [1 0 0] 1
}

func ExampleEval_2() {
	handStrings := []string{"K.9", "AQ", "T.Q"}
	fmt.Println("First hand move:")
	evalDealFromStrings(handStrings, '.', NT, FirstHand)
	fmt.Println("Second hand move:")
	evalDealFromStrings(handStrings, '.', NT, SecondHand)
	fmt.Println("Third hand move:")
	evalDealFromStrings(handStrings, '.', NT, ThirdHand)
	// Output:
	// First hand move:
	// [0 2 0] 1
	// [1 0 1] 2
	// [0 1 1] 1
	// Second hand move:
	// [1 0 1] 1
	// [0 2 0] 1
	// Third hand move:
	// [1 0 1] 2
	// [0 2 0] 1
	// [0 1 1] 1
}

func ExampleEval_time() {
	begin := time.Now()
	handStrings := []string{"♠A9♦QJ♥J7", "♠QJ♣K8♦9♥10", "♣AQ♦A1087"}
	evalDealFromStrings(handStrings, NullDelimiter, Diamonds, ThirdHand)
	fmt.Println(time.Since(begin))
	// Output:
	// long time...
}
