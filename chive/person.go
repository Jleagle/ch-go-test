package chive

import "time"

type person struct {
	Name             string
	DOB              time.Time
	Siblings         int
	HasPets          bool
	FraudScore       float32
	FavouritePlanets []string
	Transactions     map[string]string
	AmountPaid       amount
	Payments         map[string]amount
}

type amount struct {
	Currency string
	Units    int64
}
