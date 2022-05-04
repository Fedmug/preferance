package deal

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func ExampleDealMatrix_full() {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dealMatrix.String())
	fmt.Println("Suit sizes:", dealMatrix.SuitSizes())
	fmt.Println("Contingency table:", dealMatrix.ContingencyTable())
	// Output:
	// ♠AT9 ♣7 ♦QJ ♥KJ87
	// ♠QJ ♣KJ98 ♦9 ♥AQT
	// ♠K8 ♣AQ ♦AKT87 ♥9
	// Suit sizes: [7 7 8 8]
	// Contingency table: [[3 2 2] [1 4 2] [2 1 5] [4 3 1]]
}

func ExampleDealMatrix_sparse() {
	handStrings := []string{"K..J", "..Q7", "10.A"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, '.')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dealMatrix.String())
	fmt.Println("Suit sizes:", dealMatrix.SuitSizes())
	fmt.Println("Contingency table:", dealMatrix.ContingencyTable())
	// Output:
	// ♠K ♦J
	// ♦Q7
	// ♠T ♣A
	// Suit sizes: [2 1 3 0]
	// Contingency table: [[1 0 1] [0 0 1] [1 2 0] [0 0 0]]
}

func ExampleDealMatrix_moves() {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dealMatrix.GetMoves(Spades, FirstHand, All, false))
	fmt.Println(dealMatrix.GetMoves(Clubs, SecondHand, All, false))
	fmt.Println(dealMatrix.GetMoves(Hearts, ThirdHand, All, true))
	fmt.Println(dealMatrix.GetMoves(Diamonds, ThirdHand, Min, false))
	fmt.Println(dealMatrix.GetMoves(Hearts, FirstHand, Max, false))
	// Output:
	// [♠9 ♠T ♠A]
	// [♣8 ♣9 ♣J ♣K]
	// [♡9]
	// [♦7 ♦T ♦K]
	// [♥8 ♥J ♥K]
}

func BenchmarkBuildMatrix(b *testing.B) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	for i := 0; i < b.N; i++ {
		DealMatrixFromStrings(handStrings, NullDelimiter)
	}
}

func BenchmarkMovesMax(b *testing.B) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, _ := DealMatrixFromStrings(handStrings, NullDelimiter)
	for i := 0; i < b.N; i++ {
		dealMatrix.GetMoves(Suit(rand.Intn(NumberOfSuits)),
			HandIndex(rand.Intn(NumberOfHands)), Max, false)
	}
}

func BenchmarkMovesAll(b *testing.B) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, _ := DealMatrixFromStrings(handStrings, NullDelimiter)
	for i := 0; i < b.N; i++ {
		dealMatrix.GetMoves(Suit(rand.Intn(NumberOfSuits)),
			HandIndex(rand.Intn(NumberOfHands)), All, false)
	}
}

func BenchmarkMovesRandom(b *testing.B) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, _ := DealMatrixFromStrings(handStrings, NullDelimiter)
	for i := 0; i < b.N; i++ {
		dealMatrix.GetMoves(Suit(rand.Intn(NumberOfSuits)),
			HandIndex(rand.Intn(NumberOfHands)), Random, false)
	}
}

func BenchmarkSuitSizes(b *testing.B) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, _ := DealMatrixFromStrings(handStrings, NullDelimiter)
	for i := 0; i < b.N; i++ {
		dealMatrix.SuitSizes()
	}
}

func BenchmarkContingencyTable(b *testing.B) {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, _ := DealMatrixFromStrings(handStrings, NullDelimiter)
	for i := 0; i < b.N; i++ {
		dealMatrix.ContingencyTable()
	}
}
