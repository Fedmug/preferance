package deal

import "fmt"

func ExampleTrick_nt002() {
	var trick = NewTrick(NumberOfHands)
	trick.append(NewCard(Hearts, Jack, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	trick.append(NewCard(Hearts, Nine, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	trick.append(NewCard(Hearts, Queen, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	// Output:
	// ♥J
	// Taker card: ♥J
	// Taker hands: 0 1 2
	// ♥J ♥9
	// Taker card: ♥J
	// Taker hands: 0 1 2
	// ♥J ♥9 ♥Q
	// Taker card: ♥Q
	// Taker hands: 2 0 1
}

func ExampleTrick_nt000() {
	var trick = NewTrick(NumberOfHands)
	trick.append(NewCard(Spades, Ten, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	trick.append(NewCard(Clubs, King, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	trick.append(NewCard(Diamonds, Ace, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	// Output:
	// ♠T
	// Taker card: ♠T
	// Taker hands: 0 1 2
	// ♠T ♣K
	// Taker card: ♠T
	// Taker hands: 0 1 2
	// ♠T ♣K ♦A
	// Taker card: ♠T
	// Taker hands: 0 1 2
}

func ExampleTrick_t011() {
	var trick = NewTrick(NumberOfHands)
	trick.append(NewCard(Spades, Queen, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	trick.append(NewCard(Clubs, Seven, true))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	trick.append(NewCard(Spades, Ace, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	// Output:
	// ♠Q
	// Taker card: ♠Q
	// Taker hands: 0 1 2
	// ♠Q ♧7
	// Taker card: ♧7
	// Taker hands: 1 2 0
	// ♠Q ♧7 ♠A
	// Taker card: ♧7
	// Taker hands: 1 2 0
}

func ExampleTrick_t012() {
	var trick = NewTrick(NumberOfHands)
	trick.append(NewCard(Hearts, Jack, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	trick.append(NewCard(Hearts, King, false))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	trick.append(NewCard(Spades, Eight, true))
	fmt.Println(trick)
	fmt.Println("Taker card:", trick.takerCard())
	fmt.Println("Taker hands:", trick.takerHandIndex(FirstHand),
		trick.takerHandIndex(SecondHand), trick.takerHandIndex(ThirdHand))
	// Output:
	// ♥J
	// Taker card: ♥J
	// Taker hands: 0 1 2
	// ♥J ♥K
	// Taker card: ♥K
	// Taker hands: 1 2 0
	// ♥J ♥K ♤8
	// Taker card: ♤8
	// Taker hands: 2 0 1
}

func ExampleTrick_pop012() {
	var trick = NewTrick(NumberOfHands)
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
