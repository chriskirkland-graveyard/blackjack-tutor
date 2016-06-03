package main

import (
	"./blackjack"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Hello")

	// seed RNG
	rand.Seed(time.Now().UTC().UnixNano())

	// new Shoe
	s := blackjack.NewShoe(2)
	fmt.Println("Shoe:")
	fmt.Println(s)

	c := s.DealCard()
	fmt.Println("Dealt: ", c)
	c = s.DealCard()
	fmt.Println("Dealt: ", c)
	c = s.DealCard()
	fmt.Println("Dealt: ", c)
	fmt.Println("Shoe:")
	fmt.Println(s)

	s.Shuffle()
	fmt.Println("Shoe:")
	fmt.Println(s)
}
