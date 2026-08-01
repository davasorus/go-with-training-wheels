package main

import (
	"fmt"
)

type Number interface {
	int64 | float64
}

func sumInt(m map[string]int64) int64 {
	var sum int64
	for _, v := range m {
		sum += v
	}
	return sum
}

func sumFloat(m map[string]float64) float64 {
	var sum float64
	for _, v := range m {
		sum += v
	}
	return sum
}

func sumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}

func sumNumbers[K comparable, V Number](m map[K]V) V {
	var sum V
	for _, v := range m {

		sum += v

	}

	return sum
}

func main() {

	ints := map[string]int64{
		"first":  10,
		"second": 20,
		"third":  30,
	}

	floats := map[string]float64{
		"first":  10.5,
		"second": 20.5,
		"third":  30.5,
	}

	fmt.Printf("Non-Generic Sums: %v and %v\n",
		sumInt(ints),
		sumFloat(floats))

	fmt.Printf("Generic Sums: %v and %v\n",
		sumIntsOrFloats[string, int64](ints),
		sumIntsOrFloats[string, float64](floats))

	fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n",
		sumIntsOrFloats(ints),
		sumIntsOrFloats(floats))

	fmt.Printf("Generic Sums with Constraint: %v and %v\n",
		sumNumbers(ints),
		sumNumbers(floats))
}
