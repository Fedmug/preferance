package deal

import "fmt"

func ExampleCard() {
	var sevenSpades = NewCard(Spades, Seven, false)
	var aceOfHeartsTrump = NewCard(Hearts, Ace, true)
	var kingOfDiamonds = NewCard(Diamonds, King, false)
	fmt.Println(sevenSpades, sevenSpades.code)
	fmt.Println(aceOfHeartsTrump, aceOfHeartsTrump.code)
	fmt.Println(kingOfDiamonds, kingOfDiamonds.code)
	fmt.Println("♡A beats ♠7:", aceOfHeartsTrump.Beats(sevenSpades))
	fmt.Println("♦K beats ♠7:", kingOfDiamonds.Beats(sevenSpades))
	fmt.Println("♦K beats ♦J:", kingOfDiamonds.Beats(NewCard(Diamonds, Jack, false)))
	fmt.Println("♧8 beats ♦K:", NewCard(Clubs, Eight, true).Beats(kingOfDiamonds))
	// Output:
	// ♠7 0
	// ♡A 31
	// ♦K 22
	// ♡A beats ♠7: true
	// ♦K beats ♠7: false
	// ♦K beats ♦J: true
	// ♧8 beats ♦K: true
}
