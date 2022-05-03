package deal

import "fmt"

func ExampleTrick_nt002() {
	var trick = NewTrick()
	trick.Append(NewCard(Hearts, Jack, false))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	trick.Append(NewCard(Hearts, Nine, false))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	trick.Append(NewCard(Hearts, Queen, false))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	fmt.Println("Partial takers:", trick.partialTakers)
	// Output:
	// ♥J
	// Taker: ♥J
	// ♥J ♥9
	// Taker: ♥J
	// ♥J ♥9 ♥Q
	// Taker: ♥Q
	// Partial takers: [0 0 2]
}

func ExampleTrick_nt000() {
	var trick = NewTrick()
	trick.Append(NewCard(Spades, Ten, false))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	trick.Append(NewCard(Clubs, King, false))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	trick.Append(NewCard(Diamonds, Ace, false))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	fmt.Println("Partial takers:", trick.partialTakers)
	// Output:
	// ♠T
	// Taker: ♠T
	// ♠T ♣K
	// Taker: ♠T
	// ♠T ♣K ♦A
	// Taker: ♠T
	// Partial takers: [0 0 0]
}

func ExampleTrick_t011() {
	var trick = NewTrick()
	trick.Append(NewCard(Spades, Queen, false))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	trick.Append(NewCard(Clubs, Seven, true))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	trick.Append(NewCard(Spades, Ace, false))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	fmt.Println("Partial takers:", trick.partialTakers)
	// Output:
	// ♠Q
	// Taker: ♠Q
	// ♠Q ♧7
	// Taker: ♧7
	// ♠Q ♧7 ♠A
	// Taker: ♧7
	// Partial takers: [0 1 1]
}

func ExampleTrick_t012() {
	var trick = NewTrick()
	trick.Append(NewCard(Hearts, Jack, false))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	trick.Append(NewCard(Hearts, King, false))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	trick.Append(NewCard(Spades, Eight, true))
	fmt.Println(trick)
	fmt.Println("Taker:", trick.TakerCard())
	fmt.Println("Partial takers:", trick.partialTakers)
	// Output:
	// ♥J
	// Taker: ♥J
	// ♥J ♥K
	// Taker: ♥K
	// ♥J ♥K ♤8
	// Taker: ♤8
	// Partial takers: [0 1 2]
}
