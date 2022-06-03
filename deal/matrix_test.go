package deal

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
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
	fmt.Printf("Index little = %d, index big = %d\n", dealMatrix.Index(little), dealMatrix.Index(big))
	fmt.Println("Deal from index little:")
	dealFromIndex := DealMatrixFromIndex(dealMatrix.ContingencyTable(), 4, little)
	fmt.Println(dealFromIndex.String())
	fmt.Println("Deal from index big:")
	dealFromIndex = DealMatrixFromIndex(dealMatrix.ContingencyTable(), 5, big)
	fmt.Println(dealFromIndex.String())
	// Output:
	// ♠K ♦J
	// ♦Q7
	// ♠T ♣A
	// Suit sizes: [2 1 3 0]
	// Contingency table: [[1 0 1] [0 0 1] [1 2 0] [0 0 0]]
	// Index little = 4, index big = 5
	// Deal from index little:
	// ♠8 ♦8
	// ♦97
	// ♠7 ♣7
	// Deal from index big:
	// ♠A ♦K
	// ♦AQ
	// ♠K ♣A
}

func ExampleDealMatrix_squeeze_1() {
	handStrings := []string{"K..J", "..Q7", "10.A"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, '.')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Initial deal:")
	fmt.Println(dealMatrix.String())

	squeezedDeal := dealMatrix.squeeze(little)
	fmt.Println("Squeezed deal (little):")
	fmt.Println(squeezedDeal.String())
	squeezedDeal = dealMatrix.squeeze(big)
	fmt.Println("Squeezed deal (big):")
	fmt.Println(squeezedDeal.String())
	// Output:
	// Initial deal:
	// ♠K ♦J
	// ♦Q7
	// ♠T ♣A
	// Squeezed deal (little):
	// ♠8 ♦8
	// ♦97
	// ♠7 ♣7
	// Squeezed deal (big):
	// ♠A ♦K
	// ♦AQ
	// ♠K ♣A
}

