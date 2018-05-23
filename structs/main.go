package main

import (
	"fmt"
)

type contactInfo struct {
	email   string
	zipCode int
}

type person struct {
	firstName string
	lastName  string
	contactInfo
}

func (pointerToPerson *person) update(newFirstName string) {
	(*pointerToPerson).firstName = newFirstName
}

func (p person) print() {
	fmt.Printf("%+v", p)
}

func main() {
	// alex := person{firstName: "Alex", lastName: "Anderson"}

	var alex person

	alex.firstName = "aaa"
	alex.lastName = "7777777"

	alex.update("Jim")
	alex.print()
}
