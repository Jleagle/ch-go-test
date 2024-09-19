package chive

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func convert1(p person) map[string]Value {

	ret := map[string]Value{
		"name":                 {Text: p.Name},
		"dob":                  {Timestamp: p.DOB.Add(time.Hour * 8).Unix()},
		"siblings":             {Number: int64(p.Siblings)},
		"has_pets":             {Bool: p.HasPets},
		"fraud_score":          {Float: round(float64(p.FraudScore))},
		"favourite_planets":    {Array: &RepeatedValues{Strings: p.FavouritePlanets}},
		"transactions":         {Array: &RepeatedValues{KeyValue: map[string][]byte{}}},
		"payments":             {Array: &RepeatedValues{KeyValue: map[string][]byte{}}},
		"amount_paid.currency": {Text: p.AmountPaid.Currency},
		"amount_paid.units":    {Number: p.AmountPaid.Units},
	}

	for k, v := range p.Transactions {
		ret["transactions"].Array.KeyValue[k] = []byte(v)
	}
	for k, v := range p.Payments {
		ret["payments"].Array.KeyValue[k] = []byte(fmt.Sprintf("currency:%s,units:%d", v.Currency, v.Units))
	}

	return ret
}

func convert2(p person) map[string]Value {

	ret := map[string]Value{}

	value := reflect.ValueOf(p)
	typex := reflect.TypeOf(p)

	for i := 0; i < value.NumField(); i++ {
		recurse(ret, value.Field(i), typex.Field(i).Name)
	}

	return ret
}

func recurse(m map[string]Value, value reflect.Value, path ...string) {

	key := toSnake(path...)

	switch kind := value.Kind(); kind {
	case reflect.String:
		m[key] = Value{Text: value.String()}
	case reflect.Bool:
		m[key] = Value{Bool: value.Bool()}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		m[key] = Value{Number: value.Int()}
	case reflect.Float32, reflect.Float64:
		m[key] = Value{Float: round(value.Float())}
	case reflect.Slice:
		stringSlice := make([]string, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			str, ok := value.Index(i).Interface().(string)
			if ok {
				stringSlice = append(stringSlice, str)
			}
		}
		m[key] = Value{Array: &RepeatedValues{Strings: stringSlice}}
	case reflect.Map:

		if value.Type().String() == "map[string]string" {
			stringMap := map[string][]byte{}
			for _, v := range value.MapKeys() {
				stringMap[v.String()] = []byte(value.MapIndex(v).String())
			}
			m[key] = Value{Array: &RepeatedValues{KeyValue: stringMap}}
			break
		}

		if value.Type().String() == "map[string]chive.amount" {
			amountMap := map[string][]byte{}
			for _, v := range value.MapKeys() {
				a, ok := value.MapIndex(v).Interface().(amount)
				if ok {
					amountMap[v.String()] = []byte(fmt.Sprintf("currency:%s,units:%d", a.Currency, a.Units))
				}
			}
			m[key] = Value{Array: &RepeatedValues{KeyValue: amountMap}}
			break
		}

		fmt.Println("Unhandled map type: ", value.Type().String())

	case reflect.Struct:

		// time.Time
		if timeValue, ok := value.Interface().(time.Time); ok {

			m[key] = Value{Timestamp: timeValue.Add(time.Hour * 8).Unix()}
			break
		}

		// Other structs
		for i := 0; i < value.NumField(); i++ {
			recurse(m, value.Field(i), append(path, value.Type().Field(i).Name)...)
		}

	default:
		fmt.Println("Unhandled type: ", kind)
	}
}

var matchFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
var matchAllCap = regexp.MustCompile(`([a-z0-9])([A-Z])`)

func toSnake(strs ...string) string {
	var ret []string
	for _, str := range strs {
		if str == "" {
			continue
		}
		str = matchFirstCap.ReplaceAllString(str, "${1}_${2}")
		str = matchAllCap.ReplaceAllString(str, "${1}_${2}")
		ret = append(ret, strings.ToLower(str))
	}
	return strings.Join(ret, ".")
}

// Just rounding here to make the floating point match
// Otherwise 0.8 != 0.800000011920929
func round(f float64) float64 {
	return math.Round(f*10) / 10
}