func ExampleDealMatrix_squeeze_2() {
	handStrings := []string{"♠A109♣7♦QJ♥KJ87", "♠QJ♣KJ98♦9♥AQ10", "♠K8♣AQ♦AK1087♥9"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, NullDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Initial deal:")
	fmt.Println(dealMatrix.String())
	fmt.Printf("Index little = %d, index big = %d\n", dealMatrix.Index(little), dealMatrix.Index(big))

	squeezedDeal := dealMatrix.squeeze(little)
	fmt.Println("Squeezed deal (little):")
	fmt.Println(squeezedDeal.String())
	squeezedDeal = dealMatrix.squeeze(big)
	fmt.Println("Squeezed deal (big):")
	fmt.Println(squeezedDeal.String())

	fmt.Println("Deal from index little:")
	dealFromIndex := DealMatrixFromIndex(dealMatrix.ContingencyTable(), 744508437, little)
	fmt.Println(dealFromIndex.String())

	fmt.Println("Deal from index big:")
	dealFromIndex = DealMatrixFromIndex(dealMatrix.ContingencyTable(), 206763465, big)
	fmt.Println(dealFromIndex.String())
	// Output:
	// Initial deal:
	// ♠AT9 ♣7 ♦QJ ♥KJ87
	// ♠QJ ♣KJ98 ♦9 ♥AQT
	// ♠K8 ♣AQ ♦AKT87 ♥9
	// Index little = 744508437, index big = 206763465
	// Squeezed deal (little):
	// ♠K98 ♣7 ♦QJ ♥KJ87
	// ♠JT ♣QT98 ♦9 ♥AQT
	// ♠Q7 ♣KJ ♦AKT87 ♥9
	// Squeezed deal (big):
	// ♠AT9 ♣8 ♦QJ ♥KJ87
	// ♠QJ ♣KJT9 ♦9 ♥AQT
	// ♠K8 ♣AQ ♦AKT87 ♥9
	// Deal from index little:
	// ♠K98 ♣7 ♦QJ ♥KJ87
	// ♠JT ♣QT98 ♦9 ♥AQT
	// ♠Q7 ♣KJ ♦AKT87 ♥9
	// Deal from index big:
	// ♠AT9 ♣8 ♦QJ ♥KJ87
	// ♠QJ ♣KJT9 ♦9 ♥AQT
	// ♠K8 ♣AQ ♦AKT87 ♥9
}

func ExampleDealMatrix_squeeze_3() {
	handStrings := []string{"KJ.Q8.AQ", "Q9.AKJ.J", "AT.9.T8.J"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, '.')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Initial deal:")
	fmt.Println(dealMatrix.String())
	fmt.Printf("Index little = %d, index big = %d\n", dealMatrix.Index(little), dealMatrix.Index(big))

	squeezedDeal := dealMatrix.squeeze(little)
	fmt.Println("Squeezed deal (little):")
	fmt.Println(squeezedDeal.String())
	squeezedDeal = dealMatrix.squeeze(big)
	fmt.Println("Squeezed deal (big):")
	fmt.Println(squeezedDeal.String())

	fmt.Println("Deal from index little:")
	dealFromIndex := DealMatrixFromIndex(dealMatrix.ContingencyTable(), 4360, little)
	fmt.Println(dealFromIndex.String())

	fmt.Println("Deal from index big:")
	dealFromIndex = DealMatrixFromIndex(dealMatrix.ContingencyTable(), 103310, big)
	fmt.Println(dealFromIndex.String())
	// Output:
	// Initial deal:
	// ♠KJ ♣Q8 ♦AQ
	// ♠Q9 ♣AKJ ♦J
	// ♠AT ♣9 ♦T8 ♥J
	// Index little = 4360, index big = 103310
	// Squeezed deal (little):
	// ♠J9 ♣T7 ♦JT
	// ♠T7 ♣QJ9 ♦9
	// ♠Q8 ♣8 ♦87 ♥7
	// Squeezed deal (big):
	// ♠KJ ♣Q9 ♦AK
	// ♠Q9 ♣AKJ ♦Q
	// ♠AT ♣T ♦JT ♥A
	// Deal from index little:
	// ♠J9 ♣T7 ♦JT
	// ♠T7 ♣QJ9 ♦9
	// ♠Q8 ♣8 ♦87 ♥7
	// Deal from index big:
	// ♠KJ ♣Q9 ♦AK
	// ♠Q9 ♣AKJ ♦Q
	// ♠AT ♣T ♦JT ♥A
}

func ExampleDealMatrix_squeeze_4() {
	handStrings := []string{"K9.A", "7.Q7", "J.9.K"}
	dealMatrix, err := DealMatrixFromStrings(handStrings, '.')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Initial deal:")
	fmt.Println(dealMatrix.String())
	fmt.Printf("Index little = %d, index big = %d\n", dealMatrix.Index(little), dealMatrix.Index(big))

	squeezedDeal := dealMatrix.squeeze(little)
	fmt.Println("Squeezed deal (little):")
	fmt.Println(squeezedDeal.String())
	squeezedDeal = dealMatrix.squeeze(big)
	fmt.Println("Squeezed deal (big):")
	fmt.Println(squeezedDeal.String())

	fmt.Println("Deal from index little:")
	dealFromIndex := DealMatrixFromIndex(dealMatrix.ContingencyTable(), 17, little)
	fmt.Println(dealFromIndex.String())

	fmt.Println("Deal from index big:")
	dealFromIndex = DealMatrixFromIndex(dealMatrix.ContingencyTable(), 92, big)
	fmt.Println(dealFromIndex.String())
	// Output:
	// Initial deal:
	// ♠K9 ♣A
	// ♠7 ♣Q7
	// ♠J ♣9 ♦K
	// Index little = 17, index big = 92
	// Squeezed deal (little):
	// ♠T8 ♣T
	// ♠7 ♣97
	// ♠9 ♣8 ♦7
	// Squeezed deal (big):
	// ♠AQ ♣A
	// ♠J ♣KJ
	// ♠K ♣Q ♦A
	// Deal from index little:
	// ♠T8 ♣T
	// ♠7 ♣97
	// ♠9 ♣8 ♦7
	// Deal from index big:
	// ♠AQ ♣A
	// ♠J ♣KJ
	// ♠K ♣Q ♦A
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

// func TestShiftZeros(t *testing.T) {
// 	var codes = map[[3]uint8]uint8{
// 		{9, 1, 2}:   3,
// 		{1, 1, 0}:   1,
// 		{2, 0, 1}:   1,
// 		{255, 8, 0}: 255,
// 		{128, 4, 2}: 32,
// 		{18, 2, 2}:  6,
// 		{85, 5, 1}:  53,
// 		{85, 3, 1}:  45,
// 	}
// 	for params, wanted := range codes {
// 		require.EqualValues(t, wanted, squeezeZeros(params[0], params[1], params[2]))
// 	}
// }

func TestMatrixIndex(t *testing.T) {
	table := [NumberOfSuits][NumberOfHands]int8{{0, 2, 1}, {2, 2, 2}, {3, 1, 2}}
	var variants int64 = 1
	for i := 0; i < len(table); i++ {
		variants *= int64(multinomial(table[i]))
	}
	var index int64
	for index = 0; index < variants; index++ {
		dealMatrix := DealMatrixFromIndex(table, index, little)
		indexFromMatrix := dealMatrix.Index(little)
		require.EqualValues(t, index, indexFromMatrix,
			fmt.Sprintf("matrix index must equal %d, got %d", index, indexFromMatrix))
	}
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

func BenchmarkSqueezeMatrix(b *testing.B) {
	handStrings := []string{"K9.A", "7.Q7", "J.9.K"}
	dealMatrix, _ := DealMatrixFromStrings(handStrings, '.')
	for i := 0; i < b.N; i++ {
		dealMatrix.squeeze(little)
	}
}
