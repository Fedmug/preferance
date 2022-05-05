package deal

import "fmt"

func ExampleTrick_nt002() {
	var trick = NewTrick(ThirdHand, NumberOfHands)
	trick.Append(NewCard(Hearts, Jack, false))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	trick.Append(NewCard(Hearts, Nine, false))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	trick.Append(NewCard(Hearts, Queen, false))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	// Output:
	// ♥J
	// Taker card&hand: ♥J 2
	// ♥J ♥9
	// Taker card&hand: ♥9 2
	// ♥J ♥9 ♥Q
	// Taker card&hand: ♥Q 1
}

func ExampleTrick_nt000() {
	var trick = NewTrick(SecondHand, NumberOfHands)
	trick.Append(NewCard(Spades, Ten, false))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	trick.Append(NewCard(Clubs, King, false))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	trick.Append(NewCard(Diamonds, Ace, false))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	// Output:
	// ♠T
	// Taker card&hand: ♠T 1
	// ♠T ♣K
	// Taker card&hand: ♣K 1
	// ♠T ♣K ♦A
	// Taker card&hand: ♦A 1
}

func ExampleTrick_t011() {
	var trick = NewTrick(ThirdHand, NumberOfHands)
	trick.Append(NewCard(Spades, Queen, false))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	trick.Append(NewCard(Clubs, Seven, true))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	trick.Append(NewCard(Spades, Ace, false))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	// Output:
	// ♠Q
	// Taker card&hand: ♠Q 2
	// ♠Q ♧7
	// Taker card&hand: ♧7 0
	// ♠Q ♧7 ♠A
	// Taker card&hand: ♠A 0
}

func ExampleTrick_t012() {
	var trick = NewTrick(FirstHand, NumberOfHands)
	trick.Append(NewCard(Hearts, Jack, false))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	trick.Append(NewCard(Hearts, King, false))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	trick.Append(NewCard(Spades, Eight, true))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	// Output:
	// ♥J
	// Taker card&hand: ♥J 0
	// ♥J ♥K
	// Taker card&hand: ♥K 1
	// ♥J ♥K ♤8
	// Taker card&hand: ♤8 2
}

func ExampleTrick_pop012() {
	var trick = NewTrick(SecondHand, NumberOfHands)
	trick.Append(NewCard(Hearts, Jack, false))
	trick.Append(NewCard(Hearts, King, false))
	trick.Append(NewCard(Spades, Eight, true))
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	card := trick.Pop()
	fmt.Println("Popped card:", card)
	fmt.Println(trick)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	card = trick.Pop()
	fmt.Println("Popped card:", card)
	fmt.Println("Taker card&hand:", trick.TakerMove().card, trick.TakerMove().taker)
	card = trick.Pop()
	fmt.Println("Popped card:", card)
	// Output:
	// ♥J ♥K ♤8
	// Taker card&hand: ♤8 0
	// Popped card: ♤8
	// ♥J ♥K
	// Taker card&hand: ♥K 2
	// Popped card: ♥K
	// Taker card&hand: ♥J 1
	// Popped card: ♥J
}
