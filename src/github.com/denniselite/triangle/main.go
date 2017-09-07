package main

import (
	"fmt"
	"sort"
)

const (
	Right = "Right"
	Obtuse = "Obtuse"
	Acute = "Acute"
	Equilateral = "Equilateral"
	Isosceles = "Isosceles"
	Scalene = "Scalene"
)

func main() {
	fmt.Println(TriangleFormula(3, 3, 3))
}

func TriangleFormula(a float64, b float64, c float64) string {
	if (a <= 0) || (b <= 0) || (c <= 0) {
		return "Wrong input params"
	}

	// Use go routines for type computing and channel for answers
	channel := make(chan byte)
	go defineByLengthsOfSides(a, b, c, channel)
	go defineByInternalAngles(a, b, c, channel)
	type1, type2 := <- channel, <- channel
	var result string
	// Build result after multi-thread calc
	if type2 < type1 {
		result = fmt.Sprintf("%s, %s", getTriangleType(type1), getTriangleType(type2))
	} else {
		result = fmt.Sprintf("%s, %s", getTriangleType(type2), getTriangleType(type1))
	}
	return result
}

// Definitions of triangle types
// because we use bytes in channels for optimization
func getTriangleType(b byte) string {
	switch b {
	case 1: return Right
	case 2: return Obtuse
	case 3: return Acute
	case 4: return Equilateral
	case 5: return Isosceles
	case 6: return Scalene
	default: return ""
	}
}

// Define triangle type by angles.
// bytes in channels for optimization
func defineByInternalAngles(a, b, c float64, channel chan byte) {
	arr := []float64{a, b, c}
	sort.Float64s(arr)

	aSquare := arr[2] * arr[2]
	bSquare := arr[1] * arr[1]
	cSquare := arr[0] * arr[0]

	switch true {
	case aSquare == bSquare + cSquare:
		channel <- byte(1)
	case aSquare > bSquare + cSquare:
		channel <- byte(2)
	default:
		channel <- byte(3)
	}
}

// Define triangle type by sides
// bytes in channels for optimization
func defineByLengthsOfSides(a, b, c float64, channel chan byte) {
	if (a == b) && (a == c) {
		channel <- byte(4)
		return
	}
	if (a == b) || (a == c) || (b == c) {
		channel <- byte(5)
	} else {
		channel <- byte(6)
	}
}
