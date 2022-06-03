package deal

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func ExampleDealPlay() {
	handStrings := []string{"♠A109", "♠QJ♣J", "♠K8♣Q"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	dp := NewDealPlay(&dealMatrix, Diamonds, ThirdHand, ThirdHandContract)

	moves := dp.GetMoves(Max)
	fmt.Println(moves)
	dp.DoMove(moves[0])

	moves = dp.GetMoves(Max)
	fmt.Println(moves)
	dp.DoMove(moves[0])

	moves = dp.GetMoves(Max)
	fmt.Println(moves)
	dp.DoMove(moves[0])

	fmt.Printf("Trick: %s\n", dp.tricks[len(dp.tricks)-1].String())
	fmt.Println("Deal:")
	fmt.Println(dealMatrix.String())
	fmt.Println("Squeezed deal:")
	squeezedDeal := dealMatrix.squeeze(little)
	fmt.Println(squeezedDeal.String())
	// Output:
	// [♠8 ♠K ♣Q]
	// [♠T ♠A]
	// [♠Q]
	// Trick: ♠8 ♠T ♠Q
	// Deal:
	// ♠A9
	// ♠J ♣J
	// ♠K ♣Q
	// Squeezed deal:
	// ♠T7
	// ♠8 ♣7
	// ♠9 ♣8
}

func ExampleDealPlay_index() {
	handStrings := []string{"♠A109", "♠QJ♣J", "♠K8♣Q"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	dp := NewDealPlay(&dealMatrix, Diamonds, FirstHand, ThirdHandContract)
	fmt.Println("Index little:", dp.Index(little))
	fmt.Println("Index big:", dp.Index(big))

	dpFromIndexLittle := NewDealPlayFromIndex(dp.Index(little), little, ThirdHandContract)
	fmt.Println("Deal play from index little:")
	fmt.Println(dpFromIndexLittle.matrix)

	dpFromIndexBig := NewDealPlayFromIndex(dp.Index(big), big, ThirdHandContract)
	fmt.Println("Deal play from index little:")
	fmt.Println(dpFromIndexBig.matrix)
	// Output:
	// Index little: 559939863
	// Index big: 408944919
	// Deal play from index little:
	// ♠K98
	// ♠JT ♣7
	// ♠Q7 ♣8
	// Deal play from index little:
	// ♠AT9
	// ♠QJ ♣K
	// ♠K8 ♣A
}

func TestDealMatrix_from10to1(t *testing.T) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	dp := NewDealPlay(&dealMatrix, Diamonds, ThirdHand, ThirdHandContract)
	rand.Seed(time.Now().Unix())
	require.EqualValues(t, dp.matrix.DeckSize(), 30, "Deck size must equal 30")
	for stage := 10; stage > 0; stage-- {
		moves := dp.GetMoves(Max)
		firstMove := moves[rand.Intn(len(moves))]
		dp.DoMove(firstMove)
		require.EqualValues(t, dp.matrix.DeckSize(), NumberOfHands*stage-1,
			fmt.Sprintf("Deck must contain %d cards", NumberOfHands*stage-1))
		require.EqualValues(t, len(dp.tricks[len(dp.tricks)-1].moves), 1,
			fmt.Sprintf("Last trick must contain %d card", 1))

		moves = dp.GetMoves(Max)
		secondMove := moves[rand.Intn(len(moves))]
		dp.DoMove(secondMove)
		require.EqualValues(t, dp.matrix.DeckSize(), NumberOfHands*stage-2,
			fmt.Sprintf("Deck must contain %d cards", NumberOfHands*stage-2))
		require.EqualValues(t, len(dp.tricks[len(dp.tricks)-1].moves), 2,
			fmt.Sprintf("Last trick must contain %d cards", 2))

		moves = dp.GetMoves(Max)
		thirdMove := moves[rand.Intn(len(moves))]
		dp.DoMove(thirdMove)
		require.EqualValues(t, dp.matrix.DeckSize(), NumberOfHands*stage-3,
			fmt.Sprintf("Deck must contain %d cards", NumberOfHands*stage-3))
		require.EqualValues(t, len(dp.tricks[len(dp.tricks)-1].moves), 3,
			fmt.Sprintf("Last trick must contain %d cards", 3))
	}
}

func TestDealMatrix_downUp(t *testing.T) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	dp := NewDealPlay(&dealMatrix, Diamonds, ThirdHand, ThirdHandContract)
	rand.Seed(29820)
	require.EqualValues(t, dp.matrix.DeckSize(), 30, "Deck size must equal 30")
	var matrixHistory [10]DealMatrix
	for stage := 10; stage > 0; stage-- {
		matrixHistory[stage-1] = *dp.matrix
		moves := dp.GetMoves(Max)
		firstMove := moves[rand.Intn(len(moves))]
		dp.DoMove(firstMove)
		require.EqualValues(t, dp.matrix.DeckSize(), NumberOfHands*stage-1,
			fmt.Sprintf("Deck must contain %d cards", NumberOfHands*stage-1))
		require.EqualValues(t, len(dp.tricks[len(dp.tricks)-1].moves), 1,
			fmt.Sprintf("Last trick must contain %d card", 1))

		moves = dp.GetMoves(Max)
		secondMove := moves[rand.Intn(len(moves))]
		dp.DoMove(secondMove)
		require.EqualValues(t, dp.matrix.DeckSize(), NumberOfHands*stage-2,
			fmt.Sprintf("Deck must contain %d cards", NumberOfHands*stage-2))
		require.EqualValues(t, len(dp.tricks[len(dp.tricks)-1].moves), 2,
			fmt.Sprintf("Last trick must contain %d cards", 2))

		moves = dp.GetMoves(Max)
		thirdMove := moves[rand.Intn(len(moves))]
		dp.DoMove(thirdMove)
		require.EqualValues(t, dp.matrix.DeckSize(), NumberOfHands*stage-3,
			fmt.Sprintf("Deck must contain %d cards", NumberOfHands*stage-3))
		require.EqualValues(t, len(dp.tricks[len(dp.tricks)-1].moves), 3,
			fmt.Sprintf("Last trick must contain %d cards", 3))
	}
	for stage := 1; stage < 11; stage++ {
		dp.UndoMove()
		require.EqualValues(t, dp.matrix.DeckSize(), NumberOfHands*stage-2,
			fmt.Sprintf("Deck must contain %d cards", NumberOfHands*stage-2))
		require.EqualValues(t, len(dp.tricks[len(dp.tricks)-1].moves), 2,
			fmt.Sprintf("Last trick must contain %d card", 2))

		dp.UndoMove()
		require.EqualValues(t, dp.matrix.DeckSize(), NumberOfHands*stage-1,
			fmt.Sprintf("Deck must contain %d cards", NumberOfHands*stage-1))
		require.EqualValues(t, len(dp.tricks[len(dp.tricks)-1].moves), 1,
			fmt.Sprintf("Last trick must contain %d cards", 1))

		dp.UndoMove()
		require.EqualValues(t, dp.matrix.DeckSize(), NumberOfHands*stage,
			fmt.Sprintf("Deck must contain %d cards", NumberOfHands*stage))
		require.EqualValues(t, len(dp.tricks[len(dp.tricks)-1].moves), 0,
			fmt.Sprintf("Last trick must contain %d cards", 0))

		require.EqualValues(t, matrixHistory[stage-1], *dp.matrix, "Different matrices: expected\n"+
			matrixHistory[stage-1].String()+"\nactual:\n"+dp.matrix.String())
	}
	require.EqualValues(t, [NumberOfHands]int8{0, 0, 0}, dp.result,
		fmt.Sprintf("Expected zeros, got: %v", dp.result))
	dealMatrixDup, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	require.EqualValues(t, dealMatrixDup, *dp.matrix, "Matrix has changed\n"+dp.matrix.String())
}

