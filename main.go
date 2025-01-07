package main

import (
	"dragon/profile"
	"dragon/wallet"
	"fmt"
)

func main() {
	fmt.Println("Starting the Game Project")

	// Example usage of wallet functionality
	w := wallet.NewVirtualWallet()
	err := w.Connect()
	if err != nil {
		fmt.Printf("Error connecting virtual wallet: %v\n", err)
	} else {
		fmt.Println("Virtual wallet connected successfully")
	}

	// Example usage of profile functionality
	p := profile.NewProfile("JohnDoe", 1)
	p.SetLevel(5)
	fmt.Printf("Player Profile: %v\n", p)
}
