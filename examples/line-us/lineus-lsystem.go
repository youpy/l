package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"../.."
	"github.com/youpy/go-line-us"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func write(turtle *lindenmayer.Turtle, client *lineus.Client) {
	var z float64

	if turtle.Pen {
		z = 0.0
	} else {
		z = 1000.0
	}

	res, err := client.LinearInterpolation(turtle.Pos.X, turtle.Pos.Y, z)
	checkError(err)

	fmt.Print(string(res.Message))

	time.Sleep(200 * time.Millisecond)
}

func main() {
	x := flag.Float64("x", 1000.0, "initial X position")
	y := flag.Float64("y", -1300.0, "initial Y position")
	ruleF := flag.String("ruleF", "F[+F-F-F]F[--F+F+F]", "rule F")
	initialState := flag.String("initial", "F", "initial state")
	iterations := flag.Int("iterations", 4, "number of iterations")
	side := flag.Float64("side", 40.0, "side")
	angle := flag.Float64("angle", 15.0, "angle")
	hostname := flag.String("hostname", "line-us.local:1337", "hostname of line-us machine")

	flag.Parse()

	client, err := lineus.NewClient(*hostname)
	checkError(err)

	turtle := lindenmayer.NewTurtle(*x, *y)
	done := make(chan struct{})

	write(turtle, client)

	go func() {
		for _ = range turtle.Dispatcher {
			write(turtle, client)
		}

		close(done)
	}()

	lindenmayer.LSystem(turtle, *ruleF, *initialState, *iterations, *side, *angle)

	<-done
}
