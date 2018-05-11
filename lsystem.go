package lindenmayer

import (
	"math"
	"strings"
)

func LSystem(turtle *Turtle, ruleF string, initialState string, iterations int, side float64, angle float64) {
	expanded := initialState

	for i := 0; i < iterations; i += 1 {
		expanded = strings.Replace(expanded, "F", ruleF, -1)
	}

	chars := strings.Split(expanded, "")
	stack := []*Turtle{}

	for _, char := range chars {
		switch char {
		case "F":
			turtle.Forward(side)
		case "+":
			turtle.Left(angle / 180 * math.Pi)
		case "-":
			turtle.Right(angle / 180 * math.Pi)
		case "[":
			stack = append(stack, turtle.Clone())
		case "]":
			popped := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			turtle.Restore(popped)
		}
	}

	close(turtle.Dispatcher)
}
