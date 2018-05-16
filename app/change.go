package main

import "math"

// Change is coin representation of a certain amount of USD
type Change struct {
	Quarters int
	Dimes    int
	Nickles  int
	Pennies  int
}

func NewChange(amount float64) *Change {
	change := new(Change)
	if amount <= 0 {
		return change
	}

	change.Quarters = int(math.Floor(amount / 0.25))
	remainingAmount := amount - (float64(change.Quarters) * 0.25)

	change.Dimes = int(math.Floor(remainingAmount / 0.10))
	remainingAmount -= float64(change.Dimes) * 0.10

	change.Nickles = int(math.Floor(remainingAmount / 0.05))
	remainingAmount -= float64(change.Nickles) * 0.05

	change.Pennies = int(math.Floor(remainingAmount / 0.01))

	return change
}
