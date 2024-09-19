package chive

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestOutput(t *testing.T) {

	p := person{
		Name:             "John",
		DOB:              time.Date(2020, 05, 01, 1, 30, 24, 0, time.UTC),
		Siblings:         2,
		HasPets:          true,
		FraudScore:       0.8,
		FavouritePlanets: []string{"Earth", "Mars"},
		Transactions: map[string]string{
			"payment1": "XUDDHSO23i72",
			"payment2": "LJRHRJ9347JF",
		},
		AmountPaid: amount{
			Currency: "USD",
			Units:    100,
		},
		Payments: map[string]amount{
			"payment1": {Currency: "USD", Units: 50},
			"payment2": {Currency: "USD", Units: 50},
		},
	}

	// Expected Output
	expect := map[string]Value{
		"name":              {Text: "John"},
		"dob":               {Timestamp: 1588325424},
		"siblings":          {Number: 2},
		"has_pets":          {Bool: true},
		"fraud_score":       {Float: 0.8},
		"favourite_planets": {Array: &RepeatedValues{Strings: []string{"Earth", "Mars"}}},
		"transactions":      {Array: &RepeatedValues{KeyValue: map[string][]byte{"payment1": []byte("XUDDHSO23i72"), "payment2": []byte("LJRHRJ9347JF")}}},
		"payments":          {Array: &RepeatedValues{KeyValue: map[string][]byte{"payment1": []byte("currency:USD,units:50"), "payment2": []byte("currency:USD,units:50")}}},
		// Nested structs should be flattened
		"amount_paid.currency": {Text: "USD"},
		"amount_paid.units":    {Number: 100},
	}

	// Quick way - 15 mins
	result1 := convert1(p)
	if !reflect.DeepEqual(result1, expect) {
		t.Error("Expected and actual output do not match")
		log.Println(marshal(result1))
		log.Println(marshal(expect))
	}

	// Slow way - 4 hours
	result2 := convert2(p)
	if !reflect.DeepEqual(result2, expect) {
		t.Error("Expected and actual output do not match")
		log.Println(marshal(result2))
		log.Println(marshal(expect))
	}
}

func marshal(in any) string {
	b, err := json.Marshal(in)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
