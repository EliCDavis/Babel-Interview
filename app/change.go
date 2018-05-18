package main

import (
	"fmt"
	"math"
)

// Change is coin the representation of a certain amount of USD
type Change struct {
	Quarters int
	Dimes    int
	Nickles  int
	Pennies  int
}

func (change Change) String() string {
	return fmt.Sprintf("%d Quater(s), %d Dime(s), %d Nickle(s), %d Penny(s)", change.Quarters, change.Dimes, change.Nickles, change.Pennies)
}

func (change *Change) Add(v Change) {
	change.Quarters += v.Quarters
	change.Dimes += v.Dimes
	change.Nickles += v.Nickles
	change.Pennies += v.Pennies
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

	// Round in the end due to errors in floating point arithmitic
	change.Pennies = int(math.Floor(toFixed(remainingAmount, 2) / 0.01))

	return change
}

// https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision-in-golang

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