func TestDealMatrix_randDepth(t *testing.T) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	dp := NewDealPlay(&dealMatrix, Diamonds, ThirdHand, ThirdHandContract)
	rand.Seed(time.Now().Unix())
	require.EqualValues(t, dp.matrix.DeckSize(), 30, "Deck size must equal 30")
	const n = 1000
	for i := 0; i < n; i++ {
		depth := rand.Intn(31)
		for j := 0; j < depth; j++ {
			moves := dp.GetMoves(Max)
			move := moves[rand.Intn(len(moves))]
			dp.DoMove(move)
		}
		for j := 0; j < depth; j++ {
			dp.UndoMove()
		}
		dealMatrixDup, err := DealMatrixFromStrings(handStrings, NullDelimiter)
		if err != nil {
			log.Fatal(err)
		}
		require.EqualValues(t, dealMatrixDup, *dp.matrix,
			fmt.Sprintf("i=%d, depth=%d, matrix has changed\n%s", i, depth, dp.matrix.String()))
	}
}

func TestDealMatrix_result(t *testing.T) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	rand.Seed(time.Now().Unix())
	const n = 10000
	results := make(map[int8]int, n)
	for i := 0; i < n; i++ {
		dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
		if err != nil {
			log.Fatal(err)
		}
		dp := NewDealPlay(&dealMatrix, Diamonds, ThirdHand, ThirdHandContract)
		require.EqualValues(t, dp.matrix.DeckSize(), 30, "Deck size must equal 30")
		depth := 30
		for j := 0; j < depth; j++ {
			moves := dp.GetMoves(Max)
			move := moves[rand.Intn(len(moves))]
			dp.DoMove(move)
		}
		require.EqualValues(t, dp.matrix.DeckSize(), 0, "Deck must be empty")
		results[dp.result[ThirdHand]]++
	}
	/*
		var k int8
		for k = 10; k >= 0; k-- {
			if v, ok := results[k]; ok {
				fmt.Println(k, v)
			}
		} */
}

func downUp(dp *DealPlay) {
	for stage := 10; stage > 0; stage-- {
		moves := dp.GetMoves(Max)
		firstMove := moves[rand.Intn(len(moves))]
		dp.DoMove(firstMove)
		moves = dp.GetMoves(Max)
		secondMove := moves[rand.Intn(len(moves))]
		dp.DoMove(secondMove)
		moves = dp.GetMoves(Max)
		thirdMove := moves[rand.Intn(len(moves))]
		dp.DoMove(thirdMove)
	}
	for stage := 1; stage < 11; stage++ {
		dp.UndoMove()
		dp.UndoMove()
		dp.UndoMove()
	}
}

func BenchmarkDealMatrix_full_depth(b *testing.B) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	dp := NewDealPlay(&dealMatrix, Diamonds, ThirdHand, ThirdHandContract)
	rand.Seed(time.Now().Unix())
	for i := 0; i < b.N; i++ {
		downUp(dp)
	}
}

func BenchmarkDoUndoMove(b *testing.B) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	dp := NewDealPlay(&dealMatrix, Diamonds, ThirdHand, ThirdHandContract)
	moves := dealMatrix.GetMoves(Suit(rand.Intn(NumberOfSuits)),
		HandIndex(rand.Intn(NumberOfHands)), Random, false)
	for i := 0; i < b.N; i++ {
		dp.DoMove(moves[rand.Intn(len(moves))])
	}
}
