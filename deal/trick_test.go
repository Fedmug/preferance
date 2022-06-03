package deal

import "fmt"

func ExampleTrick_nt002() {
	var trick = NewTrick(ThirdHand, NumberOfHands)
	trick.append(NewCard(Hearts, Jack, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	trick.append(NewCard(Hearts, Nine, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	trick.append(NewCard(Hearts, Queen, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	// Output:
	// ♥J
	// Taker card: ♥J
	// Taker hand: 2
	// ♥J ♥9
	// Taker card: ♥J
	// Taker hand: 2
	// ♥J ♥9 ♥Q
	// Taker card: ♥Q
	// Taker hand: 1
}

func ExampleTrick_nt000() {
	var trick = NewTrick(FirstHand, NumberOfHands)
	trick.append(NewCard(Spades, Ten, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	trick.append(NewCard(Clubs, King, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	trick.append(NewCard(Diamonds, Ace, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	// Output:
	// ♠T
	// Taker card: ♠T
	// Taker hand: 0
	// ♠T ♣K
	// Taker card: ♠T
	// Taker hand: 0
	// ♠T ♣K ♦A
	// Taker card: ♠T
	// Taker hand: 0
}

func ExampleTrick_t011() {
	var trick = NewTrick(SecondHand, NumberOfHands)
	trick.append(NewCard(Spades, Queen, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	trick.append(NewCard(Clubs, Seven, true))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	trick.append(NewCard(Spades, Ace, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	// Output:
	// ♠Q
	// Taker card: ♠Q
	// Taker hand: 1
	// ♠Q ♧7
	// Taker card: ♧7
	// Taker hand: 2
	// ♠Q ♧7 ♠A
	// Taker card: ♧7
	// Taker hand: 2
}

func ExampleTrick_t012() {
	var trick = NewTrick(SecondHand, NumberOfHands)
	trick.append(NewCard(Hearts, Jack, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	trick.append(NewCard(Hearts, King, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	trick.append(NewCard(Spades, Eight, true))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hand:", trick.takerHandIndex())
	// Output:
	// ♥J
	// Taker card: ♥J
	// Taker hand: 1
	// ♥J ♥K
	// Taker card: ♥K
	// Taker hand: 2
	// ♥J ♥K ♤8
	// Taker card: ♤8
	// Taker hand: 0
}

func ExampleTrick_pop012() {
	var trick = NewTrick(SecondHand, NumberOfHands)
	trick.append(NewCard(Hearts, Jack, false))
	trick.append(NewCard(Hearts, King, false))
	trick.append(NewCard(Spades, Eight, true))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	card := trick.pop()
	fmt.Println("Popped card:", card)
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	card = trick.pop()
	fmt.Println("Popped card:", card)
	fmt.Println("Taker card:", trick.takerCard())
	card = trick.pop()
	fmt.Println("Popped card:", card)
	// Output:
	// ♥J ♥K ♤8
	// Taker card: ♤8
	// Popped card: ♤8
	// ♥J ♥K
	// Taker card: ♥K
	// Popped card: ♥K
	// Taker card: ♥J
	// Popped card: ♥J
}
